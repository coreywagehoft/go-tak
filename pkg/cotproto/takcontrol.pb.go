// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: takcontrol.proto

package cotproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// TAK Protocol control message
// This specifies to a recipient what versions
// of protocol elements this sender supports during
// decoding.
type TakControl struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Lowest TAK protocol version supported
	// If not filled in (reads as 0), version 1 is assumed
	MinProtoVersion uint32 `protobuf:"varint,1,opt,name=minProtoVersion,proto3" json:"minProtoVersion,omitempty"`
	// Highest TAK protocol version supported
	// If not filled in (reads as 0), version 1 is assumed
	MaxProtoVersion uint32 `protobuf:"varint,2,opt,name=maxProtoVersion,proto3" json:"maxProtoVersion,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *TakControl) Reset() {
	*x = TakControl{}
	mi := &file_takcontrol_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TakControl) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TakControl) ProtoMessage() {}

func (x *TakControl) ProtoReflect() protoreflect.Message {
	mi := &file_takcontrol_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TakControl.ProtoReflect.Descriptor instead.
func (*TakControl) Descriptor() ([]byte, []int) {
	return file_takcontrol_proto_rawDescGZIP(), []int{0}
}

func (x *TakControl) GetMinProtoVersion() uint32 {
	if x != nil {
		return x.MinProtoVersion
	}
	return 0
}

func (x *TakControl) GetMaxProtoVersion() uint32 {
	if x != nil {
		return x.MaxProtoVersion
	}
	return 0
}

var File_takcontrol_proto protoreflect.FileDescriptor

const file_takcontrol_proto_rawDesc = "" +
	"\n" +
	"\x10takcontrol.proto\"`\n" +
	"\n" +
	"TakControl\x12(\n" +
	"\x0fminProtoVersion\x18\x01 \x01(\rR\x0fminProtoVersion\x12(\n" +
	"\x0fmaxProtoVersion\x18\x02 \x01(\rR\x0fmaxProtoVersionB,H\x03Z(github.com/coreywagehoft/go-tak/cotprotob\x06proto3"

var (
	file_takcontrol_proto_rawDescOnce sync.Once
	file_takcontrol_proto_rawDescData []byte
)

func file_takcontrol_proto_rawDescGZIP() []byte {
	file_takcontrol_proto_rawDescOnce.Do(func() {
		file_takcontrol_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_takcontrol_proto_rawDesc), len(file_takcontrol_proto_rawDesc)))
	})
	return file_takcontrol_proto_rawDescData
}

var file_takcontrol_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_takcontrol_proto_goTypes = []any{
	(*TakControl)(nil), // 0: TakControl
}
var file_takcontrol_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_takcontrol_proto_init() }
func file_takcontrol_proto_init() {
	if File_takcontrol_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_takcontrol_proto_rawDesc), len(file_takcontrol_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_takcontrol_proto_goTypes,
		DependencyIndexes: file_takcontrol_proto_depIdxs,
		MessageInfos:      file_takcontrol_proto_msgTypes,
	}.Build()
	File_takcontrol_proto = out.File
	file_takcontrol_proto_goTypes = nil
	file_takcontrol_proto_depIdxs = nil
}
