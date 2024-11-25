// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: orden.proto

package orden

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
	CompraService_RealizarCompra_FullMethodName = "/CompraService/RealizarCompra"
)

// CompraServiceClient is the client API for CompraService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CompraServiceClient interface {
	RealizarCompra(ctx context.Context, in *Compra, opts ...grpc.CallOption) (*CompraResponse, error)
}

type compraServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCompraServiceClient(cc grpc.ClientConnInterface) CompraServiceClient {
	return &compraServiceClient{cc}
}

func (c *compraServiceClient) RealizarCompra(ctx context.Context, in *Compra, opts ...grpc.CallOption) (*CompraResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CompraResponse)
	err := c.cc.Invoke(ctx, CompraService_RealizarCompra_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompraServiceServer is the server API for CompraService service.
// All implementations must embed UnimplementedCompraServiceServer
// for forward compatibility.
type CompraServiceServer interface {
	RealizarCompra(context.Context, *Compra) (*CompraResponse, error)
	mustEmbedUnimplementedCompraServiceServer()
}

// UnimplementedCompraServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCompraServiceServer struct{}

func (UnimplementedCompraServiceServer) RealizarCompra(context.Context, *Compra) (*CompraResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RealizarCompra not implemented")
}
func (UnimplementedCompraServiceServer) mustEmbedUnimplementedCompraServiceServer() {}
func (UnimplementedCompraServiceServer) testEmbeddedByValue()                       {}

// UnsafeCompraServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CompraServiceServer will
// result in compilation errors.
type UnsafeCompraServiceServer interface {
	mustEmbedUnimplementedCompraServiceServer()
}

func RegisterCompraServiceServer(s grpc.ServiceRegistrar, srv CompraServiceServer) {
	// If the following call pancis, it indicates UnimplementedCompraServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CompraService_ServiceDesc, srv)
}

func _CompraService_RealizarCompra_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Compra)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompraServiceServer).RealizarCompra(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CompraService_RealizarCompra_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompraServiceServer).RealizarCompra(ctx, req.(*Compra))
	}
	return interceptor(ctx, in, info, handler)
}

// CompraService_ServiceDesc is the grpc.ServiceDesc for CompraService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CompraService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "CompraService",
	HandlerType: (*CompraServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RealizarCompra",
			Handler:    _CompraService_RealizarCompra_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "orden.proto",
}
