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

// WorkspaceInfoServiceClient is the client API for WorkspaceInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WorkspaceInfoServiceClient interface {
	GetWorkspaceInfo(ctx context.Context, in *WorkspaceInfoRequest, opts ...grpc.CallOption) (*WorkspaceInfo, error)
	CreateNewWorkspace(ctx context.Context, in *WorkspaceInfo, opts ...grpc.CallOption) (*WorkspaceInfo, error)
}

type workspaceInfoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkspaceInfoServiceClient(cc grpc.ClientConnInterface) WorkspaceInfoServiceClient {
	return &workspaceInfoServiceClient{cc}
}

func (c *workspaceInfoServiceClient) GetWorkspaceInfo(ctx context.Context, in *WorkspaceInfoRequest, opts ...grpc.CallOption) (*WorkspaceInfo, error) {
	out := new(WorkspaceInfo)
	err := c.cc.Invoke(ctx, "/WorkspaceInfoService/GetWorkspaceInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workspaceInfoServiceClient) CreateNewWorkspace(ctx context.Context, in *WorkspaceInfo, opts ...grpc.CallOption) (*WorkspaceInfo, error) {
	out := new(WorkspaceInfo)
	err := c.cc.Invoke(ctx, "/WorkspaceInfoService/CreateNewWorkspace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorkspaceInfoServiceServer is the server API for WorkspaceInfoService service.
// All implementations must embed UnimplementedWorkspaceInfoServiceServer
// for forward compatibility
type WorkspaceInfoServiceServer interface {
	GetWorkspaceInfo(context.Context, *WorkspaceInfoRequest) (*WorkspaceInfo, error)
	CreateNewWorkspace(context.Context, *WorkspaceInfo) (*WorkspaceInfo, error)
	mustEmbedUnimplementedWorkspaceInfoServiceServer()
}

// UnimplementedWorkspaceInfoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWorkspaceInfoServiceServer struct {
}

func (UnimplementedWorkspaceInfoServiceServer) GetWorkspaceInfo(context.Context, *WorkspaceInfoRequest) (*WorkspaceInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWorkspaceInfo not implemented")
}
func (UnimplementedWorkspaceInfoServiceServer) CreateNewWorkspace(context.Context, *WorkspaceInfo) (*WorkspaceInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNewWorkspace not implemented")
}
func (UnimplementedWorkspaceInfoServiceServer) mustEmbedUnimplementedWorkspaceInfoServiceServer() {}

// UnsafeWorkspaceInfoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WorkspaceInfoServiceServer will
// result in compilation errors.
type UnsafeWorkspaceInfoServiceServer interface {
	mustEmbedUnimplementedWorkspaceInfoServiceServer()
}

func RegisterWorkspaceInfoServiceServer(s grpc.ServiceRegistrar, srv WorkspaceInfoServiceServer) {
	s.RegisterService(&WorkspaceInfoService_ServiceDesc, srv)
}

func _WorkspaceInfoService_GetWorkspaceInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkspaceInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkspaceInfoServiceServer).GetWorkspaceInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkspaceInfoService/GetWorkspaceInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkspaceInfoServiceServer).GetWorkspaceInfo(ctx, req.(*WorkspaceInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkspaceInfoService_CreateNewWorkspace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkspaceInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkspaceInfoServiceServer).CreateNewWorkspace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkspaceInfoService/CreateNewWorkspace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkspaceInfoServiceServer).CreateNewWorkspace(ctx, req.(*WorkspaceInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// WorkspaceInfoService_ServiceDesc is the grpc.ServiceDesc for WorkspaceInfoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WorkspaceInfoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "WorkspaceInfoService",
	HandlerType: (*WorkspaceInfoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWorkspaceInfo",
			Handler:    _WorkspaceInfoService_GetWorkspaceInfo_Handler,
		},
		{
			MethodName: "CreateNewWorkspace",
			Handler:    _WorkspaceInfoService_CreateNewWorkspace_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/rendezvous.proto",
}