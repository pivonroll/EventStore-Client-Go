// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.14.0
// source: operations.proto

package operations

import (
	shared "github.com/pivonroll/EventStore-Client-Go/protos/shared"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ScavengeResp_ScavengeResult int32

const (
	ScavengeResp_Started    ScavengeResp_ScavengeResult = 0
	ScavengeResp_InProgress ScavengeResp_ScavengeResult = 1
	ScavengeResp_Stopped    ScavengeResp_ScavengeResult = 2
)

// Enum value maps for ScavengeResp_ScavengeResult.
var (
	ScavengeResp_ScavengeResult_name = map[int32]string{
		0: "Started",
		1: "InProgress",
		2: "Stopped",
	}
	ScavengeResp_ScavengeResult_value = map[string]int32{
		"Started":    0,
		"InProgress": 1,
		"Stopped":    2,
	}
)

func (x ScavengeResp_ScavengeResult) Enum() *ScavengeResp_ScavengeResult {
	p := new(ScavengeResp_ScavengeResult)
	*p = x
	return p
}

func (x ScavengeResp_ScavengeResult) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ScavengeResp_ScavengeResult) Descriptor() protoreflect.EnumDescriptor {
	return file_operations_proto_enumTypes[0].Descriptor()
}

func (ScavengeResp_ScavengeResult) Type() protoreflect.EnumType {
	return &file_operations_proto_enumTypes[0]
}

func (x ScavengeResp_ScavengeResult) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ScavengeResp_ScavengeResult.Descriptor instead.
func (ScavengeResp_ScavengeResult) EnumDescriptor() ([]byte, []int) {
	return file_operations_proto_rawDescGZIP(), []int{2, 0}
}

type StartScavengeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options *StartScavengeReq_Options `protobuf:"bytes,1,opt,name=options,proto3" json:"options,omitempty"`
}

func (x *StartScavengeReq) Reset() {
	*x = StartScavengeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartScavengeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartScavengeReq) ProtoMessage() {}

func (x *StartScavengeReq) ProtoReflect() protoreflect.Message {
	mi := &file_operations_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartScavengeReq.ProtoReflect.Descriptor instead.
func (*StartScavengeReq) Descriptor() ([]byte, []int) {
	return file_operations_proto_rawDescGZIP(), []int{0}
}

func (x *StartScavengeReq) GetOptions() *StartScavengeReq_Options {
	if x != nil {
		return x.Options
	}
	return nil
}

type StopScavengeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options *StopScavengeReq_Options `protobuf:"bytes,1,opt,name=options,proto3" json:"options,omitempty"`
}

func (x *StopScavengeReq) Reset() {
	*x = StopScavengeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopScavengeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopScavengeReq) ProtoMessage() {}

func (x *StopScavengeReq) ProtoReflect() protoreflect.Message {
	mi := &file_operations_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopScavengeReq.ProtoReflect.Descriptor instead.
func (*StopScavengeReq) Descriptor() ([]byte, []int) {
	return file_operations_proto_rawDescGZIP(), []int{1}
}

func (x *StopScavengeReq) GetOptions() *StopScavengeReq_Options {
	if x != nil {
		return x.Options
	}
	return nil
}

type ScavengeResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ScavengeId     string                      `protobuf:"bytes,1,opt,name=scavenge_id,json=scavengeId,proto3" json:"scavenge_id,omitempty"`
	ScavengeResult ScavengeResp_ScavengeResult `protobuf:"varint,2,opt,name=scavenge_result,json=scavengeResult,proto3,enum=event_store.client.operations.ScavengeResp_ScavengeResult" json:"scavenge_result,omitempty"`
}

func (x *ScavengeResp) Reset() {
	*x = ScavengeResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScavengeResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScavengeResp) ProtoMessage() {}

func (x *ScavengeResp) ProtoReflect() protoreflect.Message {
	mi := &file_operations_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScavengeResp.ProtoReflect.Descriptor instead.
func (*ScavengeResp) Descriptor() ([]byte, []int) {
	return file_operations_proto_rawDescGZIP(), []int{2}
}

func (x *ScavengeResp) GetScavengeId() string {
	if x != nil {
		return x.ScavengeId
	}
	return ""
}

func (x *ScavengeResp) GetScavengeResult() ScavengeResp_ScavengeResult {
	if x != nil {
		return x.ScavengeResult
	}
	return ScavengeResp_Started
}

type SetNodePriorityReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Priority int32 `protobuf:"varint,1,opt,name=priority,proto3" json:"priority,omitempty"`
}

func (x *SetNodePriorityReq) Reset() {
	*x = SetNodePriorityReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetNodePriorityReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetNodePriorityReq) ProtoMessage() {}

func (x *SetNodePriorityReq) ProtoReflect() protoreflect.Message {
	mi := &file_operations_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetNodePriorityReq.ProtoReflect.Descriptor instead.
func (*SetNodePriorityReq) Descriptor() ([]byte, []int) {
	return file_operations_proto_rawDescGZIP(), []int{3}
}

func (x *SetNodePriorityReq) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

type StartScavengeReq_Options struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ThreadCount    int32 `protobuf:"varint,1,opt,name=thread_count,json=threadCount,proto3" json:"thread_count,omitempty"`
	StartFromChunk int32 `protobuf:"varint,2,opt,name=start_from_chunk,json=startFromChunk,proto3" json:"start_from_chunk,omitempty"`
}

func (x *StartScavengeReq_Options) Reset() {
	*x = StartScavengeReq_Options{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartScavengeReq_Options) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartScavengeReq_Options) ProtoMessage() {}

func (x *StartScavengeReq_Options) ProtoReflect() protoreflect.Message {
	mi := &file_operations_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartScavengeReq_Options.ProtoReflect.Descriptor instead.
func (*StartScavengeReq_Options) Descriptor() ([]byte, []int) {
	return file_operations_proto_rawDescGZIP(), []int{0, 0}
}

func (x *StartScavengeReq_Options) GetThreadCount() int32 {
	if x != nil {
		return x.ThreadCount
	}
	return 0
}

func (x *StartScavengeReq_Options) GetStartFromChunk() int32 {
	if x != nil {
		return x.StartFromChunk
	}
	return 0
}

type StopScavengeReq_Options struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ScavengeId string `protobuf:"bytes,1,opt,name=scavenge_id,json=scavengeId,proto3" json:"scavenge_id,omitempty"`
}

func (x *StopScavengeReq_Options) Reset() {
	*x = StopScavengeReq_Options{}
	if protoimpl.UnsafeEnabled {
		mi := &file_operations_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopScavengeReq_Options) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopScavengeReq_Options) ProtoMessage() {}

func (x *StopScavengeReq_Options) ProtoReflect() protoreflect.Message {
	mi := &file_operations_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopScavengeReq_Options.ProtoReflect.Descriptor instead.
func (*StopScavengeReq_Options) Descriptor() ([]byte, []int) {
	return file_operations_proto_rawDescGZIP(), []int{1, 0}
}

func (x *StopScavengeReq_Options) GetScavengeId() string {
	if x != nil {
		return x.ScavengeId
	}
	return ""
}

var File_operations_proto protoreflect.FileDescriptor

var file_operations_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1d, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x1a, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xbd, 0x01, 0x0a, 0x10, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x12, 0x51, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63, 0x61, 0x76, 0x65,
	0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x56, 0x0a, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x68, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x74, 0x68, 0x72, 0x65, 0x61, 0x64,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x66,
	0x72, 0x6f, 0x6d, 0x5f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x22,
	0x8f, 0x01, 0x0a, 0x0f, 0x53, 0x74, 0x6f, 0x70, 0x53, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x12, 0x50, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x36, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x53, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x2a, 0x0a, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67, 0x65, 0x49,
	0x64, 0x22, 0xd0, 0x01, 0x0a, 0x0c, 0x53, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67,
	0x65, 0x49, 0x64, 0x12, 0x63, 0x0a, 0x0f, 0x73, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67, 0x65, 0x5f,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x3a, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53, 0x63, 0x61,
	0x76, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x2e, 0x53, 0x63, 0x61, 0x76, 0x65, 0x6e,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x0e, 0x73, 0x63, 0x61, 0x76, 0x65, 0x6e,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x3a, 0x0a, 0x0e, 0x53, 0x63, 0x61, 0x76,
	0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x65, 0x64, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x49, 0x6e, 0x50, 0x72, 0x6f,
	0x67, 0x72, 0x65, 0x73, 0x73, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x74, 0x6f, 0x70, 0x70,
	0x65, 0x64, 0x10, 0x02, 0x22, 0x30, 0x0a, 0x12, 0x53, 0x65, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x50,
	0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x32, 0xac, 0x05, 0x0a, 0x0a, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x6d, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63,
	0x61, 0x76, 0x65, 0x6e, 0x67, 0x65, 0x12, 0x2f, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x63, 0x61, 0x76,
	0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x2b, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x6b, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x70, 0x53, 0x63, 0x61, 0x76,
	0x65, 0x6e, 0x67, 0x65, 0x12, 0x2e, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x53, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x1a, 0x2b, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53, 0x63, 0x61, 0x76, 0x65, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x4e, 0x0a, 0x08, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x12, 0x20, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x20, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x12, 0x52, 0x0a, 0x0c, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65,
	0x73, 0x12, 0x20, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x20, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x50, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x4e,
	0x6f, 0x64, 0x65, 0x12, 0x20, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x20, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65,
	0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x66, 0x0a, 0x0f, 0x53, 0x65, 0x74, 0x4e, 0x6f,
	0x64, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x31, 0x2e, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x53, 0x65, 0x74, 0x4e, 0x6f,
	0x64, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x20, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x64, 0x0a, 0x1e, 0x52, 0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73,
	0x74, 0x65, 0x6e, 0x74, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x20, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x20, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x3d, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x69, 0x76, 0x6f, 0x6e, 0x72, 0x6f, 0x6c, 0x6c, 0x2f, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x2d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2d,
	0x47, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_operations_proto_rawDescOnce sync.Once
	file_operations_proto_rawDescData = file_operations_proto_rawDesc
)

func file_operations_proto_rawDescGZIP() []byte {
	file_operations_proto_rawDescOnce.Do(func() {
		file_operations_proto_rawDescData = protoimpl.X.CompressGZIP(file_operations_proto_rawDescData)
	})
	return file_operations_proto_rawDescData
}

var file_operations_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_operations_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_operations_proto_goTypes = []interface{}{
	(ScavengeResp_ScavengeResult)(0), // 0: event_store.client.operations.ScavengeResp.ScavengeResult
	(*StartScavengeReq)(nil),         // 1: event_store.client.operations.StartScavengeReq
	(*StopScavengeReq)(nil),          // 2: event_store.client.operations.StopScavengeReq
	(*ScavengeResp)(nil),             // 3: event_store.client.operations.ScavengeResp
	(*SetNodePriorityReq)(nil),       // 4: event_store.client.operations.SetNodePriorityReq
	(*StartScavengeReq_Options)(nil), // 5: event_store.client.operations.StartScavengeReq.Options
	(*StopScavengeReq_Options)(nil),  // 6: event_store.client.operations.StopScavengeReq.Options
	(*shared.Empty)(nil),             // 7: event_store.client.shared.Empty
}
var file_operations_proto_depIdxs = []int32{
	5,  // 0: event_store.client.operations.StartScavengeReq.options:type_name -> event_store.client.operations.StartScavengeReq.Options
	6,  // 1: event_store.client.operations.StopScavengeReq.options:type_name -> event_store.client.operations.StopScavengeReq.Options
	0,  // 2: event_store.client.operations.ScavengeResp.scavenge_result:type_name -> event_store.client.operations.ScavengeResp.ScavengeResult
	1,  // 3: event_store.client.operations.Operations.StartScavenge:input_type -> event_store.client.operations.StartScavengeReq
	2,  // 4: event_store.client.operations.Operations.StopScavenge:input_type -> event_store.client.operations.StopScavengeReq
	7,  // 5: event_store.client.operations.Operations.Shutdown:input_type -> event_store.client.shared.Empty
	7,  // 6: event_store.client.operations.Operations.MergeIndexes:input_type -> event_store.client.shared.Empty
	7,  // 7: event_store.client.operations.Operations.ResignNode:input_type -> event_store.client.shared.Empty
	4,  // 8: event_store.client.operations.Operations.SetNodePriority:input_type -> event_store.client.operations.SetNodePriorityReq
	7,  // 9: event_store.client.operations.Operations.RestartPersistentSubscriptions:input_type -> event_store.client.shared.Empty
	3,  // 10: event_store.client.operations.Operations.StartScavenge:output_type -> event_store.client.operations.ScavengeResp
	3,  // 11: event_store.client.operations.Operations.StopScavenge:output_type -> event_store.client.operations.ScavengeResp
	7,  // 12: event_store.client.operations.Operations.Shutdown:output_type -> event_store.client.shared.Empty
	7,  // 13: event_store.client.operations.Operations.MergeIndexes:output_type -> event_store.client.shared.Empty
	7,  // 14: event_store.client.operations.Operations.ResignNode:output_type -> event_store.client.shared.Empty
	7,  // 15: event_store.client.operations.Operations.SetNodePriority:output_type -> event_store.client.shared.Empty
	7,  // 16: event_store.client.operations.Operations.RestartPersistentSubscriptions:output_type -> event_store.client.shared.Empty
	10, // [10:17] is the sub-list for method output_type
	3,  // [3:10] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_operations_proto_init() }
func file_operations_proto_init() {
	if File_operations_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_operations_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartScavengeReq); i {
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
		file_operations_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopScavengeReq); i {
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
		file_operations_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScavengeResp); i {
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
		file_operations_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetNodePriorityReq); i {
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
		file_operations_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartScavengeReq_Options); i {
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
		file_operations_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopScavengeReq_Options); i {
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
			RawDescriptor: file_operations_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_operations_proto_goTypes,
		DependencyIndexes: file_operations_proto_depIdxs,
		EnumInfos:         file_operations_proto_enumTypes,
		MessageInfos:      file_operations_proto_msgTypes,
	}.Build()
	File_operations_proto = out.File
	file_operations_proto_rawDesc = nil
	file_operations_proto_goTypes = nil
	file_operations_proto_depIdxs = nil
}
