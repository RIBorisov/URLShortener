// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: proto/service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	URLShortenerService_Save_FullMethodName        = "/URLShortenerService/Save"
	URLShortenerService_Batch_FullMethodName       = "/URLShortenerService/Batch"
	URLShortenerService_DeleteMany_FullMethodName  = "/URLShortenerService/DeleteMany"
	URLShortenerService_Get_FullMethodName         = "/URLShortenerService/Get"
	URLShortenerService_Ping_FullMethodName        = "/URLShortenerService/Ping"
	URLShortenerService_Shorten_FullMethodName     = "/URLShortenerService/Shorten"
	URLShortenerService_Stats_FullMethodName       = "/URLShortenerService/Stats"
	URLShortenerService_SavedByUser_FullMethodName = "/URLShortenerService/SavedByUser"
)

// URLShortenerServiceClient is the client API for URLShortenerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLShortenerServiceClient interface {
	Save(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.StringValue, error)
	Batch(ctx context.Context, in *BatchRequest, opts ...grpc.CallOption) (*BatchResponse, error)
	DeleteMany(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Shorten(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error)
	Stats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error)
	SavedByUser(ctx context.Context, in *SavedByUserRequest, opts ...grpc.CallOption) (*SavedByUserResponse, error)
}

type uRLShortenerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewURLShortenerServiceClient(cc grpc.ClientConnInterface) URLShortenerServiceClient {
	return &uRLShortenerServiceClient{cc}
}

func (c *uRLShortenerServiceClient) Save(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(wrapperspb.StringValue)
	err := c.cc.Invoke(ctx, URLShortenerService_Save_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) Batch(ctx context.Context, in *BatchRequest, opts ...grpc.CallOption) (*BatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_Batch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) DeleteMany(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_DeleteMany_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) Shorten(ctx context.Context, in *ShortenRequest, opts ...grpc.CallOption) (*ShortenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShortenResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_Shorten_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) Stats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (*StatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_Stats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerServiceClient) SavedByUser(ctx context.Context, in *SavedByUserRequest, opts ...grpc.CallOption) (*SavedByUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SavedByUserResponse)
	err := c.cc.Invoke(ctx, URLShortenerService_SavedByUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerServiceServer is the server API for URLShortenerService service.
// All implementations must embed UnimplementedURLShortenerServiceServer
// for forward compatibility.
type URLShortenerServiceServer interface {
	Save(context.Context, *wrapperspb.StringValue) (*wrapperspb.StringValue, error)
	Batch(context.Context, *BatchRequest) (*BatchResponse, error)
	DeleteMany(context.Context, *DeleteRequest) (*DeleteResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Shorten(context.Context, *ShortenRequest) (*ShortenResponse, error)
	Stats(context.Context, *StatsRequest) (*StatsResponse, error)
	SavedByUser(context.Context, *SavedByUserRequest) (*SavedByUserResponse, error)
	mustEmbedUnimplementedURLShortenerServiceServer()
}

// UnimplementedURLShortenerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedURLShortenerServiceServer struct{}

func (UnimplementedURLShortenerServiceServer) Save(context.Context, *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedURLShortenerServiceServer) Batch(context.Context, *BatchRequest) (*BatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Batch not implemented")
}
func (UnimplementedURLShortenerServiceServer) DeleteMany(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMany not implemented")
}
func (UnimplementedURLShortenerServiceServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedURLShortenerServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedURLShortenerServiceServer) Shorten(context.Context, *ShortenRequest) (*ShortenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shorten not implemented")
}
func (UnimplementedURLShortenerServiceServer) Stats(context.Context, *StatsRequest) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedURLShortenerServiceServer) SavedByUser(context.Context, *SavedByUserRequest) (*SavedByUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SavedByUser not implemented")
}
func (UnimplementedURLShortenerServiceServer) mustEmbedUnimplementedURLShortenerServiceServer() {}
func (UnimplementedURLShortenerServiceServer) testEmbeddedByValue()                             {}

// UnsafeURLShortenerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerServiceServer will
// result in compilation errors.
type UnsafeURLShortenerServiceServer interface {
	mustEmbedUnimplementedURLShortenerServiceServer()
}

func RegisterURLShortenerServiceServer(s grpc.ServiceRegistrar, srv URLShortenerServiceServer) {
	// If the following call pancis, it indicates UnimplementedURLShortenerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&URLShortenerService_ServiceDesc, srv)
}

func _URLShortenerService_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_Save_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).Save(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_Batch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).Batch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_Batch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).Batch(ctx, req.(*BatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_DeleteMany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).DeleteMany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_DeleteMany_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).DeleteMany(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_Shorten_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).Shorten(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_Shorten_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).Shorten(ctx, req.(*ShortenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_Stats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).Stats(ctx, req.(*StatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerService_SavedByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SavedByUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServiceServer).SavedByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerService_SavedByUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServiceServer).SavedByUser(ctx, req.(*SavedByUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortenerService_ServiceDesc is the grpc.ServiceDesc for URLShortenerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortenerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "URLShortenerService",
	HandlerType: (*URLShortenerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Save",
			Handler:    _URLShortenerService_Save_Handler,
		},
		{
			MethodName: "Batch",
			Handler:    _URLShortenerService_Batch_Handler,
		},
		{
			MethodName: "DeleteMany",
			Handler:    _URLShortenerService_DeleteMany_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _URLShortenerService_Get_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _URLShortenerService_Ping_Handler,
		},
		{
			MethodName: "Shorten",
			Handler:    _URLShortenerService_Shorten_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _URLShortenerService_Stats_Handler,
		},
		{
			MethodName: "SavedByUser",
			Handler:    _URLShortenerService_SavedByUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}