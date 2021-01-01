// Code generated by protoc-gen-go.
// source: proto/render.proto
// DO NOT EDIT!

package renderdemo

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type RenderRequest struct {
	// GCS path to write output image into.
	GcsOutputBase string `protobuf:"bytes,1,opt,name=gcs_output_base,json=gcsOutputBase" json:"gcs_output_base,omitempty"`
	// scene object file (in .obj format) GCS path
	ObjPath string `protobuf:"bytes,2,opt,name=obj_path,json=objPath" json:"obj_path,omitempty"`
	// assets (like material files and images) to be associated with the object
	Assets []string `protobuf:"bytes,3,rep,name=assets" json:"assets,omitempty"`
	// scene rotation (in radians)
	Rotation float32 `protobuf:"fixed32,4,opt,name=rotation" json:"rotation,omitempty"`
	// num iterations
	Iterations int32 `protobuf:"varint,5,opt,name=iterations" json:"iterations,omitempty"`
}

func (m *RenderRequest) Reset()                    { *m = RenderRequest{} }
func (m *RenderRequest) String() string            { return proto.CompactTextString(m) }
func (*RenderRequest) ProtoMessage()               {}
func (*RenderRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *RenderRequest) GetGcsOutputBase() string {
	if m != nil {
		return m.GcsOutputBase
	}
	return ""
}

func (m *RenderRequest) GetObjPath() string {
	if m != nil {
		return m.ObjPath
	}
	return ""
}

func (m *RenderRequest) GetAssets() []string {
	if m != nil {
		return m.Assets
	}
	return nil
}

func (m *RenderRequest) GetRotation() float32 {
	if m != nil {
		return m.Rotation
	}
	return 0
}

func (m *RenderRequest) GetIterations() int32 {
	if m != nil {
		return m.Iterations
	}
	return 0
}

type RenderResponse struct {
	// GCS path image was written to.
	GcsOutput string `protobuf:"bytes,1,opt,name=gcs_output,json=gcsOutput" json:"gcs_output,omitempty"`
}

func (m *RenderResponse) Reset()                    { *m = RenderResponse{} }
func (m *RenderResponse) String() string            { return proto.CompactTextString(m) }
func (*RenderResponse) ProtoMessage()               {}
func (*RenderResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *RenderResponse) GetGcsOutput() string {
	if m != nil {
		return m.GcsOutput
	}
	return ""
}

func init() {
	proto.RegisterType((*RenderRequest)(nil), "renderdemo.RenderRequest")
	proto.RegisterType((*RenderResponse)(nil), "renderdemo.RenderResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Render service

type RenderClient interface {
	RenderFrame(ctx context.Context, in *RenderRequest, opts ...grpc.CallOption) (*RenderResponse, error)
}

type renderClient struct {
	cc *grpc.ClientConn
}

func NewRenderClient(cc *grpc.ClientConn) RenderClient {
	return &renderClient{cc}
}

func (c *renderClient) RenderFrame(ctx context.Context, in *RenderRequest, opts ...grpc.CallOption) (*RenderResponse, error) {
	out := new(RenderResponse)
	err := grpc.Invoke(ctx, "/renderdemo.Render/RenderFrame", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Render service

type RenderServer interface {
	RenderFrame(context.Context, *RenderRequest) (*RenderResponse, error)
}

func RegisterRenderServer(s *grpc.Server, srv RenderServer) {
	s.RegisterService(&_Render_serviceDesc, srv)
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

var _Render_serviceDesc = grpc.ServiceDesc{
	ServiceName: "renderdemo.Render",
	HandlerType: (*RenderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RenderFrame",
			Handler:    _Render_RenderFrame_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/render.proto",
}

func init() { proto.RegisterFile("proto/render.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x90, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0xc6, 0xd9, 0xd6, 0xc6, 0x66, 0xb4, 0x0a, 0x73, 0x90, 0x6d, 0x40, 0x09, 0x3d, 0x48, 0x4e,
	0x29, 0xe8, 0x1b, 0x14, 0xf1, 0xa8, 0xb2, 0x47, 0x2f, 0x61, 0xd3, 0x0e, 0xfd, 0x03, 0xcd, 0xc4,
	0x9d, 0xc9, 0x2b, 0xf9, 0x9c, 0xc2, 0xa6, 0xb6, 0x15, 0xbc, 0xed, 0x6f, 0x76, 0x97, 0xf9, 0x7d,
	0x1f, 0x60, 0x1b, 0x58, 0x79, 0x1e, 0xa8, 0x59, 0x51, 0x28, 0x23, 0x20, 0xf4, 0xb4, 0xa2, 0x3d,
	0xcf, 0xbe, 0x0d, 0x4c, 0x5c, 0x44, 0x47, 0x5f, 0x1d, 0x89, 0xe2, 0x23, 0xdc, 0xae, 0x97, 0x52,
	0x71, 0xa7, 0x6d, 0xa7, 0x55, 0xed, 0x85, 0xac, 0xc9, 0x4d, 0x91, 0xba, 0xc9, 0x7a, 0x29, 0xef,
	0x71, 0xba, 0xf0, 0x42, 0x38, 0x85, 0x31, 0xd7, 0xbb, 0xaa, 0xf5, 0xba, 0xb1, 0x83, 0xf8, 0xe0,
	0x92, 0xeb, 0xdd, 0x87, 0xd7, 0x0d, 0xde, 0x41, 0xe2, 0x45, 0x48, 0xc5, 0x0e, 0xf3, 0x61, 0x91,
	0xba, 0x03, 0x61, 0x06, 0xe3, 0xc0, 0xea, 0x75, 0xcb, 0x8d, 0xbd, 0xc8, 0x4d, 0x31, 0x70, 0x47,
	0xc6, 0x07, 0x80, 0xad, 0x52, 0x88, 0x20, 0x76, 0x94, 0x9b, 0x62, 0xe4, 0xce, 0x26, 0xb3, 0x39,
	0xdc, 0xfc, 0x7a, 0x4a, 0xcb, 0x8d, 0x10, 0xde, 0x03, 0x9c, 0x44, 0x0f, 0x8e, 0xe9, 0xd1, 0xf1,
	0xe9, 0x0d, 0x92, 0xfe, 0x03, 0xbe, 0xc0, 0x55, 0x7f, 0x7a, 0x0d, 0x7e, 0x4f, 0x38, 0x2d, 0x4f,
	0xf9, 0xcb, 0x3f, 0xd9, 0xb3, 0xec, 0xbf, 0xab, 0x7e, 0xdd, 0xe2, 0xfa, 0xf3, 0xac, 0xb7, 0x3a,
	0x89, 0x55, 0x3e, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0xa4, 0x68, 0x4a, 0x92, 0x60, 0x01, 0x00,
	0x00,
}