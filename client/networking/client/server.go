package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	newDiscovery "github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"github.com/zivkovicmilos/peer_drop/config"
	localCrypto "github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/files"
	"github.com/zivkovicmilos/peer_drop/proto"
	"github.com/zivkovicmilos/peer_drop/rest/types"
	"github.com/zivkovicmilos/peer_drop/rest/utils"
	"github.com/zivkovicmilos/peer_drop/storage"
	globalUtils "github.com/zivkovicmilos/peer_drop/utils"
	"google.golang.org/grpc"
)

type ClientServer struct {
	logger     hclog.Logger
	nodeConfig *config.NodeConfig

	// Networking metadata //
	me                  peer.ID              // the current node's peer ID
	host                host.Host            // the reference to the libp2p host
	rendezvousIDs       []peer.ID            // the peer IDs of the rendezvous nodes
	verifiedPeers       map[string][]peer.ID // the peer IDs of nodes who've passed verification for the given mnemonic
	pendingPeers        sync.Map             // peers awaiting verification. Only a single verification request is processed per user
	pendingPeersSize    int64
	kademliaDHT         *dht.IpfsDHT
	pubSub              *pubsub.PubSub                  // Reference to the pubsub service
	pubsubSubscriptions map[string]*pubsub.Subscription // In memory map of active subscriptions
	pubsubTopics        map[string]*pubsub.Topic        // In memory map of active topics
	//workspaceDirectoryMap map[string]string               // In memory map of workspace directories on disk (mnemonic -> dirName)
	pubsubSubscriptionsStop map[string]chan struct{} // Stop channel map
	pubsubTopicsStop        map[string]chan struct{} // Stop channel map
	findPeersStop           map[string]chan struct{} // Stop channel map

	// File handling //
	fileListerMap     map[string]*files.FileLister     // In memory map of file lister services (mnemonic -> fileLister)
	fileAggregatorMap map[string]*files.FileAggregator // In memory map of file aggregator services (mnemonic -> fileAggregator)

	// Workspace handler //
	newWorkspaceChannel chan *proto.WorkspaceInfo

	// Locks //
	rendezvousMux        sync.RWMutex
	verifiedPeersMuxMap  map[string]sync.RWMutex // mnemonic -> rwmutex
	fileListerMuxMap     map[string]sync.RWMutex // mnemonic -> rwmutex
	fileAggregatorMuxMap map[string]sync.RWMutex // mnemonic -> rwmutex

	// Context //
	ctx        context.Context
	cancelFunc context.CancelFunc

	// GRPC //
	proto.UnimplementedVerificationServiceServer

	// challenge id -> join request
	joinRequests map[string]*joinRequest // holds pending join requests TODO add garbage collector for stale requests
}

// NewClientServer returns a new client networking instance
func NewClientServer(
	logger hclog.Logger,
	nodeConfig *config.NodeConfig,
) *ClientServer {
	return &ClientServer{
		logger:               logger.Named("networking"),
		nodeConfig:           nodeConfig,
		rendezvousIDs:        make([]peer.ID, 0),
		verifiedPeers:        make(map[string][]peer.ID),
		joinRequests:         make(map[string]*joinRequest),
		newWorkspaceChannel:  make(chan *proto.WorkspaceInfo),
		pubsubTopics:         make(map[string]*pubsub.Topic),
		pubsubSubscriptions:  make(map[string]*pubsub.Subscription),
		fileListerMap:        make(map[string]*files.FileLister),
		fileAggregatorMap:    make(map[string]*files.FileAggregator),
		fileListerMuxMap:     make(map[string]sync.RWMutex),
		verifiedPeersMuxMap:  make(map[string]sync.RWMutex),
		fileAggregatorMuxMap: make(map[string]sync.RWMutex),

		pubsubSubscriptionsStop: make(map[string]chan struct{}),
		pubsubTopicsStop:        make(map[string]chan struct{}),
		findPeersStop:           make(map[string]chan struct{}),
	}
}

// numPendingPeers returns the number of peers pending handshake
func (cs *ClientServer) numPendingPeers() int64 {
	return atomic.LoadInt64(&cs.pendingPeersSize)
}

// isPendingPeer checks if the peer is awaiting handshake
func (cs *ClientServer) isPendingPeer(peerID peer.ID) bool {
	isPending, ok := cs.pendingPeers.Load(peerID)
	if !ok {
		return false
	}
	return isPending.(bool)
}

// removePendingPeer unsets the peer as pending handshake
func (cs *ClientServer) removePendingPeer(peerID peer.ID) {
	if _, loaded := cs.pendingPeers.LoadAndDelete(peerID); loaded {
		atomic.AddInt64(&cs.pendingPeersSize, -1)
	}
}

// setPendingPeer sets the peer as pending handshake
func (cs *ClientServer) setPendingPeer(peerID peer.ID) {
	if _, loaded := cs.pendingPeers.LoadOrStore(peerID, true); !loaded {
		atomic.AddInt64(&cs.pendingPeersSize, 1)
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

	defer cs.Stop()

	cs.host = clientHost
	cs.me = clientHost.ID()

	// Set up the local DHT
	options := []dht.Option{dht.Mode(dht.ModeServer)}

	kademliaDHT, err := dht.New(ctx, clientHost, options...)
	if err != nil {
		cs.logger.Error(fmt.Sprintf("Unable to start Kademlia DHT, %v", err))

		os.Exit(1)
	}

	cs.kademliaDHT = kademliaDHT

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

	// Set GRPC protocol handlers
	cs.setupProtocols()

	// Start the workspace handler loop
	go cs.workspaceJoinHandler()

	if startErr := cs.startPeerDropService(); startErr != nil {
		cs.logger.Error(fmt.Sprintf("Unable to start peer_drop service, %v", startErr))

		os.Exit(1)
	}

	<-closeChannel

	return
}

// Stop closes the client server and all pending services
func (cs *ClientServer) Stop() {
	cs.logger.Info("Closing client server...")
	// Cancel the context
	cs.cancelFunc()

	// Close the libp2p host
	_ = cs.host.Close()

	// Stop all find peers loops
	for _, findPeerLoop := range cs.findPeersStop {
		go func(findPeerLoop chan struct{}) {
			findPeerLoop <- struct{}{}
		}(findPeerLoop)
	}

	// Stop the topic subscribers
	for _, subscription := range cs.pubsubSubscriptionsStop {
		go func(subscription chan struct{}) {
			subscription <- struct{}{}
		}(subscription)
	}

	// Stop the topic publishers
	for _, topic := range cs.pubsubTopicsStop {
		go func(topic chan struct{}) {
			topic <- struct{}{}

		}(topic)
	}

	// Stop the file listers
	for _, fileLister := range cs.fileListerMap {
		go fileLister.Stop()
	}

	// Stop the file aggregators
	for _, fileAggregator := range cs.fileAggregatorMap {
		go fileAggregator.Stop()
	}

	// Close the workspace handler loop
	close(cs.newWorkspaceChannel)

	cs.logger.Info("Client server closed")
}

// workspaceJoinHandler handles incoming workspace join messages and starts the
// appropriate loops
func (cs *ClientServer) workspaceJoinHandler() {
	for {
		select {
		// Wait for workspace join messages
		case workspaceInfo, more := <-cs.newWorkspaceChannel:
			if more {
				// Initialize the workspace
				go func() {
					err := cs.initializeWorkspace(workspaceInfo)
					if err != nil {
						cs.logger.Error(fmt.Sprintf("Unable to initialize workspace, %v", err))

						return
					}

					cs.logger.Info(fmt.Sprintf("Successfully initialized workspace %s", workspaceInfo.Name))
				}()
			} else {
				cs.logger.Info("Closing workspace join handler...")
				return
			}
		}
	}
}

// startPeerDropService is the method that starts up the pubsub mechanism
// for sharing file lists
func (cs *ClientServer) startPeerDropService() error {
	// Grab all the available workspaces on this node
	foundWorkspaces, totalWorkspaces, findErr := storage.GetStorageHandler().GetWorkspaces(utils.NoPagination)
	if findErr != nil {
		return findErr
	}
	cs.logger.Info(fmt.Sprintf("Found a total of %d workspaces", totalWorkspaces))

	// Create the pubsub instance
	pubSub, gossipErr := pubsub.NewGossipSub(cs.ctx, cs.host)
	if gossipErr != nil {
		return gossipErr
	}
	cs.pubSub = pubSub

	// For each individual workspace, initialize it
	for _, workspaceInfo := range foundWorkspaces {
		//go func(workspaceInfo *proto.WorkspaceInfo) {
		//	err := cs.initializeWorkspace(workspaceInfo)
		//	if err != nil {
		//		cs.logger.Error(fmt.Sprintf("Unable to initialize workspace, %v", err))
		//	}
		//}(workspaceInfo)
		// TODO return to goroutine
		err := cs.initializeWorkspace(workspaceInfo)
		if err != nil {
			cs.logger.Error(fmt.Sprintf("Unable to initialize workspace, %v", err))
		}
	}

	cs.logger.Info("peer_drop service started")

	return nil
}

// initializeWorkspace initializes the publishers / listeners for an individual workspace
func (cs *ClientServer) initializeWorkspace(workspaceInfo *proto.WorkspaceInfo) error {
	mnemonic := workspaceInfo.Mnemonic
	// Create the folder structure if it doesn't exist
	if directoryErr := cs.initializeWorkspaceDirectory(workspaceInfo.Name, workspaceInfo.Mnemonic); directoryErr != nil {
		cs.logger.Error(
			fmt.Sprintf("Unable to initialize directory for %s, %v", workspaceInfo.Name, directoryErr),
		)

		return directoryErr
	}

	// Start the find peers service
	if findPeersErr := cs.findPeersWrapper(workspaceInfo); findPeersErr != nil {
		cs.logger.Error(
			fmt.Sprintf("Unable to start find peers service, %v", findPeersErr),
		)

		return findPeersErr
	}

	pubSubTopic, err := cs.pubSub.Join(workspaceInfo.Mnemonic)
	if err != nil {
		return fmt.Errorf("unable to join topic [%s], %v", workspaceInfo.Mnemonic, err)
	}

	pubSubSubscription, err := pubSubTopic.Subscribe()
	if err != nil {
		return fmt.Errorf("unable to subscribe to topic [%s], %v", workspaceInfo.Mnemonic, err)
	}

	// Update the in memory map values
	cs.pubsubSubscriptions[workspaceInfo.Mnemonic] = pubSubSubscription
	cs.pubsubTopics[workspaceInfo.Mnemonic] = pubSubTopic

	// Check if we are the owner of this workspace
	amOwner := cs.amIWorkspaceOwner(workspaceInfo)

	switch workspaceInfo.WorkspaceType {
	case config.WORKSPACE_TYPE_SEND_ONLY:
		if amOwner {
			// If we are the owner of this workspace, we only listen for messages
			go cs.startSubscriptionListener(mnemonic)
		} else {
			// If we are not the owner of this workspace, we only send messages
			go cs.startTopicPublisher(mnemonic)
		}
	case config.WORKSPACE_TYPE_RECEIVE_ONLY:
		if amOwner {
			// If we are the owner of this workspace, we only send messages
			go cs.startTopicPublisher(mnemonic)
		} else {
			// If we are not the owner of this workspace, we only receive messages
			go cs.startSubscriptionListener(mnemonic)
		}
	default:
		// Send & Receive
		go cs.startSubscriptionListener(mnemonic)
		go cs.startTopicPublisher(mnemonic)
	}

	cs.logger.Info(fmt.Sprintf("Workspace with mnemonic [%s] initialized", mnemonic))

	return nil
}

// initializeWorkspaceDirectory creates the workspace directory in the folder structure
func (cs *ClientServer) initializeWorkspaceDirectory(name string, mnemonic string) error {
	// Lowercase the directory name
	dirName := strings.ToLower(name)
	dirName = strings.Replace(dirName, " ", "-", -1)

	pathCommon := fmt.Sprintf("%s/%s/%s", cs.nodeConfig.BaseDir, config.DirectoryFiles, dirName)

	// baseDir/files/workspace-mnemonic/temp
	// Directory is used for temporary download data
	filesDirectory := fmt.Sprintf("%s/%s", pathCommon, config.DirectoryTemp)
	if createErr := globalUtils.CreateDirectory(filesDirectory); createErr != nil {
		return createErr
	}

	// baseDir/files/workspace-mnemonic/share
	// Directory is used for sharing local node files
	shareDirectory := fmt.Sprintf("%s/%s", pathCommon, config.DirectoryShare)
	if createErr := globalUtils.CreateDirectory(shareDirectory); createErr != nil {
		return createErr
	}

	// Start the file lister service for this directory
	fileLister := files.NewFileLister(
		cs.logger,
		shareDirectory,
		time.Second*5,
	)
	fileLister.Start()

	cs.registerFileLister(mnemonic, fileLister)

	return nil
}

// registerFileLister registers a new file lister
func (cs *ClientServer) registerFileLister(mnemonic string, fileLister *files.FileLister) {
	mux, _ := cs.fileListerMuxMap[mnemonic]
	mux.Lock()
	defer mux.Unlock()

	cs.fileListerMap[mnemonic] = fileLister

	fileLister.Start()
}

// unregisterFileLister unregisters a file lister
func (cs *ClientServer) unregisterFileLister(mnemonic string) {
	mux, _ := cs.fileListerMuxMap[mnemonic]
	mux.Lock()
	defer mux.Unlock()

	fileLister := cs.fileListerMap[mnemonic]
	fileLister.Stop()
	delete(cs.fileListerMap, mnemonic)
}

// findPeersWrapper is a wrapper function for starting the find peers service
// for a specific workspace
func (cs *ClientServer) findPeersWrapper(workspaceInfo *proto.WorkspaceInfo) error {
	// Contains necessary information to perform new peer handshakes
	info := &handshakeInfo{}

	workspaceCredentials, getErr := storage.GetStorageHandler().GetWorkspaceCredentials(workspaceInfo.Mnemonic)
	if getErr != nil {
		return getErr
	}

	if workspaceCredentials == nil {
		return errors.New("no workspace credentials found")
	}

	info.publicKey = workspaceCredentials.PublicKey
	info.privateKey = workspaceCredentials.PrivateKey
	info.password = workspaceCredentials.Password

	go cs.findPeers(workspaceInfo.Mnemonic, info)

	return nil
}

// unregisterFileAggregator removes / stops a file aggregator service
func (cs *ClientServer) unregisterFileAggregator(mnemonic string) {
	mux, _ := cs.fileAggregatorMuxMap[mnemonic]
	mux.Lock()
	fileAggregator := cs.fileAggregatorMap[mnemonic]

	delete(cs.fileAggregatorMap, mnemonic)
	fileAggregator.Stop()
	mux.Unlock()
}

// startSubscriptionListener starts the subscription listener for a workspace mnemonic
func (cs *ClientServer) startSubscriptionListener(mnemonic string) {
	subscription := cs.pubsubSubscriptions[mnemonic]
	subContext := context.Background()

	// Create the stop channel
	stopChannel := make(chan struct{})
	cs.pubsubSubscriptionsStop[mnemonic] = stopChannel

	// Create the file aggregator instance
	updateChannel := make(chan files.FileListWrapper)
	fileAggregator := files.NewFileAggregator(cs.logger, mnemonic, updateChannel)

	mux, _ := cs.fileAggregatorMuxMap[mnemonic]
	mux.Lock()
	cs.fileAggregatorMap[mnemonic] = fileAggregator
	mux.Unlock()

	fileAggregator.Start()

	for {
		select {
		case _ = <-stopChannel:
			cs.logger.Info(fmt.Sprintf("Stopping subscription listener for mnemonic [%s]", mnemonic))
			subscription.Cancel()
			return
		default:
		}
		fileListMessage, err := subscription.Next(subContext)
		if err != nil {
			cs.logger.Error(fmt.Sprintf("Unable to parse message, %v", err))
			return
		}
		cs.logger.Info(fmt.Sprintf("Received a new file list message for mnemonic [%s]", mnemonic))

		// Forward messages that are not from us
		if fileListMessage.ReceivedFrom == cs.me {
			cs.logger.Info("Pubsub message skipped")
			continue
		}

		peerFileList := new(proto.FileList)
		err = jsonpb.Unmarshal(bytes.NewReader(fileListMessage.Data), peerFileList)
		if err != nil {
			cs.logger.Error(fmt.Sprintf("Unmarshal error %v", err))
			continue
		}

		updateChannel <- files.FileListWrapper{
			FileList: peerFileList,
			PeerID:   fileListMessage.ReceivedFrom,
		}
	}
}

// startTopicPublisher starts up the file list sharing loop
func (cs *ClientServer) startTopicPublisher(mnemonic string) {
	topic := cs.pubsubTopics[mnemonic]
	topicContext := context.Background()
	ticker := time.NewTicker(time.Second * 5)

	// Create the stop channel
	stopChannel := make(chan struct{})
	cs.pubsubTopicsStop[mnemonic] = stopChannel

	mux, _ := cs.fileListerMuxMap[mnemonic]
	mux.RLock()
	fileLister := cs.fileListerMap[mnemonic]
	fileLister.Start()
	mux.RUnlock()

	for {
		select {
		case _ = <-stopChannel:
			ticker.Stop()
			_ = topic.Close()
			cs.logger.Info(fmt.Sprintf("Stopping topic publisher for mnemonic [%s]", mnemonic))
			return

		// Every 5 seconds, share the file list
		case _ = <-ticker.C:
			// Compile file list
			localFileList := fileLister.GetAvailableFiles()
			if len(localFileList) > 0 {

				fileList := &proto.FileList{
					FileList: localFileList,
				}

				// Share the file list to the topic
				marshaler := jsonpb.Marshaler{}
				buf := new(bytes.Buffer)
				err := marshaler.Marshal(buf, fileList)
				if err != nil {
					cs.logger.Error(fmt.Sprintf("Unable to marshal file list, %v", err))
					continue
				}

				sendErr := topic.Publish(topicContext, buf.Bytes())
				if sendErr != nil {
					cs.logger.Error(fmt.Sprintf("Unable to publish local file list, %v", sendErr))
					continue
				}

				cs.logger.Info("Workspace file list successfully published")
			}
		}
	}
}

// amIWorkspaceOwner checks if the current node has any keys which correspond
// to the given workspace owner key list
func (cs *ClientServer) amIWorkspaceOwner(info *proto.WorkspaceInfo) bool {
	// Range over the public keys of the workspace owners and check
	// if the current node is the owner of one of these keys
	identities, count, findErr := storage.GetStorageHandler().GetIdentities(utils.NoPagination, utils.DefaultSort)
	if findErr != nil {
		cs.logger.Error(fmt.Sprintf("Unable to fetch identities, %v", findErr))
	}

	if count == 0 {
		cs.logger.Error("No identity detected")
		panic("Invalid situation")
	}

	for _, workspaceOwnerPK := range info.WorkspaceOwnerPublicKeys {
		for _, identity := range identities {
			if identity.PublicKey == workspaceOwnerPK {
				return true
			}
		}
	}

	return false
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

				cs.rendezvousMux.Lock()
				cs.rendezvousIDs = append(cs.rendezvousIDs, peerinfo.ID)
				cs.rendezvousMux.Unlock()
			}
		}()
	}
	wg.Wait()
}

// GetWorkspaceInfo fetches workspace info from rendezvous nodes
func (cs *ClientServer) GetWorkspaceInfo(mnemonic string) (*proto.WorkspaceInfo, error) {
	rendezvousID, findErr := cs.findBestRendezvous()
	if findErr != nil {
		return nil, findErr
	}

	stream, err := cs.host.NewStream(cs.ctx, *rendezvousID, protocol.ID(config.WorkspaceInfoProto))
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate stream to rendezvous node, %v", err)
	}
	defer func(stream network.Stream) {
		if streamCloseErr := stream.Close(); streamCloseErr != nil {
			cs.logger.Error(fmt.Sprintf("Unable to gracefully close stream, %v", streamCloseErr))
		}
	}(stream)

	// Grab the wrapped connection
	clientConn := WrapStreamInClient(stream)

	// Instantiate the proto client
	clientProto := proto.NewWorkspaceInfoServiceClient(clientConn.(*grpc.ClientConn))

	// Call the RPC method
	return clientProto.GetWorkspaceInfo(
		context.Background(),
		&proto.WorkspaceInfoRequest{Mnemonic: mnemonic},
	)
}

// findBestRendezvous finds the best suitable rendezvous node
func (cs *ClientServer) findBestRendezvous() (*peer.ID, error) {
	cs.rendezvousMux.Lock()
	defer cs.rendezvousMux.Unlock()

	for indx, peerID := range cs.rendezvousIDs {
		if cs.host.Network().Connectedness(peerID) == network.Connected {
			return &peerID, nil
		} else {
			// Remove the unconnected rendezvous from the array
			cs.rendezvousIDs = append(cs.rendezvousIDs[:indx], cs.rendezvousIDs[indx+1:]...)
		}
	}

	return nil, errors.New("no suitable rendezvous node found")
}

// CreateWorkspace sends a new workspace request to the rendezvous nodes
func (cs *ClientServer) CreateWorkspace(workspaceRequest types.NewWorkspaceRequest) (*proto.WorkspaceInfo, error) {
	rendezvousID, findErr := cs.findBestRendezvous()
	if findErr != nil {
		return nil, findErr
	}

	stream, err := cs.host.NewStream(cs.ctx, *rendezvousID, protocol.ID(config.WorkspaceInfoProto))
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate stream to rendezvous node, %v", err)
	}
	defer func(stream network.Stream) {
		if streamCloseErr := stream.Close(); streamCloseErr != nil {
			cs.logger.Error(fmt.Sprintf("Unable to gracefully close stream, %v", streamCloseErr))
		}
	}(stream)

	// Grab the wrapped connection
	clientConn := WrapStreamInClient(stream)

	// Instantiate the proto client
	clientProto := proto.NewWorkspaceInfoServiceClient(clientConn.(*grpc.ClientConn))

	// Call the RPC method
	return clientProto.CreateNewWorkspace(
		context.Background(),
		cs.workspaceRequestToWorkspaceInfo(workspaceRequest),
	)
}

func (cs *ClientServer) workspaceRequestToWorkspaceInfo(
	workspaceRequest types.NewWorkspaceRequest,
) *proto.WorkspaceInfo {
	workspaceInfo := &proto.WorkspaceInfo{
		Name: workspaceRequest.WorkspaceName,
	}

	// Set access control type
	if workspaceRequest.WorkspaceAccessControlType == "Password" {
		workspaceInfo.SecurityType = "password"

		// Hash the password, then hash the hash
		// Bad idea, but it will work for now
		passwordHash := localCrypto.NewSHA256([]byte(workspaceRequest.WorkspaceAccessControl.Password))
		passwordHash = localCrypto.NewSHA256(passwordHash)
		workspaceInfo.SecuritySettings = &proto.WorkspaceInfo_PasswordHash{
			PasswordHash: hex.EncodeToString(passwordHash),
		}
	} else {
		contactPublicKeys := make([]string, 0)

		for _, contact := range workspaceRequest.WorkspaceAccessControl.Contacts {
			contactPublicKeys = append(contactPublicKeys, contact.PublicKey)
		}

		workspaceInfo.SecurityType = "contacts"
		workspaceInfo.SecuritySettings = &proto.WorkspaceInfo_ContactsWrapper{
			ContactsWrapper: &proto.ContactsWrapper{
				ContactPublicKeys: contactPublicKeys,
			},
		}
	}

	switch workspaceRequest.WorkspaceType {
	case "Send only":
		workspaceInfo.WorkspaceType = "send-only"
	case "Receive only":
		workspaceInfo.WorkspaceType = "receive-only"
	default:
		workspaceInfo.WorkspaceType = "send-receive"
	}

	workspaceOwners := make([]string, 0)

	workspaceOwners = append(workspaceOwners, workspaceRequest.WorkspaceOwners...)
	workspaceInfo.WorkspaceOwnerPublicKeys = workspaceOwners

	return workspaceInfo
}

// GRPC //

// findPeers is a loop function that gathers new peers for the
// specific workspace mnemonic
func (cs *ClientServer) findPeers(workspaceMnemonic string, info *handshakeInfo) {
	findPeersCtx, cancelFunc := context.WithCancel(context.Background())

	ticker := time.NewTicker(10 * time.Second)
	routingDiscovery := discovery.NewRoutingDiscovery(cs.kademliaDHT)
	discovery.Advertise(findPeersCtx, routingDiscovery, workspaceMnemonic, newDiscovery.TTL(time.Second*5))
	cs.logger.Info(fmt.Sprintf("Successfully announced workspace file request [%s]", workspaceMnemonic))

	closeChannel := make(chan struct{})
	cs.findPeersStop[workspaceMnemonic] = closeChannel

	cs.logger.Info(fmt.Sprintf("We are %s", cs.me.String()))

	// Instantiate verified peers list
	mux, _ := cs.verifiedPeersMuxMap[workspaceMnemonic]
	mux.Lock()
	verifiedPeers := make([]peer.ID, 0)
	cs.verifiedPeers[workspaceMnemonic] = verifiedPeers
	mux.Unlock()

	advertiseString := workspaceMnemonic
	for {
		select {
		case <-closeChannel:
			ticker.Stop()
			cancelFunc()
			return
		case _ = <-ticker.C:
			// Find peers
			cs.logger.Info(fmt.Sprintf("I have exactly %d peers", len(cs.host.Network().Peers())))
			cs.logger.Info(fmt.Sprintf("Searching for other peers [%s]...", advertiseString))
			peerChan, err := routingDiscovery.FindPeers(findPeersCtx, advertiseString)
			if err != nil {
				cs.logger.Error("Unable to find peers ", err)
				continue
			}

			// Attempt to verify these peers
			// TODO make sure this is not infinite
			cs.logger.Debug(fmt.Sprintf("Number of peers found %d", len(peerChan)))
			for foundPeer := range peerChan {
				if foundPeer.ID.String() == cs.me.String() {
					continue
				}

				if cs.isPendingPeer(foundPeer.ID) {
					// handshake has already started
					continue
				}

				if cs.isVerifiedPeer(foundPeer.ID, workspaceMnemonic) {
					// peer already verified
					cs.logger.Info(fmt.Sprintf("Peer %s already verified...", foundPeer.ID.String()))
					continue
				} else {
					cs.logger.Info(fmt.Sprintf("Peer %s NOT verified...", foundPeer.ID.String()))
				}

				// Mark the peer as pending for the handshake
				cs.setPendingPeer(foundPeer.ID)

				cs.logger.Info(fmt.Sprintf("Found peer %s", foundPeer.ID.Pretty()))
				cs.logger.Info(fmt.Sprintf("Attempting connection to peer %s", foundPeer.ID))

				go func(peerID peer.ID) {
					defer func() {
						// Clear the value
						if cs.isPendingPeer(peerID) {
							cs.removePendingPeer(peerID)
						}
					}()

					if handshakeErr := cs.handleHandshake(peerID, workspaceMnemonic, info); handshakeErr != nil {
						cs.logger.Error(fmt.Sprintf("Unable to perform handshake, %v", handshakeErr))
						cs.disconnectFromPeer(peerID)
					} else {
						cs.logger.Info(fmt.Sprintf("Peer verified and connection established [%s]", peerID))
						//cs.addVerifiedPeer(workspaceMnemonic, peerID)
					}
				}(foundPeer.ID)
			}
		}
	}
}

type handshakeInfo struct {
	publicKey  *string // PEM encoded
	privateKey *string // PEM encoded
	password   *string // password from DB
}

// isVerifiedPeer checks if the peer is verified for that workspace
func (cs *ClientServer) isVerifiedPeer(peerID peer.ID, mnemonic string) bool {
	mux, _ := cs.verifiedPeersMuxMap[mnemonic]
	mux.RLock()
	defer mux.RUnlock()

	verifiedPeers, ok := cs.verifiedPeers[mnemonic]
	if !ok {
		return false
	}

	for _, verifiedPeerID := range verifiedPeers {
		if verifiedPeerID == peerID {
			return true
		}
	}

	return false
}

// isRendezvousNode checks if the currently connecting peer is a rendezvous node
func (cs *ClientServer) isRendezvousNode(peerID peer.ID) bool {
	for _, nodeID := range cs.rendezvousIDs {
		if peerID == nodeID {
			return true
		}
	}

	return false
}

// disconnectFromPeer disconnects from a certain peer if there is a connection
func (cs *ClientServer) disconnectFromPeer(peerID peer.ID) {
	if cs.host.Network().Connectedness(peerID) == network.Connected {
		cs.logger.Info(fmt.Sprintf("Disconnecting from peer %s", peerID))
		_ = cs.host.Network().ClosePeer(peerID)
	}
}

// handleHandshake executes the handshake process
func (cs *ClientServer) handleHandshake(
	peerID peer.ID,
	workspaceMnemonic string,
	info *handshakeInfo,
) error {
	stream, err := cs.host.NewStream(cs.ctx, peerID, protocol.ID(config.ClientVerificationProto))
	if err != nil {
		return fmt.Errorf("unable to instantiate stream to client node, %v", err)
	}
	defer func(stream network.Stream) {
		if streamCloseErr := stream.Close(); streamCloseErr != nil {
			cs.logger.Error(fmt.Sprintf("Unable to gracefully close stream, %v", streamCloseErr))
		}
	}(stream)

	// Grab the wrapped connection
	clientConn := WrapStreamInClient(stream)

	// Instantiate the proto client
	clientProto := proto.NewVerificationServiceClient(clientConn.(*grpc.ClientConn))

	verificationRequest := &proto.VerificationRequest{}
	verificationRequest.WorkspaceMnemonic = workspaceMnemonic
	if info.publicKey != nil {
		verificationRequest.PublicKey = info.publicKey
	}

	challenge, challengeErr := clientProto.BeginVerification(
		context.Background(),
		verificationRequest,
	)
	if challengeErr != nil {
		return challengeErr
	}

	var challengeSolution *proto.ChallengeSolution
	challengeSolution = nil

	if info.publicKey != nil {
		// Public key challenge
		publicKeyChallengeSolution, solveErr := SolvePublicKeyChallenge(challenge, *info.privateKey)
		if solveErr != nil {
			return solveErr
		}

		challengeSolution = publicKeyChallengeSolution
	} else {
		// Password challenge
		passwordChallengeSolution, solveErr := SolvePasswordChallenge(challenge, *info.password)
		if solveErr != nil {
			return solveErr
		}

		challengeSolution = passwordChallengeSolution
	}

	// Now that the challenge is solved,
	// send the solution to the verifier
	verificationResponse, verificationErr := clientProto.FinishVerification(
		context.Background(),
		challengeSolution,
	)
	if verificationErr != nil {
		return verificationErr
	}

	if !verificationResponse.Confirmed {
		return errors.New("unable to pass verification")
	}

	return nil
}

// setupProtocols sets up all the supported GRPC protocols
func (cs *ClientServer) setupProtocols() {
	cs.logger.Info("Setting up GRPC protocols...")

	// Set up the verification protocol
	verificationProtocol := NewGRPCProtocol()
	proto.RegisterVerificationServiceServer(verificationProtocol.GrpcServer(), cs)
	verificationProtocol.Serve()

	cs.host.SetStreamHandler(protocol.ID(config.ClientVerificationProto), func(stream network.Stream) {
		peerID := stream.Conn().RemotePeer()
		cs.logger.Info(fmt.Sprintf("Open stream [protocol: %s], by peer %s", config.ClientVerificationProto, peerID))

		verificationProtocol.Handler()(stream)
	})

	// TODO set up the file transfer protocol

	cs.logger.Info("GRPC protocols set")

}

// joinRequest is a way for the client node to track
// any pending join requests
type joinRequest struct {
	workspaceMnemonic string
	unencryptedData   []byte
	challenge         *proto.Challenge
}

// BeginVerification starts the verification process and returns the challenge
func (cs *ClientServer) BeginVerification(
	context context.Context,
	request *proto.VerificationRequest,
) (*proto.Challenge, error) {
	// Check if the current node has the workspace info
	workspaceInfo, findErr := storage.GetStorageHandler().GetWorkspaceInfo(request.WorkspaceMnemonic)
	if findErr != nil {
		cs.logger.Error(
			"Verification requested for unknown workspace ",
			request.WorkspaceMnemonic,
			"error",
			findErr,
		)

		return nil, errors.New("unknown workspace")
	}

	workspaceCredentials, findCredErr := storage.GetStorageHandler().GetWorkspaceCredentials(request.WorkspaceMnemonic)
	if findCredErr != nil {
		cs.logger.Error(
			"Verification requested for unknown workspace ",
			request.WorkspaceMnemonic,
			"error",
			findErr,
		)

		return nil, errors.New("unknown workspace")
	}

	var challenge *proto.Challenge
	challenge = nil
	unencryptedData := []byte(uuid.New().String())

	if workspaceInfo.SecurityType == "password" {
		cs.logger.Debug("Verification with password")

		if workspaceCredentials.Password == nil {
			cs.logger.Error("Invalid credentials")

			return nil, errors.New("unknown workspace")
		}
		passwordChallenge, constructErr := ConstructPasswordChallenge(
			unencryptedData,
			*workspaceCredentials.Password,
		)
		if constructErr != nil {
			return nil, errors.New("unable to construct challenge")
		}

		challenge = passwordChallenge
	} else {
		// Use the contacts public key to construct the challenge
		cs.logger.Debug("Verification with public key")

		// Check if the contact attached a public key
		if request.PublicKey == nil {
			return nil, errors.New("invalid request")
		}

		// Search for the public key in permitted contacts
		if workspaceCredentials.PublicKey == nil {
			cs.logger.Error("Invalid credentials")

			return nil, errors.New("unknown workspace")
		}

		publicKeyChallenge, constructErr := ConstructPublicKeyChallenge(
			unencryptedData,
			*workspaceCredentials.PublicKey,
		)
		if constructErr != nil {
			return nil, errors.New("unable to construct challenge")
		}

		challenge = publicKeyChallenge
	}

	if challenge == nil {
		return nil, errors.New("unable to construct challenge")
	}

	// Save join request locally
	cs.joinRequests[challenge.ChallengeId] = &joinRequest{
		workspaceMnemonic: workspaceInfo.Mnemonic,
		challenge:         challenge,
		unencryptedData:   unencryptedData,
	}

	return challenge, nil
}

// FinishVerification solves the challenge issued in order to finish verification
func (cs *ClientServer) FinishVerification(
	context context.Context,
	request *proto.ChallengeSolution,
) (*proto.VerificationResponse, error) {
	// Check if we have the pending request
	pendingJoinRequest, found := cs.joinRequests[request.ChallengeId]
	if !found {
		// TODO close connection to peer
		return ConstructVerificationResponse("Unknown challenge", false),
			errors.New("unknown challenge")
	}

	// Verify that the time signature is correct
	timestamp := time.Unix(pendingJoinRequest.challenge.Timestamp, 0)

	if timestamp.After(time.Now().Add(time.Second * 30)) {
		// Timestamp invalid
		return ConstructVerificationResponse("Invalid timestamp", false),
			errors.New("invalid timestamp")
	}

	// Verify that the unencrypted data is correct
	if bytes.Compare(pendingJoinRequest.unencryptedData, request.DecryptedValue) != 0 {
		return ConstructVerificationResponse("Invalid decrypt data", false),
			errors.New("invalid decrypted data")
	}

	// Remove the join request
	delete(cs.joinRequests, request.ChallengeId)

	// Add the peer to verified peers
	typedContext := context.(*WrappedContext)
	cs.addVerifiedPeer(pendingJoinRequest.workspaceMnemonic, typedContext.PeerID)

	return ConstructVerificationResponse("Verification success", true), nil
}

// addVerifiedPeer adds a verified peer. [Thread safe]
func (cs *ClientServer) addVerifiedPeer(mnemonic string, newPeer peer.ID) {
	mux, _ := cs.verifiedPeersMuxMap[mnemonic]
	mux.Lock()
	defer mux.Unlock()

	verifiedPeers, ok := cs.verifiedPeers[mnemonic]
	if !ok {
		// No peers yet, create the array
		verifiedPeers = make([]peer.ID, 0)
	}

	verifiedPeers = append(verifiedPeers, newPeer)
	cs.verifiedPeers[mnemonic] = verifiedPeers
}

// removeVerifiedPeer removes a peer from the verified array
func (cs *ClientServer) removeVerifiedPeer(mnemonic string, oldPeer peer.ID) {
	mux, _ := cs.verifiedPeersMuxMap[mnemonic]
	mux.Lock()
	defer mux.Unlock()

	verifiedPeers, ok := cs.verifiedPeers[mnemonic]
	if !ok {
		// No peers yet
		return
	}

	indx := -1
	for index, verifiedPeer := range verifiedPeers {
		if verifiedPeer == oldPeer {
			indx = index
		}
	}
	if indx >= 0 {
		verifiedPeers = append(verifiedPeers[:indx], verifiedPeers[:indx+1]...)
	}

	cs.verifiedPeers[mnemonic] = verifiedPeers
}

// Workspace joining //

// JoinWorkspacePassword handles workspace join requests with passwords
func (cs *ClientServer) JoinWorkspacePassword(
	workspaceInfo *proto.WorkspaceInfo,
	password string,
) bool {
	// Prepare password hash
	passwordHash := localCrypto.NewSHA256([]byte(password))
	passwordHash = localCrypto.NewSHA256(passwordHash)

	workspacePasswordHash := workspaceInfo.SecuritySettings.(*proto.WorkspaceInfo_PasswordHash).PasswordHash

	return workspacePasswordHash == password
}

// JoinWorkspacePublicKey handles workspace join requests with public keys
func (cs *ClientServer) JoinWorkspacePublicKey(
	workspaceInfo *proto.WorkspaceInfo,
	publicKeyPEM string,
) bool {

	inSet := false

	contactPublicKeys := workspaceInfo.SecuritySettings.(*proto.WorkspaceInfo_ContactsWrapper).ContactsWrapper.ContactPublicKeys
	for _, contactPublicKey := range contactPublicKeys {
		if contactPublicKey == publicKeyPEM {
			inSet = true
			break
		}
	}

	return inSet
}

// TriggerWorkspaceInit alerts the workspace handler that a new workspace
// has been joined
func (cs *ClientServer) TriggerWorkspaceInit(workspaceInfo *proto.WorkspaceInfo) {
	cs.newWorkspaceChannel <- workspaceInfo
}
