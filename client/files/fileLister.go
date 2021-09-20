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

	sweepInterval time.Duration

	sweepInProgress atomic.Bool
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
	}
}

// Start starts the file lister
func (fl *FileLister) Start() {
	go fl.sweepDirectoryLoop()
}

// sweepDirectory goes over all the files in the sharing directory
// and updates the file map
func (fl *FileLister) sweepDirectory() {
	if fl.sweepInProgress.Load() {
		// Sweep already in progress
		return
	}
	fl.sweepInProgress.Store(true)

	var wg sync.WaitGroup

	defer fl.sweepInProgress.Store(false)

	// Sweep directory for files
	directoryFiles, err := ioutil.ReadDir("./")
	if err != nil {
		fl.logger.Error(fmt.Sprintf("Unable to read directory, %v", err))
		return
	}

	// Construct the directory files
	for _, f := range directoryFiles {
		if !f.IsDir() {
			wg.Add(1)

			go func() {
				defer wg.Done()

				filePath := filepath.Join(fl.baseDir, f.Name())
				checksum, checksumErr := fl.checksumFile(filePath)
				if checksumErr != nil {
					fl.logger.Error(fmt.Sprintf("Unable to checksum file %s", filePath))
				}

				protoFile := fileInfoToFileProto(f)
				protoFile.FileChecksum = checksum

				go fl.addFile(protoFile, checksum)
			}()
		}
	}

	wg.Wait()
}

// fileInfoToFileProto converts the regular file information into proto format
func fileInfoToFileProto(file fs.FileInfo) *proto.File {
	var protoFile *proto.File

	stringArr := strings.Split(file.Name(), ".")

	protoFile.Name = stringArr[0]
	protoFile.Extension = stringArr[1]

	protoFile.Size = file.Size()
	protoFile.DateModified = file.ModTime().Unix()

	return protoFile
}

// sweepDirectoryLoop is triggered periodically to update the file map
func (fl *FileLister) sweepDirectoryLoop() {
	//sweepContext := context.Background()
	ticker := time.NewTicker(fl.sweepInterval)

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
