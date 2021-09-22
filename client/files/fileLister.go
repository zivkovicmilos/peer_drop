package files

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/zivkovicmilos/peer_drop/proto"
	"go.uber.org/atomic"
)

// FileLister gathers all files in the workspace
// directory, and returns them as available for sharing
type FileLister struct {
	logger  hclog.Logger
	baseDir string

	fileMap    map[string]*proto.File // checksum -> file
	fileMapMux sync.RWMutex

	sweepInterval   time.Duration
	sweepInProgress atomic.Bool

	stopChannel chan struct{}
}

// NewFileLister creates a new instance of the file lister
func NewFileLister(
	logger hclog.Logger,
	baseDir string,
	sweepInterval time.Duration,
) *FileLister {
	return &FileLister{
		logger:        logger.Named(fmt.Sprintf("file-lister [%s]", baseDir)),
		baseDir:       baseDir,
		fileMap:       make(map[string]*proto.File),
		sweepInterval: sweepInterval,
		stopChannel:   make(chan struct{}),
	}
}

// Start starts the file lister
func (fl *FileLister) Start() {
	go fl.sweepDirectoryLoop()
}

// GetBaseDir gets the base directory
func (fl *FileLister) GetBaseDir() string {
	return fl.baseDir
}

func (fl *FileLister) pruneRemovedFiles(currentFiles []*proto.File) {
	fl.fileMapMux.RLock()

	fileMap := make(map[string]int)

	for _, file := range fl.fileMap {
		fileMap[file.FileChecksum]++
	}

	fl.fileMapMux.RUnlock()

	// Go through all current files
	for _, currentFile := range currentFiles {
		fileMap[currentFile.FileChecksum]--
	}

	// Prune removed files

	for prunedChecksum, pruned := range fileMap {
		if pruned > 0 {
			fl.removeFile(prunedChecksum)
		}
	}
}

// sweepDirectory goes over all the files in the sharing directory
// and updates the file map
func (fl *FileLister) sweepDirectory() {
	if fl.sweepInProgress.Load() {
		// Sweep already in progress
		return
	}
	fl.sweepInProgress.Store(true)
	fl.logger.Info("Directory sweep started")
	var filesFound atomic.Int64

	var wg sync.WaitGroup

	defer func() {
		fl.logger.Info(fmt.Sprintf("Directory sweep finished with %d files", filesFound.Load()))
		fl.sweepInProgress.Store(false)
	}()

	// Sweep directory for files
	directoryFiles, err := ioutil.ReadDir(fl.baseDir)
	if err != nil {
		fl.logger.Error(fmt.Sprintf("Unable to read directory, %v", err))
		return
	}

	currentFiles := make([]*proto.File, 0)
	// Construct the directory files
	for _, f := range directoryFiles {
		if !f.IsDir() {
			wg.Add(1)

			filePath := filepath.Join(fl.baseDir, f.Name())
			checksum, checksumErr := fl.checksumFile(filePath)
			if checksumErr != nil {
				fl.logger.Error(fmt.Sprintf("Unable to checksum file %s", filePath))
			}

			protoFile := fileInfoToFileProto(f)
			protoFile.FileChecksum = checksum

			fl.addFile(protoFile, checksum)
			currentFiles = append(currentFiles, protoFile)
			filesFound.Inc()

			wg.Done()
		}
	}

	fl.pruneRemovedFiles(currentFiles)

	wg.Wait()
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// fileInfoToFileProto converts the regular file information into proto format
func fileInfoToFileProto(file fs.FileInfo) *proto.File {
	protoFile := &proto.File{}

	protoFile.Name = fileNameWithoutExtension(file.Name())
	protoFile.Extension = filepath.Ext(file.Name())

	protoFile.Size = file.Size()
	protoFile.DateModified = file.ModTime().Unix()

	return protoFile
}

// Stop stops the file lister service
func (fl *FileLister) Stop() {
	fl.stopChannel <- struct{}{}
}

// sweepDirectoryLoop is triggered periodically to update the file map
func (fl *FileLister) sweepDirectoryLoop() {
	// Do an initial sweep
	fl.sweepDirectory()

	// Start the loop
	ticker := time.NewTicker(fl.sweepInterval)

	go func() {
		<-fl.stopChannel
		ticker.Stop()

		return
	}()
	for {
		select {
		case _ = <-ticker.C:
			go fl.sweepDirectory()
		}
	}
}

// checksumFile checksums the file at the given path
// Generates a SHA256 hash of the file
func (fl *FileLister) checksumFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// GetFileInfo returns the file information
func (fl *FileLister) GetFileInfo(checksum string) (*proto.File, error) {
	fl.fileMapMux.RLock()
	defer fl.fileMapMux.RUnlock()

	file, _ := fl.fileMap[checksum]

	return file, nil
}

// addFile adds a file to the file map
func (fl *FileLister) addFile(file *proto.File, checksum string) {
	fl.fileMapMux.Lock()
	defer fl.fileMapMux.Unlock()

	fl.fileMap[checksum] = file
}

// removeFile removes a file from the file map
func (fl *FileLister) removeFile(checksum string) {
	fl.fileMapMux.Lock()
	defer fl.fileMapMux.Unlock()

	delete(fl.fileMap, checksum)
}

// GetAvailableFiles returns the available files for sharing in the workspace
func (fl *FileLister) GetAvailableFiles() []*proto.File {
	fl.fileMapMux.RLock()
	defer fl.fileMapMux.RUnlock()

	fileList := make([]*proto.File, 0)

	for _, file := range fl.fileMap {
		fileList = append(fileList, file)
	}

	return fileList
}
