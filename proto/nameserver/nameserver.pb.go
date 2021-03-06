// Code generated by protoc-gen-go.
// source: nameserver.proto
// DO NOT EDIT!

/*
Package nameserver is a generated protocol buffer package.

It is generated from these files:
	nameserver.proto

It has these top-level messages:
	Req
	JoinReq
	JoinResp
	Worker
*/
package nameserver

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
const _ = proto.ProtoPackageIsVersion1

type Req struct {
}

func (m *Req) Reset()                    { *m = Req{} }
func (m *Req) String() string            { return proto.CompactTextString(m) }
func (*Req) ProtoMessage()               {}
func (*Req) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type JoinReq struct {
	WorkerId        uint64  `protobuf:"varint,1,opt,name=WorkerId,json=workerId" json:"WorkerId,omitempty"`
	Host            string  `protobuf:"bytes,2,opt,name=Host,json=host" json:"Host,omitempty"`
	Port            int32   `protobuf:"varint,3,opt,name=Port,json=port" json:"Port,omitempty"`
	CurrConnection  uint32  `protobuf:"varint,4,opt,name=CurrConnection,json=currConnection" json:"CurrConnection,omitempty"`
	CloseConnection uint32  `protobuf:"varint,5,opt,name=CloseConnection,json=closeConnection" json:"CloseConnection,omitempty"`
	CpuUsage        float64 `protobuf:"fixed64,6,opt,name=CpuUsage,json=cpuUsage" json:"CpuUsage,omitempty"`
	Version         string  `protobuf:"bytes,7,opt,name=Version,json=version" json:"Version,omitempty"`
}

func (m *JoinReq) Reset()                    { *m = JoinReq{} }
func (m *JoinReq) String() string            { return proto.CompactTextString(m) }
func (*JoinReq) ProtoMessage()               {}
func (*JoinReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type JoinResp struct {
	Success  bool   `protobuf:"varint,1,opt,name=Success,json=success" json:"Success,omitempty"`
	WorkerId uint64 `protobuf:"varint,2,opt,name=WorkerId,json=workerId" json:"WorkerId,omitempty"`
	ErrMsg   string `protobuf:"bytes,3,opt,name=ErrMsg,json=errMsg" json:"ErrMsg,omitempty"`
}

func (m *JoinResp) Reset()                    { *m = JoinResp{} }
func (m *JoinResp) String() string            { return proto.CompactTextString(m) }
func (*JoinResp) ProtoMessage()               {}
func (*JoinResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Worker struct {
	WorkerId        uint64  `protobuf:"varint,1,opt,name=WorkerId,json=workerId" json:"WorkerId,omitempty"`
	ListenAddr      string  `protobuf:"bytes,2,opt,name=ListenAddr,json=listenAddr" json:"ListenAddr,omitempty"`
	CurrConnection  uint32  `protobuf:"varint,3,opt,name=CurrConnection,json=currConnection" json:"CurrConnection,omitempty"`
	CloseConnection uint32  `protobuf:"varint,4,opt,name=CloseConnection,json=closeConnection" json:"CloseConnection,omitempty"`
	CpuUsage        float64 `protobuf:"fixed64,5,opt,name=CpuUsage,json=cpuUsage" json:"CpuUsage,omitempty"`
	LastAlive       int64   `protobuf:"varint,6,opt,name=LastAlive,json=lastAlive" json:"LastAlive,omitempty"`
	Version         string  `protobuf:"bytes,7,opt,name=Version,json=version" json:"Version,omitempty"`
}

func (m *Worker) Reset()                    { *m = Worker{} }
func (m *Worker) String() string            { return proto.CompactTextString(m) }
func (*Worker) ProtoMessage()               {}
func (*Worker) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*Req)(nil), "nameserver.Req")
	proto.RegisterType((*JoinReq)(nil), "nameserver.JoinReq")
	proto.RegisterType((*JoinResp)(nil), "nameserver.JoinResp")
	proto.RegisterType((*Worker)(nil), "nameserver.Worker")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for NameService service

type NameServiceClient interface {
	// 每5分钟向服务器(NameService)报活
	WorkerJoin(ctx context.Context, in *JoinReq, opts ...grpc.CallOption) (*JoinResp, error)
	ListWorkers(ctx context.Context, in *Req, opts ...grpc.CallOption) (NameService_ListWorkersClient, error)
}

type nameServiceClient struct {
	cc *grpc.ClientConn
}

func NewNameServiceClient(cc *grpc.ClientConn) NameServiceClient {
	return &nameServiceClient{cc}
}

func (c *nameServiceClient) WorkerJoin(ctx context.Context, in *JoinReq, opts ...grpc.CallOption) (*JoinResp, error) {
	out := new(JoinResp)
	err := grpc.Invoke(ctx, "/nameserver.NameService/WorkerJoin", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nameServiceClient) ListWorkers(ctx context.Context, in *Req, opts ...grpc.CallOption) (NameService_ListWorkersClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_NameService_serviceDesc.Streams[0], c.cc, "/nameserver.NameService/ListWorkers", opts...)
	if err != nil {
		return nil, err
	}
	x := &nameServiceListWorkersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NameService_ListWorkersClient interface {
	Recv() (*Worker, error)
	grpc.ClientStream
}

type nameServiceListWorkersClient struct {
	grpc.ClientStream
}

func (x *nameServiceListWorkersClient) Recv() (*Worker, error) {
	m := new(Worker)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for NameService service

type NameServiceServer interface {
	// 每5分钟向服务器(NameService)报活
	WorkerJoin(context.Context, *JoinReq) (*JoinResp, error)
	ListWorkers(*Req, NameService_ListWorkersServer) error
}

func RegisterNameServiceServer(s *grpc.Server, srv NameServiceServer) {
	s.RegisterService(&_NameService_serviceDesc, srv)
}

func _NameService_WorkerJoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(JoinReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(NameServiceServer).WorkerJoin(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _NameService_ListWorkers_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Req)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NameServiceServer).ListWorkers(m, &nameServiceListWorkersServer{stream})
}

type NameService_ListWorkersServer interface {
	Send(*Worker) error
	grpc.ServerStream
}

type nameServiceListWorkersServer struct {
	grpc.ServerStream
}

func (x *nameServiceListWorkersServer) Send(m *Worker) error {
	return x.ServerStream.SendMsg(m)
}

var _NameService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "nameserver.NameService",
	HandlerType: (*NameServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WorkerJoin",
			Handler:    _NameService_WorkerJoin_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListWorkers",
			Handler:       _NameService_ListWorkers_Handler,
			ServerStreams: true,
		},
	},
}

var fileDescriptor0 = []byte{
	// 391 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x92, 0xcd, 0x8a, 0xdb, 0x30,
	0x14, 0x85, 0xa3, 0xf8, 0xff, 0x86, 0xc4, 0x45, 0x2d, 0xc5, 0x84, 0x52, 0x8a, 0x17, 0x25, 0x2b,
	0x53, 0x5a, 0x28, 0x74, 0x99, 0x98, 0x42, 0x5b, 0xd2, 0x12, 0x14, 0xfa, 0xb3, 0x4d, 0x1c, 0x91,
	0x9a, 0xf1, 0x48, 0x1e, 0x49, 0xc9, 0x64, 0x39, 0x2f, 0x39, 0x2f, 0x31, 0x4f, 0x31, 0xb2, 0xe5,
	0xcc, 0x38, 0x99, 0xdf, 0x9d, 0xee, 0xd1, 0xb9, 0x70, 0xcf, 0xc7, 0x81, 0x17, 0x6c, 0x71, 0x4a,
	0x25, 0x15, 0x5b, 0x2a, 0x92, 0x52, 0x70, 0xc5, 0x31, 0xdc, 0x2a, 0xb1, 0x03, 0x16, 0xa1, 0x67,
	0xf1, 0x25, 0x02, 0xef, 0x07, 0xcf, 0x99, 0x7e, 0xe3, 0x21, 0xf8, 0x7f, 0xb9, 0x38, 0xa1, 0xe2,
	0xfb, 0x2a, 0x42, 0xef, 0xd0, 0xc8, 0x26, 0xfe, 0x79, 0x33, 0x63, 0x0c, 0xf6, 0x37, 0x2e, 0x55,
	0xd4, 0xd5, 0x7a, 0x40, 0xec, 0xff, 0xfa, 0x5d, 0x69, 0x33, 0x2e, 0x54, 0x64, 0x69, 0xcd, 0x21,
	0x76, 0xa9, 0xdf, 0xf8, 0x3d, 0x0c, 0xd2, 0x8d, 0x10, 0x29, 0x67, 0x8c, 0x66, 0x2a, 0xe7, 0x2c,
	0xb2, 0xf5, 0x6f, 0x9f, 0x0c, 0xb2, 0x03, 0x15, 0x8f, 0x20, 0x4c, 0x0b, 0x2e, 0x69, 0xcb, 0xe8,
	0xd4, 0xc6, 0x30, 0x3b, 0x94, 0xab, 0xab, 0xd2, 0x72, 0xf3, 0x5b, 0x2e, 0xd6, 0x34, 0x72, 0xb5,
	0x05, 0x11, 0x3f, 0x6b, 0x66, 0x1c, 0x81, 0xf7, 0x87, 0x0a, 0x59, 0x6d, 0x7b, 0xf5, 0x61, 0xde,
	0xd6, 0x8c, 0xf1, 0x3f, 0xf0, 0x4d, 0x2c, 0x59, 0x56, 0xae, 0xf9, 0x26, 0xcb, 0xa8, 0x94, 0x75,
	0x2c, 0x9f, 0x78, 0xd2, 0x8c, 0x07, 0x89, 0xbb, 0x47, 0x89, 0x5f, 0x83, 0xfb, 0x55, 0x88, 0x9f,
	0x72, 0x5d, 0xe7, 0x0b, 0x88, 0x4b, 0xeb, 0x29, 0xbe, 0x42, 0xe0, 0x9a, 0xa5, 0x47, 0x81, 0xbd,
	0x05, 0x98, 0xe6, 0x52, 0x51, 0x36, 0x5e, 0xad, 0x44, 0x83, 0x0d, 0x8a, 0x1b, 0xe5, 0x1e, 0x50,
	0xd6, 0x73, 0x41, 0xd9, 0x4f, 0x83, 0x72, 0x8e, 0x40, 0xbd, 0x81, 0x60, 0xba, 0x90, 0x6a, 0x5c,
	0xe4, 0x5b, 0x43, 0xd1, 0x22, 0x41, 0xb1, 0x17, 0x1e, 0xc6, 0xf8, 0xf1, 0x02, 0x41, 0xef, 0x97,
	0x2e, 0xcd, 0x5c, 0x97, 0x26, 0xcf, 0x28, 0xfe, 0x02, 0x60, 0x12, 0x57, 0x70, 0xf1, 0xcb, 0xa4,
	0x55, 0xb1, 0xa6, 0x45, 0xc3, 0x57, 0x77, 0x45, 0x59, 0xc6, 0x1d, 0xfc, 0x19, 0x7a, 0x15, 0x10,
	0xb3, 0x2e, 0x71, 0xd8, 0xb6, 0x55, 0x7b, 0xb8, 0x2d, 0x18, 0x57, 0xdc, 0xf9, 0x80, 0x26, 0x31,
	0x84, 0x39, 0x4f, 0x14, 0xe7, 0x4c, 0x71, 0xb6, 0x4e, 0x96, 0x7c, 0x37, 0xe9, 0x4f, 0xf8, 0x6e,
	0x7f, 0x15, 0x15, 0x33, 0xb4, 0x74, 0xeb, 0x7e, 0x7f, 0xba, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xea,
	0x2e, 0x0b, 0xaa, 0xf3, 0x02, 0x00, 0x00,
}
