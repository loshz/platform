// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: proto/v1/eventd.proto

package apiv1

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

const (
	EventService_RegisterHost_FullMethodName = "/proto.v1.EventService/RegisterHost"
	EventService_SendEvent_FullMethodName    = "/proto.v1.EventService/SendEvent"
)

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventServiceClient interface {
	RegisterHost(ctx context.Context, in *RegisterHostRequest, opts ...grpc.CallOption) (*RegisterHostResponse, error)
	SendEvent(ctx context.Context, opts ...grpc.CallOption) (EventService_SendEventClient, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) RegisterHost(ctx context.Context, in *RegisterHostRequest, opts ...grpc.CallOption) (*RegisterHostResponse, error) {
	out := new(RegisterHostResponse)
	err := c.cc.Invoke(ctx, EventService_RegisterHost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) SendEvent(ctx context.Context, opts ...grpc.CallOption) (EventService_SendEventClient, error) {
	stream, err := c.cc.NewStream(ctx, &EventService_ServiceDesc.Streams[0], EventService_SendEvent_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &eventServiceSendEventClient{stream}
	return x, nil
}

type EventService_SendEventClient interface {
	Send(*SendEventRequest) error
	CloseAndRecv() (*SendEventResponse, error)
	grpc.ClientStream
}

type eventServiceSendEventClient struct {
	grpc.ClientStream
}

func (x *eventServiceSendEventClient) Send(m *SendEventRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *eventServiceSendEventClient) CloseAndRecv() (*SendEventResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(SendEventResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventServiceServer is the server API for EventService service.
// All implementations must embed UnimplementedEventServiceServer
// for forward compatibility
type EventServiceServer interface {
	RegisterHost(context.Context, *RegisterHostRequest) (*RegisterHostResponse, error)
	SendEvent(EventService_SendEventServer) error
	mustEmbedUnimplementedEventServiceServer()
}

// UnimplementedEventServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEventServiceServer struct {
}

func (UnimplementedEventServiceServer) RegisterHost(context.Context, *RegisterHostRequest) (*RegisterHostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterHost not implemented")
}
func (UnimplementedEventServiceServer) SendEvent(EventService_SendEventServer) error {
	return status.Errorf(codes.Unimplemented, "method SendEvent not implemented")
}
func (UnimplementedEventServiceServer) mustEmbedUnimplementedEventServiceServer() {}

// UnsafeEventServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceServer will
// result in compilation errors.
type UnsafeEventServiceServer interface {
	mustEmbedUnimplementedEventServiceServer()
}

func RegisterEventServiceServer(s grpc.ServiceRegistrar, srv EventServiceServer) {
	s.RegisterService(&EventService_ServiceDesc, srv)
}

func _EventService_RegisterHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).RegisterHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EventService_RegisterHost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).RegisterHost(ctx, req.(*RegisterHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_SendEvent_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventServiceServer).SendEvent(&eventServiceSendEventServer{stream})
}

type EventService_SendEventServer interface {
	SendAndClose(*SendEventResponse) error
	Recv() (*SendEventRequest, error)
	grpc.ServerStream
}

type eventServiceSendEventServer struct {
	grpc.ServerStream
}

func (x *eventServiceSendEventServer) SendAndClose(m *SendEventResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *eventServiceSendEventServer) Recv() (*SendEventRequest, error) {
	m := new(SendEventRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventService_ServiceDesc is the grpc.ServiceDesc for EventService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.v1.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterHost",
			Handler:    _EventService_RegisterHost_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendEvent",
			Handler:       _EventService_SendEvent_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "proto/v1/eventd.proto",
}
