// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: cross/v1test/cross.proto

package crossv1test

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number int64                `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Sleep  *durationpb.Duration `protobuf:"bytes,2,opt,name=sleep,proto3" json:"sleep,omitempty"`
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{0}
}

func (x *PingRequest) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *PingRequest) GetSleep() *durationpb.Duration {
	if x != nil {
		return x.Sleep
	}
	return nil
}

type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number int64 `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{1}
}

func (x *PingResponse) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

type FailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *FailRequest) Reset() {
	*x = FailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailRequest) ProtoMessage() {}

func (x *FailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailRequest.ProtoReflect.Descriptor instead.
func (*FailRequest) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{2}
}

func (x *FailRequest) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type FailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FailResponse) Reset() {
	*x = FailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailResponse) ProtoMessage() {}

func (x *FailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailResponse.ProtoReflect.Descriptor instead.
func (*FailResponse) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{3}
}

type SumRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number int64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *SumRequest) Reset() {
	*x = SumRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SumRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SumRequest) ProtoMessage() {}

func (x *SumRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SumRequest.ProtoReflect.Descriptor instead.
func (*SumRequest) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{4}
}

func (x *SumRequest) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

type SumResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sum int64 `protobuf:"varint,1,opt,name=sum,proto3" json:"sum,omitempty"`
}

func (x *SumResponse) Reset() {
	*x = SumResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SumResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SumResponse) ProtoMessage() {}

func (x *SumResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SumResponse.ProtoReflect.Descriptor instead.
func (*SumResponse) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{5}
}

func (x *SumResponse) GetSum() int64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

type CountUpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number int64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *CountUpRequest) Reset() {
	*x = CountUpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CountUpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CountUpRequest) ProtoMessage() {}

func (x *CountUpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CountUpRequest.ProtoReflect.Descriptor instead.
func (*CountUpRequest) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{6}
}

func (x *CountUpRequest) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

type CountUpResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number int64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *CountUpResponse) Reset() {
	*x = CountUpResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CountUpResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CountUpResponse) ProtoMessage() {}

func (x *CountUpResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CountUpResponse.ProtoReflect.Descriptor instead.
func (*CountUpResponse) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{7}
}

func (x *CountUpResponse) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

type CumSumRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number int64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *CumSumRequest) Reset() {
	*x = CumSumRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CumSumRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CumSumRequest) ProtoMessage() {}

func (x *CumSumRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CumSumRequest.ProtoReflect.Descriptor instead.
func (*CumSumRequest) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{8}
}

func (x *CumSumRequest) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

type CumSumResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sum int64 `protobuf:"varint,1,opt,name=sum,proto3" json:"sum,omitempty"`
}

func (x *CumSumResponse) Reset() {
	*x = CumSumResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cross_v1test_cross_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CumSumResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CumSumResponse) ProtoMessage() {}

func (x *CumSumResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cross_v1test_cross_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CumSumResponse.ProtoReflect.Descriptor instead.
func (*CumSumResponse) Descriptor() ([]byte, []int) {
	return file_cross_v1test_cross_proto_rawDescGZIP(), []int{9}
}

func (x *CumSumResponse) GetSum() int64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

var File_cross_v1test_cross_proto protoreflect.FileDescriptor

var file_cross_v1test_cross_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x63,
	0x72, 0x6f, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x63, 0x72, 0x6f, 0x73,
	0x73, 0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x0b, 0x50, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x2f, 0x0a, 0x05, 0x73, 0x6c, 0x65, 0x65, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x73, 0x6c, 0x65, 0x65, 0x70,
	0x22, 0x26, 0x0a, 0x0c, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x21, 0x0a, 0x0b, 0x46, 0x61, 0x69, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x0e, 0x0a, 0x0c, 0x46,
	0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x0a, 0x0a, 0x53,
	0x75, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x22, 0x1f, 0x0a, 0x0b, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x73, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x73,
	0x75, 0x6d, 0x22, 0x28, 0x0a, 0x0e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x55, 0x70, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x29, 0x0a, 0x0f,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x55, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x27, 0x0a, 0x0d, 0x43, 0x75, 0x6d, 0x53, 0x75,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x22, 0x22, 0x0a, 0x0e, 0x43, 0x75, 0x6d, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x03, 0x73, 0x75, 0x6d, 0x32, 0xe7, 0x02, 0x0a, 0x0c, 0x43, 0x72, 0x6f, 0x73, 0x73, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3f, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x19, 0x2e,
	0x63, 0x72, 0x6f, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73,
	0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x04, 0x46, 0x61, 0x69, 0x6c, 0x12, 0x19,
	0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x61,
	0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x63, 0x72, 0x6f, 0x73,
	0x73, 0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x03, 0x53, 0x75, 0x6d, 0x12, 0x18,
	0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x75,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73,
	0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x75, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x4a, 0x0a, 0x07, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x55, 0x70, 0x12, 0x1c, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x74, 0x65, 0x73,
	0x74, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x55, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2e,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x55, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x30, 0x01, 0x12, 0x49, 0x0a, 0x06, 0x43, 0x75, 0x6d, 0x53, 0x75, 0x6d, 0x12, 0x1b, 0x2e,
	0x63, 0x72, 0x6f, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x75, 0x6d,
	0x53, 0x75, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x63, 0x72, 0x6f,
	0x73, 0x73, 0x2e, 0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x75, 0x6d, 0x53, 0x75, 0x6d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0xc0,
	0x01, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x74,
	0x65, 0x73, 0x74, 0x42, 0x0a, 0x43, 0x72, 0x6f, 0x73, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x4f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65,
	0x72, 0x70, 0x63, 0x2f, 0x72, 0x65, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2f,
	0x76, 0x31, 0x74, 0x65, 0x73, 0x74, 0x3b, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x76, 0x31, 0x74, 0x65,
	0x73, 0x74, 0xa2, 0x02, 0x03, 0x43, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x43, 0x72, 0x6f, 0x73, 0x73,
	0x2e, 0x56, 0x31, 0x74, 0x65, 0x73, 0x74, 0xca, 0x02, 0x0c, 0x43, 0x72, 0x6f, 0x73, 0x73, 0x5c,
	0x56, 0x31, 0x74, 0x65, 0x73, 0x74, 0xe2, 0x02, 0x18, 0x43, 0x72, 0x6f, 0x73, 0x73, 0x5c, 0x56,
	0x31, 0x74, 0x65, 0x73, 0x74, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x0d, 0x43, 0x72, 0x6f, 0x73, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x74, 0x65, 0x73,
	0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cross_v1test_cross_proto_rawDescOnce sync.Once
	file_cross_v1test_cross_proto_rawDescData = file_cross_v1test_cross_proto_rawDesc
)

func file_cross_v1test_cross_proto_rawDescGZIP() []byte {
	file_cross_v1test_cross_proto_rawDescOnce.Do(func() {
		file_cross_v1test_cross_proto_rawDescData = protoimpl.X.CompressGZIP(file_cross_v1test_cross_proto_rawDescData)
	})
	return file_cross_v1test_cross_proto_rawDescData
}

var file_cross_v1test_cross_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_cross_v1test_cross_proto_goTypes = []interface{}{
	(*PingRequest)(nil),         // 0: cross.v1test.PingRequest
	(*PingResponse)(nil),        // 1: cross.v1test.PingResponse
	(*FailRequest)(nil),         // 2: cross.v1test.FailRequest
	(*FailResponse)(nil),        // 3: cross.v1test.FailResponse
	(*SumRequest)(nil),          // 4: cross.v1test.SumRequest
	(*SumResponse)(nil),         // 5: cross.v1test.SumResponse
	(*CountUpRequest)(nil),      // 6: cross.v1test.CountUpRequest
	(*CountUpResponse)(nil),     // 7: cross.v1test.CountUpResponse
	(*CumSumRequest)(nil),       // 8: cross.v1test.CumSumRequest
	(*CumSumResponse)(nil),      // 9: cross.v1test.CumSumResponse
	(*durationpb.Duration)(nil), // 10: google.protobuf.Duration
}
var file_cross_v1test_cross_proto_depIdxs = []int32{
	10, // 0: cross.v1test.PingRequest.sleep:type_name -> google.protobuf.Duration
	0,  // 1: cross.v1test.CrossService.Ping:input_type -> cross.v1test.PingRequest
	2,  // 2: cross.v1test.CrossService.Fail:input_type -> cross.v1test.FailRequest
	4,  // 3: cross.v1test.CrossService.Sum:input_type -> cross.v1test.SumRequest
	6,  // 4: cross.v1test.CrossService.CountUp:input_type -> cross.v1test.CountUpRequest
	8,  // 5: cross.v1test.CrossService.CumSum:input_type -> cross.v1test.CumSumRequest
	1,  // 6: cross.v1test.CrossService.Ping:output_type -> cross.v1test.PingResponse
	3,  // 7: cross.v1test.CrossService.Fail:output_type -> cross.v1test.FailResponse
	5,  // 8: cross.v1test.CrossService.Sum:output_type -> cross.v1test.SumResponse
	7,  // 9: cross.v1test.CrossService.CountUp:output_type -> cross.v1test.CountUpResponse
	9,  // 10: cross.v1test.CrossService.CumSum:output_type -> cross.v1test.CumSumResponse
	6,  // [6:11] is the sub-list for method output_type
	1,  // [1:6] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_cross_v1test_cross_proto_init() }
func file_cross_v1test_cross_proto_init() {
	if File_cross_v1test_cross_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cross_v1test_cross_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cross_v1test_cross_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cross_v1test_cross_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FailRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cross_v1test_cross_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FailResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cross_v1test_cross_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SumRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cross_v1test_cross_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SumResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cross_v1test_cross_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CountUpRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cross_v1test_cross_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CountUpResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cross_v1test_cross_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CumSumRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cross_v1test_cross_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CumSumResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cross_v1test_cross_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cross_v1test_cross_proto_goTypes,
		DependencyIndexes: file_cross_v1test_cross_proto_depIdxs,
		MessageInfos:      file_cross_v1test_cross_proto_msgTypes,
	}.Build()
	File_cross_v1test_cross_proto = out.File
	file_cross_v1test_cross_proto_rawDesc = nil
	file_cross_v1test_cross_proto_goTypes = nil
	file_cross_v1test_cross_proto_depIdxs = nil
}