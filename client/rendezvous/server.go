package rendezvous

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"github.com/zivkovicmilos/peer_drop/config"
	localCrypto "github.com/zivkovicmilos/peer_drop/crypto"
)

type RendezvousServer struct {
	logger           hclog.Logger
	nodeConfig       *config.NodeConfig
	rendezvousConfig *config.RendezvousConfig
}

// NewRendezvousServer returns a new instance of the rendezvous server
func NewRendezvousServer(
	logger hclog.Logger,
	nodeConfig *config.NodeConfig,
	rendezvousConfig *config.RendezvousConfig,
) *RendezvousServer {
	return &RendezvousServer{
		logger:           logger,
		nodeConfig:       nodeConfig,
		rendezvousConfig: rendezvousConfig,
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

	ctx := context.Background()

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ip4/%s/tcp/%d", r.nodeConfig.HostAddress, r.nodeConfig.Libp2pPort),
	)

	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(libp2pKey),
	)

	if err != nil {
		panic(err)
	}

	_, err = dht.New(ctx, host)
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", r.nodeConfig.HostAddress, r.nodeConfig.Libp2pPort, host.ID().Pretty())
	fmt.Println("")
}
