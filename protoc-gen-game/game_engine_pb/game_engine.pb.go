// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: game_engine.proto

package game_engine_pb

import (
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

type Rule_Single_Operator int32

const (
	Rule_Single_NO_OP Rule_Single_Operator = 0
	Rule_Single_EQ    Rule_Single_Operator = 1
	Rule_Single_NEQ   Rule_Single_Operator = 2
	Rule_Single_LT    Rule_Single_Operator = 3
	Rule_Single_LTE   Rule_Single_Operator = 4
	Rule_Single_GT    Rule_Single_Operator = 5
	Rule_Single_GTE   Rule_Single_Operator = 6
)

// Enum value maps for Rule_Single_Operator.
var (
	Rule_Single_Operator_name = map[int32]string{
		0: "NO_OP",
		1: "EQ",
		2: "NEQ",
		3: "LT",
		4: "LTE",
		5: "GT",
		6: "GTE",
	}
	Rule_Single_Operator_value = map[string]int32{
		"NO_OP": 0,
		"EQ":    1,
		"NEQ":   2,
		"LT":    3,
		"LTE":   4,
		"GT":    5,
		"GTE":   6,
	}
)

func (x Rule_Single_Operator) Enum() *Rule_Single_Operator {
	p := new(Rule_Single_Operator)
	*p = x
	return p
}

func (x Rule_Single_Operator) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Rule_Single_Operator) Descriptor() protoreflect.EnumDescriptor {
	return file_game_engine_proto_enumTypes[0].Descriptor()
}

func (Rule_Single_Operator) Type() protoreflect.EnumType {
	return &file_game_engine_proto_enumTypes[0]
}

func (x Rule_Single_Operator) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Rule_Single_Operator.Descriptor instead.
func (Rule_Single_Operator) EnumDescriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{2, 0, 0}
}

type Action struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Effect []*Effect `protobuf:"bytes,1,rep,name=effect,proto3" json:"effect,omitempty"`
	Rule   *Rule     `protobuf:"bytes,2,opt,name=rule,proto3" json:"rule,omitempty"`
	Error  *Error    `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Action) Reset() {
	*x = Action{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action) ProtoMessage() {}

func (x *Action) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action.ProtoReflect.Descriptor instead.
func (*Action) Descriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{0}
}

func (x *Action) GetEffect() []*Effect {
	if x != nil {
		return x.Effect
	}
	return nil
}

func (x *Action) GetRule() *Rule {
	if x != nil {
		return x.Rule
	}
	return nil
}

func (x *Action) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

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
		mi := &file_game_engine_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Effect) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Effect) ProtoMessage() {}

func (x *Effect) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[1]
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
	return file_game_engine_proto_rawDescGZIP(), []int{1}
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

type Rule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to RuleField:
	//	*Rule_Single_
	//	*Rule_And_
	//	*Rule_Or_
	RuleField isRule_RuleField `protobuf_oneof:"rule_field"`
}

func (x *Rule) Reset() {
	*x = Rule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule) ProtoMessage() {}

func (x *Rule) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule.ProtoReflect.Descriptor instead.
func (*Rule) Descriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{2}
}

func (m *Rule) GetRuleField() isRule_RuleField {
	if m != nil {
		return m.RuleField
	}
	return nil
}

func (x *Rule) GetSingle() *Rule_Single {
	if x, ok := x.GetRuleField().(*Rule_Single_); ok {
		return x.Single
	}
	return nil
}

func (x *Rule) GetAnd() *Rule_And {
	if x, ok := x.GetRuleField().(*Rule_And_); ok {
		return x.And
	}
	return nil
}

func (x *Rule) GetOr() *Rule_Or {
	if x, ok := x.GetRuleField().(*Rule_Or_); ok {
		return x.Or
	}
	return nil
}

type isRule_RuleField interface {
	isRule_RuleField()
}

type Rule_Single_ struct {
	Single *Rule_Single `protobuf:"bytes,1,opt,name=single,proto3,oneof"`
}

type Rule_And_ struct {
	And *Rule_And `protobuf:"bytes,2,opt,name=and,proto3,oneof"`
}

type Rule_Or_ struct {
	Or *Rule_Or `protobuf:"bytes,3,opt,name=or,proto3,oneof"`
}

func (*Rule_Single_) isRule_RuleField() {}

func (*Rule_And_) isRule_RuleField() {}

func (*Rule_Or_) isRule_RuleField() {}

type Operand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Operand:
	//	*Operand_Value
	//	*Operand_Prop
	//	*Operand_Input
	Operand isOperand_Operand `protobuf_oneof:"operand"`
}

func (x *Operand) Reset() {
	*x = Operand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Operand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Operand) ProtoMessage() {}

func (x *Operand) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Operand.ProtoReflect.Descriptor instead.
func (*Operand) Descriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{3}
}

func (m *Operand) GetOperand() isOperand_Operand {
	if m != nil {
		return m.Operand
	}
	return nil
}

func (x *Operand) GetValue() *Value {
	if x, ok := x.GetOperand().(*Operand_Value); ok {
		return x.Value
	}
	return nil
}

func (x *Operand) GetProp() *Path {
	if x, ok := x.GetOperand().(*Operand_Prop); ok {
		return x.Prop
	}
	return nil
}

func (x *Operand) GetInput() *Path {
	if x, ok := x.GetOperand().(*Operand_Input); ok {
		return x.Input
	}
	return nil
}

type isOperand_Operand interface {
	isOperand_Operand()
}

type Operand_Value struct {
	Value *Value `protobuf:"bytes,1,opt,name=value,proto3,oneof"`
}

type Operand_Prop struct {
	Prop *Path `protobuf:"bytes,2,opt,name=prop,proto3,oneof"`
}

type Operand_Input struct {
	Input *Path `protobuf:"bytes,3,opt,name=input,proto3,oneof"`
}

func (*Operand_Value) isOperand_Operand() {}

func (*Operand_Prop) isOperand_Operand() {}

func (*Operand_Input) isOperand_Operand() {}

type Path struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path []string `protobuf:"bytes,1,rep,name=path,proto3" json:"path,omitempty"` // TODO: rename to 'segments'
}

func (x *Path) Reset() {
	*x = Path{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Path) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Path) ProtoMessage() {}

func (x *Path) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Path.ProtoReflect.Descriptor instead.
func (*Path) Descriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{4}
}

func (x *Path) GetPath() []string {
	if x != nil {
		return x.Path
	}
	return nil
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//	*Value_Bool
	//	*Value_Integer
	//	*Value_Float
	//	*Value_String_
	Value isValue_Value `protobuf_oneof:"value"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{5}
}

func (m *Value) GetValue() isValue_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *Value) GetBool() bool {
	if x, ok := x.GetValue().(*Value_Bool); ok {
		return x.Bool
	}
	return false
}

func (x *Value) GetInteger() int64 {
	if x, ok := x.GetValue().(*Value_Integer); ok {
		return x.Integer
	}
	return 0
}

func (x *Value) GetFloat() float32 {
	if x, ok := x.GetValue().(*Value_Float); ok {
		return x.Float
	}
	return 0
}

func (x *Value) GetString_() string {
	if x, ok := x.GetValue().(*Value_String_); ok {
		return x.String_
	}
	return ""
}

type isValue_Value interface {
	isValue_Value()
}

type Value_Bool struct {
	Bool bool `protobuf:"varint,1,opt,name=bool,proto3,oneof"`
}

type Value_Integer struct {
	Integer int64 `protobuf:"varint,2,opt,name=integer,proto3,oneof"`
}

type Value_Float struct {
	Float float32 `protobuf:"fixed32,3,opt,name=float,proto3,oneof"`
}

type Value_String_ struct {
	String_ string `protobuf:"bytes,4,opt,name=string,proto3,oneof"`
}

func (*Value_Bool) isValue_Value() {}

func (*Value_Integer) isValue_Value() {}

func (*Value_Float) isValue_Value() {}

func (*Value_String_) isValue_Value() {}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Code    string `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	Msg     string `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{6}
}

func (x *Error) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *Error) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Error) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type Effect_Update struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dest *Path    `protobuf:"bytes,1,opt,name=dest,proto3" json:"dest,omitempty"`
	Src  *Operand `protobuf:"bytes,2,opt,name=src,proto3" json:"src,omitempty"`
}

func (x *Effect_Update) Reset() {
	*x = Effect_Update{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Effect_Update) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Effect_Update) ProtoMessage() {}

func (x *Effect_Update) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[7]
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
	return file_game_engine_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Effect_Update) GetDest() *Path {
	if x != nil {
		return x.Dest
	}
	return nil
}

func (x *Effect_Update) GetSrc() *Operand {
	if x != nil {
		return x.Src
	}
	return nil
}

type Rule_Single struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Operator Rule_Single_Operator `protobuf:"varint,1,opt,name=operator,proto3,enum=game_engine.Rule_Single_Operator" json:"operator,omitempty"`
	Left     *Operand             `protobuf:"bytes,2,opt,name=left,proto3" json:"left,omitempty"`
	Right    *Operand             `protobuf:"bytes,3,opt,name=right,proto3" json:"right,omitempty"`
}

func (x *Rule_Single) Reset() {
	*x = Rule_Single{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule_Single) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule_Single) ProtoMessage() {}

func (x *Rule_Single) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule_Single.ProtoReflect.Descriptor instead.
func (*Rule_Single) Descriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Rule_Single) GetOperator() Rule_Single_Operator {
	if x != nil {
		return x.Operator
	}
	return Rule_Single_NO_OP
}

func (x *Rule_Single) GetLeft() *Operand {
	if x != nil {
		return x.Left
	}
	return nil
}

func (x *Rule_Single) GetRight() *Operand {
	if x != nil {
		return x.Right
	}
	return nil
}

type Rule_And struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rules []*Rule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (x *Rule_And) Reset() {
	*x = Rule_And{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule_And) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule_And) ProtoMessage() {}

func (x *Rule_And) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule_And.ProtoReflect.Descriptor instead.
func (*Rule_And) Descriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{2, 1}
}

func (x *Rule_And) GetRules() []*Rule {
	if x != nil {
		return x.Rules
	}
	return nil
}

type Rule_Or struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rules []*Rule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (x *Rule_Or) Reset() {
	*x = Rule_Or{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_engine_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule_Or) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule_Or) ProtoMessage() {}

func (x *Rule_Or) ProtoReflect() protoreflect.Message {
	mi := &file_game_engine_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule_Or.ProtoReflect.Descriptor instead.
func (*Rule_Or) Descriptor() ([]byte, []int) {
	return file_game_engine_proto_rawDescGZIP(), []int{2, 2}
}

func (x *Rule_Or) GetRules() []*Rule {
	if x != nil {
		return x.Rules
	}
	return nil
}

var file_game_engine_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.ServiceOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50000,
		Name:          "game_engine.is_action_service",
		Tag:           "varint,50000,opt,name=is_action_service",
		Filename:      "game_engine.proto",
	},
	{
		ExtendedType:  (*descriptor.MethodOptions)(nil),
		ExtensionType: (*Action)(nil),
		Field:         50001,
		Name:          "game_engine.action",
		Tag:           "bytes,50001,opt,name=action",
		Filename:      "game_engine.proto",
	},
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50002,
		Name:          "game_engine.is_game_state",
		Tag:           "varint,50002,opt,name=is_game_state",
		Filename:      "game_engine.proto",
	},
}

// Extension fields to descriptor.ServiceOptions.
var (
	// optional bool is_action_service = 50000;
	E_IsActionService = &file_game_engine_proto_extTypes[0] // TODO: deliberate these numbers
)

// Extension fields to descriptor.MethodOptions.
var (
	// optional game_engine.Action action = 50001;
	E_Action = &file_game_engine_proto_extTypes[1] // TODO: deliberate these numbers
)

// Extension fields to descriptor.MessageOptions.
var (
	// optional bool is_game_state = 50002;
	E_IsGameState = &file_game_engine_proto_extTypes[2] // TODO: deliberate these numbers
)

var File_game_engine_proto protoreflect.FileDescriptor

var file_game_engine_proto_rawDesc = []byte{
	0x0a, 0x11, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x86, 0x01, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2b, 0x0a,
	0x06, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x45, 0x66, 0x66, 0x65,
	0x63, 0x74, 0x52, 0x06, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x12, 0x25, 0x0a, 0x04, 0x72, 0x75,
	0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x75, 0x6c,
	0x65, 0x12, 0x28, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0xa4, 0x01, 0x0a, 0x06,
	0x45, 0x66, 0x66, 0x65, 0x63, 0x74, 0x12, 0x34, 0x0a, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x2e, 0x45, 0x66, 0x66, 0x65, 0x63, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x48, 0x00, 0x52, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x1a, 0x57, 0x0a, 0x06,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2e, 0x50, 0x61, 0x74, 0x68, 0x52, 0x04, 0x64, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a,
	0x03, 0x73, 0x72, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x61, 0x6d,
	0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x6e, 0x64,
	0x52, 0x03, 0x73, 0x72, 0x63, 0x42, 0x0b, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0xe4, 0x03, 0x0a, 0x04, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x32, 0x0a, 0x06, 0x73,
	0x69, 0x6e, 0x67, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x61,
	0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x2e, 0x53,
	0x69, 0x6e, 0x67, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x06, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x12,
	0x29, 0x0a, 0x03, 0x61, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x67,
	0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x2e,
	0x41, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x03, 0x61, 0x6e, 0x64, 0x12, 0x26, 0x0a, 0x02, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x2e, 0x4f, 0x72, 0x48, 0x00, 0x52, 0x02,
	0x6f, 0x72, 0x1a, 0xe7, 0x01, 0x0a, 0x06, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x12, 0x3d, 0x0a,
	0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x21, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x75,
	0x6c, 0x65, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x28, 0x0a, 0x04,
	0x6c, 0x65, 0x66, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x61, 0x6d,
	0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x6e, 0x64,
	0x52, 0x04, 0x6c, 0x65, 0x66, 0x74, 0x12, 0x2a, 0x0a, 0x05, 0x72, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x6e, 0x64, 0x52, 0x05, 0x72, 0x69, 0x67,
	0x68, 0x74, 0x22, 0x48, 0x0a, 0x08, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x09,
	0x0a, 0x05, 0x4e, 0x4f, 0x5f, 0x4f, 0x50, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x45, 0x51, 0x10,
	0x01, 0x12, 0x07, 0x0a, 0x03, 0x4e, 0x45, 0x51, 0x10, 0x02, 0x12, 0x06, 0x0a, 0x02, 0x4c, 0x54,
	0x10, 0x03, 0x12, 0x07, 0x0a, 0x03, 0x4c, 0x54, 0x45, 0x10, 0x04, 0x12, 0x06, 0x0a, 0x02, 0x47,
	0x54, 0x10, 0x05, 0x12, 0x07, 0x0a, 0x03, 0x47, 0x54, 0x45, 0x10, 0x06, 0x1a, 0x2e, 0x0a, 0x03,
	0x41, 0x6e, 0x64, 0x12, 0x27, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x1a, 0x2d, 0x0a, 0x02,
	0x4f, 0x72, 0x12, 0x27, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e,
	0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x42, 0x0c, 0x0a, 0x0a, 0x72,
	0x75, 0x6c, 0x65, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x22, 0x94, 0x01, 0x0a, 0x07, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x6e, 0x64, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x27, 0x0a, 0x04, 0x70, 0x72, 0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x50, 0x61,
	0x74, 0x68, 0x48, 0x00, 0x52, 0x04, 0x70, 0x72, 0x6f, 0x70, 0x12, 0x29, 0x0a, 0x05, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x61, 0x6d, 0x65,
	0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x50, 0x61, 0x74, 0x68, 0x48, 0x00, 0x52, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x6e, 0x64,
	0x22, 0x1a, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x74, 0x0a, 0x05,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x04, 0x62, 0x6f, 0x6f, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x04, 0x62, 0x6f, 0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x07, 0x69,
	0x6e, 0x74, 0x65, 0x67, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x07,
	0x69, 0x6e, 0x74, 0x65, 0x67, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x05, 0x66, 0x6c, 0x6f, 0x61, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x48, 0x00, 0x52, 0x05, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x12,
	0x18, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x47, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x3a, 0x50, 0x0a, 0x11, 0x69,
	0x73, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x69, 0x73, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x50, 0x0a,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x3a,
	0x48, 0x0a, 0x0d, 0x69, 0x73, 0x5f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xd2, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x47, 0x61, 0x6d,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x42, 0x39, 0x5a, 0x37, 0x61, 0x6e, 0x67,
	0x65, 0x6c, 0x62, 0x65, 0x6c, 0x74, 0x72, 0x61, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e,
	0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_game_engine_proto_rawDescOnce sync.Once
	file_game_engine_proto_rawDescData = file_game_engine_proto_rawDesc
)

func file_game_engine_proto_rawDescGZIP() []byte {
	file_game_engine_proto_rawDescOnce.Do(func() {
		file_game_engine_proto_rawDescData = protoimpl.X.CompressGZIP(file_game_engine_proto_rawDescData)
	})
	return file_game_engine_proto_rawDescData
}

var file_game_engine_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_game_engine_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_game_engine_proto_goTypes = []interface{}{
	(Rule_Single_Operator)(0),         // 0: game_engine.Rule.Single.Operator
	(*Action)(nil),                    // 1: game_engine.Action
	(*Effect)(nil),                    // 2: game_engine.Effect
	(*Rule)(nil),                      // 3: game_engine.Rule
	(*Operand)(nil),                   // 4: game_engine.Operand
	(*Path)(nil),                      // 5: game_engine.Path
	(*Value)(nil),                     // 6: game_engine.Value
	(*Error)(nil),                     // 7: game_engine.Error
	(*Effect_Update)(nil),             // 8: game_engine.Effect.Update
	(*Rule_Single)(nil),               // 9: game_engine.Rule.Single
	(*Rule_And)(nil),                  // 10: game_engine.Rule.And
	(*Rule_Or)(nil),                   // 11: game_engine.Rule.Or
	(*descriptor.ServiceOptions)(nil), // 12: google.protobuf.ServiceOptions
	(*descriptor.MethodOptions)(nil),  // 13: google.protobuf.MethodOptions
	(*descriptor.MessageOptions)(nil), // 14: google.protobuf.MessageOptions
}
var file_game_engine_proto_depIdxs = []int32{
	2,  // 0: game_engine.Action.effect:type_name -> game_engine.Effect
	3,  // 1: game_engine.Action.rule:type_name -> game_engine.Rule
	7,  // 2: game_engine.Action.error:type_name -> game_engine.Error
	8,  // 3: game_engine.Effect.update:type_name -> game_engine.Effect.Update
	9,  // 4: game_engine.Rule.single:type_name -> game_engine.Rule.Single
	10, // 5: game_engine.Rule.and:type_name -> game_engine.Rule.And
	11, // 6: game_engine.Rule.or:type_name -> game_engine.Rule.Or
	6,  // 7: game_engine.Operand.value:type_name -> game_engine.Value
	5,  // 8: game_engine.Operand.prop:type_name -> game_engine.Path
	5,  // 9: game_engine.Operand.input:type_name -> game_engine.Path
	5,  // 10: game_engine.Effect.Update.dest:type_name -> game_engine.Path
	4,  // 11: game_engine.Effect.Update.src:type_name -> game_engine.Operand
	0,  // 12: game_engine.Rule.Single.operator:type_name -> game_engine.Rule.Single.Operator
	4,  // 13: game_engine.Rule.Single.left:type_name -> game_engine.Operand
	4,  // 14: game_engine.Rule.Single.right:type_name -> game_engine.Operand
	3,  // 15: game_engine.Rule.And.rules:type_name -> game_engine.Rule
	3,  // 16: game_engine.Rule.Or.rules:type_name -> game_engine.Rule
	12, // 17: game_engine.is_action_service:extendee -> google.protobuf.ServiceOptions
	13, // 18: game_engine.action:extendee -> google.protobuf.MethodOptions
	14, // 19: game_engine.is_game_state:extendee -> google.protobuf.MessageOptions
	1,  // 20: game_engine.action:type_name -> game_engine.Action
	21, // [21:21] is the sub-list for method output_type
	21, // [21:21] is the sub-list for method input_type
	20, // [20:21] is the sub-list for extension type_name
	17, // [17:20] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_game_engine_proto_init() }
func file_game_engine_proto_init() {
	if File_game_engine_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_game_engine_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action); i {
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
		file_game_engine_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_game_engine_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rule); i {
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
		file_game_engine_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Operand); i {
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
		file_game_engine_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Path); i {
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
		file_game_engine_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
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
		file_game_engine_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_game_engine_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
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
		file_game_engine_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rule_Single); i {
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
		file_game_engine_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rule_And); i {
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
		file_game_engine_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rule_Or); i {
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
	file_game_engine_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Effect_Update_)(nil),
	}
	file_game_engine_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Rule_Single_)(nil),
		(*Rule_And_)(nil),
		(*Rule_Or_)(nil),
	}
	file_game_engine_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Operand_Value)(nil),
		(*Operand_Prop)(nil),
		(*Operand_Input)(nil),
	}
	file_game_engine_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*Value_Bool)(nil),
		(*Value_Integer)(nil),
		(*Value_Float)(nil),
		(*Value_String_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_game_engine_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 3,
			NumServices:   0,
		},
		GoTypes:           file_game_engine_proto_goTypes,
		DependencyIndexes: file_game_engine_proto_depIdxs,
		EnumInfos:         file_game_engine_proto_enumTypes,
		MessageInfos:      file_game_engine_proto_msgTypes,
		ExtensionInfos:    file_game_engine_proto_extTypes,
	}.Build()
	File_game_engine_proto = out.File
	file_game_engine_proto_rawDesc = nil
	file_game_engine_proto_goTypes = nil
	file_game_engine_proto_depIdxs = nil
}
