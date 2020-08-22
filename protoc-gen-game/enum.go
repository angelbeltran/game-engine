package main

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
)

// NOTE: enum_key option value must be a FULLY QUALIFIED enum name.
func getMessagesWithEnumKeys(files []*desc.FileDescriptor) (map[*desc.MessageDescriptor]*desc.EnumDescriptor, error) {
	msgs, err := getMessagesAndExtensions(files, isEnumKeyFieldNumber)
	if err != nil {
		return nil, err
	}

	enumPerMsg := make(map[*desc.MessageDescriptor]*desc.EnumDescriptor, len(msgs))
	for msg, ext := range msgs {
		enumNamePtr, ok := ext.(*string)
		if !ok {
			return nil, fmt.Errorf("enum key option should have been a *string, but got a %T", ext)
		}

		enumName := *enumNamePtr

		// NOTE enumName must be fully qualified
		enums := getEnumsWithName(files, enumName)

		switch len(enums) {
		case 0:
			return nil, fmt.Errorf("invalid enum_key options: no enum named %q exists", enumName)
		case 1:
			enumPerMsg[msg] = enums[0]
		default:
			return nil, fmt.Errorf("invalid enum_key options: found multiple enums with the name %q", enumName)
		}
	}

	return enumPerMsg, nil
}

func getEnumsWithName(files []*desc.FileDescriptor, fullyQualifiedName string) []*desc.EnumDescriptor {
	var matches []*desc.EnumDescriptor

	traversed := make(map[string]bool)

	for _, file := range files {
		for _, file := range append([]*desc.FileDescriptor{file}, file.GetDependencies()...) {
			if traversed[file.GetFullyQualifiedName()] {
				continue
			}

			if enumDesc := file.FindEnum(fullyQualifiedName); enumDesc != nil {
				matches = append(matches, enumDesc)
			}

			traversed[file.GetFullyQualifiedName()] = true
		}
	}

	return matches
}

func validateMessagesWithEnums(m map[*desc.MessageDescriptor]*desc.EnumDescriptor) error {
	for msg, enum := range m {
		vals := make(map[int32]bool)

		for _, v := range enum.GetValues() {
			if i := v.GetNumber(); i != 0 {
				vals[i] = true
			}
		}

		fields := make(map[int32]*desc.FieldDescriptor)
		for _, f := range msg.GetFields() {
			if i := f.GetNumber(); vals[i] {
				fields[i] = f
			}
		}

		if len(fields) != len(vals) {
			return fmt.Errorf("not all values for enum %s are present on %s message type", enum.GetFullyQualifiedName(), msg.GetFullyQualifiedName())
		}

		var first *desc.FieldDescriptor
		for _, field := range fields {
			if first == nil {
				first = field
			} else {
				if !hasSameType(first, field) {
					return fmt.Errorf("fields corresponding to enum values must all have the same type: %q and %q of %q do not have the same type", first, field, msg)
				}
			}
		}
	}

	return nil
}

/*
func (fd *FieldDescriptor) GetEnumType() *EnumDescriptor
func (fd *FieldDescriptor) GetMapKeyType() *FieldDescriptor
func (fd *FieldDescriptor) GetMapValueType() *FieldDescriptor
func (fd *FieldDescriptor) GetMessageType() *MessageDescriptor
func (fd *FieldDescriptor) GetType() dpb.FieldDescriptorProto_Type
func (fd *FieldDescriptor) IsExtension() bool
func (fd *FieldDescriptor) IsMap() bool
func (fd *FieldDescriptor) IsRepeated() bool
*/

func hasSameType(a, b *desc.FieldDescriptor) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if a.GetType() != b.GetType() {
		return false
	}

	if a.IsExtension() != b.IsExtension() {
		return false
	}

	if a.IsMap() != b.IsMap() {
		return false
	}

	if a.IsRepeated() != b.IsRepeated() {
		return false
	}

	if a.IsMap() {
		if !hasSameType(a.GetMapKeyType(), b.GetMapKeyType()) {
			return false
		}

		if !hasSameType(a.GetMapValueType(), b.GetMapValueType()) {
			return false
		}
	}

	aet := a.GetEnumType()
	bet := b.GetEnumType()

	if (aet == nil) != (bet == nil) {
		return false
	}
	if aet != nil {
		return aet.GetFullyQualifiedName() == bet.GetFullyQualifiedName()
	}

	amt := a.GetMessageType()
	bmt := b.GetMessageType()

	if (amt == nil) != (bmt == nil) {
		return false
	}

	return amt == nil || amt.GetFullyQualifiedName() == bmt.GetFullyQualifiedName()
}
