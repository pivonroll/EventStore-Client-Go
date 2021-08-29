// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gossip

import (
	context "context"
	shared "github.com/EventStore/EventStore-Client-Go/protos/shared"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GossipClient is the client API for Gossip service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GossipClient interface {
	Read(ctx context.Context, in *shared.Empty, opts ...grpc.CallOption) (*ClusterInfo, error)
}

type gossipClient struct {
	cc grpc.ClientConnInterface
}

func NewGossipClient(cc grpc.ClientConnInterface) GossipClient {
	return &gossipClient{cc}
}

func (c *gossipClient) Read(ctx context.Context, in *shared.Empty, opts ...grpc.CallOption) (*ClusterInfo, error) {
	out := new(ClusterInfo)
	err := c.cc.Invoke(ctx, "/event_store.client.gossip.Gossip/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GossipServer is the server API for Gossip service.
// All implementations must embed UnimplementedGossipServer
// for forward compatibility
type GossipServer interface {
	Read(context.Context, *shared.Empty) (*ClusterInfo, error)
	mustEmbedUnimplementedGossipServer()
}

// UnimplementedGossipServer must be embedded to have forward compatible implementations.
type UnimplementedGossipServer struct {
}

func (UnimplementedGossipServer) Read(context.Context, *shared.Empty) (*ClusterInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedGossipServer) mustEmbedUnimplementedGossipServer() {}

// UnsafeGossipServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GossipServer will
// result in compilation errors.
type UnsafeGossipServer interface {
	mustEmbedUnimplementedGossipServer()
}

func RegisterGossipServer(s grpc.ServiceRegistrar, srv GossipServer) {
	s.RegisterService(&Gossip_ServiceDesc, srv)
}

func _Gossip_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(shared.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GossipServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event_store.client.gossip.Gossip/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GossipServer).Read(ctx, req.(*shared.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Gossip_ServiceDesc is the grpc.ServiceDesc for Gossip service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gossip_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event_store.client.gossip.Gossip",
	HandlerType: (*GossipServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Read",
			Handler:    _Gossip_Read_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gossip.proto",
}
