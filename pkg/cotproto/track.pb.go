// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: track.proto

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

// All items are required unless otherwise noted!
// "required" means if they are missing on send, the conversion
// to the message format will be rejected and fall back to opaque
// XML representation
type Track struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Speed         float64                `protobuf:"fixed64,1,opt,name=speed,proto3" json:"speed,omitempty"`   // speed=
	Course        float64                `protobuf:"fixed64,2,opt,name=course,proto3" json:"course,omitempty"` // course=
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Track) Reset() {
	*x = Track{}
	mi := &file_track_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Track) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Track) ProtoMessage() {}

func (x *Track) ProtoReflect() protoreflect.Message {
	mi := &file_track_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Track.ProtoReflect.Descriptor instead.
func (*Track) Descriptor() ([]byte, []int) {
	return file_track_proto_rawDescGZIP(), []int{0}
}

func (x *Track) GetSpeed() float64 {
	if x != nil {
		return x.Speed
	}
	return 0
}

func (x *Track) GetCourse() float64 {
	if x != nil {
		return x.Course
	}
	return 0
}

var File_track_proto protoreflect.FileDescriptor

const file_track_proto_rawDesc = "" +
	"\n" +
	"\vtrack.proto\"5\n" +
	"\x05Track\x12\x14\n" +
	"\x05speed\x18\x01 \x01(\x01R\x05speed\x12\x16\n" +
	"\x06course\x18\x02 \x01(\x01R\x06courseB,H\x03Z(github.com/coreywagehoft/go-tak/cotprotob\x06proto3"

var (
	file_track_proto_rawDescOnce sync.Once
	file_track_proto_rawDescData []byte
)

func file_track_proto_rawDescGZIP() []byte {
	file_track_proto_rawDescOnce.Do(func() {
		file_track_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_track_proto_rawDesc), len(file_track_proto_rawDesc)))
	})
	return file_track_proto_rawDescData
}

var file_track_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_track_proto_goTypes = []any{
	(*Track)(nil), // 0: Track
}
var file_track_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_track_proto_init() }
func file_track_proto_init() {
	if File_track_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_track_proto_rawDesc), len(file_track_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_track_proto_goTypes,
		DependencyIndexes: file_track_proto_depIdxs,
		MessageInfos:      file_track_proto_msgTypes,
	}.Build()
	File_track_proto = out.File
	file_track_proto_goTypes = nil
	file_track_proto_depIdxs = nil
}
