package config

// NodeConfig is the wrapper for the node configuration
type NodeConfig struct {
	HostAddress string
	HttpPort    int
	GrpcPort    int
	Libp2pPort  int
	BaseDir     string
}

// RendezvousConfig contains rendezvous nodes to which other rendezvous nodes
// should initially connect to
type RendezvousConfig struct {
	RendezvousNodes []string
}

// Default values
var (
	ServerHost       = "0.0.0.0"
	ServerHTTPPort   = 5000 // Used for the UI -> Client REST communication
	ServerGRPCPort   = 5001 // Used for Client <-> Client RPC communication
	ServerLibp2pPort = 5002 // Used for Client <-> Client network communication
)

// Directory names
var (
	DirectoryFiles   = "files"
	DirectoryLibp2p  = "libp2p"
	DirectoryStorage = "storage"
	DirectoryBase    = "app_data"

	// Client local //
	DirectoryTemp  = "temp"
	DirectoryShare = "share"
)

var (
	WORKSPACE_TYPE_SEND_ONLY    = "send-only"
	WORKSPACE_TYPE_RECEIVE_ONLY = "receive-only"
	WORKSPACE_TYPE_SEND_RECEIVE = "send-receive"
)

// Default rendezvous nodes that are already up and running
var (
	DefaultRendezvousNodes = []string{
		"/ip4/0.0.0.0/tcp/10002/p2p/QmX7WapncWyxiyuUfR33r9WiqR9brPuMiFhJ5v3hK67C2y", // TODO set
	}
)

var (
	// WorkspaceInfoProto is the protocol name used between Client and Rendezvous nodes
	// to exchange workspace information
	WorkspaceInfoProto = "/workspace-info-proto/0.1"

	// ClientVerificationProto is the protocol name used between Client nodes
	// to perform a handshake verification
	ClientVerificationProto = "/client-verification-proto/0.1"

	// FileSharingProto is the protocol used for file sharing
	FileSharingProto = "/file-share-proto/0.1"
)
