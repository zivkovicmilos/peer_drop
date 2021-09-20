package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"github.com/zivkovicmilos/peer_drop/config"
	localCrypto "github.com/zivkovicmilos/peer_drop/crypto"
	"github.com/zivkovicmilos/peer_drop/proto"
	"github.com/zivkovicmilos/peer_drop/rest/types"
	"github.com/zivkovicmilos/peer_drop/storage"
	"google.golang.org/grpc"
)

type ClientServer struct {
	logger       hclog.Logger
	nodeConfig   *config.NodeConfig
	closeChannel chan struct{}

	// Networking metadata //
	me               peer.ID              // the current node's peer ID
	host             host.Host            // the reference to the libp2p host
	rendezvousIDs    []peer.ID            // the peer IDs of the rendezvous nodes
	verifiedPeers    map[string][]peer.ID // the peer IDs of nodes who've passed verification for the given mnemonic
	pendingPeers     sync.Map             // peers awaiting verification. Only a single verification request is processed per user
	pendingPeersSize int64
	kademliaDHT      *dht.IpfsDHT

	// Locks //
	rendezvousMux    sync.Mutex
	verifiedPeersMux sync.Mutex // a mux map would be a better solution

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
		logger:        logger.Named("networking"),
		nodeConfig:    nodeConfig,
		rendezvousIDs: make([]peer.ID, 0),
		verifiedPeers: make(map[string][]peer.ID),
		joinRequests:  make(map[string]*joinRequest),
	}
}

func (cs *ClientServer) numPendingPeers() int64 {
	return atomic.LoadInt64(&cs.pendingPeersSize)
}

func (cs *ClientServer) isPendingPeer(peerID peer.ID) bool {
	isPending, ok := cs.pendingPeers.Load(peerID)
	if !ok {
		return false
	}
	return isPending.(bool)
}

func (cs *ClientServer) removePendingPeer(peerID peer.ID) {
	if _, loaded := cs.pendingPeers.LoadAndDelete(peerID); loaded {
		atomic.AddInt64(&cs.pendingPeersSize, -1)
	}
}

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
				cs.rendezvousIDs = append(cs.rendezvousIDs, peerinfo.ID)
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

		workspaceInfo.SecuritySettings = &proto.WorkspaceInfo_PasswordHash{
			PasswordHash: workspaceRequest.WorkspaceAccessControl.Password,
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

	ticker := time.NewTicker(5 * time.Second)
	routingDiscovery := discovery.NewRoutingDiscovery(cs.kademliaDHT)
	discovery.Advertise(findPeersCtx, routingDiscovery, workspaceMnemonic)
	cs.logger.Info(fmt.Sprintf("Successfully announced workspace file request [%s]", workspaceMnemonic))

	for {
		select {
		case <-cs.closeChannel:
			ticker.Stop()
			cancelFunc()
			return
		case _ = <-ticker.C:
			// Find peers
			cs.logger.Info("Searching for other peers...")
			peerChan, err := routingDiscovery.FindPeers(findPeersCtx, workspaceMnemonic)
			if err != nil {
				cs.logger.Error("Unable to find peers ", err)
				continue
			}

			// Attempt to verify these peers
			for foundPeer := range peerChan {
				if foundPeer.ID == cs.me {
					continue
				}

				if cs.isPendingPeer(foundPeer.ID) {
					// handshake has already started
					continue
				}

				if cs.isVerifiedPeer(foundPeer.ID, workspaceMnemonic) {
					// peer already verified
					continue
				}

				// Mark the peer as pending for the handshake
				cs.setPendingPeer(foundPeer.ID)

				cs.logger.Info("Found peer", foundPeer.ID)
				cs.logger.Info("Attempting connection to peer", foundPeer.ID)

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
	cs.verifiedPeersMux.Lock()
	defer cs.verifiedPeersMux.Unlock()

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
		cs.logger.Info("Disconnecting from peer", peerID)
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

	var verificationRequest *proto.VerificationRequest
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

		return nil, errors.New("unknown workspace") // TODO conform grpc error
	}

	var challenge *proto.Challenge
	challenge = nil
	unencryptedData := []byte(uuid.New().String())

	if workspaceInfo.SecurityType == "password" {
		passwordChallenge, constructErr := ConstructPasswordChallenge(
			unencryptedData,
			"myPassword", // TODO change password
		)
		if constructErr != nil {
			return nil, errors.New("unable to construct challenge")
		}

		challenge = passwordChallenge
	} else {
		// Use the contacts public key to construct the challenge

		// Check if the contact attached a public key
		if request.PublicKey == nil {
			return nil, errors.New("invalid request")
		}

		// Search for the public key in permitted contacts
		contactsSecurity := workspaceInfo.SecuritySettings.(*proto.WorkspaceInfo_ContactsWrapper)

		for _, publicKey := range contactsSecurity.ContactsWrapper.ContactPublicKeys {
			if publicKey == *request.PublicKey {
				// Key found, create the challenge
				publicKeyChallenge, constructErr := ConstructPublicKeyChallenge(
					unencryptedData,
					publicKey,
				)
				if constructErr != nil {
					return nil, errors.New("unable to construct challenge")
				}

				challenge = publicKeyChallenge
				break
			}
		}
	}

	if challenge == nil {
		return nil, errors.New("invalid request")
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
	typedContext := context.(WrappedContext)
	cs.addVerifiedPeer(pendingJoinRequest.workspaceMnemonic, typedContext.PeerID)

	return ConstructVerificationResponse("Verification success", true), nil
}

// addVerifiedPeer adds a verified peer. [Thread safe]
func (cs *ClientServer) addVerifiedPeer(mnemonic string, newPeer peer.ID) {
	cs.verifiedPeersMux.Lock()
	defer cs.verifiedPeersMux.Unlock()

	verifiedPeers, ok := cs.verifiedPeers[mnemonic]
	if !ok {
		// No peers yet, create the array
		verifiedPeers = make([]peer.ID, 0)
	}

	verifiedPeers = append(verifiedPeers, newPeer)
	cs.verifiedPeers[mnemonic] = verifiedPeers
}

func (cs *ClientServer) removeVerifiedPeer(mnemonic string, oldPeer peer.ID) {
	cs.verifiedPeersMux.Lock()
	defer cs.verifiedPeersMux.Unlock()

	verifiedPeers, ok := cs.verifiedPeers[mnemonic]
	if !ok {
		// No peers yet
		return
	}

	indx := -1
	for index, peer := range verifiedPeers {
		if peer == oldPeer {
			indx = index
		}
	}
	if indx >= 0 {
		verifiedPeers = append(verifiedPeers[:indx], verifiedPeers[:indx+1]...)
	}

	cs.verifiedPeers[mnemonic] = verifiedPeers
}
