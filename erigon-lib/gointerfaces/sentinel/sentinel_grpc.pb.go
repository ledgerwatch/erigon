// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.2
// source: p2psentinel/sentinel.proto

package sentinel

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
	Sentinel_SubscribeGossip_FullMethodName = "/sentinel.Sentinel/SubscribeGossip"
	Sentinel_SendRequest_FullMethodName     = "/sentinel.Sentinel/SendRequest"
	Sentinel_SetStatus_FullMethodName       = "/sentinel.Sentinel/SetStatus"
	Sentinel_GetPeers_FullMethodName        = "/sentinel.Sentinel/GetPeers"
	Sentinel_BanPeer_FullMethodName         = "/sentinel.Sentinel/BanPeer"
	Sentinel_UnbanPeer_FullMethodName       = "/sentinel.Sentinel/UnbanPeer"
	Sentinel_PenalizePeer_FullMethodName    = "/sentinel.Sentinel/PenalizePeer"
	Sentinel_RewardPeer_FullMethodName      = "/sentinel.Sentinel/RewardPeer"
	Sentinel_PublishGossip_FullMethodName   = "/sentinel.Sentinel/PublishGossip"
	Sentinel_Identity_FullMethodName        = "/sentinel.Sentinel/Identity"
	Sentinel_PeersInfo_FullMethodName       = "/sentinel.Sentinel/PeersInfo"
)

// SentinelClient is the client API for Sentinel service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SentinelClient interface {
	SubscribeGossip(ctx context.Context, in *SubscriptionData, opts ...grpc.CallOption) (Sentinel_SubscribeGossipClient, error)
	SendRequest(ctx context.Context, in *RequestData, opts ...grpc.CallOption) (*ResponseData, error)
	SetStatus(ctx context.Context, in *Status, opts ...grpc.CallOption) (*EmptyMessage, error)
	GetPeers(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*PeerCount, error)
	BanPeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*EmptyMessage, error)
	UnbanPeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*EmptyMessage, error)
	PenalizePeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*EmptyMessage, error)
	RewardPeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*EmptyMessage, error)
	PublishGossip(ctx context.Context, in *GossipData, opts ...grpc.CallOption) (*EmptyMessage, error)
	Identity(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*IdentityResponse, error)
	PeersInfo(ctx context.Context, in *PeersInfoRequest, opts ...grpc.CallOption) (*PeersInfoResponse, error)
}

type sentinelClient struct {
	cc grpc.ClientConnInterface
}

func NewSentinelClient(cc grpc.ClientConnInterface) SentinelClient {
	return &sentinelClient{cc}
}

func (c *sentinelClient) SubscribeGossip(ctx context.Context, in *SubscriptionData, opts ...grpc.CallOption) (Sentinel_SubscribeGossipClient, error) {
	stream, err := c.cc.NewStream(ctx, &Sentinel_ServiceDesc.Streams[0], Sentinel_SubscribeGossip_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &sentinelSubscribeGossipClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Sentinel_SubscribeGossipClient interface {
	Recv() (*GossipData, error)
	grpc.ClientStream
}

type sentinelSubscribeGossipClient struct {
	grpc.ClientStream
}

func (x *sentinelSubscribeGossipClient) Recv() (*GossipData, error) {
	m := new(GossipData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sentinelClient) SendRequest(ctx context.Context, in *RequestData, opts ...grpc.CallOption) (*ResponseData, error) {
	out := new(ResponseData)
	err := c.cc.Invoke(ctx, Sentinel_SendRequest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentinelClient) SetStatus(ctx context.Context, in *Status, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, Sentinel_SetStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentinelClient) GetPeers(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*PeerCount, error) {
	out := new(PeerCount)
	err := c.cc.Invoke(ctx, Sentinel_GetPeers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentinelClient) BanPeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, Sentinel_BanPeer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentinelClient) UnbanPeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, Sentinel_UnbanPeer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentinelClient) PenalizePeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, Sentinel_PenalizePeer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentinelClient) RewardPeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, Sentinel_RewardPeer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentinelClient) PublishGossip(ctx context.Context, in *GossipData, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, Sentinel_PublishGossip_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentinelClient) Identity(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*IdentityResponse, error) {
	out := new(IdentityResponse)
	err := c.cc.Invoke(ctx, Sentinel_Identity_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sentinelClient) PeersInfo(ctx context.Context, in *PeersInfoRequest, opts ...grpc.CallOption) (*PeersInfoResponse, error) {
	out := new(PeersInfoResponse)
	err := c.cc.Invoke(ctx, Sentinel_PeersInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SentinelServer is the server API for Sentinel service.
// All implementations must embed UnimplementedSentinelServer
// for forward compatibility
type SentinelServer interface {
	SubscribeGossip(*SubscriptionData, Sentinel_SubscribeGossipServer) error
	SendRequest(context.Context, *RequestData) (*ResponseData, error)
	SetStatus(context.Context, *Status) (*EmptyMessage, error)
	GetPeers(context.Context, *EmptyMessage) (*PeerCount, error)
	BanPeer(context.Context, *Peer) (*EmptyMessage, error)
	UnbanPeer(context.Context, *Peer) (*EmptyMessage, error)
	PenalizePeer(context.Context, *Peer) (*EmptyMessage, error)
	RewardPeer(context.Context, *Peer) (*EmptyMessage, error)
	PublishGossip(context.Context, *GossipData) (*EmptyMessage, error)
	Identity(context.Context, *EmptyMessage) (*IdentityResponse, error)
	PeersInfo(context.Context, *PeersInfoRequest) (*PeersInfoResponse, error)
	mustEmbedUnimplementedSentinelServer()
}

// UnimplementedSentinelServer must be embedded to have forward compatible implementations.
type UnimplementedSentinelServer struct {
}

func (UnimplementedSentinelServer) SubscribeGossip(*SubscriptionData, Sentinel_SubscribeGossipServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeGossip not implemented")
}
func (UnimplementedSentinelServer) SendRequest(context.Context, *RequestData) (*ResponseData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRequest not implemented")
}
func (UnimplementedSentinelServer) SetStatus(context.Context, *Status) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStatus not implemented")
}
func (UnimplementedSentinelServer) GetPeers(context.Context, *EmptyMessage) (*PeerCount, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPeers not implemented")
}
func (UnimplementedSentinelServer) BanPeer(context.Context, *Peer) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BanPeer not implemented")
}
func (UnimplementedSentinelServer) UnbanPeer(context.Context, *Peer) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnbanPeer not implemented")
}
func (UnimplementedSentinelServer) PenalizePeer(context.Context, *Peer) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PenalizePeer not implemented")
}
func (UnimplementedSentinelServer) RewardPeer(context.Context, *Peer) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RewardPeer not implemented")
}
func (UnimplementedSentinelServer) PublishGossip(context.Context, *GossipData) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishGossip not implemented")
}
func (UnimplementedSentinelServer) Identity(context.Context, *EmptyMessage) (*IdentityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Identity not implemented")
}
func (UnimplementedSentinelServer) PeersInfo(context.Context, *PeersInfoRequest) (*PeersInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PeersInfo not implemented")
}
func (UnimplementedSentinelServer) mustEmbedUnimplementedSentinelServer() {}

// UnsafeSentinelServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SentinelServer will
// result in compilation errors.
type UnsafeSentinelServer interface {
	mustEmbedUnimplementedSentinelServer()
}

func RegisterSentinelServer(s grpc.ServiceRegistrar, srv SentinelServer) {
	s.RegisterService(&Sentinel_ServiceDesc, srv)
}

func _Sentinel_SubscribeGossip_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscriptionData)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SentinelServer).SubscribeGossip(m, &sentinelSubscribeGossipServer{stream})
}

type Sentinel_SubscribeGossipServer interface {
	Send(*GossipData) error
	grpc.ServerStream
}

type sentinelSubscribeGossipServer struct {
	grpc.ServerStream
}

func (x *sentinelSubscribeGossipServer) Send(m *GossipData) error {
	return x.ServerStream.SendMsg(m)
}

func _Sentinel_SendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).SendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_SendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).SendRequest(ctx, req.(*RequestData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sentinel_SetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Status)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).SetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_SetStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).SetStatus(ctx, req.(*Status))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sentinel_GetPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).GetPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_GetPeers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).GetPeers(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sentinel_BanPeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).BanPeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_BanPeer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).BanPeer(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sentinel_UnbanPeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).UnbanPeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_UnbanPeer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).UnbanPeer(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sentinel_PenalizePeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).PenalizePeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_PenalizePeer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).PenalizePeer(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sentinel_RewardPeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).RewardPeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_RewardPeer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).RewardPeer(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sentinel_PublishGossip_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GossipData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).PublishGossip(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_PublishGossip_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).PublishGossip(ctx, req.(*GossipData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sentinel_Identity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).Identity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_Identity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).Identity(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sentinel_PeersInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PeersInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SentinelServer).PeersInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sentinel_PeersInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SentinelServer).PeersInfo(ctx, req.(*PeersInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Sentinel_ServiceDesc is the grpc.ServiceDesc for Sentinel service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sentinel_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sentinel.Sentinel",
	HandlerType: (*SentinelServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendRequest",
			Handler:    _Sentinel_SendRequest_Handler,
		},
		{
			MethodName: "SetStatus",
			Handler:    _Sentinel_SetStatus_Handler,
		},
		{
			MethodName: "GetPeers",
			Handler:    _Sentinel_GetPeers_Handler,
		},
		{
			MethodName: "BanPeer",
			Handler:    _Sentinel_BanPeer_Handler,
		},
		{
			MethodName: "UnbanPeer",
			Handler:    _Sentinel_UnbanPeer_Handler,
		},
		{
			MethodName: "PenalizePeer",
			Handler:    _Sentinel_PenalizePeer_Handler,
		},
		{
			MethodName: "RewardPeer",
			Handler:    _Sentinel_RewardPeer_Handler,
		},
		{
			MethodName: "PublishGossip",
			Handler:    _Sentinel_PublishGossip_Handler,
		},
		{
			MethodName: "Identity",
			Handler:    _Sentinel_Identity_Handler,
		},
		{
			MethodName: "PeersInfo",
			Handler:    _Sentinel_PeersInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeGossip",
			Handler:       _Sentinel_SubscribeGossip_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "p2psentinel/sentinel.proto",
}
