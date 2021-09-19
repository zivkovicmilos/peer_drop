package rendezvous

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"github.com/zivkovicmilos/peer_drop/config"
	localCrypto "github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/mnemonic"
	"github.com/zivkovicmilos/peer_drop/proto"
	"github.com/zivkovicmilos/peer_drop/storage"

	localGRPC "github.com/zivkovicmilos/peer_drop/networking/client"
)

var (
	rendezvousTopic = "workspace-info"
)

type RendezvousServer struct {
	logger           hclog.Logger
	nodeConfig       *config.NodeConfig
	rendezvousConfig *config.RendezvousConfig

	// Networking metadata //
	me   peer.ID
	host host.Host

	// Message handling //
	workspaceInfoMsgQueue chan *proto.WorkspaceInfo

	// Pubsub //
	pubSub             *pubsub.PubSub       // Reference to the main pubsub instance
	pubSubTopic        *pubsub.Topic        // Reference to the workspace info sharing topic
	pubSubSubscription *pubsub.Subscription // Reference to the workspace info sharing subscription

	// Context //
	ctx        context.Context
	cancelFunc context.CancelFunc

	// GRPC //
	proto.UnimplementedWorkspaceInfoServiceServer
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
		filepath.Join(r.nodeConfig.BaseDir, config.DirectoryLibp2p),
		"libp2p_key_rendezvous.asc",
	)
	if keyError != nil {
		os.Exit(1)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	r.ctx = ctx
	r.cancelFunc = cancelFunc

	sourceMultiAddr, multiAddrErr := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ip4/%s/tcp/%d", r.nodeConfig.HostAddress, r.nodeConfig.Libp2pPort),
	)
	if multiAddrErr != nil {
		r.logger.Error(fmt.Sprintf("Unable to create source multiaddr, %v", multiAddrErr))

		cancelFunc()
		os.Exit(1)
	}

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

	defer func() {
		r.cancelFunc()
		_ = rendezvousHost.Close()
	}()

	r.host = rendezvousHost

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

	// Connect to other rendezvous nodes
	r.connectToRendezvousPeers()

	// Set up the pubsub service
	pubsubError := r.setupPubsub()
	if pubsubError != nil {
		r.logger.Error(fmt.Sprintf("Unable to setup the pubsub instance, %v", pubsubError))

		cancelFunc()
		os.Exit(1)
	}

	// Set up the GRPC protocol handlers
	r.setupGRPCProtocols()

	// Wait for a close signal
	<-closeChannel
}

func (r *RendezvousServer) sendDummyMessage() {
	m := &proto.WorkspaceInfo{
		Mnemonic:                 "mymnemonic",
		WorkspaceOwnerPublicKeys: []string{"123", "456"},
		SecurityType:             "password",
		SecuritySettings:         &proto.WorkspaceInfo_PasswordHash{PasswordHash: "passwordHash"},
	}

	marshaler := jsonpb.Marshaler{}
	buf := new(bytes.Buffer)
	err := marshaler.Marshal(buf, m)
	if err != nil {
		r.logger.Info("INVALID MARSHAL")
	}

	sendErr := r.pubSubTopic.Publish(r.ctx, buf.Bytes())
	if sendErr != nil {
		r.logger.Info("INVALID SEND")
	}

	r.logger.Info("MESSAGE SENT")
}

// connectToRendezvousPeers attempts to connect to other rendezvous nodes
func (r *RendezvousServer) connectToRendezvousPeers() {
	r.logger.Info("Attempting connection to other rendezvous nodes...")
	addedPeers := 0

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
			if connectErr := r.host.Connect(r.ctx, *peerinfo); connectErr != nil {
				r.logger.Error(fmt.Sprintf("Unable to connect to peer: %s", peerinfo.String()))
			} else {
				r.logger.Info(fmt.Sprintf("Connection established to rendezvous node: %s", peerinfo.String()))
				addedPeers = addedPeers + 1
			}
		}()
	}
	wg.Wait()

	if addedPeers > 0 {
		r.logger.Info(fmt.Sprintf("Connection established to %d rendezvous nodes", addedPeers))
	} else {
		r.logger.Info("No rendezvous peers found")
	}
}

// setupPubsub sets up the pubsub service
func (r *RendezvousServer) setupPubsub() error {
	pubSub, gossipErr := pubsub.NewGossipSub(r.ctx, r.host)
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

	// TODO add better closing management with channels

	for {
		workspaceInfoMsg, err := r.pubSubSubscription.Next(r.ctx)
		if err != nil {
			r.logger.Error(fmt.Sprintf("Unable to parse message, %v", err))

			close(r.workspaceInfoMsgQueue)
			return
		}

		// Forward messages that are not from us
		if workspaceInfoMsg.ReceivedFrom == r.me {
			continue
		}

		workspaceInfo := new(proto.WorkspaceInfo)
		err = jsonpb.Unmarshal(bytes.NewReader(workspaceInfoMsg.Data), workspaceInfo)
		if err != nil {
			r.logger.Error(fmt.Sprintf("Unmarshal error %v", err))
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

// GRPC //

// setupGRPCProtocols sets up the supported handlers for GRPC protocols
func (r *RendezvousServer) setupGRPCProtocols() {
	r.logger.Info("Setting up GRPC protocols...")

	workspaceInfoProtocol := localGRPC.NewGRPCProtocol()
	proto.RegisterWorkspaceInfoServiceServer(workspaceInfoProtocol.GrpcServer(), r)
	workspaceInfoProtocol.Serve()

	r.host.SetStreamHandler(protocol.ID(config.WorkspaceInfoProto), func(stream network.Stream) {
		peerID := stream.Conn().RemotePeer()
		r.logger.Info(fmt.Sprintf("Open stream [protocol: %s], by peer %s", config.WorkspaceInfoProto, peerID))

		workspaceInfoProtocol.Handler()(stream)
	})

	r.logger.Info("GRPC protocols set")
}

// GetWorkspaceInfo returns the workspace info found in a local DB
func (r *RendezvousServer) GetWorkspaceInfo(
	context context.Context,
	request *proto.WorkspaceInfoRequest,
) (*proto.WorkspaceInfo, error) {
	foundWorkspaceInfo, findErr := storage.GetStorageHandler().GetWorkspaceInfo(request.Mnemonic)
	if findErr != nil {
		return nil, fmt.Errorf("err") // TODO conform grpc
	}

	return foundWorkspaceInfo, nil
}

// CreateNewWorkspace generates mnemonic for the workspace and gossips it
func (r *RendezvousServer) CreateNewWorkspace(
	context context.Context,
	workspaceInfo *proto.WorkspaceInfo,
) (*proto.WorkspaceInfo, error) {
	mg := mnemonic.MnemonicGenerator{NumWords: 6}

	generatedMnemonic, generateErr := mg.GenerateMnemonic()
	if generateErr != nil {
		return workspaceInfo, generateErr
	}

	workspaceInfo.Mnemonic = generatedMnemonic
	r.handleNewWorkspace(workspaceInfo) // TODO make this a go routine?

	return workspaceInfo, nil
}

func (r *RendezvousServer) handleNewWorkspace(workspaceInfo *proto.WorkspaceInfo) {
	r.logger.Info("Attempting to publish workspace info")

	storeErr := storage.GetStorageHandler().CreateWorkspaceInfo(workspaceInfo)
	if storeErr != nil {
		r.logger.Error(fmt.Sprintf("Unable to store workspace info, %v", storeErr))
	}

	marshaler := jsonpb.Marshaler{}
	buf := new(bytes.Buffer)
	err := marshaler.Marshal(buf, workspaceInfo)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Unable to marshal workspace info, %v", err))
	}

	sendErr := r.pubSubTopic.Publish(r.ctx, buf.Bytes())
	if sendErr != nil {
		r.logger.Error(fmt.Sprintf("Unable to publish workspace info, %v", sendErr))
	}

	r.logger.Info("Workspace info successfully published")
}
