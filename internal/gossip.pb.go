// Code generated by protoc-gen-go.
// source: gossip.proto
// DO NOT EDIT!

/*
Package internal is a generated protocol buffer package.

It is generated from these files:
	gossip.proto

It has these top-level messages:
	Peer
	JoinRequest
	JoinResponse
	UnjoinRequest
	UnjoinResponse
	PingRequest
	PingResponse
	WriteRequest
	WriteResponse
	EraseRequest
	EraseResponse
*/
package internal

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

type Peer struct {
	Key  string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Addr string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
}

func (m *Peer) Reset()                    { *m = Peer{} }
func (m *Peer) String() string            { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()               {}
func (*Peer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type JoinRequest struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Laddr string `protobuf:"bytes,2,opt,name=laddr" json:"laddr,omitempty"`
}

func (m *JoinRequest) Reset()                    { *m = JoinRequest{} }
func (m *JoinRequest) String() string            { return proto.CompactTextString(m) }
func (*JoinRequest) ProtoMessage()               {}
func (*JoinRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type JoinResponse struct {
	Key   string  `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Peers []*Peer `protobuf:"bytes,2,rep,name=peers" json:"peers,omitempty"`
}

func (m *JoinResponse) Reset()                    { *m = JoinResponse{} }
func (m *JoinResponse) String() string            { return proto.CompactTextString(m) }
func (*JoinResponse) ProtoMessage()               {}
func (*JoinResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *JoinResponse) GetPeers() []*Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

type UnjoinRequest struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *UnjoinRequest) Reset()                    { *m = UnjoinRequest{} }
func (m *UnjoinRequest) String() string            { return proto.CompactTextString(m) }
func (*UnjoinRequest) ProtoMessage()               {}
func (*UnjoinRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type UnjoinResponse struct {
}

func (m *UnjoinResponse) Reset()                    { *m = UnjoinResponse{} }
func (m *UnjoinResponse) String() string            { return proto.CompactTextString(m) }
func (*UnjoinResponse) ProtoMessage()               {}
func (*UnjoinResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type PingRequest struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Laddr string `protobuf:"bytes,2,opt,name=laddr" json:"laddr,omitempty"`
}

func (m *PingRequest) Reset()                    { *m = PingRequest{} }
func (m *PingRequest) String() string            { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()               {}
func (*PingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type PingResponse struct {
}

func (m *PingResponse) Reset()                    { *m = PingResponse{} }
func (m *PingResponse) String() string            { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()               {}
func (*PingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type WriteRequest struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Msg []byte `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (m *WriteRequest) Reset()                    { *m = WriteRequest{} }
func (m *WriteRequest) String() string            { return proto.CompactTextString(m) }
func (*WriteRequest) ProtoMessage()               {}
func (*WriteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type WriteResponse struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *WriteResponse) Reset()                    { *m = WriteResponse{} }
func (m *WriteResponse) String() string            { return proto.CompactTextString(m) }
func (*WriteResponse) ProtoMessage()               {}
func (*WriteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type EraseRequest struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *EraseRequest) Reset()                    { *m = EraseRequest{} }
func (m *EraseRequest) String() string            { return proto.CompactTextString(m) }
func (*EraseRequest) ProtoMessage()               {}
func (*EraseRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type EraseResponse struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *EraseResponse) Reset()                    { *m = EraseResponse{} }
func (m *EraseResponse) String() string            { return proto.CompactTextString(m) }
func (*EraseResponse) ProtoMessage()               {}
func (*EraseResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func init() {
	proto.RegisterType((*Peer)(nil), "internal.Peer")
	proto.RegisterType((*JoinRequest)(nil), "internal.JoinRequest")
	proto.RegisterType((*JoinResponse)(nil), "internal.JoinResponse")
	proto.RegisterType((*UnjoinRequest)(nil), "internal.UnjoinRequest")
	proto.RegisterType((*UnjoinResponse)(nil), "internal.UnjoinResponse")
	proto.RegisterType((*PingRequest)(nil), "internal.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "internal.PingResponse")
	proto.RegisterType((*WriteRequest)(nil), "internal.WriteRequest")
	proto.RegisterType((*WriteResponse)(nil), "internal.WriteResponse")
	proto.RegisterType((*EraseRequest)(nil), "internal.EraseRequest")
	proto.RegisterType((*EraseResponse)(nil), "internal.EraseResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Gossip service

type GossipClient interface {
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error)
	Unjoin(ctx context.Context, in *UnjoinRequest, opts ...grpc.CallOption) (*UnjoinResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	Erase(ctx context.Context, in *EraseRequest, opts ...grpc.CallOption) (*EraseResponse, error)
}

type gossipClient struct {
	cc *grpc.ClientConn
}

func NewGossipClient(cc *grpc.ClientConn) GossipClient {
	return &gossipClient{cc}
}

func (c *gossipClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error) {
	out := new(JoinResponse)
	err := grpc.Invoke(ctx, "/internal.Gossip/Join", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gossipClient) Unjoin(ctx context.Context, in *UnjoinRequest, opts ...grpc.CallOption) (*UnjoinResponse, error) {
	out := new(UnjoinResponse)
	err := grpc.Invoke(ctx, "/internal.Gossip/Unjoin", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gossipClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := grpc.Invoke(ctx, "/internal.Gossip/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gossipClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := grpc.Invoke(ctx, "/internal.Gossip/Write", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gossipClient) Erase(ctx context.Context, in *EraseRequest, opts ...grpc.CallOption) (*EraseResponse, error) {
	out := new(EraseResponse)
	err := grpc.Invoke(ctx, "/internal.Gossip/Erase", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Gossip service

type GossipServer interface {
	Join(context.Context, *JoinRequest) (*JoinResponse, error)
	Unjoin(context.Context, *UnjoinRequest) (*UnjoinResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Write(context.Context, *WriteRequest) (*WriteResponse, error)
	Erase(context.Context, *EraseRequest) (*EraseResponse, error)
}

func RegisterGossipServer(s *grpc.Server, srv GossipServer) {
	s.RegisterService(&_Gossip_serviceDesc, srv)
}

func _Gossip_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(GossipServer).Join(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Gossip_Unjoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(UnjoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(GossipServer).Unjoin(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Gossip_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(GossipServer).Ping(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Gossip_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(GossipServer).Write(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Gossip_Erase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(EraseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(GossipServer).Erase(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Gossip_serviceDesc = grpc.ServiceDesc{
	ServiceName: "internal.Gossip",
	HandlerType: (*GossipServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _Gossip_Join_Handler,
		},
		{
			MethodName: "Unjoin",
			Handler:    _Gossip_Unjoin_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Gossip_Ping_Handler,
		},
		{
			MethodName: "Write",
			Handler:    _Gossip_Write_Handler,
		},
		{
			MethodName: "Erase",
			Handler:    _Gossip_Erase_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 304 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0xbb, 0xad, 0x2d, 0xfa, 0x9a, 0x15, 0x09, 0xe8, 0x4a, 0x55, 0xd0, 0x9c, 0xea, 0xa5,
	0x87, 0x79, 0x10, 0x06, 0x1e, 0x45, 0xf0, 0x34, 0x04, 0xf1, 0x5c, 0x59, 0x28, 0xd5, 0x99, 0xd4,
	0xa4, 0x1e, 0xfc, 0xb3, 0xfc, 0x0f, 0x6d, 0xd2, 0xce, 0xbd, 0xb8, 0x0a, 0x3b, 0x26, 0xef, 0xfb,
	0xf2, 0x7d, 0xef, 0xd7, 0x02, 0x29, 0xa5, 0xd6, 0x55, 0x9d, 0xd7, 0x4a, 0x36, 0x92, 0x1e, 0x54,
	0xa2, 0xe1, 0x4a, 0x14, 0x6b, 0x76, 0x09, 0xfe, 0x92, 0x73, 0x45, 0x23, 0x98, 0xbc, 0xf1, 0xaf,
	0x64, 0x74, 0x31, 0xca, 0x0e, 0x29, 0x01, 0xbf, 0x58, 0xad, 0x54, 0x32, 0x36, 0x27, 0x76, 0x05,
	0xd1, 0x83, 0xac, 0xc4, 0x23, 0xff, 0xf8, 0xe4, 0xba, 0x71, 0x95, 0x53, 0x08, 0xd6, 0x48, 0xba,
	0x00, 0xd2, 0x49, 0x75, 0x2d, 0x85, 0xe6, 0xae, 0xf6, 0x1c, 0x82, 0xba, 0x8d, 0xd2, 0xad, 0x76,
	0x92, 0x45, 0xf3, 0x38, 0xdf, 0x94, 0xc8, 0x4d, 0x03, 0x76, 0x06, 0xd3, 0x27, 0xf1, 0xfa, 0x4f,
	0x10, 0x3b, 0x82, 0x78, 0x33, 0xed, 0xde, 0x36, 0xb5, 0x96, 0x95, 0x28, 0xf7, 0xa9, 0x15, 0x03,
	0xe9, 0xa4, 0xbd, 0x35, 0x03, 0xf2, 0xac, 0xaa, 0x86, 0x0f, 0x7a, 0xdb, 0xc3, 0xbb, 0x2e, 0xad,
	0x93, 0x98, 0x52, 0xbd, 0x72, 0x60, 0x23, 0x76, 0x0a, 0xe4, 0x4e, 0x15, 0x7a, 0xf0, 0x1d, 0x63,
	0xed, 0x87, 0x03, 0xd6, 0xf9, 0xf7, 0x18, 0xc2, 0x7b, 0xfb, 0x49, 0xe8, 0x0d, 0xf8, 0x06, 0x1a,
	0x3d, 0xde, 0x02, 0x41, 0xbc, 0xd3, 0x93, 0xbf, 0xd7, 0xfd, 0x12, 0x1e, 0xbd, 0x85, 0xb0, 0x63,
	0x42, 0x67, 0x5b, 0x8d, 0xc3, 0x30, 0x4d, 0x76, 0x07, 0xbf, 0xf6, 0x36, 0xd7, 0x50, 0xc1, 0xb9,
	0x08, 0x28, 0xce, 0x75, 0xe0, 0x79, 0x74, 0x01, 0x81, 0x85, 0x42, 0x91, 0x04, 0xf3, 0x4c, 0x67,
	0x3b, 0xf7, 0xd8, 0x6b, 0xa9, 0x60, 0x2f, 0x66, 0x88, 0xbd, 0x0e, 0x3e, 0xe6, 0xbd, 0x84, 0xf6,
	0xe7, 0xbd, 0xfe, 0x09, 0x00, 0x00, 0xff, 0xff, 0x11, 0x28, 0x50, 0x44, 0xcc, 0x02, 0x00, 0x00,
}