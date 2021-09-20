package client

import (
	"context"
	"io"
	"net"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	manet "github.com/multiformats/go-multiaddr/net"
	"google.golang.org/grpc"
	grpcPeer "google.golang.org/grpc/peer"
)

type GRPCProtocol struct {
	ctx      context.Context
	streamCh chan network.Stream

	grpcServer *grpc.Server
}

func NewGRPCProtocol() *GRPCProtocol {
	g := &GRPCProtocol{
		ctx:        context.Background(),
		streamCh:   make(chan network.Stream),
		grpcServer: grpc.NewServer(grpc.UnaryInterceptor(interceptor)),
	}

	return g
}

type WrappedContext struct {
	context.Context
	PeerID peer.ID
}

func interceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	peer, _ := grpcPeer.FromContext(ctx)

	// we expect our libp2p wrapper
	addr := peer.Addr.(*libp2pInfoWrapper)

	ctx2 := &WrappedContext{
		Context: ctx,
		// Wrap the context so the peer id can be extracted
		// from any handler
		PeerID: addr.peerID,
	}
	h, err := handler(ctx2, req)

	return h, err
}

func (g *GRPCProtocol) Client(stream network.Stream) interface{} {
	return WrapStreamInClient(stream)
}

func (g *GRPCProtocol) Serve() {
	go g.grpcServer.Serve(g)
}

func (g *GRPCProtocol) Accept() (net.Conn, error) {
	select {
	case <-g.ctx.Done():
		return nil, io.EOF
	case stream := <-g.streamCh:
		return &streamWrapper{Stream: stream}, nil
	}
}

// Addr implements the net.Listener interface
func (g *GRPCProtocol) Addr() net.Addr {
	return fakeLocalAddr()
}

func (g *GRPCProtocol) Close() error {
	return nil
}

func (g *GRPCProtocol) Handler() func(network.Stream) {
	return func(stream network.Stream) {
		select {
		case <-g.ctx.Done():
			return
		case g.streamCh <- stream:
		}
	}
}

func (g *GRPCProtocol) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	g.grpcServer.RegisterService(sd, ss)
}

func (g *GRPCProtocol) GrpcServer() *grpc.Server {
	return g.grpcServer
}

func WrapStreamInClient(s network.Stream) interface{} { // *grpc.ClientConn
	opts := grpc.WithContextDialer(func(ctx context.Context, peerIdStr string) (net.Conn, error) {
		return &streamWrapper{s}, nil
	})
	conn, err := grpc.Dial("", grpc.WithInsecure(), opts)
	if err != nil {
		return nil
	}

	return conn
}

type streamWrapper struct {
	network.Stream
}

type libp2pInfoWrapper struct {
	peerID peer.ID
	net.Addr
}

// LocalAddr returns the local address
func (c *streamWrapper) LocalAddr() net.Addr {
	addr, err := manet.ToNetAddr(c.Stream.Conn().LocalMultiaddr())
	if err != nil {
		return fakeLocalAddr()
	}
	return &libp2pInfoWrapper{Addr: addr, peerID: c.Stream.Conn().LocalPeer()}
}

// RemoteAddr returns the remote address
func (c *streamWrapper) RemoteAddr() net.Addr {
	addr, err := manet.ToNetAddr(c.Stream.Conn().RemoteMultiaddr())
	if err != nil {
		return fakeRemoteAddr()
	}
	return &libp2pInfoWrapper{Addr: addr, peerID: c.Stream.Conn().RemotePeer()}
}

var _ net.Conn = &streamWrapper{}

// fakeLocalAddr returns a dummy local address
func fakeLocalAddr() net.Addr {
	localIP := net.ParseIP("127.0.0.1")
	return &net.TCPAddr{IP: localIP, Port: 0}
}

// fakeRemoteAddr returns a dummy remote address
func fakeRemoteAddr() net.Addr {
	remoteIP := net.ParseIP("127.1.0.1")
	return &net.TCPAddr{IP: remoteIP, Port: 0}
}
