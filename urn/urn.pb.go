// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: urn/urn.proto

package urn

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

// A ResourceURN represents a unique identifier of a resource in a service.
type ResourceURN struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// canonical means the full representation of the URN
	Canonical string `protobuf:"bytes,1,opt,name=canonical,proto3" json:"canonical,omitempty"`
	// namespace is the namespace name, which reflects a collection
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// partition is the partition name, which reflects a collection
	Partition string `protobuf:"bytes,3,opt,name=partition,proto3" json:"partition,omitempty"`
	// service is the service name, which reflects a collection
	Service string `protobuf:"bytes,4,opt,name=service,proto3" json:"service,omitempty"`
	// region is the service region, which reflects a collection
	Region string `protobuf:"bytes,5,opt,name=region,proto3" json:"region,omitempty"`
	// identifier is the identifier of the resource within a collection
	Identifier string `protobuf:"bytes,6,opt,name=identifier,proto3" json:"identifier,omitempty"`
	// resource can be an associated resource of the URN
	Resource string `protobuf:"bytes,7,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *ResourceURN) Reset() {
	*x = ResourceURN{}
	if protoimpl.UnsafeEnabled {
		mi := &file_urn_urn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceURN) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceURN) ProtoMessage() {}

func (x *ResourceURN) ProtoReflect() protoreflect.Message {
	mi := &file_urn_urn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceURN.ProtoReflect.Descriptor instead.
func (*ResourceURN) Descriptor() ([]byte, []int) {
	return file_urn_urn_proto_rawDescGZIP(), []int{0}
}

func (x *ResourceURN) GetCanonical() string {
	if x != nil {
		return x.Canonical
	}
	return ""
}

func (x *ResourceURN) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *ResourceURN) GetPartition() string {
	if x != nil {
		return x.Partition
	}
	return ""
}

func (x *ResourceURN) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *ResourceURN) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *ResourceURN) GetIdentifier() string {
	if x != nil {
		return x.Identifier
	}
	return ""
}

func (x *ResourceURN) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

var File_urn_urn_proto protoreflect.FileDescriptor

var file_urn_urn_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x75, 0x72, 0x6e, 0x2f, 0x75, 0x72, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x13, 0x6b, 0x61, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x78, 0x69, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x22, 0xd5, 0x01, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x55, 0x52, 0x4e, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63,
	0x61, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x1f, 0x5a, 0x1d,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x61, 0x74, 0x61, 0x6c,
	0x6c, 0x61, 0x78, 0x69, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x75, 0x72, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_urn_urn_proto_rawDescOnce sync.Once
	file_urn_urn_proto_rawDescData = file_urn_urn_proto_rawDesc
)

func file_urn_urn_proto_rawDescGZIP() []byte {
	file_urn_urn_proto_rawDescOnce.Do(func() {
		file_urn_urn_proto_rawDescData = protoimpl.X.CompressGZIP(file_urn_urn_proto_rawDescData)
	})
	return file_urn_urn_proto_rawDescData
}

var file_urn_urn_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_urn_urn_proto_goTypes = []interface{}{
	(*ResourceURN)(nil), // 0: katallaxie.protobuf.ResourceURN
}
var file_urn_urn_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_urn_urn_proto_init() }
func file_urn_urn_proto_init() {
	if File_urn_urn_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_urn_urn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceURN); i {
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
			RawDescriptor: file_urn_urn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_urn_urn_proto_goTypes,
		DependencyIndexes: file_urn_urn_proto_depIdxs,
		MessageInfos:      file_urn_urn_proto_msgTypes,
	}.Build()
	File_urn_urn_proto = out.File
	file_urn_urn_proto_rawDesc = nil
	file_urn_urn_proto_goTypes = nil
	file_urn_urn_proto_depIdxs = nil
}
