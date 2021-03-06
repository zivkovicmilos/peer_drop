// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FileSharingClient is the client API for FileSharing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileSharingClient interface {
	RequestFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileDownloadMetadata, error)
	DownloadFile(ctx context.Context, in *FileRequestID, opts ...grpc.CallOption) (FileSharing_DownloadFileClient, error)
}

type fileSharingClient struct {
	cc grpc.ClientConnInterface
}

func NewFileSharingClient(cc grpc.ClientConnInterface) FileSharingClient {
	return &fileSharingClient{cc}
}

func (c *fileSharingClient) RequestFile(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileDownloadMetadata, error) {
	out := new(FileDownloadMetadata)
	err := c.cc.Invoke(ctx, "/FileSharing/RequestFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileSharingClient) DownloadFile(ctx context.Context, in *FileRequestID, opts ...grpc.CallOption) (FileSharing_DownloadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileSharing_ServiceDesc.Streams[0], "/FileSharing/DownloadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileSharingDownloadFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileSharing_DownloadFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type fileSharingDownloadFileClient struct {
	grpc.ClientStream
}

func (x *fileSharingDownloadFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileSharingServer is the server API for FileSharing service.
// All implementations must embed UnimplementedFileSharingServer
// for forward compatibility
type FileSharingServer interface {
	RequestFile(context.Context, *FileRequest) (*FileDownloadMetadata, error)
	DownloadFile(*FileRequestID, FileSharing_DownloadFileServer) error
	mustEmbedUnimplementedFileSharingServer()
}

// UnimplementedFileSharingServer must be embedded to have forward compatible implementations.
type UnimplementedFileSharingServer struct {
}

func (UnimplementedFileSharingServer) RequestFile(context.Context, *FileRequest) (*FileDownloadMetadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestFile not implemented")
}
func (UnimplementedFileSharingServer) DownloadFile(*FileRequestID, FileSharing_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedFileSharingServer) mustEmbedUnimplementedFileSharingServer() {}

// UnsafeFileSharingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileSharingServer will
// result in compilation errors.
type UnsafeFileSharingServer interface {
	mustEmbedUnimplementedFileSharingServer()
}

func RegisterFileSharingServer(s grpc.ServiceRegistrar, srv FileSharingServer) {
	s.RegisterService(&FileSharing_ServiceDesc, srv)
}

func _FileSharing_RequestFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileSharingServer).RequestFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FileSharing/RequestFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileSharingServer).RequestFile(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileSharing_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileRequestID)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileSharingServer).DownloadFile(m, &fileSharingDownloadFileServer{stream})
}

type FileSharing_DownloadFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type fileSharingDownloadFileServer struct {
	grpc.ServerStream
}

func (x *fileSharingDownloadFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

// FileSharing_ServiceDesc is the grpc.ServiceDesc for FileSharing service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileSharing_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FileSharing",
	HandlerType: (*FileSharingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestFile",
			Handler:    _FileSharing_RequestFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadFile",
			Handler:       _FileSharing_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/fileSharing.proto",
}
