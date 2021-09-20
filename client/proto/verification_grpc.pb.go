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

// VerificationServiceClient is the client API for VerificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VerificationServiceClient interface {
	// BeginVerification starts the verification process by issuing a challenge
	// that the recipient needs to solve
	BeginVerification(ctx context.Context, in *VerificationRequest, opts ...grpc.CallOption) (*Challenge, error)
	// FinishVerification finishes the verification process by sending
	// the completed challenge and returning the status of the verification
	FinishVerification(ctx context.Context, in *ChallengeSolution, opts ...grpc.CallOption) (*VerificationResponse, error)
}

type verificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVerificationServiceClient(cc grpc.ClientConnInterface) VerificationServiceClient {
	return &verificationServiceClient{cc}
}

func (c *verificationServiceClient) BeginVerification(ctx context.Context, in *VerificationRequest, opts ...grpc.CallOption) (*Challenge, error) {
	out := new(Challenge)
	err := c.cc.Invoke(ctx, "/VerificationService/BeginVerification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *verificationServiceClient) FinishVerification(ctx context.Context, in *ChallengeSolution, opts ...grpc.CallOption) (*VerificationResponse, error) {
	out := new(VerificationResponse)
	err := c.cc.Invoke(ctx, "/VerificationService/FinishVerification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VerificationServiceServer is the server API for VerificationService service.
// All implementations must embed UnimplementedVerificationServiceServer
// for forward compatibility
type VerificationServiceServer interface {
	// BeginVerification starts the verification process by issuing a challenge
	// that the recipient needs to solve
	BeginVerification(context.Context, *VerificationRequest) (*Challenge, error)
	// FinishVerification finishes the verification process by sending
	// the completed challenge and returning the status of the verification
	FinishVerification(context.Context, *ChallengeSolution) (*VerificationResponse, error)
	mustEmbedUnimplementedVerificationServiceServer()
}

// UnimplementedVerificationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVerificationServiceServer struct {
}

func (UnimplementedVerificationServiceServer) BeginVerification(context.Context, *VerificationRequest) (*Challenge, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginVerification not implemented")
}
func (UnimplementedVerificationServiceServer) FinishVerification(context.Context, *ChallengeSolution) (*VerificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FinishVerification not implemented")
}
func (UnimplementedVerificationServiceServer) mustEmbedUnimplementedVerificationServiceServer() {}

// UnsafeVerificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VerificationServiceServer will
// result in compilation errors.
type UnsafeVerificationServiceServer interface {
	mustEmbedUnimplementedVerificationServiceServer()
}

func RegisterVerificationServiceServer(s grpc.ServiceRegistrar, srv VerificationServiceServer) {
	s.RegisterService(&VerificationService_ServiceDesc, srv)
}

func _VerificationService_BeginVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerificationServiceServer).BeginVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VerificationService/BeginVerification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerificationServiceServer).BeginVerification(ctx, req.(*VerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VerificationService_FinishVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChallengeSolution)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VerificationServiceServer).FinishVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VerificationService/FinishVerification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VerificationServiceServer).FinishVerification(ctx, req.(*ChallengeSolution))
	}
	return interceptor(ctx, in, info, handler)
}

// VerificationService_ServiceDesc is the grpc.ServiceDesc for VerificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VerificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "VerificationService",
	HandlerType: (*VerificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BeginVerification",
			Handler:    _VerificationService_BeginVerification_Handler,
		},
		{
			MethodName: "FinishVerification",
			Handler:    _VerificationService_FinishVerification_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/verification.proto",
}