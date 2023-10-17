// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.15.8
// source: mp3_SDFS_message.proto

package protobuf

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

type SDFSMessageType int32

const (
	SDFSMessageType_DELETE SDFSMessageType = 0
	SDFSMessageType_UPDATE SDFSMessageType = 1
)

// Enum value maps for SDFSMessageType.
var (
	SDFSMessageType_name = map[int32]string{
		0: "DELETE",
		1: "UPDATE",
	}
	SDFSMessageType_value = map[string]int32{
		"DELETE": 0,
		"UPDATE": 1,
	}
)

func (x SDFSMessageType) Enum() *SDFSMessageType {
	p := new(SDFSMessageType)
	*p = x
	return p
}

func (x SDFSMessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SDFSMessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_mp3_SDFS_message_proto_enumTypes[0].Descriptor()
}

func (SDFSMessageType) Type() protoreflect.EnumType {
	return &file_mp3_SDFS_message_proto_enumTypes[0]
}

func (x SDFSMessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SDFSMessageType.Descriptor instead.
func (SDFSMessageType) EnumDescriptor() ([]byte, []int) {
	return file_mp3_SDFS_message_proto_rawDescGZIP(), []int{0}
}

type SDFSMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageType  SDFSMessageType `protobuf:"varint,1,opt,name=messageType,proto3,enum=cs425_mp3.SDFSMessageType" json:"messageType,omitempty"`
	Replicas     []string        `protobuf:"bytes,2,rep,name=replicas,proto3" json:"replicas,omitempty"`
	SdfsFileName string          `protobuf:"bytes,3,opt,name=sdfsFileName,proto3" json:"sdfsFileName,omitempty"`
}

func (x *SDFSMessage) Reset() {
	*x = SDFSMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp3_SDFS_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SDFSMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SDFSMessage) ProtoMessage() {}

func (x *SDFSMessage) ProtoReflect() protoreflect.Message {
	mi := &file_mp3_SDFS_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SDFSMessage.ProtoReflect.Descriptor instead.
func (*SDFSMessage) Descriptor() ([]byte, []int) {
	return file_mp3_SDFS_message_proto_rawDescGZIP(), []int{0}
}

func (x *SDFSMessage) GetMessageType() SDFSMessageType {
	if x != nil {
		return x.MessageType
	}
	return SDFSMessageType_DELETE
}

func (x *SDFSMessage) GetReplicas() []string {
	if x != nil {
		return x.Replicas
	}
	return nil
}

func (x *SDFSMessage) GetSdfsFileName() string {
	if x != nil {
		return x.SdfsFileName
	}
	return ""
}

var File_mp3_SDFS_message_proto protoreflect.FileDescriptor

var file_mp3_SDFS_message_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6d, 0x70, 0x33, 0x5f, 0x53, 0x44, 0x46, 0x53, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x73, 0x34, 0x32, 0x35, 0x5f,
	0x6d, 0x70, 0x33, 0x22, 0x8b, 0x01, 0x0a, 0x0b, 0x53, 0x44, 0x46, 0x53, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x3c, 0x0a, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x63, 0x73, 0x34, 0x32, 0x35,
	0x5f, 0x6d, 0x70, 0x33, 0x2e, 0x53, 0x44, 0x46, 0x53, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x12, 0x22, 0x0a,
	0x0c, 0x73, 0x64, 0x66, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x64, 0x66, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x2a, 0x29, 0x0a, 0x0f, 0x53, 0x44, 0x46, 0x53, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x01, 0x42, 0x0c, 0x5a, 0x0a,
	0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_mp3_SDFS_message_proto_rawDescOnce sync.Once
	file_mp3_SDFS_message_proto_rawDescData = file_mp3_SDFS_message_proto_rawDesc
)

func file_mp3_SDFS_message_proto_rawDescGZIP() []byte {
	file_mp3_SDFS_message_proto_rawDescOnce.Do(func() {
		file_mp3_SDFS_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_mp3_SDFS_message_proto_rawDescData)
	})
	return file_mp3_SDFS_message_proto_rawDescData
}

var file_mp3_SDFS_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_mp3_SDFS_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_mp3_SDFS_message_proto_goTypes = []interface{}{
	(SDFSMessageType)(0), // 0: cs425_mp3.SDFSMessageType
	(*SDFSMessage)(nil),  // 1: cs425_mp3.SDFSMessage
}
var file_mp3_SDFS_message_proto_depIdxs = []int32{
	0, // 0: cs425_mp3.SDFSMessage.messageType:type_name -> cs425_mp3.SDFSMessageType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_mp3_SDFS_message_proto_init() }
func file_mp3_SDFS_message_proto_init() {
	if File_mp3_SDFS_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mp3_SDFS_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SDFSMessage); i {
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
			RawDescriptor: file_mp3_SDFS_message_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mp3_SDFS_message_proto_goTypes,
		DependencyIndexes: file_mp3_SDFS_message_proto_depIdxs,
		EnumInfos:         file_mp3_SDFS_message_proto_enumTypes,
		MessageInfos:      file_mp3_SDFS_message_proto_msgTypes,
	}.Build()
	File_mp3_SDFS_message_proto = out.File
	file_mp3_SDFS_message_proto_rawDesc = nil
	file_mp3_SDFS_message_proto_goTypes = nil
	file_mp3_SDFS_message_proto_depIdxs = nil
}
