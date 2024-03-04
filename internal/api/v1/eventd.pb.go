// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: proto/v1/eventd.proto

package apiv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EventType int32

const (
	EventType_EVENT_TYPE_UNSPECIFIED EventType = 0
	EventType_EVENT_TYPE_NETWORK     EventType = 1
)

// Enum value maps for EventType.
var (
	EventType_name = map[int32]string{
		0: "EVENT_TYPE_UNSPECIFIED",
		1: "EVENT_TYPE_NETWORK",
	}
	EventType_value = map[string]int32{
		"EVENT_TYPE_UNSPECIFIED": 0,
		"EVENT_TYPE_NETWORK":     1,
	}
)

func (x EventType) Enum() *EventType {
	p := new(EventType)
	*p = x
	return p
}

func (x EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_v1_eventd_proto_enumTypes[0].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_proto_v1_eventd_proto_enumTypes[0]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_proto_v1_eventd_proto_rawDescGZIP(), []int{0}
}

type NetworkTransport int32

const (
	NetworkTransport_NETWORK_TRANSPORT_UNSPECIFIED NetworkTransport = 0
	NetworkTransport_NETWORK_TRANSPORT_TCP         NetworkTransport = 1
	NetworkTransport_NETWORK_TRANSPORT_UDP         NetworkTransport = 2
)

// Enum value maps for NetworkTransport.
var (
	NetworkTransport_name = map[int32]string{
		0: "NETWORK_TRANSPORT_UNSPECIFIED",
		1: "NETWORK_TRANSPORT_TCP",
		2: "NETWORK_TRANSPORT_UDP",
	}
	NetworkTransport_value = map[string]int32{
		"NETWORK_TRANSPORT_UNSPECIFIED": 0,
		"NETWORK_TRANSPORT_TCP":         1,
		"NETWORK_TRANSPORT_UDP":         2,
	}
)

func (x NetworkTransport) Enum() *NetworkTransport {
	p := new(NetworkTransport)
	*p = x
	return p
}

func (x NetworkTransport) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NetworkTransport) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_v1_eventd_proto_enumTypes[1].Descriptor()
}

func (NetworkTransport) Type() protoreflect.EnumType {
	return &file_proto_v1_eventd_proto_enumTypes[1]
}

func (x NetworkTransport) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NetworkTransport.Descriptor instead.
func (NetworkTransport) EnumDescriptor() ([]byte, []int) {
	return file_proto_v1_eventd_proto_rawDescGZIP(), []int{1}
}

type Host struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MachineId string `protobuf:"bytes,1,opt,name=machine_id,json=machineId,proto3" json:"machine_id,omitempty"`
	Hostname  string `protobuf:"bytes,2,opt,name=hostname,proto3" json:"hostname,omitempty"`
}

func (x *Host) Reset() {
	*x = Host{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_eventd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Host) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Host) ProtoMessage() {}

func (x *Host) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_eventd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Host.ProtoReflect.Descriptor instead.
func (*Host) Descriptor() ([]byte, []int) {
	return file_proto_v1_eventd_proto_rawDescGZIP(), []int{0}
}

func (x *Host) GetMachineId() string {
	if x != nil {
		return x.MachineId
	}
	return ""
}

func (x *Host) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

type RegisterHostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host      *Host `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Timestamp int64 `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *RegisterHostRequest) Reset() {
	*x = RegisterHostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_eventd_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterHostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterHostRequest) ProtoMessage() {}

func (x *RegisterHostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_eventd_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterHostRequest.ProtoReflect.Descriptor instead.
func (*RegisterHostRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_eventd_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterHostRequest) GetHost() *Host {
	if x != nil {
		return x.Host
	}
	return nil
}

func (x *RegisterHostRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type RegisterHostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MachineId string `protobuf:"bytes,1,opt,name=machine_id,json=machineId,proto3" json:"machine_id,omitempty"`
}

func (x *RegisterHostResponse) Reset() {
	*x = RegisterHostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_eventd_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterHostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterHostResponse) ProtoMessage() {}

func (x *RegisterHostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_eventd_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterHostResponse.ProtoReflect.Descriptor instead.
func (*RegisterHostResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_eventd_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterHostResponse) GetMachineId() string {
	if x != nil {
		return x.MachineId
	}
	return ""
}

type SendEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      EventType `protobuf:"varint,1,opt,name=type,proto3,enum=proto.v1.EventType" json:"type,omitempty"`
	MachineId string    `protobuf:"bytes,2,opt,name=machine_id,json=machineId,proto3" json:"machine_id,omitempty"`
	Data      []byte    `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *SendEventRequest) Reset() {
	*x = SendEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_eventd_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendEventRequest) ProtoMessage() {}

func (x *SendEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_eventd_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendEventRequest.ProtoReflect.Descriptor instead.
func (*SendEventRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_eventd_proto_rawDescGZIP(), []int{3}
}

func (x *SendEventRequest) GetType() EventType {
	if x != nil {
		return x.Type
	}
	return EventType_EVENT_TYPE_UNSPECIFIED
}

func (x *SendEventRequest) GetMachineId() string {
	if x != nil {
		return x.MachineId
	}
	return ""
}

func (x *SendEventRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type SendEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventsTotal uint64 `protobuf:"varint,1,opt,name=events_total,json=eventsTotal,proto3" json:"events_total,omitempty"`
}

func (x *SendEventResponse) Reset() {
	*x = SendEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_eventd_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendEventResponse) ProtoMessage() {}

func (x *SendEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_eventd_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendEventResponse.ProtoReflect.Descriptor instead.
func (*SendEventResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_eventd_proto_rawDescGZIP(), []int{4}
}

func (x *SendEventResponse) GetEventsTotal() uint64 {
	if x != nil {
		return x.EventsTotal
	}
	return 0
}

type NetworkEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transport  NetworkTransport `protobuf:"varint,1,opt,name=transport,proto3,enum=proto.v1.NetworkTransport" json:"transport,omitempty"`
	SourceIp   *IP              `protobuf:"bytes,2,opt,name=source_ip,json=sourceIp,proto3" json:"source_ip,omitempty"`
	SourcePort uint32           `protobuf:"varint,3,opt,name=source_port,json=sourcePort,proto3" json:"source_port,omitempty"`
	DestIp     *IP              `protobuf:"bytes,4,opt,name=dest_ip,json=destIp,proto3" json:"dest_ip,omitempty"`
	DestPort   uint32           `protobuf:"varint,5,opt,name=dest_port,json=destPort,proto3" json:"dest_port,omitempty"`
}

func (x *NetworkEvent) Reset() {
	*x = NetworkEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_eventd_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkEvent) ProtoMessage() {}

func (x *NetworkEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_eventd_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkEvent.ProtoReflect.Descriptor instead.
func (*NetworkEvent) Descriptor() ([]byte, []int) {
	return file_proto_v1_eventd_proto_rawDescGZIP(), []int{5}
}

func (x *NetworkEvent) GetTransport() NetworkTransport {
	if x != nil {
		return x.Transport
	}
	return NetworkTransport_NETWORK_TRANSPORT_UNSPECIFIED
}

func (x *NetworkEvent) GetSourceIp() *IP {
	if x != nil {
		return x.SourceIp
	}
	return nil
}

func (x *NetworkEvent) GetSourcePort() uint32 {
	if x != nil {
		return x.SourcePort
	}
	return 0
}

func (x *NetworkEvent) GetDestIp() *IP {
	if x != nil {
		return x.DestIp
	}
	return nil
}

func (x *NetworkEvent) GetDestPort() uint32 {
	if x != nil {
		return x.DestPort
	}
	return 0
}

type IP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Address:
	//
	//	*IP_V4
	//	*IP_V6
	Address isIP_Address `protobuf_oneof:"address"`
}

func (x *IP) Reset() {
	*x = IP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_eventd_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IP) ProtoMessage() {}

func (x *IP) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_eventd_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IP.ProtoReflect.Descriptor instead.
func (*IP) Descriptor() ([]byte, []int) {
	return file_proto_v1_eventd_proto_rawDescGZIP(), []int{6}
}

func (m *IP) GetAddress() isIP_Address {
	if m != nil {
		return m.Address
	}
	return nil
}

func (x *IP) GetV4() uint32 {
	if x, ok := x.GetAddress().(*IP_V4); ok {
		return x.V4
	}
	return 0
}

func (x *IP) GetV6() []byte {
	if x, ok := x.GetAddress().(*IP_V6); ok {
		return x.V6
	}
	return nil
}

type isIP_Address interface {
	isIP_Address()
}

type IP_V4 struct {
	// IPv4 [4]byte
	V4 uint32 `protobuf:"fixed32,1,opt,name=v4,proto3,oneof"`
}

type IP_V6 struct {
	// IPv6 [16]byte
	V6 []byte `protobuf:"bytes,2,opt,name=v6,proto3,oneof"`
}

func (*IP_V4) isIP_Address() {}

func (*IP_V6) isIP_Address() {}

var File_proto_v1_eventd_proto protoreflect.FileDescriptor

var file_proto_v1_eventd_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76,
	0x31, 0x22, 0x41, 0x0a, 0x04, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x61, 0x63,
	0x68, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d,
	0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x57, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x04, 0x68,
	0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12,
	0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x35, 0x0a,
	0x14, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x61, 0x63, 0x68, 0x69,
	0x6e, 0x65, 0x49, 0x64, 0x22, 0x6e, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76,
	0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x36, 0x0a, 0x11, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x5f, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0xd8, 0x01, 0x0a,
	0x0c, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x38, 0x0a,
	0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x09, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x29, 0x0a, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x50, 0x52, 0x08, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x49, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50,
	0x6f, 0x72, 0x74, 0x12, 0x25, 0x0a, 0x07, 0x64, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x70, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e,
	0x49, 0x50, 0x52, 0x06, 0x64, 0x65, 0x73, 0x74, 0x49, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65,
	0x73, 0x74, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x64,
	0x65, 0x73, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x33, 0x0a, 0x02, 0x49, 0x50, 0x12, 0x10, 0x0a,
	0x02, 0x76, 0x34, 0x18, 0x01, 0x20, 0x01, 0x28, 0x07, 0x48, 0x00, 0x52, 0x02, 0x76, 0x34, 0x12,
	0x10, 0x0a, 0x02, 0x76, 0x36, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x02, 0x76,
	0x36, 0x42, 0x09, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2a, 0x3f, 0x0a, 0x09,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x45, 0x56, 0x45,
	0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x4e, 0x45, 0x54, 0x57, 0x4f, 0x52, 0x4b, 0x10, 0x01, 0x2a, 0x6b, 0x0a,
	0x10, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72,
	0x74, 0x12, 0x21, 0x0a, 0x1d, 0x4e, 0x45, 0x54, 0x57, 0x4f, 0x52, 0x4b, 0x5f, 0x54, 0x52, 0x41,
	0x4e, 0x53, 0x50, 0x4f, 0x52, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x19, 0x0a, 0x15, 0x4e, 0x45, 0x54, 0x57, 0x4f, 0x52, 0x4b, 0x5f,
	0x54, 0x52, 0x41, 0x4e, 0x53, 0x50, 0x4f, 0x52, 0x54, 0x5f, 0x54, 0x43, 0x50, 0x10, 0x01, 0x12,
	0x19, 0x0a, 0x15, 0x4e, 0x45, 0x54, 0x57, 0x4f, 0x52, 0x4b, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53,
	0x50, 0x4f, 0x52, 0x54, 0x5f, 0x55, 0x44, 0x50, 0x10, 0x02, 0x32, 0xa9, 0x01, 0x0a, 0x0c, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x0c, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x1d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x48,
	0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x48, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x09,
	0x53, 0x65, 0x6e, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x6f, 0x73, 0x68, 0x7a, 0x2f, 0x70, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_v1_eventd_proto_rawDescOnce sync.Once
	file_proto_v1_eventd_proto_rawDescData = file_proto_v1_eventd_proto_rawDesc
)

func file_proto_v1_eventd_proto_rawDescGZIP() []byte {
	file_proto_v1_eventd_proto_rawDescOnce.Do(func() {
		file_proto_v1_eventd_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_eventd_proto_rawDescData)
	})
	return file_proto_v1_eventd_proto_rawDescData
}

var file_proto_v1_eventd_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_proto_v1_eventd_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_v1_eventd_proto_goTypes = []interface{}{
	(EventType)(0),               // 0: proto.v1.EventType
	(NetworkTransport)(0),        // 1: proto.v1.NetworkTransport
	(*Host)(nil),                 // 2: proto.v1.Host
	(*RegisterHostRequest)(nil),  // 3: proto.v1.RegisterHostRequest
	(*RegisterHostResponse)(nil), // 4: proto.v1.RegisterHostResponse
	(*SendEventRequest)(nil),     // 5: proto.v1.SendEventRequest
	(*SendEventResponse)(nil),    // 6: proto.v1.SendEventResponse
	(*NetworkEvent)(nil),         // 7: proto.v1.NetworkEvent
	(*IP)(nil),                   // 8: proto.v1.IP
}
var file_proto_v1_eventd_proto_depIdxs = []int32{
	2, // 0: proto.v1.RegisterHostRequest.host:type_name -> proto.v1.Host
	0, // 1: proto.v1.SendEventRequest.type:type_name -> proto.v1.EventType
	1, // 2: proto.v1.NetworkEvent.transport:type_name -> proto.v1.NetworkTransport
	8, // 3: proto.v1.NetworkEvent.source_ip:type_name -> proto.v1.IP
	8, // 4: proto.v1.NetworkEvent.dest_ip:type_name -> proto.v1.IP
	3, // 5: proto.v1.EventService.RegisterHost:input_type -> proto.v1.RegisterHostRequest
	5, // 6: proto.v1.EventService.SendEvent:input_type -> proto.v1.SendEventRequest
	4, // 7: proto.v1.EventService.RegisterHost:output_type -> proto.v1.RegisterHostResponse
	6, // 8: proto.v1.EventService.SendEvent:output_type -> proto.v1.SendEventResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_v1_eventd_proto_init() }
func file_proto_v1_eventd_proto_init() {
	if File_proto_v1_eventd_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_eventd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Host); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_v1_eventd_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterHostRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_v1_eventd_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterHostResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_v1_eventd_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendEventRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_v1_eventd_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendEventResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_v1_eventd_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkEvent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_v1_eventd_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IP); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_proto_v1_eventd_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*IP_V4)(nil),
		(*IP_V6)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_v1_eventd_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_eventd_proto_goTypes,
		DependencyIndexes: file_proto_v1_eventd_proto_depIdxs,
		EnumInfos:         file_proto_v1_eventd_proto_enumTypes,
		MessageInfos:      file_proto_v1_eventd_proto_msgTypes,
	}.Build()
	File_proto_v1_eventd_proto = out.File
	file_proto_v1_eventd_proto_rawDesc = nil
	file_proto_v1_eventd_proto_goTypes = nil
	file_proto_v1_eventd_proto_depIdxs = nil
}
