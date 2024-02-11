// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: proto/shortener.proto

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
	Shortener_AddURL_FullMethodName         = "/shortener.Shortener/AddURL"
	Shortener_GetuRLByUserID_FullMethodName = "/shortener.Shortener/GetuRLByUserID"
	Shortener_DeleteURLs_FullMethodName     = "/shortener.Shortener/DeleteURLs"
	Shortener_Batch_FullMethodName          = "/shortener.Shortener/Batch"
	Shortener_GetOriginalURL_FullMethodName = "/shortener.Shortener/GetOriginalURL"
	Shortener_GetStats_FullMethodName       = "/shortener.Shortener/GetStats"
	Shortener_Ping_FullMethodName           = "/shortener.Shortener/Ping"
)

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	AddURL(ctx context.Context, in *AddURLRequest, opts ...grpc.CallOption) (*AddURLResponse, error)
	GetuRLByUserID(ctx context.Context, in *NullRequest, opts ...grpc.CallOption) (*GetURLByUserResponse, error)
	DeleteURLs(ctx context.Context, in *DeleteURLsRequest, opts ...grpc.CallOption) (*DeleteURLsResponse, error)
	Batch(ctx context.Context, in *BatchRequest, opts ...grpc.CallOption) (*BatchResponse, error)
	GetOriginalURL(ctx context.Context, in *GetOriginalURLRequest, opts ...grpc.CallOption) (*GetOriginalURLResponse, error)
	GetStats(ctx context.Context, in *NullRequest, opts ...grpc.CallOption) (*GetStatsResponse, error)
	Ping(ctx context.Context, in *NullRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) AddURL(ctx context.Context, in *AddURLRequest, opts ...grpc.CallOption) (*AddURLResponse, error) {
	out := new(AddURLResponse)
	err := c.cc.Invoke(ctx, Shortener_AddURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetuRLByUserID(ctx context.Context, in *NullRequest, opts ...grpc.CallOption) (*GetURLByUserResponse, error) {
	out := new(GetURLByUserResponse)
	err := c.cc.Invoke(ctx, Shortener_GetuRLByUserID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) DeleteURLs(ctx context.Context, in *DeleteURLsRequest, opts ...grpc.CallOption) (*DeleteURLsResponse, error) {
	out := new(DeleteURLsResponse)
	err := c.cc.Invoke(ctx, Shortener_DeleteURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) Batch(ctx context.Context, in *BatchRequest, opts ...grpc.CallOption) (*BatchResponse, error) {
	out := new(BatchResponse)
	err := c.cc.Invoke(ctx, Shortener_Batch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetOriginalURL(ctx context.Context, in *GetOriginalURLRequest, opts ...grpc.CallOption) (*GetOriginalURLResponse, error) {
	out := new(GetOriginalURLResponse)
	err := c.cc.Invoke(ctx, Shortener_GetOriginalURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetStats(ctx context.Context, in *NullRequest, opts ...grpc.CallOption) (*GetStatsResponse, error) {
	out := new(GetStatsResponse)
	err := c.cc.Invoke(ctx, Shortener_GetStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) Ping(ctx context.Context, in *NullRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, Shortener_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	AddURL(context.Context, *AddURLRequest) (*AddURLResponse, error)
	GetuRLByUserID(context.Context, *NullRequest) (*GetURLByUserResponse, error)
	DeleteURLs(context.Context, *DeleteURLsRequest) (*DeleteURLsResponse, error)
	Batch(context.Context, *BatchRequest) (*BatchResponse, error)
	GetOriginalURL(context.Context, *GetOriginalURLRequest) (*GetOriginalURLResponse, error)
	GetStats(context.Context, *NullRequest) (*GetStatsResponse, error)
	Ping(context.Context, *NullRequest) (*PingResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

func (UnimplementedShortenerServer) AddURL(context.Context, *AddURLRequest) (*AddURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddURL not implemented")
}
func (UnimplementedShortenerServer) GetuRLByUserID(context.Context, *NullRequest) (*GetURLByUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetuRLByUserID not implemented")
}
func (UnimplementedShortenerServer) DeleteURLs(context.Context, *DeleteURLsRequest) (*DeleteURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURLs not implemented")
}
func (UnimplementedShortenerServer) Batch(context.Context, *BatchRequest) (*BatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Batch not implemented")
}
func (UnimplementedShortenerServer) GetOriginalURL(context.Context, *GetOriginalURLRequest) (*GetOriginalURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalURL not implemented")
}
func (UnimplementedShortenerServer) GetStats(context.Context, *NullRequest) (*GetStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedShortenerServer) Ping(context.Context, *NullRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_AddURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).AddURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_AddURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).AddURL(ctx, req.(*AddURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetuRLByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NullRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetuRLByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_GetuRLByUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetuRLByUserID(ctx, req.(*NullRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_DeleteURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).DeleteURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_DeleteURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).DeleteURLs(ctx, req.(*DeleteURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_Batch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Batch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_Batch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Batch(ctx, req.(*BatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetOriginalURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOriginalURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetOriginalURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_GetOriginalURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetOriginalURL(ctx, req.(*GetOriginalURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NullRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_GetStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetStats(ctx, req.(*NullRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NullRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Ping(ctx, req.(*NullRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddURL",
			Handler:    _Shortener_AddURL_Handler,
		},
		{
			MethodName: "GetuRLByUserID",
			Handler:    _Shortener_GetuRLByUserID_Handler,
		},
		{
			MethodName: "DeleteURLs",
			Handler:    _Shortener_DeleteURLs_Handler,
		},
		{
			MethodName: "Batch",
			Handler:    _Shortener_Batch_Handler,
		},
		{
			MethodName: "GetOriginalURL",
			Handler:    _Shortener_GetOriginalURL_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _Shortener_GetStats_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Shortener_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortener.proto",
}