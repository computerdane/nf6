// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: nf6.proto

package nf6

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
	Nf6Public_CreateAccount_FullMethodName = "/nf6.Nf6Public/CreateAccount"
	Nf6Public_GetCaCert_FullMethodName     = "/nf6.Nf6Public/GetCaCert"
)

// Nf6PublicClient is the client API for Nf6Public service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type Nf6PublicClient interface {
	CreateAccount(ctx context.Context, in *CreateAccount_Request, opts ...grpc.CallOption) (*None, error)
	GetCaCert(ctx context.Context, in *None, opts ...grpc.CallOption) (*GetCaCert_Reply, error)
}

type nf6PublicClient struct {
	cc grpc.ClientConnInterface
}

func NewNf6PublicClient(cc grpc.ClientConnInterface) Nf6PublicClient {
	return &nf6PublicClient{cc}
}

func (c *nf6PublicClient) CreateAccount(ctx context.Context, in *CreateAccount_Request, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, Nf6Public_CreateAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6PublicClient) GetCaCert(ctx context.Context, in *None, opts ...grpc.CallOption) (*GetCaCert_Reply, error) {
	out := new(GetCaCert_Reply)
	err := c.cc.Invoke(ctx, Nf6Public_GetCaCert_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Nf6PublicServer is the server API for Nf6Public service.
// All implementations must embed UnimplementedNf6PublicServer
// for forward compatibility
type Nf6PublicServer interface {
	CreateAccount(context.Context, *CreateAccount_Request) (*None, error)
	GetCaCert(context.Context, *None) (*GetCaCert_Reply, error)
	mustEmbedUnimplementedNf6PublicServer()
}

// UnimplementedNf6PublicServer must be embedded to have forward compatible implementations.
type UnimplementedNf6PublicServer struct {
}

func (UnimplementedNf6PublicServer) CreateAccount(context.Context, *CreateAccount_Request) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedNf6PublicServer) GetCaCert(context.Context, *None) (*GetCaCert_Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCaCert not implemented")
}
func (UnimplementedNf6PublicServer) mustEmbedUnimplementedNf6PublicServer() {}

// UnsafeNf6PublicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Nf6PublicServer will
// result in compilation errors.
type UnsafeNf6PublicServer interface {
	mustEmbedUnimplementedNf6PublicServer()
}

func RegisterNf6PublicServer(s grpc.ServiceRegistrar, srv Nf6PublicServer) {
	s.RegisterService(&Nf6Public_ServiceDesc, srv)
}

func _Nf6Public_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccount_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6PublicServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6Public_CreateAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6PublicServer).CreateAccount(ctx, req.(*CreateAccount_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6Public_GetCaCert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(None)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6PublicServer).GetCaCert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6Public_GetCaCert_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6PublicServer).GetCaCert(ctx, req.(*None))
	}
	return interceptor(ctx, in, info, handler)
}

// Nf6Public_ServiceDesc is the grpc.ServiceDesc for Nf6Public service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Nf6Public_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nf6.Nf6Public",
	HandlerType: (*Nf6PublicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAccount",
			Handler:    _Nf6Public_CreateAccount_Handler,
		},
		{
			MethodName: "GetCaCert",
			Handler:    _Nf6Public_GetCaCert_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nf6.proto",
}

const (
	Nf6_GetAccount_FullMethodName           = "/nf6.Nf6/GetAccount"
	Nf6_UpdateAccount_FullMethodName        = "/nf6.Nf6/UpdateAccount"
	Nf6_CreateHost_FullMethodName           = "/nf6.Nf6/CreateHost"
	Nf6_GetHost_FullMethodName              = "/nf6.Nf6/GetHost"
	Nf6_ListHosts_FullMethodName            = "/nf6.Nf6/ListHosts"
	Nf6_UpdateHost_FullMethodName           = "/nf6.Nf6/UpdateHost"
	Nf6_CreateRepo_FullMethodName           = "/nf6.Nf6/CreateRepo"
	Nf6_GetRepo_FullMethodName              = "/nf6.Nf6/GetRepo"
	Nf6_ListRepos_FullMethodName            = "/nf6.Nf6/ListRepos"
	Nf6_UpdateRepo_FullMethodName           = "/nf6.Nf6/UpdateRepo"
	Nf6_GitServer_GetAccount_FullMethodName = "/nf6.Nf6/GitServer_GetAccount"
	Nf6_GitServer_ListRepos_FullMethodName  = "/nf6.Nf6/GitServer_ListRepos"
	Nf6_WgServer_ListHosts_FullMethodName   = "/nf6.Nf6/WgServer_ListHosts"
)

// Nf6Client is the client API for Nf6 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type Nf6Client interface {
	// Users
	GetAccount(ctx context.Context, in *None, opts ...grpc.CallOption) (*GetAccount_Reply, error)
	UpdateAccount(ctx context.Context, in *UpdateAccount_Request, opts ...grpc.CallOption) (*None, error)
	CreateHost(ctx context.Context, in *CreateHost_Request, opts ...grpc.CallOption) (*None, error)
	GetHost(ctx context.Context, in *GetHost_Request, opts ...grpc.CallOption) (*GetHost_Reply, error)
	ListHosts(ctx context.Context, in *None, opts ...grpc.CallOption) (*ListHosts_Reply, error)
	UpdateHost(ctx context.Context, in *UpdateHost_Request, opts ...grpc.CallOption) (*None, error)
	CreateRepo(ctx context.Context, in *CreateRepo_Request, opts ...grpc.CallOption) (*None, error)
	GetRepo(ctx context.Context, in *GetRepo_Request, opts ...grpc.CallOption) (*GetRepo_Reply, error)
	ListRepos(ctx context.Context, in *None, opts ...grpc.CallOption) (*ListRepos_Reply, error)
	UpdateRepo(ctx context.Context, in *UpdateRepo_Request, opts ...grpc.CallOption) (*None, error)
	// Git server
	GitServer_GetAccount(ctx context.Context, in *GitServer_GetAccount_Request, opts ...grpc.CallOption) (*GitServer_GetAccount_Reply, error)
	GitServer_ListRepos(ctx context.Context, in *GitServer_ListRepos_Request, opts ...grpc.CallOption) (*GitServer_ListRepos_Reply, error)
	// WireGuard server
	WgServer_ListHosts(ctx context.Context, in *None, opts ...grpc.CallOption) (*WgServer_ListHosts_Reply, error)
}

type nf6Client struct {
	cc grpc.ClientConnInterface
}

func NewNf6Client(cc grpc.ClientConnInterface) Nf6Client {
	return &nf6Client{cc}
}

func (c *nf6Client) GetAccount(ctx context.Context, in *None, opts ...grpc.CallOption) (*GetAccount_Reply, error) {
	out := new(GetAccount_Reply)
	err := c.cc.Invoke(ctx, Nf6_GetAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) UpdateAccount(ctx context.Context, in *UpdateAccount_Request, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, Nf6_UpdateAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) CreateHost(ctx context.Context, in *CreateHost_Request, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, Nf6_CreateHost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) GetHost(ctx context.Context, in *GetHost_Request, opts ...grpc.CallOption) (*GetHost_Reply, error) {
	out := new(GetHost_Reply)
	err := c.cc.Invoke(ctx, Nf6_GetHost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) ListHosts(ctx context.Context, in *None, opts ...grpc.CallOption) (*ListHosts_Reply, error) {
	out := new(ListHosts_Reply)
	err := c.cc.Invoke(ctx, Nf6_ListHosts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) UpdateHost(ctx context.Context, in *UpdateHost_Request, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, Nf6_UpdateHost_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) CreateRepo(ctx context.Context, in *CreateRepo_Request, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, Nf6_CreateRepo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) GetRepo(ctx context.Context, in *GetRepo_Request, opts ...grpc.CallOption) (*GetRepo_Reply, error) {
	out := new(GetRepo_Reply)
	err := c.cc.Invoke(ctx, Nf6_GetRepo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) ListRepos(ctx context.Context, in *None, opts ...grpc.CallOption) (*ListRepos_Reply, error) {
	out := new(ListRepos_Reply)
	err := c.cc.Invoke(ctx, Nf6_ListRepos_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) UpdateRepo(ctx context.Context, in *UpdateRepo_Request, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, Nf6_UpdateRepo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) GitServer_GetAccount(ctx context.Context, in *GitServer_GetAccount_Request, opts ...grpc.CallOption) (*GitServer_GetAccount_Reply, error) {
	out := new(GitServer_GetAccount_Reply)
	err := c.cc.Invoke(ctx, Nf6_GitServer_GetAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) GitServer_ListRepos(ctx context.Context, in *GitServer_ListRepos_Request, opts ...grpc.CallOption) (*GitServer_ListRepos_Reply, error) {
	out := new(GitServer_ListRepos_Reply)
	err := c.cc.Invoke(ctx, Nf6_GitServer_ListRepos_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nf6Client) WgServer_ListHosts(ctx context.Context, in *None, opts ...grpc.CallOption) (*WgServer_ListHosts_Reply, error) {
	out := new(WgServer_ListHosts_Reply)
	err := c.cc.Invoke(ctx, Nf6_WgServer_ListHosts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Nf6Server is the server API for Nf6 service.
// All implementations must embed UnimplementedNf6Server
// for forward compatibility
type Nf6Server interface {
	// Users
	GetAccount(context.Context, *None) (*GetAccount_Reply, error)
	UpdateAccount(context.Context, *UpdateAccount_Request) (*None, error)
	CreateHost(context.Context, *CreateHost_Request) (*None, error)
	GetHost(context.Context, *GetHost_Request) (*GetHost_Reply, error)
	ListHosts(context.Context, *None) (*ListHosts_Reply, error)
	UpdateHost(context.Context, *UpdateHost_Request) (*None, error)
	CreateRepo(context.Context, *CreateRepo_Request) (*None, error)
	GetRepo(context.Context, *GetRepo_Request) (*GetRepo_Reply, error)
	ListRepos(context.Context, *None) (*ListRepos_Reply, error)
	UpdateRepo(context.Context, *UpdateRepo_Request) (*None, error)
	// Git server
	GitServer_GetAccount(context.Context, *GitServer_GetAccount_Request) (*GitServer_GetAccount_Reply, error)
	GitServer_ListRepos(context.Context, *GitServer_ListRepos_Request) (*GitServer_ListRepos_Reply, error)
	// WireGuard server
	WgServer_ListHosts(context.Context, *None) (*WgServer_ListHosts_Reply, error)
	mustEmbedUnimplementedNf6Server()
}

// UnimplementedNf6Server must be embedded to have forward compatible implementations.
type UnimplementedNf6Server struct {
}

func (UnimplementedNf6Server) GetAccount(context.Context, *None) (*GetAccount_Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedNf6Server) UpdateAccount(context.Context, *UpdateAccount_Request) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAccount not implemented")
}
func (UnimplementedNf6Server) CreateHost(context.Context, *CreateHost_Request) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateHost not implemented")
}
func (UnimplementedNf6Server) GetHost(context.Context, *GetHost_Request) (*GetHost_Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHost not implemented")
}
func (UnimplementedNf6Server) ListHosts(context.Context, *None) (*ListHosts_Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListHosts not implemented")
}
func (UnimplementedNf6Server) UpdateHost(context.Context, *UpdateHost_Request) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHost not implemented")
}
func (UnimplementedNf6Server) CreateRepo(context.Context, *CreateRepo_Request) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRepo not implemented")
}
func (UnimplementedNf6Server) GetRepo(context.Context, *GetRepo_Request) (*GetRepo_Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRepo not implemented")
}
func (UnimplementedNf6Server) ListRepos(context.Context, *None) (*ListRepos_Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRepos not implemented")
}
func (UnimplementedNf6Server) UpdateRepo(context.Context, *UpdateRepo_Request) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRepo not implemented")
}
func (UnimplementedNf6Server) GitServer_GetAccount(context.Context, *GitServer_GetAccount_Request) (*GitServer_GetAccount_Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GitServer_GetAccount not implemented")
}
func (UnimplementedNf6Server) GitServer_ListRepos(context.Context, *GitServer_ListRepos_Request) (*GitServer_ListRepos_Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GitServer_ListRepos not implemented")
}
func (UnimplementedNf6Server) WgServer_ListHosts(context.Context, *None) (*WgServer_ListHosts_Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WgServer_ListHosts not implemented")
}
func (UnimplementedNf6Server) mustEmbedUnimplementedNf6Server() {}

// UnsafeNf6Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Nf6Server will
// result in compilation errors.
type UnsafeNf6Server interface {
	mustEmbedUnimplementedNf6Server()
}

func RegisterNf6Server(s grpc.ServiceRegistrar, srv Nf6Server) {
	s.RegisterService(&Nf6_ServiceDesc, srv)
}

func _Nf6_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(None)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_GetAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).GetAccount(ctx, req.(*None))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_UpdateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAccount_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).UpdateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_UpdateAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).UpdateAccount(ctx, req.(*UpdateAccount_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_CreateHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateHost_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).CreateHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_CreateHost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).CreateHost(ctx, req.(*CreateHost_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_GetHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHost_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).GetHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_GetHost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).GetHost(ctx, req.(*GetHost_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_ListHosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(None)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).ListHosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_ListHosts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).ListHosts(ctx, req.(*None))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_UpdateHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateHost_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).UpdateHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_UpdateHost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).UpdateHost(ctx, req.(*UpdateHost_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_CreateRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRepo_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).CreateRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_CreateRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).CreateRepo(ctx, req.(*CreateRepo_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_GetRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRepo_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).GetRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_GetRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).GetRepo(ctx, req.(*GetRepo_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_ListRepos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(None)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).ListRepos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_ListRepos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).ListRepos(ctx, req.(*None))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_UpdateRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRepo_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).UpdateRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_UpdateRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).UpdateRepo(ctx, req.(*UpdateRepo_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_GitServer_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GitServer_GetAccount_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).GitServer_GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_GitServer_GetAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).GitServer_GetAccount(ctx, req.(*GitServer_GetAccount_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_GitServer_ListRepos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GitServer_ListRepos_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).GitServer_ListRepos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_GitServer_ListRepos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).GitServer_ListRepos(ctx, req.(*GitServer_ListRepos_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Nf6_WgServer_ListHosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(None)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6Server).WgServer_ListHosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6_WgServer_ListHosts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6Server).WgServer_ListHosts(ctx, req.(*None))
	}
	return interceptor(ctx, in, info, handler)
}

// Nf6_ServiceDesc is the grpc.ServiceDesc for Nf6 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Nf6_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nf6.Nf6",
	HandlerType: (*Nf6Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAccount",
			Handler:    _Nf6_GetAccount_Handler,
		},
		{
			MethodName: "UpdateAccount",
			Handler:    _Nf6_UpdateAccount_Handler,
		},
		{
			MethodName: "CreateHost",
			Handler:    _Nf6_CreateHost_Handler,
		},
		{
			MethodName: "GetHost",
			Handler:    _Nf6_GetHost_Handler,
		},
		{
			MethodName: "ListHosts",
			Handler:    _Nf6_ListHosts_Handler,
		},
		{
			MethodName: "UpdateHost",
			Handler:    _Nf6_UpdateHost_Handler,
		},
		{
			MethodName: "CreateRepo",
			Handler:    _Nf6_CreateRepo_Handler,
		},
		{
			MethodName: "GetRepo",
			Handler:    _Nf6_GetRepo_Handler,
		},
		{
			MethodName: "ListRepos",
			Handler:    _Nf6_ListRepos_Handler,
		},
		{
			MethodName: "UpdateRepo",
			Handler:    _Nf6_UpdateRepo_Handler,
		},
		{
			MethodName: "GitServer_GetAccount",
			Handler:    _Nf6_GitServer_GetAccount_Handler,
		},
		{
			MethodName: "GitServer_ListRepos",
			Handler:    _Nf6_GitServer_ListRepos_Handler,
		},
		{
			MethodName: "WgServer_ListHosts",
			Handler:    _Nf6_WgServer_ListHosts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nf6.proto",
}

const (
	Nf6Wg_CreateRoute_FullMethodName = "/nf6.Nf6Wg/CreateRoute"
)

// Nf6WgClient is the client API for Nf6Wg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type Nf6WgClient interface {
	CreateRoute(ctx context.Context, in *CreateRoute_Request, opts ...grpc.CallOption) (*None, error)
}

type nf6WgClient struct {
	cc grpc.ClientConnInterface
}

func NewNf6WgClient(cc grpc.ClientConnInterface) Nf6WgClient {
	return &nf6WgClient{cc}
}

func (c *nf6WgClient) CreateRoute(ctx context.Context, in *CreateRoute_Request, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, Nf6Wg_CreateRoute_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Nf6WgServer is the server API for Nf6Wg service.
// All implementations must embed UnimplementedNf6WgServer
// for forward compatibility
type Nf6WgServer interface {
	CreateRoute(context.Context, *CreateRoute_Request) (*None, error)
	mustEmbedUnimplementedNf6WgServer()
}

// UnimplementedNf6WgServer must be embedded to have forward compatible implementations.
type UnimplementedNf6WgServer struct {
}

func (UnimplementedNf6WgServer) CreateRoute(context.Context, *CreateRoute_Request) (*None, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoute not implemented")
}
func (UnimplementedNf6WgServer) mustEmbedUnimplementedNf6WgServer() {}

// UnsafeNf6WgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Nf6WgServer will
// result in compilation errors.
type UnsafeNf6WgServer interface {
	mustEmbedUnimplementedNf6WgServer()
}

func RegisterNf6WgServer(s grpc.ServiceRegistrar, srv Nf6WgServer) {
	s.RegisterService(&Nf6Wg_ServiceDesc, srv)
}

func _Nf6Wg_CreateRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoute_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Nf6WgServer).CreateRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Nf6Wg_CreateRoute_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Nf6WgServer).CreateRoute(ctx, req.(*CreateRoute_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Nf6Wg_ServiceDesc is the grpc.ServiceDesc for Nf6Wg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Nf6Wg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nf6.Nf6Wg",
	HandlerType: (*Nf6WgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRoute",
			Handler:    _Nf6Wg_CreateRoute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nf6.proto",
}
