// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: proto/NetLocker.proto

package NetLockerController

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
	NetLockerService_TryAndLock_FullMethodName = "/NetLocker.NetLockerService/TryAndLock"
	NetLockerService_Unlock_FullMethodName     = "/NetLocker.NetLockerService/Unlock"
)

// NetLockerServiceClient is the client API for NetLockerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetLockerServiceClient interface {
	TryAndLock(ctx context.Context, in *NetLockRequest, opts ...grpc.CallOption) (*NetLockerResponse, error)
	Unlock(ctx context.Context, in *NetUnlockRequest, opts ...grpc.CallOption) (*NetLockerResponse, error)
}

type netLockerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNetLockerServiceClient(cc grpc.ClientConnInterface) NetLockerServiceClient {
	return &netLockerServiceClient{cc}
}

func (c *netLockerServiceClient) TryAndLock(ctx context.Context, in *NetLockRequest, opts ...grpc.CallOption) (*NetLockerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NetLockerResponse)
	err := c.cc.Invoke(ctx, NetLockerService_TryAndLock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *netLockerServiceClient) Unlock(ctx context.Context, in *NetUnlockRequest, opts ...grpc.CallOption) (*NetLockerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NetLockerResponse)
	err := c.cc.Invoke(ctx, NetLockerService_Unlock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetLockerServiceServer is the server API for NetLockerService service.
// All implementations must embed UnimplementedNetLockerServiceServer
// for forward compatibility.
type NetLockerServiceServer interface {
	TryAndLock(context.Context, *NetLockRequest) (*NetLockerResponse, error)
	Unlock(context.Context, *NetUnlockRequest) (*NetLockerResponse, error)
	mustEmbedUnimplementedNetLockerServiceServer()
}

// UnimplementedNetLockerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNetLockerServiceServer struct{}

func (UnimplementedNetLockerServiceServer) TryAndLock(context.Context, *NetLockRequest) (*NetLockerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryAndLock not implemented")
}
func (UnimplementedNetLockerServiceServer) Unlock(context.Context, *NetUnlockRequest) (*NetLockerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unlock not implemented")
}
func (UnimplementedNetLockerServiceServer) mustEmbedUnimplementedNetLockerServiceServer() {}
func (UnimplementedNetLockerServiceServer) testEmbeddedByValue()                          {}

// UnsafeNetLockerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetLockerServiceServer will
// result in compilation errors.
type UnsafeNetLockerServiceServer interface {
	mustEmbedUnimplementedNetLockerServiceServer()
}

func RegisterNetLockerServiceServer(s grpc.ServiceRegistrar, srv NetLockerServiceServer) {
	// If the following call pancis, it indicates UnimplementedNetLockerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NetLockerService_ServiceDesc, srv)
}

func _NetLockerService_TryAndLock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetLockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetLockerServiceServer).TryAndLock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NetLockerService_TryAndLock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetLockerServiceServer).TryAndLock(ctx, req.(*NetLockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetLockerService_Unlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetUnlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetLockerServiceServer).Unlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NetLockerService_Unlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetLockerServiceServer).Unlock(ctx, req.(*NetUnlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NetLockerService_ServiceDesc is the grpc.ServiceDesc for NetLockerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetLockerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NetLocker.NetLockerService",
	HandlerType: (*NetLockerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TryAndLock",
			Handler:    _NetLockerService_TryAndLock_Handler,
		},
		{
			MethodName: "Unlock",
			Handler:    _NetLockerService_Unlock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/NetLocker.proto",
}
