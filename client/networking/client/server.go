package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"github.com/zivkovicmilos/peer_drop/config"
	localCrypto "github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/proto"
	"github.com/zivkovicmilos/peer_drop/storage"
	"google.golang.org/grpc"
)

type ClientServer struct {
	logger     hclog.Logger
	nodeConfig *config.NodeConfig

	// Networking metadata //
	me   peer.ID
	host host.Host

	// Context //
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewClientServer returns a new client networking instance
func NewClientServer(
	logger hclog.Logger,
	nodeConfig *config.NodeConfig,
) *ClientServer {
	return &ClientServer{
		logger:     logger.Named("client-node"),
		nodeConfig: nodeConfig,
	}
}

// Start starts the client networking instance
func (cs *ClientServer) Start(closeChannel chan struct{}) {
	// Set up the rendezvous nodes
	if setErr := cs.setupRendezvous(); setErr != nil {
		cs.logger.Error(fmt.Sprintf("Unable to set up rendezvous nodes, %v", setErr))

		os.Exit(1)
	}

	// Base libp2p setup
	libp2pKey, keyError := localCrypto.ReadLibp2pKey(
		filepath.Join(cs.nodeConfig.BaseDir, config.DirectoryLibp2p),
		"libp2p_key_client.asc",
	)
	if keyError != nil {
		os.Exit(1)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cs.ctx = ctx
	cs.cancelFunc = cancelFunc

	sourceMultiAddr, multiAddrErr := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ip4/%s/tcp/%d", cs.nodeConfig.HostAddress, cs.nodeConfig.Libp2pPort),
	)
	if multiAddrErr != nil {
		cs.logger.Error(fmt.Sprintf("Unable to create source multiaddr, %v", multiAddrErr))

		cancelFunc()
		os.Exit(1)
	}

	clientHost, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(libp2pKey),
	)
	if err != nil {
		cs.logger.Error(fmt.Sprintf("Unable to start libp2p service, %v", err))

		cancelFunc()
		os.Exit(1)
	}

	defer func() {
		cs.cancelFunc()
		_ = clientHost.Close()
	}()

	cs.host = clientHost

	// Set up the local DHT
	kademliaDHT, err := dht.New(ctx, clientHost)
	if err != nil {
		cs.logger.Error(fmt.Sprintf("Unable to start Kademlia DHT, %v", err))

		os.Exit(1)
	}

	// Bootstrap the DHT. This spawns a Background
	// thread that will refresh the peer table every five minutes
	cs.logger.Info("Bootstrapping the DHT...")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		cs.logger.Error(fmt.Sprintf("Unable to bootstrap the DHT, %v", err))

		cancelFunc()
		os.Exit(1)
	}

	// Dial the rendezvous nodes
	cs.dialRendezvous()

	<-closeChannel
}

// setupRendezvous sets the initial rendezvous nodes
func (cs *ClientServer) setupRendezvous() (storeErr error) {
	foundRendezvous, _ := storage.GetStorageHandler().GetRendezvousNodes()
	if len(foundRendezvous) < 1 {
		// No rendezvous nodes set, use the default ones
		storeErr = storage.GetStorageHandler().SetRendezvousNodes(config.DefaultRendezvousNodes)
	}

	return
}

// dialRendezvous dials the set rendezvous nodes
// to enable peer discovery
func (cs *ClientServer) dialRendezvous() {
	cs.logger.Info("Attempting connection to rendezvous nodes...")

	rendezvousNodes, _ := storage.GetStorageHandler().GetRendezvousNodes()

	var wg sync.WaitGroup
	for _, peerAddr := range rendezvousNodes {
		mAddr, mAddrErr := multiaddr.NewMultiaddr(peerAddr)
		if mAddrErr != nil {
			cs.logger.Error(fmt.Sprintf("Unable to create multiaddr from value, %v", mAddrErr))
		}

		peerinfo, _ := peer.AddrInfoFromP2pAddr(mAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := cs.host.Connect(cs.ctx, *peerinfo); err != nil {
				cs.logger.Error(fmt.Sprintf("Unable to connect to rendezvous node (%s), %v", *peerinfo, err))
			} else {
				cs.logger.Info(fmt.Sprintf("Successfully connected to rendezvous node %s", *peerinfo))

				cs.initiateStreamToRendezvous(peerinfo.ID)
			}
		}()
	}
	wg.Wait()
}

// initiateStreamToRendezvous creates a communication stream between the rendezvous node
// and the current node
func (cs *ClientServer) initiateStreamToRendezvous(peerID peer.ID) {
	stream, err := cs.host.NewStream(cs.ctx, peerID, protocol.ID(config.WorkspaceInfoProto))
	if err != nil {
		cs.logger.Error(fmt.Sprintf("Unable to instantiate stream to rendezvous node, %v", err))
	}

	clientConn := WrapStreamInClient(stream)

	// TODO persist this in memory
	_ = proto.NewWorkspaceInfoServiceClient(clientConn.(*grpc.ClientConn))

	//workspaceInfo, err := clientProto.GetWorkspaceInfo(
	//	context.Background(),
	//	&proto.WorkspaceInfoRequest{Mnemonic: "dummy"},
	//)
	//if err != nil {
	//	cs.logger.Error(fmt.Sprintf("Unable to get workspace info, %v", workspaceInfo))
	//} else {
	//	cs.logger.Info(fmt.Sprintf("Received workspace info with mnemonic: %s", workspaceInfo.Mnemonic))
	//}
}
