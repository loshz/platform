// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: proto/v1/service.proto

package pbv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PlatformService_Status_FullMethodName = "/service.v1.PlatformService/Status"
)

// PlatformServiceClient is the client API for PlatformService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlatformServiceClient interface {
	// Status...
	Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PlatformServiceStatusResponse, error)
}

type platformServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPlatformServiceClient(cc grpc.ClientConnInterface) PlatformServiceClient {
	return &platformServiceClient{cc}
}

func (c *platformServiceClient) Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PlatformServiceStatusResponse, error) {
	out := new(PlatformServiceStatusResponse)
	err := c.cc.Invoke(ctx, PlatformService_Status_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlatformServiceServer is the server API for PlatformService service.
// All implementations must embed UnimplementedPlatformServiceServer
// for forward compatibility
type PlatformServiceServer interface {
	// Status...
	Status(context.Context, *emptypb.Empty) (*PlatformServiceStatusResponse, error)
	mustEmbedUnimplementedPlatformServiceServer()
}

// UnimplementedPlatformServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPlatformServiceServer struct {
}

func (UnimplementedPlatformServiceServer) Status(context.Context, *emptypb.Empty) (*PlatformServiceStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedPlatformServiceServer) mustEmbedUnimplementedPlatformServiceServer() {}

// UnsafePlatformServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlatformServiceServer will
// result in compilation errors.
type UnsafePlatformServiceServer interface {
	mustEmbedUnimplementedPlatformServiceServer()
}

func RegisterPlatformServiceServer(s grpc.ServiceRegistrar, srv PlatformServiceServer) {
	s.RegisterService(&PlatformService_ServiceDesc, srv)
}

func _PlatformService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlatformServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PlatformService_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlatformServiceServer).Status(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// PlatformService_ServiceDesc is the grpc.ServiceDesc for PlatformService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlatformService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.v1.PlatformService",
	HandlerType: (*PlatformServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _PlatformService_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/service.proto",
}
