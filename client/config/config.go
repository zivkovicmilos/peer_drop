package config

// NodeConfig is the wrapper for the node configuration
type NodeConfig struct {
	HostAddress string
	HttpPort    int
	GrpcPort    int
	Libp2pPort  int
	BaseDir     string
}

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
