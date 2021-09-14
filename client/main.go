package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/hashicorp/go-hclog"
	"github.com/zivkovicmilos/peer_drop/rest"
	"github.com/zivkovicmilos/peer_drop/storage"
)

func main() {
	// Create the logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "peer_drop",
		Level: hclog.LevelFromString("DEBUG"),
	})

	// Set up the storage
	storageHandler := storage.GetStorageHandler()
	storageHandler.SetLogger(logger)
	storageErr := storageHandler.OpenDB(filepath.Join("/", "peerDropStorage"))
	if storageErr != nil {
		os.Exit(1)
	}

	// Set up the close mechanism
	closeChannel := make(chan os.Signal, 1)
	signal.Notify(closeChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Set up the http dispatcher
	dispatcher := rest.NewDispatcher(logger, storageHandler)
	dispatcher.Start(closeChannel)
}
