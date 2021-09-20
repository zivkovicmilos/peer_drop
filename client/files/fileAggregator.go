package files

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zivkovicmilos/peer_drop/proto"
)

type FileListWrapper struct {
	FileList *proto.FileList
	PeerID   peer.ID
}

// FileAggregator aggregates different file lists for workspaces
type FileAggregator struct {
	logger hclog.Logger

	updateChannel chan FileListWrapper // Update channel that's filled by clientServer
	fileMap       map[string][]peer.ID // Map indicating which peers have a certain file (checksum -> []peerID)
	peerFileArray map[peer.ID][]string // Map indicating all files that the peer is offering (peerID -> []checksum)
	fileArray     []*proto.File        // All files available to the client in the workspace

	fileMapMuxMap map[string]sync.RWMutex
	fileArrayMux  sync.RWMutex
}

// NewFileAggregator creates a new instance of the file aggregator
func NewFileAggregator(
	logger hclog.Logger,
	workspaceName string,
	updateChannel chan FileListWrapper,
) *FileAggregator {
	return &FileAggregator{
		logger:        logger.Named(fmt.Sprintf("file-aggregator [%s]", workspaceName)),
		fileMap:       make(map[string][]peer.ID),
		peerFileArray: make(map[peer.ID][]string),
		fileArray:     make([]*proto.File, 0),
		updateChannel: updateChannel,
	}
}

// Start starts the File aggregator loop
func (fa *FileAggregator) Start() {
	go fa.aggregateFilesLoop()
}

// Stop stops the file aggregator service
func (fa *FileAggregator) Stop() {
	close(fa.updateChannel)
}

// aggregateFilesLoop listens for new file list events
func (fa *FileAggregator) aggregateFilesLoop() {
	for {
		fileListWrapper, more := <-fa.updateChannel
		if more {
			// Parse the file list

			// Find the differences
			fa.fileArrayMux.RLock()
			fileDifferences := fa.findFileDifference(fa.fileArray, fileListWrapper.FileList.FileList)
			fa.fileArrayMux.RUnlock()

			// Update all the relevant structures
			fa.pruneFileMap(fileDifferences, fileListWrapper.PeerID)
		} else {
			// exit signal caught
			return
		}
	}
}

// findFileDifference finds which files are different between the arrays
func (fa *FileAggregator) findFileDifference(a, b []*proto.File) []*proto.File {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x.FileChecksum] = struct{}{}
	}

	var diff []*proto.File
	for _, x := range a {
		if _, found := mb[x.FileChecksum]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// pruneFileMap updates file structures based on the file differences
func (fa *FileAggregator) pruneFileMap(fileDifferences []*proto.File, peerID peer.ID) {
	for _, file := range fileDifferences {
		mux, _ := fa.fileMapMuxMap[file.FileChecksum]
		mux.Lock()
		// Check if the file mapping to peerID exists.
		peerArray, ok := fa.fileMap[file.FileChecksum]
		if !ok {
			// If it doesn't exist, that means the file was added
			fa.addFileToFileArray(file)
			newArray := make([]peer.ID, 0)
			newArray = append(newArray, peerID)

			fa.fileMap[file.FileChecksum] = newArray
		} else {
			// If it exists, that means the file was removed
			newArray := fa.pruneFromPeerArray(peerArray, peerID)
			fa.fileMap[file.FileChecksum] = newArray
			if len(newArray) == 0 {
				// Remove the file from the global array as nobody serves it
				fa.pruneFileFromFileArray(file)
			}

		}

		mux.Unlock()
	}
}

// addFileToFileArray adds a new file to the global file array. [Thread safe]
func (fa *FileAggregator) addFileToFileArray(file *proto.File) {
	fa.fileArrayMux.Lock()
	defer fa.fileArrayMux.Unlock()

	fa.fileArray = append(fa.fileArray, file)
}

// pruneFileFromFileArray removes a specific file from the global file array/ [Thread safe]
func (fa *FileAggregator) pruneFileFromFileArray(file *proto.File) {
	fa.fileArrayMux.Lock()
	defer fa.fileArrayMux.Unlock()

	index := -1
	for searchIndex, searchFile := range fa.fileArray {
		if searchFile.FileChecksum == file.FileChecksum {
			index = searchIndex
			break
		}
	}

	if index >= 0 {
		fa.fileArray = append(fa.fileArray[:index], fa.fileArray[index+1:]...)
	}
}

// pruneFromPeerArray removes a peer ID from the peer array
func (fa *FileAggregator) pruneFromPeerArray(peerArray []peer.ID, peerID peer.ID) []peer.ID {
	index := -1
	for searchIndex, searchPeerID := range peerArray {
		if searchPeerID == peerID {
			index = searchIndex
			break
		}
	}

	if index >= 0 {
		return append(peerArray[:index], peerArray[index+1:]...)
	}

	return peerArray
}

// GetFileList returns the available file list
func (fa *FileAggregator) GetFileList() []*proto.File {
	fa.fileArrayMux.RLock()
	defer fa.fileArrayMux.RUnlock()

	return fa.fileArray
}
