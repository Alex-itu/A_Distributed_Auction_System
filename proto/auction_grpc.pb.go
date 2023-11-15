// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: proto/auction.proto

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

const (
	AuctionService_Bid_FullMethodName              = "/proto.AuctionService/Bid"
	AuctionService_Result_FullMethodName           = "/proto.AuctionService/Result"
	AuctionService_ConnectionStream_FullMethodName = "/proto.AuctionService/connectionStream"
)

// AuctionServiceClient is the client API for AuctionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuctionServiceClient interface {
	Bid(ctx context.Context, in *BidAmount, opts ...grpc.CallOption) (*Ack, error)
	Result(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Outcome, error)
	ConnectionStream(ctx context.Context, opts ...grpc.CallOption) (AuctionService_ConnectionStreamClient, error)
}

type auctionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuctionServiceClient(cc grpc.ClientConnInterface) AuctionServiceClient {
	return &auctionServiceClient{cc}
}

func (c *auctionServiceClient) Bid(ctx context.Context, in *BidAmount, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, AuctionService_Bid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) Result(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Outcome, error) {
	out := new(Outcome)
	err := c.cc.Invoke(ctx, AuctionService_Result_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) ConnectionStream(ctx context.Context, opts ...grpc.CallOption) (AuctionService_ConnectionStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &AuctionService_ServiceDesc.Streams[0], AuctionService_ConnectionStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &auctionServiceConnectionStreamClient{stream}
	return x, nil
}

type AuctionService_ConnectionStreamClient interface {
	Send(*BackupStream) error
	Recv() (*BackupStream, error)
	grpc.ClientStream
}

type auctionServiceConnectionStreamClient struct {
	grpc.ClientStream
}

func (x *auctionServiceConnectionStreamClient) Send(m *BackupStream) error {
	return x.ClientStream.SendMsg(m)
}

func (x *auctionServiceConnectionStreamClient) Recv() (*BackupStream, error) {
	m := new(BackupStream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AuctionServiceServer is the server API for AuctionService service.
// All implementations must embed UnimplementedAuctionServiceServer
// for forward compatibility
type AuctionServiceServer interface {
	Bid(context.Context, *BidAmount) (*Ack, error)
	Result(context.Context, *Void) (*Outcome, error)
	ConnectionStream(AuctionService_ConnectionStreamServer) error
	mustEmbedUnimplementedAuctionServiceServer()
}

// UnimplementedAuctionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuctionServiceServer struct {
}

func (UnimplementedAuctionServiceServer) Bid(context.Context, *BidAmount) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bid not implemented")
}
func (UnimplementedAuctionServiceServer) Result(context.Context, *Void) (*Outcome, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedAuctionServiceServer) ConnectionStream(AuctionService_ConnectionStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ConnectionStream not implemented")
}
func (UnimplementedAuctionServiceServer) mustEmbedUnimplementedAuctionServiceServer() {}

// UnsafeAuctionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuctionServiceServer will
// result in compilation errors.
type UnsafeAuctionServiceServer interface {
	mustEmbedUnimplementedAuctionServiceServer()
}

func RegisterAuctionServiceServer(s grpc.ServiceRegistrar, srv AuctionServiceServer) {
	s.RegisterService(&AuctionService_ServiceDesc, srv)
}

func _AuctionService_Bid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BidAmount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).Bid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuctionService_Bid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).Bid(ctx, req.(*BidAmount))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_Result_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).Result(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuctionService_Result_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).Result(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_ConnectionStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AuctionServiceServer).ConnectionStream(&auctionServiceConnectionStreamServer{stream})
}

type AuctionService_ConnectionStreamServer interface {
	Send(*BackupStream) error
	Recv() (*BackupStream, error)
	grpc.ServerStream
}

type auctionServiceConnectionStreamServer struct {
	grpc.ServerStream
}

func (x *auctionServiceConnectionStreamServer) Send(m *BackupStream) error {
	return x.ServerStream.SendMsg(m)
}

func (x *auctionServiceConnectionStreamServer) Recv() (*BackupStream, error) {
	m := new(BackupStream)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AuctionService_ServiceDesc is the grpc.ServiceDesc for AuctionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuctionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.AuctionService",
	HandlerType: (*AuctionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bid",
			Handler:    _AuctionService_Bid_Handler,
		},
		{
			MethodName: "Result",
			Handler:    _AuctionService_Result_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "connectionStream",
			Handler:       _AuctionService_ConnectionStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/auction.proto",
}