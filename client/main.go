package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hashicorp/go-hclog"
	"github.com/multiformats/go-multiaddr"
	"github.com/zivkovicmilos/peer_drop/config"
	"github.com/zivkovicmilos/peer_drop/rendezvous"
	"github.com/zivkovicmilos/peer_drop/rest"
	"github.com/zivkovicmilos/peer_drop/storage"
)

type ServiceHandler struct {
	logger       hclog.Logger
	closeChannel chan os.Signal

	serviceListeners map[string]chan struct{}
}

var serviceHandler *ServiceHandler

type RendezvousNodes []string

func (rn *RendezvousNodes) String() string {
	return formatArray(*rn)
}

func formatArray(array []string) string {
	return "(" + strings.Join(array, ",") + ")"
}

func (rn *RendezvousNodes) Set(value string) error {
	if _, err := multiaddr.NewMultiaddr(value); err != nil {
		return err
	}
	*rn = append(*rn, value)
	return nil
}

func main() {
	// ===== BASE SETUP ===== //

	// Set the flags
	hostPtr := flag.String("host", config.ServerHost,
		fmt.Sprintf("Host address of the client. Default %s", config.ServerHost),
	)
	httpPortPtr := flag.Int("http-port", config.ServerHTTPPort,
		fmt.Sprintf("HTTP port of the client. Default %d", config.ServerHTTPPort),
	)
	grpcPortPtr := flag.Int("grpc-port", config.ServerGRPCPort,
		fmt.Sprintf("GRPC port of the client. Defualt %d", config.ServerGRPCPort),
	)
	libp2pPortPtr := flag.Int("libp2p-port", config.ServerLibp2pPort,
		fmt.Sprintf("GRPC port of the client. Defualt %d", config.ServerLibp2pPort),
	)
	rendezvousMode := flag.Bool("rendezvous", false,
		fmt.Sprintf("server mode of the client. Default %t", false),
	)
	var rendezvousNodes RendezvousNodes
	flag.Var(&rendezvousNodes, "rendezvous-node",
		fmt.Sprintf("server mode of the client. Default %t", false),
	)

	flag.Parse()

	// Set up the base directory
	directoryError := createDirectory(fmt.Sprintf("%s/%s", config.DirectoryBase, config.DirectoryStorage))
	if directoryError != nil {
		os.Exit(1)
	}

	// Set up the storage
	storageHandler := storage.GetStorageHandler()
	storageHandler.SetLogger(serviceHandler.logger)
	storageHandler.SetCloseChannel(serviceHandler.registerCloseListener("storageHandler"))

	if storageErr := storage.GetStorageHandler().OpenDB(
		fmt.Sprintf("%s/%s/", config.DirectoryBase, config.DirectoryStorage),
	); storageErr != nil {
		os.Exit(1)
	}

	// Initialize the service handler
	serviceHandler = &ServiceHandler{
		serviceListeners: make(map[string]chan struct{}),
	}

	// Set up the logger
	serviceHandler.logger = hclog.New(&hclog.LoggerOptions{
		Name:  "peer_drop",
		Level: hclog.LevelFromString("DEBUG"),
	})

	// Set up the close mechanism
	serviceHandler.closeChannel = make(chan os.Signal, 1)
	signal.Notify(serviceHandler.closeChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the broadcast notifier
	go serviceHandler.broadcastNotifier()

	// ===== CLIENT SPECIFIC SETUP ===== //

	// Start up the appropriate client instance
	nodeConfig := &config.NodeConfig{
		HostAddress: *hostPtr,
		HttpPort:    *httpPortPtr,
		GrpcPort:    *grpcPortPtr,
		Libp2pPort:  *libp2pPortPtr,
	}

	if *rendezvousMode {
		setupAsRendezvous(nodeConfig, &config.RendezvousConfig{
			RendezvousNodes: rendezvousNodes,
		},
		)
	} else {
		setupAsClient(nodeConfig)
	}
}

// setupAsClient sets up the current peer_drop. node as a client
func setupAsClient(nodeConfig *config.NodeConfig) {
	// Set up the directories
	if directoryError := createDirectory(
		fmt.Sprintf("%s/%s", config.DirectoryBase, config.DirectoryFiles),
	); directoryError != nil {
		os.Exit(1)
	}

	// Set up the http dispatcher
	dispatcher := rest.NewDispatcher(
		serviceHandler.logger,
		nodeConfig,
	)

	dispatcher.Start(serviceHandler.registerCloseListener("dispatcher"))
}

// setupAsRendezvous sets up the current peer_drop. node as a rendezvous server
func setupAsRendezvous(nodeConfig *config.NodeConfig, rendezvousConfig *config.RendezvousConfig) {
	// Set up the directories
	if directoryError := createDirectory(
		fmt.Sprintf("%s/%s", config.DirectoryBase, config.DirectoryLibp2p),
	); directoryError != nil {
		os.Exit(1)
	}

	// Set up the rendezvous server
	rendezvousServer := rendezvous.NewRendezvousServer(
		serviceHandler.logger,
		nodeConfig,
		rendezvousConfig,
	)

	rendezvousServer.Start(serviceHandler.registerCloseListener("rendezvous"))
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
