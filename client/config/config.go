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
)

// Default rendezvous nodes that are already up and running
var (
	DefaultRendezvousNodes = []string{
		"/ip4/0.0.0.0/tcp/10002/p2p/QmNmc9kho9D6ji1Curf2e9Xsi8h7PAHZPpr8r7BcaevnTx", // TODO set
	}
)

var (
	// WorkspaceInfoProto is the protocol name used between Client and Rendezvous nodes
	// to exchange workspace information
	WorkspaceInfoProto = "/workspace-info-proto/0.1"
)
