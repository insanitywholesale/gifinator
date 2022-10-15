// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gifcreator.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Product int32

const (
	Product_UNKNOWN_PRODUCT Product = 0
	Product_GRPC            Product = 1
	Product_KUBERNETES      Product = 2
	Product_GO              Product = 3
)

var Product_name = map[int32]string{
	0: "UNKNOWN_PRODUCT",
	1: "GRPC",
	2: "KUBERNETES",
	3: "GO",
}

var Product_value = map[string]int32{
	"UNKNOWN_PRODUCT": 0,
	"GRPC":            1,
	"KUBERNETES":      2,
	"GO":              3,
}

func (x Product) String() string {
	return proto.EnumName(Product_name, int32(x))
}

func (Product) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_52bf8112deee5818, []int{0}
}

type GetJobResponse_Status int32

const (
	GetJobResponse_UNKNOWN_STATUS GetJobResponse_Status = 0
	GetJobResponse_PENDING        GetJobResponse_Status = 1
	GetJobResponse_DONE           GetJobResponse_Status = 2
	GetJobResponse_FAILED         GetJobResponse_Status = 3
)

var GetJobResponse_Status_name = map[int32]string{
	0: "UNKNOWN_STATUS",
	1: "PENDING",
	2: "DONE",
	3: "FAILED",
}

var GetJobResponse_Status_value = map[string]int32{
	"UNKNOWN_STATUS": 0,
	"PENDING":        1,
	"DONE":           2,
	"FAILED":         3,
}

func (x GetJobResponse_Status) String() string {
	return proto.EnumName(GetJobResponse_Status_name, int32(x))
}

func (GetJobResponse_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_52bf8112deee5818, []int{3, 0}
}

type StartJobRequest struct {
	// TODO(light): what scene parameters do we want to give?
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ProductToPlug        Product  `protobuf:"varint,2,opt,name=product_to_plug,json=productToPlug,proto3,enum=renderdemo.Product" json:"product_to_plug,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartJobRequest) Reset()         { *m = StartJobRequest{} }
func (m *StartJobRequest) String() string { return proto.CompactTextString(m) }
func (*StartJobRequest) ProtoMessage()    {}
func (*StartJobRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_52bf8112deee5818, []int{0}
}

func (m *StartJobRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartJobRequest.Unmarshal(m, b)
}
func (m *StartJobRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartJobRequest.Marshal(b, m, deterministic)
}
func (m *StartJobRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartJobRequest.Merge(m, src)
}
func (m *StartJobRequest) XXX_Size() int {
	return xxx_messageInfo_StartJobRequest.Size(m)
}
func (m *StartJobRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StartJobRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StartJobRequest proto.InternalMessageInfo

func (m *StartJobRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *StartJobRequest) GetProductToPlug() Product {
	if m != nil {
		return m.ProductToPlug
	}
	return Product_UNKNOWN_PRODUCT
}

type StartJobResponse struct {
	JobId                string   `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartJobResponse) Reset()         { *m = StartJobResponse{} }
func (m *StartJobResponse) String() string { return proto.CompactTextString(m) }
func (*StartJobResponse) ProtoMessage()    {}
func (*StartJobResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_52bf8112deee5818, []int{1}
}

func (m *StartJobResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartJobResponse.Unmarshal(m, b)
}
func (m *StartJobResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartJobResponse.Marshal(b, m, deterministic)
}
func (m *StartJobResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartJobResponse.Merge(m, src)
}
func (m *StartJobResponse) XXX_Size() int {
	return xxx_messageInfo_StartJobResponse.Size(m)
}
func (m *StartJobResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StartJobResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StartJobResponse proto.InternalMessageInfo

func (m *StartJobResponse) GetJobId() string {
	if m != nil {
		return m.JobId
	}
	return ""
}

type GetJobRequest struct {
	JobId                string   `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetJobRequest) Reset()         { *m = GetJobRequest{} }
func (m *GetJobRequest) String() string { return proto.CompactTextString(m) }
func (*GetJobRequest) ProtoMessage()    {}
func (*GetJobRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_52bf8112deee5818, []int{2}
}

func (m *GetJobRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetJobRequest.Unmarshal(m, b)
}
func (m *GetJobRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetJobRequest.Marshal(b, m, deterministic)
}
func (m *GetJobRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetJobRequest.Merge(m, src)
}
func (m *GetJobRequest) XXX_Size() int {
	return xxx_messageInfo_GetJobRequest.Size(m)
}
func (m *GetJobRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetJobRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetJobRequest proto.InternalMessageInfo

func (m *GetJobRequest) GetJobId() string {
	if m != nil {
		return m.JobId
	}
	return ""
}

type GetJobResponse struct {
	Status GetJobResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=renderdemo.GetJobResponse_Status" json:"status,omitempty"`
	// World-readable URL for created image.
	ImageUrl             string   `protobuf:"bytes,2,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetJobResponse) Reset()         { *m = GetJobResponse{} }
func (m *GetJobResponse) String() string { return proto.CompactTextString(m) }
func (*GetJobResponse) ProtoMessage()    {}
func (*GetJobResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_52bf8112deee5818, []int{3}
}

func (m *GetJobResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetJobResponse.Unmarshal(m, b)
}
func (m *GetJobResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetJobResponse.Marshal(b, m, deterministic)
}
func (m *GetJobResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetJobResponse.Merge(m, src)
}
func (m *GetJobResponse) XXX_Size() int {
	return xxx_messageInfo_GetJobResponse.Size(m)
}
func (m *GetJobResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetJobResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetJobResponse proto.InternalMessageInfo

func (m *GetJobResponse) GetStatus() GetJobResponse_Status {
	if m != nil {
		return m.Status
	}
	return GetJobResponse_UNKNOWN_STATUS
}

func (m *GetJobResponse) GetImageUrl() string {
	if m != nil {
		return m.ImageUrl
	}
	return ""
}

func init() {
	proto.RegisterEnum("renderdemo.Product", Product_name, Product_value)
	proto.RegisterEnum("renderdemo.GetJobResponse_Status", GetJobResponse_Status_name, GetJobResponse_Status_value)
	proto.RegisterType((*StartJobRequest)(nil), "renderdemo.StartJobRequest")
	proto.RegisterType((*StartJobResponse)(nil), "renderdemo.StartJobResponse")
	proto.RegisterType((*GetJobRequest)(nil), "renderdemo.GetJobRequest")
	proto.RegisterType((*GetJobResponse)(nil), "renderdemo.GetJobResponse")
}

func init() { proto.RegisterFile("gifcreator.proto", fileDescriptor_52bf8112deee5818) }

var fileDescriptor_52bf8112deee5818 = []byte{
	// 418 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0xeb, 0xb4, 0xb8, 0xc9, 0xa0, 0x3a, 0xab, 0xa9, 0x90, 0x4a, 0xcb, 0xa1, 0xe4, 0x80,
	0x0a, 0x42, 0x8e, 0x14, 0x4e, 0x88, 0x43, 0x69, 0x13, 0x63, 0x85, 0x22, 0xc7, 0x5a, 0xdb, 0x42,
	0xe2, 0x62, 0xad, 0xe3, 0x8d, 0xd9, 0xca, 0xf6, 0x9a, 0xf5, 0x5a, 0x88, 0xf7, 0xe0, 0x25, 0x78,
	0x4b, 0x54, 0xdb, 0x51, 0x03, 0x6a, 0x4f, 0xb6, 0x76, 0xbf, 0xf9, 0xe7, 0x9f, 0x7f, 0x16, 0x48,
	0x26, 0x36, 0x6b, 0xc5, 0x99, 0x96, 0xca, 0xae, 0x94, 0xd4, 0x12, 0x41, 0xf1, 0x32, 0xe5, 0x2a,
	0xe5, 0x85, 0x9c, 0x24, 0x30, 0x0e, 0x34, 0x53, 0xfa, 0xb3, 0x4c, 0x28, 0xff, 0xd1, 0xf0, 0x5a,
	0x23, 0xc2, 0x41, 0xc9, 0x0a, 0x7e, 0x62, 0x9c, 0x1b, 0x17, 0x23, 0xda, 0xfe, 0xe3, 0x07, 0x18,
	0x57, 0x4a, 0xa6, 0xcd, 0x5a, 0xc7, 0x5a, 0xc6, 0x55, 0xde, 0x64, 0x27, 0x83, 0x73, 0xe3, 0xc2,
	0x9a, 0x1d, 0xdb, 0xf7, 0x62, 0xb6, 0xdf, 0x21, 0xf4, 0xa8, 0x67, 0x43, 0xe9, 0xe7, 0x4d, 0x36,
	0x79, 0x0d, 0xe4, 0xbe, 0x47, 0x5d, 0xc9, 0xb2, 0xe6, 0xf8, 0x0c, 0xcc, 0x5b, 0x99, 0xc4, 0x22,
	0xed, 0xdb, 0x3c, 0xb9, 0x95, 0xc9, 0x32, 0x9d, 0xbc, 0x82, 0x23, 0x97, 0xef, 0x9a, 0x79, 0x84,
	0xfb, 0x63, 0x80, 0xb5, 0x05, 0x7b, 0xc5, 0xf7, 0x60, 0xd6, 0x9a, 0xe9, 0xa6, 0x6e, 0x49, 0x6b,
	0xf6, 0x72, 0xd7, 0xd9, 0xbf, 0xac, 0x1d, 0xb4, 0x20, 0xed, 0x0b, 0xf0, 0x0c, 0x46, 0xa2, 0x60,
	0x19, 0x8f, 0x1b, 0x95, 0xb7, 0x73, 0x8d, 0xe8, 0xb0, 0x3d, 0x88, 0x54, 0x3e, 0xb9, 0x04, 0xb3,
	0xc3, 0x11, 0xc1, 0x8a, 0xbc, 0x1b, 0x6f, 0xf5, 0xd5, 0x8b, 0x83, 0xf0, 0x2a, 0x8c, 0x02, 0xb2,
	0x87, 0x4f, 0xe1, 0xd0, 0x77, 0xbc, 0xc5, 0xd2, 0x73, 0x89, 0x81, 0x43, 0x38, 0x58, 0xac, 0x3c,
	0x87, 0x0c, 0x10, 0xc0, 0xfc, 0x74, 0xb5, 0xfc, 0xe2, 0x2c, 0xc8, 0xfe, 0x9b, 0x8f, 0x70, 0xd8,
	0x07, 0x83, 0xc7, 0x30, 0xde, 0x2a, 0xf8, 0x74, 0xb5, 0x88, 0xe6, 0x21, 0xd9, 0xbb, 0xab, 0x72,
	0xa9, 0x3f, 0x27, 0x06, 0x5a, 0x00, 0x37, 0xd1, 0xb5, 0x43, 0x3d, 0x27, 0x74, 0x02, 0x32, 0x40,
	0x13, 0x06, 0xee, 0x8a, 0xec, 0xcf, 0x7e, 0x1b, 0x00, 0xae, 0xd8, 0xcc, 0xbb, 0x2d, 0xa2, 0x03,
	0xc3, 0x6d, 0x9e, 0x78, 0xb6, 0x3b, 0xe5, 0x7f, 0x9b, 0x3c, 0x7d, 0xf1, 0xf0, 0x65, 0x1f, 0xd8,
	0x25, 0x98, 0x5d, 0x2c, 0xf8, 0xfc, 0xa1, 0xa8, 0x3a, 0x89, 0xd3, 0xc7, 0x53, 0xbc, 0xb6, 0xbf,
	0xbd, 0xcd, 0x84, 0xce, 0x59, 0x62, 0xaf, 0x65, 0x31, 0x15, 0x65, 0xcd, 0x4a, 0xa1, 0x7f, 0xfd,
	0xfc, 0x2e, 0x73, 0x5e, 0xb3, 0x9c, 0x4f, 0x33, 0xb1, 0x11, 0xe5, 0x9d, 0xe3, 0x69, 0xfb, 0xee,
	0x12, 0xb3, 0xfd, 0xbc, 0xfb, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x07, 0xb3, 0xe3, 0x92, 0x02,
	0x00, 0x00,
}
