// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: missionannouncement.proto

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

// Container for MissionAnnouncement and metadata
type MissionAnnouncement struct {
	state                   protoimpl.MessageState `protogen:"open.v1"`
	Payload                 *TakMessage            `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	MissionName             string                 `protobuf:"bytes,2,opt,name=missionName,proto3" json:"missionName,omitempty"`
	MissionAnnouncementType string                 `protobuf:"bytes,3,opt,name=missionAnnouncementType,proto3" json:"missionAnnouncementType,omitempty"`
	CreatorUid              string                 `protobuf:"bytes,4,opt,name=creatorUid,proto3" json:"creatorUid,omitempty"`
	GroupVector             string                 `protobuf:"bytes,5,opt,name=groupVector,proto3" json:"groupVector,omitempty"`
	ClientUid               string                 `protobuf:"bytes,6,opt,name=clientUid,proto3" json:"clientUid,omitempty"`
	Uids                    []string               `protobuf:"bytes,7,rep,name=uids,proto3" json:"uids,omitempty"`
	unknownFields           protoimpl.UnknownFields
	sizeCache               protoimpl.SizeCache
}

func (x *MissionAnnouncement) Reset() {
	*x = MissionAnnouncement{}
	mi := &file_missionannouncement_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MissionAnnouncement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MissionAnnouncement) ProtoMessage() {}

func (x *MissionAnnouncement) ProtoReflect() protoreflect.Message {
	mi := &file_missionannouncement_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MissionAnnouncement.ProtoReflect.Descriptor instead.
func (*MissionAnnouncement) Descriptor() ([]byte, []int) {
	return file_missionannouncement_proto_rawDescGZIP(), []int{0}
}

func (x *MissionAnnouncement) GetPayload() *TakMessage {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *MissionAnnouncement) GetMissionName() string {
	if x != nil {
		return x.MissionName
	}
	return ""
}

func (x *MissionAnnouncement) GetMissionAnnouncementType() string {
	if x != nil {
		return x.MissionAnnouncementType
	}
	return ""
}

func (x *MissionAnnouncement) GetCreatorUid() string {
	if x != nil {
		return x.CreatorUid
	}
	return ""
}

func (x *MissionAnnouncement) GetGroupVector() string {
	if x != nil {
		return x.GroupVector
	}
	return ""
}

func (x *MissionAnnouncement) GetClientUid() string {
	if x != nil {
		return x.ClientUid
	}
	return ""
}

func (x *MissionAnnouncement) GetUids() []string {
	if x != nil {
		return x.Uids
	}
	return nil
}

var File_missionannouncement_proto protoreflect.FileDescriptor

const file_missionannouncement_proto_rawDesc = "" +
	"\n" +
	"\x19missionannouncement.proto\x1a\x10takmessage.proto\"\x8c\x02\n" +
	"\x13MissionAnnouncement\x12%\n" +
	"\apayload\x18\x01 \x01(\v2\v.TakMessageR\apayload\x12 \n" +
	"\vmissionName\x18\x02 \x01(\tR\vmissionName\x128\n" +
	"\x17missionAnnouncementType\x18\x03 \x01(\tR\x17missionAnnouncementType\x12\x1e\n" +
	"\n" +
	"creatorUid\x18\x04 \x01(\tR\n" +
	"creatorUid\x12 \n" +
	"\vgroupVector\x18\x05 \x01(\tR\vgroupVector\x12\x1c\n" +
	"\tclientUid\x18\x06 \x01(\tR\tclientUid\x12\x12\n" +
	"\x04uids\x18\a \x03(\tR\x04uidsB,H\x03Z(github.com/coreywagehoft/go-tak/cotprotob\x06proto3"

var (
	file_missionannouncement_proto_rawDescOnce sync.Once
	file_missionannouncement_proto_rawDescData []byte
)

func file_missionannouncement_proto_rawDescGZIP() []byte {
	file_missionannouncement_proto_rawDescOnce.Do(func() {
		file_missionannouncement_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_missionannouncement_proto_rawDesc), len(file_missionannouncement_proto_rawDesc)))
	})
	return file_missionannouncement_proto_rawDescData
}

var file_missionannouncement_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_missionannouncement_proto_goTypes = []any{
	(*MissionAnnouncement)(nil), // 0: MissionAnnouncement
	(*TakMessage)(nil),          // 1: TakMessage
}
var file_missionannouncement_proto_depIdxs = []int32{
	1, // 0: MissionAnnouncement.payload:type_name -> TakMessage
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_missionannouncement_proto_init() }
func file_missionannouncement_proto_init() {
	if File_missionannouncement_proto != nil {
		return
	}
	file_takmessage_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_missionannouncement_proto_rawDesc), len(file_missionannouncement_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_missionannouncement_proto_goTypes,
		DependencyIndexes: file_missionannouncement_proto_depIdxs,
		MessageInfos:      file_missionannouncement_proto_msgTypes,
	}.Build()
	File_missionannouncement_proto = out.File
	file_missionannouncement_proto_goTypes = nil
	file_missionannouncement_proto_depIdxs = nil
}
