package rendezvous

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"github.com/zivkovicmilos/peer_drop/config"
	localCrypto "github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/proto"
	"github.com/zivkovicmilos/peer_drop/storage"
)

var (
	rendezvousTopic = "workspace-info"
)

type RendezvousServer struct {
	logger           hclog.Logger
	nodeConfig       *config.NodeConfig
	rendezvousConfig *config.RendezvousConfig

	// Networking metadata //
	me peer.ID

	// Message handling //
	workspaceInfoMsgQueue chan *proto.WorkspaceInfo

	// Pubsub //
	pubSub             *pubsub.PubSub       // Reference to the main pubsub instance
	pubSubTopic        *pubsub.Topic        // Reference to the workspace info sharing topic
	pubSubSubscription *pubsub.Subscription // Reference to the workspace info sharing subscription

	// Context //
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewRendezvousServer returns a new instance of the rendezvous server
func NewRendezvousServer(
	logger hclog.Logger,
	nodeConfig *config.NodeConfig,
	rendezvousConfig *config.RendezvousConfig,
) *RendezvousServer {
	return &RendezvousServer{
		logger:                logger.Named("rendezvous"),
		nodeConfig:            nodeConfig,
		rendezvousConfig:      rendezvousConfig,
		workspaceInfoMsgQueue: make(chan *proto.WorkspaceInfo),
	}
}

// Start starts the rendezvous server
func (r *RendezvousServer) Start(closeChannel chan struct{}) {
	libp2pKey, keyError := localCrypto.ReadLibp2pKey(
		filepath.Join(config.DirectoryBase, config.DirectoryLibp2p),
		"libp2p_key.asc",
	)
	if keyError != nil {
		os.Exit(1)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	r.ctx = ctx
	r.cancelFunc = cancelFunc

	defer r.cancelFunc()

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ip4/%s/tcp/%d", r.nodeConfig.HostAddress, r.nodeConfig.Libp2pPort),
	)

	rendezvousHost, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(libp2pKey),
	)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Unable to start libp2p service, %v", err))

		cancelFunc()
		os.Exit(1)
	}

	_, err = dht.New(ctx, rendezvousHost)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Unable to start DHT service, %v", err))

		cancelFunc()
		os.Exit(1)
	}

	r.logger.Info(
		fmt.Sprintf("Rendezvous node started with ID: /ip4/%s/tcp/%v/p2p/%s",
			r.nodeConfig.HostAddress, r.nodeConfig.Libp2pPort, rendezvousHost.ID().Pretty(),
		),
	)
	r.me = rendezvousHost.ID()

	// TODO set stream handlers

	// Connect to other rendezvous nodes
	r.connectToRendezvousPeers(rendezvousHost, ctx)

	// Set up the pubsub service
	pubsubError := r.setupPubsub(rendezvousHost, ctx)
	if pubsubError != nil {
		r.logger.Error(fmt.Sprintf("Unable to setup the pubsub instance, %v", pubsubError))

		cancelFunc()
		os.Exit(1)
	}

	// Wait for a close signal
	<-closeChannel
}

// connectToRendezvousPeers attempts to connect to other rendezvous nodes
func (r *RendezvousServer) connectToRendezvousPeers(host host.Host, ctx context.Context) {
	r.logger.Info("Attempting connection to other rendezvous nodes...")

	var wg sync.WaitGroup
	for _, peerAddr := range r.rendezvousConfig.RendezvousNodes {
		mAddr, mAddrErr := multiaddr.NewMultiaddr(peerAddr)
		if mAddrErr != nil {
			r.logger.Error(fmt.Sprintf("Unable to create multiaddr from value, %v", mAddrErr))
		}

		peerinfo, _ := peer.AddrInfoFromP2pAddr(mAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if connectErr := host.Connect(ctx, *peerinfo); connectErr != nil {
				r.logger.Error(fmt.Sprintf("Unable to connect to peer: %s", peerinfo.String()))
			} else {
				r.logger.Info(fmt.Sprintf("Connection extablished to rendezvous node: %s", peerinfo.String()))
			}
		}()
	}
	wg.Wait()
}

// setupPubsub sets up the pubsub service
func (r *RendezvousServer) setupPubsub(host host.Host, ctx context.Context) error {
	pubSub, gossipErr := pubsub.NewGossipSub(ctx, host)
	if gossipErr != nil {
		return gossipErr
	}
	r.pubSub = pubSub

	// Join to the workspace info sharing topic so the node
	// can publish new messages
	pubSubTopic, err := r.pubSub.Join(rendezvousTopic)
	if err != nil {
		return err
	}

	// Subscribe to the workspace info sharing topic so the node
	// can get new messages
	pubSubSubscription, err := pubSubTopic.Subscribe()
	if err != nil {
		return err
	}

	r.pubSubTopic = pubSubTopic
	r.pubSubSubscription = pubSubSubscription

	go r.readPubsubLoop()

	return nil
}

// readPubsubLoop runs the pubsub service loop.
// It grabs new messages as they arrive, and updates the node storage
func (r *RendezvousServer) readPubsubLoop() {
	go r.storageUpdateListener()

	for {
		workspaceInfoMsg, err := r.pubSubSubscription.Next(r.ctx)
		if err != nil {
			close(r.workspaceInfoMsgQueue)
			return
		}

		// Forward messages that are not from us
		if workspaceInfoMsg.ReceivedFrom == r.me {
			continue
		}

		workspaceInfo := new(proto.WorkspaceInfo)
		err = json.Unmarshal(workspaceInfoMsg.Data, workspaceInfo)
		if err != nil {
			continue
		}

		// Send the workspace info to be stored
		r.workspaceInfoMsgQueue <- workspaceInfo
	}
}

// storageUpdateListener listens for new workspace info messages
// and updates the storage
func (r *RendezvousServer) storageUpdateListener() {
	for {
		workspaceInfo, more := <-r.workspaceInfoMsgQueue
		if more {
			r.logger.Info("Storage update listener new message")

			storeErr := storage.GetStorageHandler().CreateWorkspaceInfo(workspaceInfo)
			if storeErr != nil {
				r.logger.Error(fmt.Sprintf("Unable to store workspace info, %v", storeErr))
				continue
			}

			r.logger.Info("Successfully stored workspace info")
		} else {
			// Close message received
			r.logger.Info("Storage update listener received stop signal...")
		}
	}
}
