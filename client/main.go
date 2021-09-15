package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-hclog"
	"github.com/zivkovicmilos/peer_drop/rest"
	"github.com/zivkovicmilos/peer_drop/storage"
)

var (
	directoryFiles   = "files"
	directoryStorage = "storage"
	directoryBase    = "app_data"
)

type ServiceHandler struct {
	logger       hclog.Logger
	closeChannel chan os.Signal

	serviceListeners map[string]chan struct{}
}

func main() {
	// Initialize the service handler
	serviceHandler := &ServiceHandler{
		serviceListeners: make(map[string]chan struct{}),
	}

	// Create the logger
	serviceHandler.logger = hclog.New(&hclog.LoggerOptions{
		Name:  "peer_drop",
		Level: hclog.LevelFromString("DEBUG"),
	})

	// Set up the close mechanism
	serviceHandler.closeChannel = make(chan os.Signal, 1)
	signal.Notify(serviceHandler.closeChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the broadcast notifier
	go serviceHandler.broadcastNotifier()

	// Set up the directories
	directoryError := setupDirectories()
	if directoryError != nil {
		os.Exit(1)
	}

	// Set up the storage
	storageHandler := storage.GetStorageHandler()
	storageHandler.SetLogger(serviceHandler.logger)
	storageHandler.SetCloseChannel(serviceHandler.registerCloseListener("storageHandler"))

	storageErr := storageHandler.OpenDB(
		fmt.Sprintf("%s/%s/", directoryBase, directoryStorage),
	)
	if storageErr != nil {
		os.Exit(1)
	}

	// Set up the http dispatcher
	dispatcher := rest.NewDispatcher(serviceHandler.logger, storageHandler)

	dispatcher.Start(serviceHandler.registerCloseListener("dispatcher"))
}

// createDirectory creates a single directory
func createDirectory(path string) error {
	err := os.MkdirAll(path, os.ModeDir)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

// setupDirectories sets up the base peer_drop folder structure
func setupDirectories() error {
	dirError := createDirectory(fmt.Sprintf("%s/%s", directoryBase, directoryFiles))
	if dirError != nil {
		return dirError
	}

	return createDirectory(fmt.Sprintf("%s/%s", directoryBase, directoryStorage))
}

// broadcastNotifier waits for a term signal and alerts all listening services
func (sh *ServiceHandler) broadcastNotifier() {
	<-sh.closeChannel

	// Alert all registered services
	for _, serviceChannel := range sh.serviceListeners {
		serviceChannel <- struct{}{}
	}
}

// registerCloseListener adds a new listener for Service Handler's close events
func (sh *ServiceHandler) registerCloseListener(serviceName string) chan struct{} {
	channel := make(chan struct{}, 1)
	sh.serviceListeners[serviceName] = channel

	return sh.serviceListeners[serviceName]
}
