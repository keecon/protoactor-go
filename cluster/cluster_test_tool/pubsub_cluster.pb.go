// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: pubsub_cluster.proto

package cluster_test_tool

import (
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

type DataPublished struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data int32 `protobuf:"varint,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *DataPublished) Reset() {
	*x = DataPublished{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pubsub_cluster_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataPublished) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataPublished) ProtoMessage() {}

func (x *DataPublished) ProtoReflect() protoreflect.Message {
	mi := &file_pubsub_cluster_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataPublished.ProtoReflect.Descriptor instead.
func (*DataPublished) Descriptor() ([]byte, []int) {
	return file_pubsub_cluster_proto_rawDescGZIP(), []int{0}
}

func (x *DataPublished) GetData() int32 {
	if x != nil {
		return x.Data
	}
	return 0
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pubsub_cluster_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_pubsub_cluster_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_pubsub_cluster_proto_rawDescGZIP(), []int{1}
}

var File_pubsub_cluster_proto protoreflect.FileDescriptor

var file_pubsub_cluster_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x75, 0x62, 0x73, 0x75, 0x62, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f,
	0x74, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x22, 0x23, 0x0a, 0x0d, 0x44, 0x61, 0x74,
	0x61, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x0a,
	0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3e, 0x5a, 0x3c, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x73, 0x79, 0x6e, 0x6b, 0x72,
	0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2d, 0x67, 0x6f,
	0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x5f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x6f, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pubsub_cluster_proto_rawDescOnce sync.Once
	file_pubsub_cluster_proto_rawDescData = file_pubsub_cluster_proto_rawDesc
)

func file_pubsub_cluster_proto_rawDescGZIP() []byte {
	file_pubsub_cluster_proto_rawDescOnce.Do(func() {
		file_pubsub_cluster_proto_rawDescData = protoimpl.X.CompressGZIP(file_pubsub_cluster_proto_rawDescData)
	})
	return file_pubsub_cluster_proto_rawDescData
}

var file_pubsub_cluster_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pubsub_cluster_proto_goTypes = []interface{}{
	(*DataPublished)(nil), // 0: cluster_test_tool.DataPublished
	(*Response)(nil),      // 1: cluster_test_tool.Response
}
var file_pubsub_cluster_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pubsub_cluster_proto_init() }
func file_pubsub_cluster_proto_init() {
	if File_pubsub_cluster_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pubsub_cluster_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataPublished); i {
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
		file_pubsub_cluster_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_pubsub_cluster_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pubsub_cluster_proto_goTypes,
		DependencyIndexes: file_pubsub_cluster_proto_depIdxs,
		MessageInfos:      file_pubsub_cluster_proto_msgTypes,
	}.Build()
	File_pubsub_cluster_proto = out.File
	file_pubsub_cluster_proto_rawDesc = nil
	file_pubsub_cluster_proto_goTypes = nil
	file_pubsub_cluster_proto_depIdxs = nil
}
