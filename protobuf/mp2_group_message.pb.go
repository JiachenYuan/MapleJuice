// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.15.8
// source: mp2_group_message.proto

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

type NodeInfoRow_NodeStatus int32

const (
	NodeInfoRow_Alive     NodeInfoRow_NodeStatus = 0
	NodeInfoRow_Suspected NodeInfoRow_NodeStatus = 1
	NodeInfoRow_Failed    NodeInfoRow_NodeStatus = 2 // TODO: might delete it since it should never be sent
)

// Enum value maps for NodeInfoRow_NodeStatus.
var (
	NodeInfoRow_NodeStatus_name = map[int32]string{
		0: "Alive",
		1: "Suspected",
		2: "Failed",
	}
	NodeInfoRow_NodeStatus_value = map[string]int32{
		"Alive":     0,
		"Suspected": 1,
		"Failed":    2,
	}
)

func (x NodeInfoRow_NodeStatus) Enum() *NodeInfoRow_NodeStatus {
	p := new(NodeInfoRow_NodeStatus)
	*p = x
	return p
}

func (x NodeInfoRow_NodeStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NodeInfoRow_NodeStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_mp2_group_message_proto_enumTypes[0].Descriptor()
}

func (NodeInfoRow_NodeStatus) Type() protoreflect.EnumType {
	return &file_mp2_group_message_proto_enumTypes[0]
}

func (x NodeInfoRow_NodeStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NodeInfoRow_NodeStatus.Descriptor instead.
func (NodeInfoRow_NodeStatus) EnumDescriptor() ([]byte, []int) {
	return file_mp2_group_message_proto_rawDescGZIP(), []int{0, 0}
}

type GroupMessage_MessageType int32

const (
	GroupMessage_JOIN   GroupMessage_MessageType = 0
	GroupMessage_GOSSIP GroupMessage_MessageType = 2
)

// Enum value maps for GroupMessage_MessageType.
var (
	GroupMessage_MessageType_name = map[int32]string{
		0: "JOIN",
		2: "GOSSIP",
	}
	GroupMessage_MessageType_value = map[string]int32{
		"JOIN":   0,
		"GOSSIP": 2,
	}
)

func (x GroupMessage_MessageType) Enum() *GroupMessage_MessageType {
	p := new(GroupMessage_MessageType)
	*p = x
	return p
}

func (x GroupMessage_MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GroupMessage_MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_mp2_group_message_proto_enumTypes[1].Descriptor()
}

func (GroupMessage_MessageType) Type() protoreflect.EnumType {
	return &file_mp2_group_message_proto_enumTypes[1]
}

func (x GroupMessage_MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GroupMessage_MessageType.Descriptor instead.
func (GroupMessage_MessageType) EnumDescriptor() ([]byte, []int) {
	return file_mp2_group_message_proto_rawDescGZIP(), []int{2, 0}
}

type NodeInfoRow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeID string                 `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	SeqNum int32                  `protobuf:"varint,2,opt,name=seqNum,proto3" json:"seqNum,omitempty"`
	Status NodeInfoRow_NodeStatus `protobuf:"varint,3,opt,name=status,proto3,enum=cs425_mp2.NodeInfoRow_NodeStatus" json:"status,omitempty"`
}

func (x *NodeInfoRow) Reset() {
	*x = NodeInfoRow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp2_group_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeInfoRow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeInfoRow) ProtoMessage() {}

func (x *NodeInfoRow) ProtoReflect() protoreflect.Message {
	mi := &file_mp2_group_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeInfoRow.ProtoReflect.Descriptor instead.
func (*NodeInfoRow) Descriptor() ([]byte, []int) {
	return file_mp2_group_message_proto_rawDescGZIP(), []int{0}
}

func (x *NodeInfoRow) GetNodeID() string {
	if x != nil {
		return x.NodeID
	}
	return ""
}

func (x *NodeInfoRow) GetSeqNum() int32 {
	if x != nil {
		return x.SeqNum
	}
	return 0
}

func (x *NodeInfoRow) GetStatus() NodeInfoRow_NodeStatus {
	if x != nil {
		return x.Status
	}
	return NodeInfoRow_Alive
}

type NodeInfoList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rows []*NodeInfoRow `protobuf:"bytes,1,rep,name=rows,proto3" json:"rows,omitempty"`
}

func (x *NodeInfoList) Reset() {
	*x = NodeInfoList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp2_group_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeInfoList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeInfoList) ProtoMessage() {}

func (x *NodeInfoList) ProtoReflect() protoreflect.Message {
	mi := &file_mp2_group_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeInfoList.ProtoReflect.Descriptor instead.
func (*NodeInfoList) Descriptor() ([]byte, []int) {
	return file_mp2_group_message_proto_rawDescGZIP(), []int{1}
}

func (x *NodeInfoList) GetRows() []*NodeInfoRow {
	if x != nil {
		return x.Rows
	}
	return nil
}

type GroupMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type         GroupMessage_MessageType `protobuf:"varint,1,opt,name=type,proto3,enum=cs425_mp2.GroupMessage_MessageType" json:"type,omitempty"`
	NodeInfoList *NodeInfoList            `protobuf:"bytes,2,opt,name=node_info_list,json=nodeInfoList,proto3" json:"node_info_list,omitempty"`
}

func (x *GroupMessage) Reset() {
	*x = GroupMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp2_group_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupMessage) ProtoMessage() {}

func (x *GroupMessage) ProtoReflect() protoreflect.Message {
	mi := &file_mp2_group_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupMessage.ProtoReflect.Descriptor instead.
func (*GroupMessage) Descriptor() ([]byte, []int) {
	return file_mp2_group_message_proto_rawDescGZIP(), []int{2}
}

func (x *GroupMessage) GetType() GroupMessage_MessageType {
	if x != nil {
		return x.Type
	}
	return GroupMessage_JOIN
}

func (x *GroupMessage) GetNodeInfoList() *NodeInfoList {
	if x != nil {
		return x.NodeInfoList
	}
	return nil
}

var File_mp2_group_message_proto protoreflect.FileDescriptor

var file_mp2_group_message_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6d, 0x70, 0x32, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x73, 0x34, 0x32, 0x35,
	0x5f, 0x6d, 0x70, 0x32, 0x22, 0xac, 0x01, 0x0a, 0x0b, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x6f, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x65, 0x71, 0x4e, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x65,
	0x71, 0x4e, 0x75, 0x6d, 0x12, 0x39, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e, 0x63, 0x73, 0x34, 0x32, 0x35, 0x5f, 0x6d, 0x70, 0x32,
	0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x6f, 0x77, 0x2e, 0x4e, 0x6f, 0x64,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x32, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x09, 0x0a,
	0x05, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x75, 0x73, 0x70,
	0x65, 0x63, 0x74, 0x65, 0x64, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x61, 0x69, 0x6c, 0x65,
	0x64, 0x10, 0x02, 0x22, 0x3a, 0x0a, 0x0c, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x04, 0x72, 0x6f, 0x77, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x63, 0x73, 0x34, 0x32, 0x35, 0x5f, 0x6d, 0x70, 0x32, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x6f, 0x77, 0x52, 0x04, 0x72, 0x6f, 0x77, 0x73, 0x22,
	0xab, 0x01, 0x0a, 0x0c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x37, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23,
	0x2e, 0x63, 0x73, 0x34, 0x32, 0x35, 0x5f, 0x6d, 0x70, 0x32, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x3d, 0x0a, 0x0e, 0x6e, 0x6f, 0x64,
	0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x63, 0x73, 0x34, 0x32, 0x35, 0x5f, 0x6d, 0x70, 0x32, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x0c, 0x6e, 0x6f, 0x64, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x23, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4a, 0x4f, 0x49, 0x4e, 0x10,
	0x00, 0x12, 0x0a, 0x0a, 0x06, 0x47, 0x4f, 0x53, 0x53, 0x49, 0x50, 0x10, 0x02, 0x42, 0x0c, 0x5a,
	0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_mp2_group_message_proto_rawDescOnce sync.Once
	file_mp2_group_message_proto_rawDescData = file_mp2_group_message_proto_rawDesc
)

func file_mp2_group_message_proto_rawDescGZIP() []byte {
	file_mp2_group_message_proto_rawDescOnce.Do(func() {
		file_mp2_group_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_mp2_group_message_proto_rawDescData)
	})
	return file_mp2_group_message_proto_rawDescData
}

var file_mp2_group_message_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_mp2_group_message_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_mp2_group_message_proto_goTypes = []interface{}{
	(NodeInfoRow_NodeStatus)(0),   // 0: cs425_mp2.NodeInfoRow.NodeStatus
	(GroupMessage_MessageType)(0), // 1: cs425_mp2.GroupMessage.MessageType
	(*NodeInfoRow)(nil),           // 2: cs425_mp2.NodeInfoRow
	(*NodeInfoList)(nil),          // 3: cs425_mp2.NodeInfoList
	(*GroupMessage)(nil),          // 4: cs425_mp2.GroupMessage
}
var file_mp2_group_message_proto_depIdxs = []int32{
	0, // 0: cs425_mp2.NodeInfoRow.status:type_name -> cs425_mp2.NodeInfoRow.NodeStatus
	2, // 1: cs425_mp2.NodeInfoList.rows:type_name -> cs425_mp2.NodeInfoRow
	1, // 2: cs425_mp2.GroupMessage.type:type_name -> cs425_mp2.GroupMessage.MessageType
	3, // 3: cs425_mp2.GroupMessage.node_info_list:type_name -> cs425_mp2.NodeInfoList
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_mp2_group_message_proto_init() }
func file_mp2_group_message_proto_init() {
	if File_mp2_group_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mp2_group_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeInfoRow); i {
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
		file_mp2_group_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeInfoList); i {
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
		file_mp2_group_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GroupMessage); i {
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
			RawDescriptor: file_mp2_group_message_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mp2_group_message_proto_goTypes,
		DependencyIndexes: file_mp2_group_message_proto_depIdxs,
		EnumInfos:         file_mp2_group_message_proto_enumTypes,
		MessageInfos:      file_mp2_group_message_proto_msgTypes,
	}.Build()
	File_mp2_group_message_proto = out.File
	file_mp2_group_message_proto_rawDesc = nil
	file_mp2_group_message_proto_goTypes = nil
	file_mp2_group_message_proto_depIdxs = nil
}