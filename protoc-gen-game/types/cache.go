package types

import "github.com/jhump/protoreflect/desc"

var (
	messageCache = make(map[string]Type)
	fieldCache   = make(map[string]Type)
	enumCache    = make(map[string]Type)
	oneOfCache   = make(map[string]Type)
)

func loadMessageType(d *desc.MessageDescriptor) (Type, bool) {
	t, ok := messageCache[d.GetFullyQualifiedName()]
	return t, ok
}

func storeMessageType(d *desc.MessageDescriptor, t Type) {
	messageCache[d.GetFullyQualifiedName()] = t
}

func loadFieldType(d *desc.FieldDescriptor) (Type, bool) {
	t, ok := fieldCache[d.GetFullyQualifiedName()]
	return t, ok
}

func storeFieldType(d *desc.FieldDescriptor, t Type) {
	fieldCache[d.GetFullyQualifiedName()] = t
}

func loadEnumType(d *desc.EnumDescriptor) (Type, bool) {
	t, ok := enumCache[d.GetFullyQualifiedName()]
	return t, ok
}

func storeEnumType(d *desc.EnumDescriptor, t Type) {
	enumCache[d.GetFullyQualifiedName()] = t
}

func loadOneOfType(d *desc.OneOfDescriptor) (Type, bool) {
	t, ok := oneOfCache[d.GetFullyQualifiedName()]
	return t, ok
}

func storeOneOfType(d *desc.OneOfDescriptor, t Type) {
	oneOfCache[d.GetFullyQualifiedName()] = t
}
