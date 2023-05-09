// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: render.proto

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

// RenderClient is the client API for Render service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RenderClient interface {
	RenderFrame(ctx context.Context, in *RenderRequest, opts ...grpc.CallOption) (*RenderResponse, error)
}

type renderClient struct {
	cc grpc.ClientConnInterface
}

func NewRenderClient(cc grpc.ClientConnInterface) RenderClient {
	return &renderClient{cc}
}

func (c *renderClient) RenderFrame(ctx context.Context, in *RenderRequest, opts ...grpc.CallOption) (*RenderResponse, error) {
	out := new(RenderResponse)
	err := c.cc.Invoke(ctx, "/renderdemo.Render/RenderFrame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RenderServer is the server API for Render service.
// All implementations must embed UnimplementedRenderServer
// for forward compatibility
type RenderServer interface {
	RenderFrame(context.Context, *RenderRequest) (*RenderResponse, error)
	mustEmbedUnimplementedRenderServer()
}

// UnimplementedRenderServer must be embedded to have forward compatible implementations.
type UnimplementedRenderServer struct {
}

func (UnimplementedRenderServer) RenderFrame(context.Context, *RenderRequest) (*RenderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenderFrame not implemented")
}
func (UnimplementedRenderServer) mustEmbedUnimplementedRenderServer() {}

// UnsafeRenderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RenderServer will
// result in compilation errors.
type UnsafeRenderServer interface {
	mustEmbedUnimplementedRenderServer()
}

func RegisterRenderServer(s grpc.ServiceRegistrar, srv RenderServer) {
	s.RegisterService(&Render_ServiceDesc, srv)
}

func _Render_RenderFrame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RenderServer).RenderFrame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/renderdemo.Render/RenderFrame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RenderServer).RenderFrame(ctx, req.(*RenderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Render_ServiceDesc is the grpc.ServiceDesc for Render service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Render_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "renderdemo.Render",
	HandlerType: (*RenderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RenderFrame",
			Handler:    _Render_RenderFrame_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "render.proto",
}
