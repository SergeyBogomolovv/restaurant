// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: sso.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SSO_RegisterCustomer_FullMethodName = "/sso.SSO/RegisterCustomer"
	SSO_RegisterWaiter_FullMethodName   = "/sso.SSO/RegisterWaiter"
	SSO_RegisterAdmin_FullMethodName    = "/sso.SSO/RegisterAdmin"
	SSO_LoginCustomer_FullMethodName    = "/sso.SSO/LoginCustomer"
	SSO_LoginWaiter_FullMethodName      = "/sso.SSO/LoginWaiter"
	SSO_LoginAdmin_FullMethodName       = "/sso.SSO/LoginAdmin"
	SSO_Refresh_FullMethodName          = "/sso.SSO/Refresh"
	SSO_Logout_FullMethodName           = "/sso.SSO/Logout"
)

// SSOClient is the client API for SSO service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SSOClient interface {
	RegisterCustomer(ctx context.Context, in *RegisterCustomerRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	RegisterWaiter(ctx context.Context, in *RegisterWaiterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	RegisterAdmin(ctx context.Context, in *RegisterAdminRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	LoginCustomer(ctx context.Context, in *LoginCustomerRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	LoginWaiter(ctx context.Context, in *LoginEmployeeRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	LoginAdmin(ctx context.Context, in *LoginEmployeeRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
}

type sSOClient struct {
	cc grpc.ClientConnInterface
}

func NewSSOClient(cc grpc.ClientConnInterface) SSOClient {
	return &sSOClient{cc}
}

func (c *sSOClient) RegisterCustomer(ctx context.Context, in *RegisterCustomerRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, SSO_RegisterCustomer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSOClient) RegisterWaiter(ctx context.Context, in *RegisterWaiterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, SSO_RegisterWaiter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSOClient) RegisterAdmin(ctx context.Context, in *RegisterAdminRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, SSO_RegisterAdmin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSOClient) LoginCustomer(ctx context.Context, in *LoginCustomerRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, SSO_LoginCustomer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSOClient) LoginWaiter(ctx context.Context, in *LoginEmployeeRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, SSO_LoginWaiter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSOClient) LoginAdmin(ctx context.Context, in *LoginEmployeeRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, SSO_LoginAdmin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSOClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RefreshResponse)
	err := c.cc.Invoke(ctx, SSO_Refresh_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSOClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogoutResponse)
	err := c.cc.Invoke(ctx, SSO_Logout_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SSOServer is the server API for SSO service.
// All implementations must embed UnimplementedSSOServer
// for forward compatibility.
type SSOServer interface {
	RegisterCustomer(context.Context, *RegisterCustomerRequest) (*RegisterResponse, error)
	RegisterWaiter(context.Context, *RegisterWaiterRequest) (*RegisterResponse, error)
	RegisterAdmin(context.Context, *RegisterAdminRequest) (*RegisterResponse, error)
	LoginCustomer(context.Context, *LoginCustomerRequest) (*LoginResponse, error)
	LoginWaiter(context.Context, *LoginEmployeeRequest) (*LoginResponse, error)
	LoginAdmin(context.Context, *LoginEmployeeRequest) (*LoginResponse, error)
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	mustEmbedUnimplementedSSOServer()
}

// UnimplementedSSOServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSSOServer struct{}

func (UnimplementedSSOServer) RegisterCustomer(context.Context, *RegisterCustomerRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterCustomer not implemented")
}
func (UnimplementedSSOServer) RegisterWaiter(context.Context, *RegisterWaiterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterWaiter not implemented")
}
func (UnimplementedSSOServer) RegisterAdmin(context.Context, *RegisterAdminRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterAdmin not implemented")
}
func (UnimplementedSSOServer) LoginCustomer(context.Context, *LoginCustomerRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginCustomer not implemented")
}
func (UnimplementedSSOServer) LoginWaiter(context.Context, *LoginEmployeeRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginWaiter not implemented")
}
func (UnimplementedSSOServer) LoginAdmin(context.Context, *LoginEmployeeRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginAdmin not implemented")
}
func (UnimplementedSSOServer) Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (UnimplementedSSOServer) Logout(context.Context, *LogoutRequest) (*LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedSSOServer) mustEmbedUnimplementedSSOServer() {}
func (UnimplementedSSOServer) testEmbeddedByValue()             {}

// UnsafeSSOServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SSOServer will
// result in compilation errors.
type UnsafeSSOServer interface {
	mustEmbedUnimplementedSSOServer()
}

func RegisterSSOServer(s grpc.ServiceRegistrar, srv SSOServer) {
	// If the following call pancis, it indicates UnimplementedSSOServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SSO_ServiceDesc, srv)
}

func _SSO_RegisterCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServer).RegisterCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SSO_RegisterCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServer).RegisterCustomer(ctx, req.(*RegisterCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSO_RegisterWaiter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterWaiterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServer).RegisterWaiter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SSO_RegisterWaiter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServer).RegisterWaiter(ctx, req.(*RegisterWaiterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSO_RegisterAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterAdminRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServer).RegisterAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SSO_RegisterAdmin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServer).RegisterAdmin(ctx, req.(*RegisterAdminRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSO_LoginCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServer).LoginCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SSO_LoginCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServer).LoginCustomer(ctx, req.(*LoginCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSO_LoginWaiter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginEmployeeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServer).LoginWaiter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SSO_LoginWaiter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServer).LoginWaiter(ctx, req.(*LoginEmployeeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSO_LoginAdmin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginEmployeeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServer).LoginAdmin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SSO_LoginAdmin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServer).LoginAdmin(ctx, req.(*LoginEmployeeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSO_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SSO_Refresh_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSO_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SSO_Logout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SSO_ServiceDesc is the grpc.ServiceDesc for SSO service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SSO_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sso.SSO",
	HandlerType: (*SSOServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterCustomer",
			Handler:    _SSO_RegisterCustomer_Handler,
		},
		{
			MethodName: "RegisterWaiter",
			Handler:    _SSO_RegisterWaiter_Handler,
		},
		{
			MethodName: "RegisterAdmin",
			Handler:    _SSO_RegisterAdmin_Handler,
		},
		{
			MethodName: "LoginCustomer",
			Handler:    _SSO_LoginCustomer_Handler,
		},
		{
			MethodName: "LoginWaiter",
			Handler:    _SSO_LoginWaiter_Handler,
		},
		{
			MethodName: "LoginAdmin",
			Handler:    _SSO_LoginAdmin_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _SSO_Refresh_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _SSO_Logout_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sso.proto",
}