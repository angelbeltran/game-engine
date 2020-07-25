// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: github.com/angelbeltran/game-engine/examples/puerto-rico/good.proto

package puerto_rico_pb

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type GoodID int32

const (
	GoodID_NO_GOOD GoodID = 0
	GoodID_CORN    GoodID = 1
	GoodID_INDIGO  GoodID = 2
	GoodID_SUGAR   GoodID = 3
	GoodID_TOBACCO GoodID = 4
	GoodID_COFFEE  GoodID = 5
)

// Enum value maps for GoodID.
var (
	GoodID_name = map[int32]string{
		0: "NO_GOOD",
		1: "CORN",
		2: "INDIGO",
		3: "SUGAR",
		4: "TOBACCO",
		5: "COFFEE",
	}
	GoodID_value = map[string]int32{
		"NO_GOOD": 0,
		"CORN":    1,
		"INDIGO":  2,
		"SUGAR":   3,
		"TOBACCO": 4,
		"COFFEE":  5,
	}
)

func (x GoodID) Enum() *GoodID {
	p := new(GoodID)
	*p = x
	return p
}

func (x GoodID) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GoodID) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_enumTypes[0].Descriptor()
}

func (GoodID) Type() protoreflect.EnumType {
	return &file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_enumTypes[0]
}

func (x GoodID) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GoodID.Descriptor instead.
func (GoodID) EnumDescriptor() ([]byte, []int) {
	return file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDescGZIP(), []int{0}
}

var File_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto protoreflect.FileDescriptor

var file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDesc = []byte{
	0x0a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x67,
	0x65, 0x6c, 0x62, 0x65, 0x6c, 0x74, 0x72, 0x61, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x70,
	0x75, 0x65, 0x72, 0x74, 0x6f, 0x2d, 0x72, 0x69, 0x63, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x70, 0x75, 0x65, 0x72, 0x74, 0x6f, 0x5f, 0x72, 0x69,
	0x63, 0x6f, 0x5f, 0x67, 0x61, 0x6d, 0x65, 0x2a, 0x4f, 0x0a, 0x06, 0x47, 0x6f, 0x6f, 0x64, 0x49,
	0x44, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x4f, 0x5f, 0x47, 0x4f, 0x4f, 0x44, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x43, 0x4f, 0x52, 0x4e, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x49, 0x4e, 0x44, 0x49,
	0x47, 0x4f, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x55, 0x47, 0x41, 0x52, 0x10, 0x03, 0x12,
	0x0b, 0x0a, 0x07, 0x54, 0x4f, 0x42, 0x41, 0x43, 0x43, 0x4f, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06,
	0x43, 0x4f, 0x46, 0x46, 0x45, 0x45, 0x10, 0x05, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x67, 0x65, 0x6c, 0x62, 0x65, 0x6c, 0x74,
	0x72, 0x61, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x70, 0x75, 0x65, 0x72, 0x74, 0x6f, 0x2d,
	0x72, 0x69, 0x63, 0x6f, 0x2f, 0x70, 0x75, 0x65, 0x72, 0x74, 0x6f, 0x5f, 0x72, 0x69, 0x63, 0x6f,
	0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDescOnce sync.Once
	file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDescData = file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDesc
)

func file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDescGZIP() []byte {
	file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDescOnce.Do(func() {
		file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDescData)
	})
	return file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDescData
}

var file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_goTypes = []interface{}{
	(GoodID)(0), // 0: puerto_rico_game.GoodID
}
var file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_init() }
func file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_init() {
	if File_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_goTypes,
		DependencyIndexes: file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_depIdxs,
		EnumInfos:         file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_enumTypes,
	}.Build()
	File_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto = out.File
	file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_rawDesc = nil
	file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_goTypes = nil
	file_github_com_angelbeltran_game_engine_examples_puerto_rico_good_proto_depIdxs = nil
}
