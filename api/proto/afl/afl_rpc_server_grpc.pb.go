// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: afl_rpc_server.proto

package brownlowdev

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

// AflClient is the client API for Afl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AflClient interface {
	CreateUmpire(ctx context.Context, in *CreateUmpireRequest, opts ...grpc.CallOption) (*CreateUmpireResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}

type aflClient struct {
	cc grpc.ClientConnInterface
}

func NewAflClient(cc grpc.ClientConnInterface) AflClient {
	return &aflClient{cc}
}

func (c *aflClient) CreateUmpire(ctx context.Context, in *CreateUmpireRequest, opts ...grpc.CallOption) (*CreateUmpireResponse, error) {
	out := new(CreateUmpireResponse)
	err := c.cc.Invoke(ctx, "/grpc.Afl/CreateUmpire", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aflClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/grpc.Afl/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AflServer is the server API for Afl service.
// All implementations must embed UnimplementedAflServer
// for forward compatibility
type AflServer interface {
	CreateUmpire(context.Context, *CreateUmpireRequest) (*CreateUmpireResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	mustEmbedUnimplementedAflServer()
}

// UnimplementedAflServer must be embedded to have forward compatible implementations.
type UnimplementedAflServer struct {
}

func (UnimplementedAflServer) CreateUmpire(context.Context, *CreateUmpireRequest) (*CreateUmpireResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUmpire not implemented")
}
func (UnimplementedAflServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAflServer) mustEmbedUnimplementedAflServer() {}

// UnsafeAflServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AflServer will
// result in compilation errors.
type UnsafeAflServer interface {
	mustEmbedUnimplementedAflServer()
}

func RegisterAflServer(s grpc.ServiceRegistrar, srv AflServer) {
	s.RegisterService(&Afl_ServiceDesc, srv)
}

func _Afl_CreateUmpire_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUmpireRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AflServer).CreateUmpire(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Afl/CreateUmpire",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AflServer).CreateUmpire(ctx, req.(*CreateUmpireRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Afl_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AflServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Afl/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AflServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Afl_ServiceDesc is the grpc.ServiceDesc for Afl service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Afl_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Afl",
	HandlerType: (*AflServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUmpire",
			Handler:    _Afl_CreateUmpire_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Afl_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "afl_rpc_server.proto",
}
