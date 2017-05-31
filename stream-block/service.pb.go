// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

/*
Package streamblock is a generated protocol buffer package.

It is generated from these files:
	service.proto

It has these top-level messages:
	ObserveRequest
	ObserveResponse
*/
package streamblock

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ObserveRequest struct {
	Data string `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *ObserveRequest) Reset()                    { *m = ObserveRequest{} }
func (m *ObserveRequest) String() string            { return proto.CompactTextString(m) }
func (*ObserveRequest) ProtoMessage()               {}
func (*ObserveRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ObserveRequest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type ObserveResponse struct {
	Data string `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *ObserveResponse) Reset()                    { *m = ObserveResponse{} }
func (m *ObserveResponse) String() string            { return proto.CompactTextString(m) }
func (*ObserveResponse) ProtoMessage()               {}
func (*ObserveResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ObserveResponse) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*ObserveRequest)(nil), "streamblock.ObserveRequest")
	proto.RegisterType((*ObserveResponse)(nil), "streamblock.ObserveResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Observer service

type ObserverClient interface {
	Observe(ctx context.Context, in *ObserveRequest, opts ...grpc.CallOption) (Observer_ObserveClient, error)
}

type observerClient struct {
	cc *grpc.ClientConn
}

func NewObserverClient(cc *grpc.ClientConn) ObserverClient {
	return &observerClient{cc}
}

func (c *observerClient) Observe(ctx context.Context, in *ObserveRequest, opts ...grpc.CallOption) (Observer_ObserveClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Observer_serviceDesc.Streams[0], c.cc, "/streamblock.Observer/Observe", opts...)
	if err != nil {
		return nil, err
	}
	x := &observerObserveClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Observer_ObserveClient interface {
	Recv() (*ObserveResponse, error)
	grpc.ClientStream
}

type observerObserveClient struct {
	grpc.ClientStream
}

func (x *observerObserveClient) Recv() (*ObserveResponse, error) {
	m := new(ObserveResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Observer service

type ObserverServer interface {
	Observe(*ObserveRequest, Observer_ObserveServer) error
}

func RegisterObserverServer(s *grpc.Server, srv ObserverServer) {
	s.RegisterService(&_Observer_serviceDesc, srv)
}

func _Observer_Observe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ObserveRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ObserverServer).Observe(m, &observerObserveServer{stream})
}

type Observer_ObserveServer interface {
	Send(*ObserveResponse) error
	grpc.ServerStream
}

type observerObserveServer struct {
	grpc.ServerStream
}

func (x *observerObserveServer) Send(m *ObserveResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _Observer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "streamblock.Observer",
	HandlerType: (*ObserverServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Observe",
			Handler:       _Observer_Observe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2e, 0x2e, 0x29, 0x4a, 0x4d,
	0xcc, 0x4d, 0xca, 0xc9, 0x4f, 0xce, 0x56, 0x52, 0xe1, 0xe2, 0xf3, 0x4f, 0x02, 0xc9, 0xa7, 0x06,
	0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x09, 0x71, 0xb1, 0xa4, 0x24, 0x96, 0x24, 0x4a, 0x30,
	0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x4a, 0xaa, 0x5c, 0xfc, 0x70, 0x55, 0xc5, 0x05, 0xf9,
	0x79, 0xc5, 0xa9, 0xd8, 0x94, 0x19, 0x85, 0x70, 0x71, 0x40, 0x95, 0x15, 0x09, 0x79, 0x70, 0xb1,
	0x43, 0xd9, 0x42, 0xd2, 0x7a, 0x48, 0x36, 0xea, 0xa1, 0x5a, 0x27, 0x25, 0x83, 0x5d, 0x12, 0x62,
	0x8b, 0x12, 0x83, 0x01, 0x63, 0x12, 0x1b, 0xd8, 0xd9, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xce, 0x28, 0xa5, 0x12, 0xc7, 0x00, 0x00, 0x00,
}