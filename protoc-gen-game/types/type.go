package types

import (
	"sort"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
)

// TODO: handle infinite structures.

type (
	// The core of the package.
	Type interface {
		String() string
		IsSameType(Type) bool
		IsType()
	}

	Bool    struct{}
	Integer struct{}
	Float   struct{}
	String  struct{}
	Bytes   struct{}
	Enum    struct {
		Name   string
		Values map[int64]string
	}
	OneOf map[string]Type
	List  struct {
		Value Type
	}
	Structured map[string]Type
	Map        struct {
		Key   Type
		Value Type
	}
)

// FromMessage will create or lookup a cached Type describing the message
// descriptor. Results are cached for future calls.
func FromMessage(d *desc.MessageDescriptor) Type {
	if t, found := loadMessageType(d); found {
		return t
	}

	m := make(Structured)

	for _, field := range d.GetFields() {
		m[field.GetName()] = FromField(field)
	}

	storeMessageType(d, m)

	return m
}

// FromField will create or lookup a cached Type describing the field
// descriptor. Results are cached for future calls.
func FromField(d *desc.FieldDescriptor) Type {
	if t, found := loadFieldType(d); found {
		return t
	}

	var t Type

	if ed := d.GetEnumType(); ed != nil {
		t = FromEnum(ed)
	} else if od := d.GetOneOf(); od != nil {
		t = FromOneOf(od)
	} else if d.IsMap() {
		t = Map{
			Key:   FromField(d.GetMapKeyType()),
			Value: FromField(d.GetMapValueType()),
		}
	} else if md := d.GetMessageType(); md != nil {
		t = FromMessage(md)
	} else {
		switch d.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_BOOL:
			t = Bool{}

		case descriptor.FieldDescriptorProto_TYPE_STRING:
			t = String{}

		case descriptor.FieldDescriptorProto_TYPE_BYTES:
			t = Bytes{}

		case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
			descriptor.FieldDescriptorProto_TYPE_FLOAT:
			t = Float{}

		case descriptor.FieldDescriptorProto_TYPE_INT64,
			descriptor.FieldDescriptorProto_TYPE_UINT64,
			descriptor.FieldDescriptorProto_TYPE_INT32,
			descriptor.FieldDescriptorProto_TYPE_FIXED64,
			descriptor.FieldDescriptorProto_TYPE_FIXED32,
			descriptor.FieldDescriptorProto_TYPE_UINT32,
			descriptor.FieldDescriptorProto_TYPE_SFIXED32,
			descriptor.FieldDescriptorProto_TYPE_SFIXED64,
			descriptor.FieldDescriptorProto_TYPE_SINT32,
			descriptor.FieldDescriptorProto_TYPE_SINT64:
			t = Integer{}
		}
	}

	if d.IsRepeated() {
		t = List{t}
	}

	storeFieldType(d, t)

	return t
}

// FromEnum will create or lookup a cached Type describing the enum descriptor.
// Results are cached for future calls.
func FromEnum(d *desc.EnumDescriptor) Type {
	if t, found := loadEnumType(d); found {
		return t
	}

	t := Enum{
		Name:   d.GetFullyQualifiedName(),
		Values: make(map[int64]string),
	}

	for _, val := range d.GetValues() {
		t.Values[int64(val.GetNumber())] = val.GetName()
	}

	storeEnumType(d, t)

	return t
}

// FromOneOf will create or lookup a cached Type describing the oneof
// descriptor. Results are cached for future calls.
func FromOneOf(d *desc.OneOfDescriptor) Type {
	if t, found := loadOneOfType(d); found {
		return t
	}

	oo := make(OneOf)

	for _, field := range d.GetChoices() {
		oo[field.GetName()] = FromField(field)
	}

	storeOneOfType(d, oo)

	return oo
}

func (Bool) String() string {
	return "bool"
}

func (Bool) IsSameType(t Type) bool {
	_, ok := t.(Bool)
	return ok
}

func (Bool) IsType() {
}

func (Integer) String() string {
	return "integer"
}

func (Integer) IsSameType(t Type) bool {
	_, ok := t.(Integer)
	return ok
}

func (Integer) IsType() {
}

func (Float) String() string {
	return "float"
}

func (Float) IsSameType(t Type) bool {
	_, ok := t.(Float)
	return ok
}

func (Float) IsType() {
}

func (String) String() string {
	return "string"
}

func (String) IsSameType(t Type) bool {
	_, ok := t.(String)
	return ok
}

func (String) IsType() {
}

func (Bytes) String() string {
	return "bytes"
}

func (Bytes) IsSameType(t Type) bool {
	_, ok := t.(Bytes)
	return ok
}

func (Bytes) IsType() {
}

func (t Enum) String() string {
	return "enum[" + t.Name + "]"
}

func (t Enum) IsSameType(u Type) bool {
	v, ok := u.(Enum)
	if !ok || t.Name != v.Name || len(t.Values) != len(v.Values) {
		return false
	}

	for i, name := range t.Values {
		if other, ok := v.Values[i]; !ok || name != other {
			return false
		}
	}

	return true
}

func (Enum) IsType() {
}

func (t OneOf) String() string {
	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	pairs := make([]string, len(keys))
	for i, k := range keys {
		pairs[i] = k + ": " + t[k].String()
	}

	return "oneof[" + strings.Join(pairs, ", ") + "]"
}

func (t OneOf) IsSameType(u Type) bool {
	v, ok := u.(OneOf)
	if !ok || len(t) != len(v) {
		return false
	}

	for name, field := range t {
		if other, ok := v[name]; !ok || !field.IsSameType(other) {
			return false
		}
	}

	return true
}

func (OneOf) IsType() {
}

func (t List) String() string {
	return "list[" + t.Value.String() + "]"
}

func (t List) IsSameType(u Type) bool {
	v, ok := u.(List)
	return ok && t.Value.IsSameType(v.Value)
}

func (List) IsType() {
}

func (t Structured) String() string {
	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	pairs := make([]string, len(keys))
	for i, k := range keys {
		pairs[i] = k + ": " + t[k].String()
	}

	return "structured{" + strings.Join(pairs, ", ") + "}"
}

func (t Structured) IsSameType(u Type) bool {
	v, ok := u.(Structured)
	if !ok || len(t) != len(v) {
		return false
	}

	for name, field := range t {
		if other, ok := v[name]; !ok || !field.IsSameType(other) {
			return false
		}
	}

	return true
}

func (Structured) IsType() {
}

func (t Map) String() string {
	return "map<" + t.Key.String() + ", " + t.Value.String() + ">"
}

func (t Map) IsSameType(u Type) bool {
	v, ok := u.(Map)

	return ok && t.Key.IsSameType(v.Key) && t.Value.IsSameType(v.Value)
}

func (Map) IsType() {
}
