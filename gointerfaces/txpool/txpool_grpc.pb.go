// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package txpool

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

// TxpoolClient is the client API for Txpool service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TxpoolClient interface {
	FindUnknownTransactions(ctx context.Context, in *TxHashes, opts ...grpc.CallOption) (*TxHashes, error)
	ImportTransactions(ctx context.Context, in *ImportRequest, opts ...grpc.CallOption) (*ImportReply, error)
	GetTransactions(ctx context.Context, in *GetTransactionsRequest, opts ...grpc.CallOption) (*GetTransactionsReply, error)
}

type txpoolClient struct {
	cc grpc.ClientConnInterface
}

func NewTxpoolClient(cc grpc.ClientConnInterface) TxpoolClient {
	return &txpoolClient{cc}
}

func (c *txpoolClient) FindUnknownTransactions(ctx context.Context, in *TxHashes, opts ...grpc.CallOption) (*TxHashes, error) {
	out := new(TxHashes)
	err := c.cc.Invoke(ctx, "/txpool.Txpool/FindUnknownTransactions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *txpoolClient) ImportTransactions(ctx context.Context, in *ImportRequest, opts ...grpc.CallOption) (*ImportReply, error) {
	out := new(ImportReply)
	err := c.cc.Invoke(ctx, "/txpool.Txpool/ImportTransactions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *txpoolClient) GetTransactions(ctx context.Context, in *GetTransactionsRequest, opts ...grpc.CallOption) (*GetTransactionsReply, error) {
	out := new(GetTransactionsReply)
	err := c.cc.Invoke(ctx, "/txpool.Txpool/GetTransactions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TxpoolServer is the server API for Txpool service.
// All implementations must embed UnimplementedTxpoolServer
// for forward compatibility
type TxpoolServer interface {
	FindUnknownTransactions(context.Context, *TxHashes) (*TxHashes, error)
	ImportTransactions(context.Context, *ImportRequest) (*ImportReply, error)
	GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsReply, error)
	mustEmbedUnimplementedTxpoolServer()
}

// UnimplementedTxpoolServer must be embedded to have forward compatible implementations.
type UnimplementedTxpoolServer struct {
}

func (UnimplementedTxpoolServer) FindUnknownTransactions(context.Context, *TxHashes) (*TxHashes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindUnknownTransactions not implemented")
}
func (UnimplementedTxpoolServer) ImportTransactions(context.Context, *ImportRequest) (*ImportReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportTransactions not implemented")
}
func (UnimplementedTxpoolServer) GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactions not implemented")
}
func (UnimplementedTxpoolServer) mustEmbedUnimplementedTxpoolServer() {}

// UnsafeTxpoolServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TxpoolServer will
// result in compilation errors.
type UnsafeTxpoolServer interface {
	mustEmbedUnimplementedTxpoolServer()
}

func RegisterTxpoolServer(s grpc.ServiceRegistrar, srv TxpoolServer) {
	s.RegisterService(&Txpool_ServiceDesc, srv)
}

func _Txpool_FindUnknownTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxHashes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TxpoolServer).FindUnknownTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/txpool.Txpool/FindUnknownTransactions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TxpoolServer).FindUnknownTransactions(ctx, req.(*TxHashes))
	}
	return interceptor(ctx, in, info, handler)
}

func _Txpool_ImportTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TxpoolServer).ImportTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/txpool.Txpool/ImportTransactions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TxpoolServer).ImportTransactions(ctx, req.(*ImportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Txpool_GetTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TxpoolServer).GetTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/txpool.Txpool/GetTransactions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TxpoolServer).GetTransactions(ctx, req.(*GetTransactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Txpool_ServiceDesc is the grpc.ServiceDesc for Txpool service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Txpool_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "txpool.Txpool",
	HandlerType: (*TxpoolServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindUnknownTransactions",
			Handler:    _Txpool_FindUnknownTransactions_Handler,
		},
		{
			MethodName: "ImportTransactions",
			Handler:    _Txpool_ImportTransactions_Handler,
		},
		{
			MethodName: "GetTransactions",
			Handler:    _Txpool_GetTransactions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "txpool/txpool.proto",
}
