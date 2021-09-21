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
	"github.com/zivkovicmilos/peer_drop/networking/client"
	"github.com/zivkovicmilos/peer_drop/rendezvous"
	"github.com/zivkovicmilos/peer_drop/rest/dispatcher"
	servicehandler "github.com/zivkovicmilos/peer_drop/service-handler"
	"github.com/zivkovicmilos/peer_drop/storage"
	"github.com/zivkovicmilos/peer_drop/utils"
)

type RendezvousNodes []string

func (rn *RendezvousNodes) String() string {
	return formatArray(*rn)
}

func formatArray(array RendezvousNodes) string {
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
	baseDirPtr := flag.String("directory", config.DirectoryBase,
		fmt.Sprintf("Base directory for the client data. Default %s", config.DirectoryBase),
	)
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
	directoryError := utils.CreateDirectory(fmt.Sprintf("%s/%s", *baseDirPtr, config.DirectoryStorage))
	if directoryError != nil {
		os.Exit(1)
	}

	// Set up the logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "peer_drop",
		Level: hclog.LevelFromString("DEBUG"),
	})

	// Set up the storage
	storageHandler := storage.GetStorageHandler()
	storageHandler.SetLogger(logger)
	storageHandler.SetCloseChannel(
		servicehandler.GetServiceHandler().RegisterCloseListener("storageHandler"),
	)

	// Open DB
	if storageErr := storage.GetStorageHandler().OpenDB(
		fmt.Sprintf("%s/%s/", *baseDirPtr, config.DirectoryStorage),
	); storageErr != nil {
		os.Exit(1)
	}

	// Set up the close mechanism
	closeChannel := make(chan os.Signal, 1)
	signal.Notify(closeChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	servicehandler.GetServiceHandler().SetCloseChannel(closeChannel)

	// Start the broadcast notifier
	go servicehandler.GetServiceHandler().BroadcastNotifier()

	// ===== CLIENT / RENDEZVOUS SPECIFIC SETUP ===== //

	// Start up the appropriate client instance
	nodeConfig := &config.NodeConfig{
		HostAddress: *hostPtr,
		HttpPort:    *httpPortPtr,
		GrpcPort:    *grpcPortPtr,
		Libp2pPort:  *libp2pPortPtr,
		BaseDir:     *baseDirPtr,
	}

	if *rendezvousMode {
		setupAsRendezvous(
			logger,
			nodeConfig,
			&config.RendezvousConfig{
				RendezvousNodes: rendezvousNodes,
			},
		)
	} else {
		setupAsClient(logger, nodeConfig)
	}
}

// setupAsClient sets up the current peer_drop. node as a client
func setupAsClient(logger hclog.Logger, nodeConfig *config.NodeConfig) {
	// Set up the directories
	if directoryError := utils.CreateDirectory(
		fmt.Sprintf("%s/%s", nodeConfig.BaseDir, config.DirectoryFiles),
	); directoryError != nil {
		os.Exit(1)
	}

	if directoryError := utils.CreateDirectory(
		fmt.Sprintf("%s/%s", nodeConfig.BaseDir, config.DirectoryLibp2p),
	); directoryError != nil {
		os.Exit(1)
	}

	// Set up the http dispatcher
	createdDispatcher := dispatcher.NewDispatcher(
		logger,
		nodeConfig,
	)

	go createdDispatcher.Start(servicehandler.GetServiceHandler().RegisterCloseListener("dispatcher"))

	// Set up the networking layer
	clientServer := client.NewClientServer(logger, nodeConfig)
	servicehandler.GetServiceHandler().SetClientServer(clientServer)

	clientServer.Start(servicehandler.GetServiceHandler().RegisterCloseListener("client-server")) // TODO start as goroutine
}

// setupAsRendezvous sets up the current peer_drop. node as a rendezvous server
func setupAsRendezvous(
	logger hclog.Logger,
	nodeConfig *config.NodeConfig,
	rendezvousConfig *config.RendezvousConfig,
) {
	// Set up the directories
	if directoryError := utils.CreateDirectory(
		fmt.Sprintf("%s/%s", nodeConfig.BaseDir, config.DirectoryLibp2p),
	); directoryError != nil {
		os.Exit(1)
	}

	// Set up the rendezvous server
	rendezvousServer := rendezvous.NewRendezvousServer(
		logger,
		nodeConfig,
		rendezvousConfig,
	)

	rendezvousServer.Start(servicehandler.GetServiceHandler().RegisterCloseListener("rendezvous"))
}
