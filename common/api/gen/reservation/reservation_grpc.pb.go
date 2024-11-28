// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: reservation.proto

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
	Reservation_CreateReservation_FullMethodName = "/reservation.Reservation/CreateReservation"
	Reservation_CancelReservation_FullMethodName = "/reservation.Reservation/CancelReservation"
	Reservation_CloseReservation_FullMethodName  = "/reservation.Reservation/CloseReservation"
)

// ReservationClient is the client API for Reservation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReservationClient interface {
	CreateReservation(ctx context.Context, in *CreateReservationRequest, opts ...grpc.CallOption) (*CreateReservationResponse, error)
	CancelReservation(ctx context.Context, in *CancelReservationRequest, opts ...grpc.CallOption) (*CancelReservationResponse, error)
	CloseReservation(ctx context.Context, in *CloseReservationRequest, opts ...grpc.CallOption) (*CloseReservationResponse, error)
}

type reservationClient struct {
	cc grpc.ClientConnInterface
}

func NewReservationClient(cc grpc.ClientConnInterface) ReservationClient {
	return &reservationClient{cc}
}

func (c *reservationClient) CreateReservation(ctx context.Context, in *CreateReservationRequest, opts ...grpc.CallOption) (*CreateReservationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateReservationResponse)
	err := c.cc.Invoke(ctx, Reservation_CreateReservation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationClient) CancelReservation(ctx context.Context, in *CancelReservationRequest, opts ...grpc.CallOption) (*CancelReservationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CancelReservationResponse)
	err := c.cc.Invoke(ctx, Reservation_CancelReservation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reservationClient) CloseReservation(ctx context.Context, in *CloseReservationRequest, opts ...grpc.CallOption) (*CloseReservationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseReservationResponse)
	err := c.cc.Invoke(ctx, Reservation_CloseReservation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReservationServer is the server API for Reservation service.
// All implementations must embed UnimplementedReservationServer
// for forward compatibility.
type ReservationServer interface {
	CreateReservation(context.Context, *CreateReservationRequest) (*CreateReservationResponse, error)
	CancelReservation(context.Context, *CancelReservationRequest) (*CancelReservationResponse, error)
	CloseReservation(context.Context, *CloseReservationRequest) (*CloseReservationResponse, error)
	mustEmbedUnimplementedReservationServer()
}

// UnimplementedReservationServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedReservationServer struct{}

func (UnimplementedReservationServer) CreateReservation(context.Context, *CreateReservationRequest) (*CreateReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateReservation not implemented")
}
func (UnimplementedReservationServer) CancelReservation(context.Context, *CancelReservationRequest) (*CancelReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelReservation not implemented")
}
func (UnimplementedReservationServer) CloseReservation(context.Context, *CloseReservationRequest) (*CloseReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseReservation not implemented")
}
func (UnimplementedReservationServer) mustEmbedUnimplementedReservationServer() {}
func (UnimplementedReservationServer) testEmbeddedByValue()                     {}

// UnsafeReservationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReservationServer will
// result in compilation errors.
type UnsafeReservationServer interface {
	mustEmbedUnimplementedReservationServer()
}

func RegisterReservationServer(s grpc.ServiceRegistrar, srv ReservationServer) {
	// If the following call pancis, it indicates UnimplementedReservationServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Reservation_ServiceDesc, srv)
}

func _Reservation_CreateReservation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateReservationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServer).CreateReservation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Reservation_CreateReservation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServer).CreateReservation(ctx, req.(*CreateReservationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Reservation_CancelReservation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelReservationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServer).CancelReservation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Reservation_CancelReservation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServer).CancelReservation(ctx, req.(*CancelReservationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Reservation_CloseReservation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseReservationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReservationServer).CloseReservation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Reservation_CloseReservation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReservationServer).CloseReservation(ctx, req.(*CloseReservationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Reservation_ServiceDesc is the grpc.ServiceDesc for Reservation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Reservation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "reservation.Reservation",
	HandlerType: (*ReservationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateReservation",
			Handler:    _Reservation_CreateReservation_Handler,
		},
		{
			MethodName: "CancelReservation",
			Handler:    _Reservation_CancelReservation_Handler,
		},
		{
			MethodName: "CloseReservation",
			Handler:    _Reservation_CloseReservation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "reservation.proto",
}