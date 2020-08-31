// Generated by protoc-gen-game/generation. DO NOT EDIT.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: effect.proto

package game_engine_pb

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

type Effect struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Operation:
	//	*Effect_Update_
	Operation isEffect_Operation `protobuf_oneof:"operation"`
}

func (x *Effect) Reset() {
	*x = Effect{}
	if protoimpl.UnsafeEnabled {
		mi := &file_effect_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Effect) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Effect) ProtoMessage() {}

func (x *Effect) ProtoReflect() protoreflect.Message {
	mi := &file_effect_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Effect.ProtoReflect.Descriptor instead.
func (*Effect) Descriptor() ([]byte, []int) {
	return file_effect_proto_rawDescGZIP(), []int{0}
}

func (m *Effect) GetOperation() isEffect_Operation {
	if m != nil {
		return m.Operation
	}
	return nil
}

func (x *Effect) GetUpdate() *Effect_Update {
	if x, ok := x.GetOperation().(*Effect_Update_); ok {
		return x.Update
	}
	return nil
}

type isEffect_Operation interface {
	isEffect_Operation()
}

type Effect_Update_ struct {
	Update *Effect_Update `protobuf:"bytes,1,opt,name=update,proto3,oneof"`
}

func (*Effect_Update_) isEffect_Operation() {}

type Effect_Update struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State *Reference `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	Value *Value     `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Effect_Update) Reset() {
	*x = Effect_Update{}
	if protoimpl.UnsafeEnabled {
		mi := &file_effect_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Effect_Update) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Effect_Update) ProtoMessage() {}

func (x *Effect_Update) ProtoReflect() protoreflect.Message {
	mi := &file_effect_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Effect_Update.ProtoReflect.Descriptor instead.
func (*Effect_Update) Descriptor() ([]byte, []int) {
	return file_effect_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Effect_Update) GetState() *Reference {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *Effect_Update) GetValue() *Value {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_effect_proto protoreflect.FileDescriptor

var file_effect_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b,
	0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x1a, 0x0f, 0x72, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x5f, 0x61, 0x6e, 0x64, 0x5f, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xad, 0x01, 0x0a, 0x06, 0x45, 0x66, 0x66,
	0x65, 0x63, 0x74, 0x12, 0x34, 0x0a, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2e, 0x45, 0x66, 0x66, 0x65, 0x63, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x48,
	0x00, 0x52, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x1a, 0x60, 0x0a, 0x06, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x67, 0x65, 0x6c, 0x62, 0x65, 0x6c, 0x74,
	0x72, 0x61, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2f,
	0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5f, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_effect_proto_rawDescOnce sync.Once
	file_effect_proto_rawDescData = file_effect_proto_rawDesc
)

func file_effect_proto_rawDescGZIP() []byte {
	file_effect_proto_rawDescOnce.Do(func() {
		file_effect_proto_rawDescData = protoimpl.X.CompressGZIP(file_effect_proto_rawDescData)
	})
	return file_effect_proto_rawDescData
}

var file_effect_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_effect_proto_goTypes = []interface{}{
	(*Effect)(nil),        // 0: game_engine.Effect
	(*Effect_Update)(nil), // 1: game_engine.Effect.Update
	(*Reference)(nil),     // 2: game_engine.Reference
	(*Value)(nil),         // 3: game_engine.Value
}
var file_effect_proto_depIdxs = []int32{
	1, // 0: game_engine.Effect.update:type_name -> game_engine.Effect.Update
	2, // 1: game_engine.Effect.Update.state:type_name -> game_engine.Reference
	3, // 2: game_engine.Effect.Update.value:type_name -> game_engine.Value
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_effect_proto_init() }
func file_effect_proto_init() {
	if File_effect_proto != nil {
		return
	}
	file_reference_proto_init()
	file_values_and_functions_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_effect_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Effect); i {
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
		file_effect_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Effect_Update); i {
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
	file_effect_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Effect_Update_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_effect_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_effect_proto_goTypes,
		DependencyIndexes: file_effect_proto_depIdxs,
		MessageInfos:      file_effect_proto_msgTypes,
	}.Build()
	File_effect_proto = out.File
	file_effect_proto_rawDesc = nil
	file_effect_proto_goTypes = nil
	file_effect_proto_depIdxs = nil
}