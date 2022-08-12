// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.1
// source: api/v1/database/errors.proto

package database

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

type DBErrorType int32

const (
	DBErrorType_SESSION      DBErrorType = 0
	DBErrorType_KEY_VALUE    DBErrorType = 1
	DBErrorType_RAFT_CONTROL DBErrorType = 2
	DBErrorType_RAFT_CLUSTER DBErrorType = 3
)

// Enum value maps for DBErrorType.
var (
	DBErrorType_name = map[int32]string{
		0: "SESSION",
		1: "KEY_VALUE",
		2: "RAFT_CONTROL",
		3: "RAFT_CLUSTER",
	}
	DBErrorType_value = map[string]int32{
		"SESSION":      0,
		"KEY_VALUE":    1,
		"RAFT_CONTROL": 2,
		"RAFT_CLUSTER": 3,
	}
)

func (x DBErrorType) Enum() *DBErrorType {
	p := new(DBErrorType)
	*p = x
	return p
}

func (x DBErrorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DBErrorType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_v1_database_errors_proto_enumTypes[0].Descriptor()
}

func (DBErrorType) Type() protoreflect.EnumType {
	return &file_api_v1_database_errors_proto_enumTypes[0]
}

func (x DBErrorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DBErrorType.Descriptor instead.
func (DBErrorType) EnumDescriptor() ([]byte, []int) {
	return file_api_v1_database_errors_proto_rawDescGZIP(), []int{0}
}

type DBError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    DBErrorType `protobuf:"varint,1,opt,name=Type,proto3,enum=database.DBErrorType" json:"Type,omitempty"`
	Message string      `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *DBError) Reset() {
	*x = DBError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_database_errors_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBError) ProtoMessage() {}

func (x *DBError) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_database_errors_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBError.ProtoReflect.Descriptor instead.
func (*DBError) Descriptor() ([]byte, []int) {
	return file_api_v1_database_errors_proto_rawDescGZIP(), []int{0}
}

func (x *DBError) GetType() DBErrorType {
	if x != nil {
		return x.Type
	}
	return DBErrorType_SESSION
}

func (x *DBError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_api_v1_database_errors_proto protoreflect.FileDescriptor

var file_api_v1_database_errors_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73,
	0x65, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x22, 0x4e, 0x0a, 0x07, 0x44, 0x42, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x12, 0x29, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x15, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x44, 0x42, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x4d, 0x0a, 0x0b, 0x44, 0x42, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x45, 0x53, 0x53, 0x49,
	0x4f, 0x4e, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x4b, 0x45, 0x59, 0x5f, 0x56, 0x41, 0x4c, 0x55,
	0x45, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x43, 0x4f, 0x4e, 0x54,
	0x52, 0x4f, 0x4c, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x43, 0x4c,
	0x55, 0x53, 0x54, 0x45, 0x52, 0x10, 0x03, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x64, 0x61, 0x74,
	0x61, 0x62, 0x61, 0x73, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_database_errors_proto_rawDescOnce sync.Once
	file_api_v1_database_errors_proto_rawDescData = file_api_v1_database_errors_proto_rawDesc
)

func file_api_v1_database_errors_proto_rawDescGZIP() []byte {
	file_api_v1_database_errors_proto_rawDescOnce.Do(func() {
		file_api_v1_database_errors_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_database_errors_proto_rawDescData)
	})
	return file_api_v1_database_errors_proto_rawDescData
}

var file_api_v1_database_errors_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_v1_database_errors_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_v1_database_errors_proto_goTypes = []interface{}{
	(DBErrorType)(0), // 0: database.DBErrorType
	(*DBError)(nil),  // 1: database.DBError
}
var file_api_v1_database_errors_proto_depIdxs = []int32{
	0, // 0: database.DBError.Type:type_name -> database.DBErrorType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_v1_database_errors_proto_init() }
func file_api_v1_database_errors_proto_init() {
	if File_api_v1_database_errors_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_database_errors_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBError); i {
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
			RawDescriptor: file_api_v1_database_errors_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_database_errors_proto_goTypes,
		DependencyIndexes: file_api_v1_database_errors_proto_depIdxs,
		EnumInfos:         file_api_v1_database_errors_proto_enumTypes,
		MessageInfos:      file_api_v1_database_errors_proto_msgTypes,
	}.Build()
	File_api_v1_database_errors_proto = out.File
	file_api_v1_database_errors_proto_rawDesc = nil
	file_api_v1_database_errors_proto_goTypes = nil
	file_api_v1_database_errors_proto_depIdxs = nil
}