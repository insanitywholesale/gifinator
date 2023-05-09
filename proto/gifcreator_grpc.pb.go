// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: gifcreator.proto

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

// GifCreatorClient is the client API for GifCreator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GifCreatorClient interface {
	StartJob(ctx context.Context, in *StartJobRequest, opts ...grpc.CallOption) (*StartJobResponse, error)
	GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*GetJobResponse, error)
}

type gifCreatorClient struct {
	cc grpc.ClientConnInterface
}

func NewGifCreatorClient(cc grpc.ClientConnInterface) GifCreatorClient {
	return &gifCreatorClient{cc}
}

func (c *gifCreatorClient) StartJob(ctx context.Context, in *StartJobRequest, opts ...grpc.CallOption) (*StartJobResponse, error) {
	out := new(StartJobResponse)
	err := c.cc.Invoke(ctx, "/renderdemo.GifCreator/StartJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gifCreatorClient) GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*GetJobResponse, error) {
	out := new(GetJobResponse)
	err := c.cc.Invoke(ctx, "/renderdemo.GifCreator/GetJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GifCreatorServer is the server API for GifCreator service.
// All implementations must embed UnimplementedGifCreatorServer
// for forward compatibility
type GifCreatorServer interface {
	StartJob(context.Context, *StartJobRequest) (*StartJobResponse, error)
	GetJob(context.Context, *GetJobRequest) (*GetJobResponse, error)
	mustEmbedUnimplementedGifCreatorServer()
}

// UnimplementedGifCreatorServer must be embedded to have forward compatible implementations.
type UnimplementedGifCreatorServer struct {
}

func (UnimplementedGifCreatorServer) StartJob(context.Context, *StartJobRequest) (*StartJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartJob not implemented")
}
func (UnimplementedGifCreatorServer) GetJob(context.Context, *GetJobRequest) (*GetJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJob not implemented")
}
func (UnimplementedGifCreatorServer) mustEmbedUnimplementedGifCreatorServer() {}

// UnsafeGifCreatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GifCreatorServer will
// result in compilation errors.
type UnsafeGifCreatorServer interface {
	mustEmbedUnimplementedGifCreatorServer()
}

func RegisterGifCreatorServer(s grpc.ServiceRegistrar, srv GifCreatorServer) {
	s.RegisterService(&GifCreator_ServiceDesc, srv)
}

func _GifCreator_StartJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GifCreatorServer).StartJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/renderdemo.GifCreator/StartJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GifCreatorServer).StartJob(ctx, req.(*StartJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GifCreator_GetJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GifCreatorServer).GetJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/renderdemo.GifCreator/GetJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GifCreatorServer).GetJob(ctx, req.(*GetJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GifCreator_ServiceDesc is the grpc.ServiceDesc for GifCreator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GifCreator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "renderdemo.GifCreator",
	HandlerType: (*GifCreatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartJob",
			Handler:    _GifCreator_StartJob_Handler,
		},
		{
			MethodName: "GetJob",
			Handler:    _GifCreator_GetJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gifcreator.proto",
}