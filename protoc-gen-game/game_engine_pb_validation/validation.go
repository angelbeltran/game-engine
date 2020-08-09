package game_engine_pb_validation

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	pb "github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
)

func ValidateValue(val *pb.Value, state, input *desc.MessageDescriptor) (Type, error) {
	if val.Value == nil {
		return "", fmt.Errorf("empty value")
	}

	switch v := val.Value.(type) {
	case *pb.Value_Bool:
		return TypeBool, ValidateBoolValueReferences(v.Bool, state, input)
	case *pb.Value_Int:
		return TypeInt, ValidateIntValueReferences(v.Int, state, input)
	case *pb.Value_Float:
		return TypeFloat, ValidateFloatValueReferences(v.Float, state, input)
	case *pb.Value_String_:
		return TypeString, ValidateStringValueReferences(v.String_, state, input)
	default:
		return "", fmt.Errorf("unrecognized value type: %T", v)
	}
}

type Type string

const (
	TypeNone Type = ""
	TypeBool Type = "bool"
	TypeInt Type = "int"
	TypeFloat Type = "float"
	TypeString Type = "string"
)

// Generic tools for verifying Reference values.

func ValidateReference(ref *pb.Reference, md *desc.MessageDescriptor, t Type) error {
	return VerifyPathType(ref.Path, md, t)
}

func VerifyPathType(path []string, md *desc.MessageDescriptor, t Type) error {
	return VerifyPath(path, md, func(fd *desc.FieldDescriptor) error {
		u, err := FieldDescriptorTypeToSupportedType(fd.GetType())
		if err != nil {
			return fmt.Errorf("invalid field %s: %w", fd.GetFullyQualifiedName(), err)
		}

		if t != u {
			return fmt.Errorf("unexpected type at field %s: expected %s but found %s", fd.GetFullyQualifiedName(), t, u)
		}

		return nil
	})
}

func VerifyPath(path []string, md *desc.MessageDescriptor, validators ...func (*desc.FieldDescriptor) error) error {
	if md == nil {
		return fmt.Errorf("nil message descriptor")
	}
	if len(path) == 0 {
		return fmt.Errorf("no path specified")
	}

	fd := md.FindFieldByName(path[0])
	if fd == nil {
		return fmt.Errorf("no field named %s", path[0])
	}
	if len(path) == 1 {
		if md = fd.GetMessageType(); md != nil {
			return fmt.Errorf("field named %s is a message type", path[0])
		}

		for _, f := range validators {
			if err := f(fd); err != nil {
				return err
			}
		}

		return nil
	}

	if md = fd.GetMessageType(); md == nil {
		return fmt.Errorf("field %s is not a message type", path[0])
	}
	if err := VerifyPath(path[1:], md, validators...); err != nil {
		return fmt.Errorf("error verifying field beyond field %s: %w", path[0], err)
	}

	return nil
}

func FieldDescriptorTypeToSupportedType(t dpb.FieldDescriptorProto_Type) (Type, error) {
	switch t {
	case 1:
		return TypeFloat, nil
	case 13:
		return TypeInt, nil
	case 14:
		return TypeString, nil
	case 15:
		return TypeInt, nil
	case 16:
		return TypeInt, nil
	case 17:
		return TypeInt, nil
	case 18:
		return TypeInt, nil
	case 2:
		return TypeFloat, nil
	case 3:
		return TypeInt, nil
	case 4:
		return TypeInt, nil
	case 5:
		return TypeInt, nil
	case 6:
		return TypeInt, nil
	case 7:
		return TypeInt, nil
	case 8:
		return TypeBool, nil
	case 9:
		return TypeString, nil
	}
	return TypeNone, fmt.Errorf("unsupported field type: %s", t)
}

// Protobuf Type-Specific Validation Tools.

// Value Validation

func ValidateBoolValueReferences(val *pb.BoolValue, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolValue value missing")
	}

	switch v := val.Value.(type) {
	case *pb.BoolValue_Constant:
		return nil
	case *pb.BoolValue_Input:
		return ValidateReference(v.Input, input, TypeBool)
	case *pb.BoolValue_State:
		return ValidateReference(v.State, state, TypeBool)
	case *pb.BoolValue_BoolFunc:
		return ValidateBoolToBoolFunctionReferences(v.BoolFunc, state, input)
	case *pb.BoolValue_IntFunc:
		return ValidateIntToBoolFunctionReferences(v.IntFunc, state, input)
	case *pb.BoolValue_FloatFunc:
		return ValidateFloatToBoolFunctionReferences(v.FloatFunc, state, input)
	case *pb.BoolValue_StringFunc:
		return ValidateStringToBoolFunctionReferences(v.StringFunc, state, input)
	case *pb.BoolValue_BoolBoolFunc:
		return ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.BoolValue_BoolIntFunc:
		return ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.BoolValue_BoolFloatFunc:
		return ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.BoolValue_BoolStringFunc:
		return ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.BoolValue_IntBoolFunc:
		return ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.BoolValue_IntIntFunc:
		return ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc, state, input)
	case *pb.BoolValue_IntFloatFunc:
		return ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.BoolValue_IntStringFunc:
		return ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc, state, input)
	case *pb.BoolValue_FloatBoolFunc:
		return ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.BoolValue_FloatIntFunc:
		return ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.BoolValue_FloatFloatFunc:
		return ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.BoolValue_FloatStringFunc:
		return ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.BoolValue_StringBoolFunc:
		return ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.BoolValue_StringIntFunc:
		return ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc, state, input)
	case *pb.BoolValue_StringFloatFunc:
		return ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.BoolValue_StringStringFunc:
		return ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc, state, input)
	case *pb.BoolValue_If:
		return ValidateBoolValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized value type from BoolValue: %T", v)
	}

	return nil
}

func ValidateIntValueReferences(val *pb.IntValue, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntValue value missing")
	}

	switch v := val.Value.(type) {
	case *pb.IntValue_Constant:
		return nil
	case *pb.IntValue_Input:
		return ValidateReference(v.Input, input, TypeInt)
	case *pb.IntValue_State:
		return ValidateReference(v.State, state, TypeInt)
	case *pb.IntValue_BoolFunc:
		return ValidateBoolToIntFunctionReferences(v.BoolFunc, state, input)
	case *pb.IntValue_IntFunc:
		return ValidateIntToIntFunctionReferences(v.IntFunc, state, input)
	case *pb.IntValue_FloatFunc:
		return ValidateFloatToIntFunctionReferences(v.FloatFunc, state, input)
	case *pb.IntValue_StringFunc:
		return ValidateStringToIntFunctionReferences(v.StringFunc, state, input)
	case *pb.IntValue_BoolBoolFunc:
		return ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.IntValue_BoolIntFunc:
		return ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.IntValue_BoolFloatFunc:
		return ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.IntValue_BoolStringFunc:
		return ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.IntValue_IntBoolFunc:
		return ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.IntValue_IntIntFunc:
		return ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc, state, input)
	case *pb.IntValue_IntFloatFunc:
		return ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.IntValue_IntStringFunc:
		return ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc, state, input)
	case *pb.IntValue_FloatBoolFunc:
		return ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.IntValue_FloatIntFunc:
		return ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.IntValue_FloatFloatFunc:
		return ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.IntValue_FloatStringFunc:
		return ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.IntValue_StringBoolFunc:
		return ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.IntValue_StringIntFunc:
		return ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc, state, input)
	case *pb.IntValue_StringFloatFunc:
		return ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.IntValue_StringStringFunc:
		return ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc, state, input)
	case *pb.IntValue_If:
		return ValidateIntValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized value type from IntValue: %T", v)
	}

	return nil
}

func ValidateFloatValueReferences(val *pb.FloatValue, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatValue value missing")
	}

	switch v := val.Value.(type) {
	case *pb.FloatValue_Constant:
		return nil
	case *pb.FloatValue_Input:
		return ValidateReference(v.Input, input, TypeFloat)
	case *pb.FloatValue_State:
		return ValidateReference(v.State, state, TypeFloat)
	case *pb.FloatValue_BoolFunc:
		return ValidateBoolToFloatFunctionReferences(v.BoolFunc, state, input)
	case *pb.FloatValue_IntFunc:
		return ValidateIntToFloatFunctionReferences(v.IntFunc, state, input)
	case *pb.FloatValue_FloatFunc:
		return ValidateFloatToFloatFunctionReferences(v.FloatFunc, state, input)
	case *pb.FloatValue_StringFunc:
		return ValidateStringToFloatFunctionReferences(v.StringFunc, state, input)
	case *pb.FloatValue_BoolBoolFunc:
		return ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.FloatValue_BoolIntFunc:
		return ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.FloatValue_BoolFloatFunc:
		return ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.FloatValue_BoolStringFunc:
		return ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.FloatValue_IntBoolFunc:
		return ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.FloatValue_IntIntFunc:
		return ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc, state, input)
	case *pb.FloatValue_IntFloatFunc:
		return ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.FloatValue_IntStringFunc:
		return ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc, state, input)
	case *pb.FloatValue_FloatBoolFunc:
		return ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.FloatValue_FloatIntFunc:
		return ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.FloatValue_FloatFloatFunc:
		return ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.FloatValue_FloatStringFunc:
		return ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.FloatValue_StringBoolFunc:
		return ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.FloatValue_StringIntFunc:
		return ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc, state, input)
	case *pb.FloatValue_StringFloatFunc:
		return ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.FloatValue_StringStringFunc:
		return ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc, state, input)
	case *pb.FloatValue_If:
		return ValidateFloatValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized value type from FloatValue: %T", v)
	}

	return nil
}

func ValidateStringValueReferences(val *pb.StringValue, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringValue value missing")
	}

	switch v := val.Value.(type) {
	case *pb.StringValue_Constant:
		return nil
	case *pb.StringValue_Input:
		return ValidateReference(v.Input, input, TypeString)
	case *pb.StringValue_State:
		return ValidateReference(v.State, state, TypeString)
	case *pb.StringValue_BoolFunc:
		return ValidateBoolToStringFunctionReferences(v.BoolFunc, state, input)
	case *pb.StringValue_IntFunc:
		return ValidateIntToStringFunctionReferences(v.IntFunc, state, input)
	case *pb.StringValue_FloatFunc:
		return ValidateFloatToStringFunctionReferences(v.FloatFunc, state, input)
	case *pb.StringValue_StringFunc:
		return ValidateStringToStringFunctionReferences(v.StringFunc, state, input)
	case *pb.StringValue_BoolBoolFunc:
		return ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.StringValue_BoolIntFunc:
		return ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.StringValue_BoolFloatFunc:
		return ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.StringValue_BoolStringFunc:
		return ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.StringValue_IntBoolFunc:
		return ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.StringValue_IntIntFunc:
		return ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc, state, input)
	case *pb.StringValue_IntFloatFunc:
		return ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.StringValue_IntStringFunc:
		return ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc, state, input)
	case *pb.StringValue_FloatBoolFunc:
		return ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.StringValue_FloatIntFunc:
		return ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.StringValue_FloatFloatFunc:
		return ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.StringValue_FloatStringFunc:
		return ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.StringValue_StringBoolFunc:
		return ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.StringValue_StringIntFunc:
		return ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc, state, input)
	case *pb.StringValue_StringFloatFunc:
		return ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.StringValue_StringStringFunc:
		return ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc, state, input)
	case *pb.StringValue_If:
		return ValidateStringValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized value type from StringValue: %T", v)
	}

	return nil
}

// Unary Function Validation

func ValidateBoolToBoolFunctionReferences(val *pb.BoolToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolToBoolFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.BoolToBoolFunction_Input:
		return ValidateReference(v.Input, input, TypeBool)
	case *pb.BoolToBoolFunction_State:
		return ValidateReference(v.State, state, TypeBool)
	case *pb.BoolToBoolFunction_BoolFunc:
		return ValidateBoolToBoolFunctionReferences(v.BoolFunc, state, input)
	case *pb.BoolToBoolFunction_IntFunc:
		return ValidateIntToBoolFunctionReferences(v.IntFunc, state, input)
	case *pb.BoolToBoolFunction_FloatFunc:
		return ValidateFloatToBoolFunctionReferences(v.FloatFunc, state, input)
	case *pb.BoolToBoolFunction_StringFunc:
		return ValidateStringToBoolFunctionReferences(v.StringFunc, state, input)
	case *pb.BoolToBoolFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.BoolToBoolFunction_BoolIntFunc:
		return ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.BoolToBoolFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.BoolToBoolFunction_BoolStringFunc:
		return ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.BoolToBoolFunction_IntBoolFunc:
		return ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.BoolToBoolFunction_IntIntFunc:
		return ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc, state, input)
	case *pb.BoolToBoolFunction_IntFloatFunc:
		return ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.BoolToBoolFunction_IntStringFunc:
		return ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc, state, input)
	case *pb.BoolToBoolFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.BoolToBoolFunction_FloatIntFunc:
		return ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.BoolToBoolFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.BoolToBoolFunction_FloatStringFunc:
		return ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.BoolToBoolFunction_StringBoolFunc:
		return ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.BoolToBoolFunction_StringIntFunc:
		return ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc, state, input)
	case *pb.BoolToBoolFunction_StringFloatFunc:
		return ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.BoolToBoolFunction_StringStringFunc:
		return ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc, state, input)
	case *pb.BoolToBoolFunction_If:
		return ValidateBoolValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from BoolToBoolFunction: %T", v)
	}

	return nil
}

func ValidateBoolToIntFunctionReferences(val *pb.BoolToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolToIntFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.BoolToIntFunction_Input:
		return ValidateReference(v.Input, input, TypeBool)
	case *pb.BoolToIntFunction_State:
		return ValidateReference(v.State, state, TypeBool)
	case *pb.BoolToIntFunction_BoolFunc:
		return ValidateBoolToBoolFunctionReferences(v.BoolFunc, state, input)
	case *pb.BoolToIntFunction_IntFunc:
		return ValidateIntToBoolFunctionReferences(v.IntFunc, state, input)
	case *pb.BoolToIntFunction_FloatFunc:
		return ValidateFloatToBoolFunctionReferences(v.FloatFunc, state, input)
	case *pb.BoolToIntFunction_StringFunc:
		return ValidateStringToBoolFunctionReferences(v.StringFunc, state, input)
	case *pb.BoolToIntFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.BoolToIntFunction_BoolIntFunc:
		return ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.BoolToIntFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.BoolToIntFunction_BoolStringFunc:
		return ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.BoolToIntFunction_IntBoolFunc:
		return ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.BoolToIntFunction_IntIntFunc:
		return ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc, state, input)
	case *pb.BoolToIntFunction_IntFloatFunc:
		return ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.BoolToIntFunction_IntStringFunc:
		return ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc, state, input)
	case *pb.BoolToIntFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.BoolToIntFunction_FloatIntFunc:
		return ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.BoolToIntFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.BoolToIntFunction_FloatStringFunc:
		return ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.BoolToIntFunction_StringBoolFunc:
		return ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.BoolToIntFunction_StringIntFunc:
		return ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc, state, input)
	case *pb.BoolToIntFunction_StringFloatFunc:
		return ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.BoolToIntFunction_StringStringFunc:
		return ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc, state, input)
	case *pb.BoolToIntFunction_If:
		return ValidateBoolValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from BoolToIntFunction: %T", v)
	}

	return nil
}

func ValidateBoolToFloatFunctionReferences(val *pb.BoolToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolToFloatFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.BoolToFloatFunction_Input:
		return ValidateReference(v.Input, input, TypeBool)
	case *pb.BoolToFloatFunction_State:
		return ValidateReference(v.State, state, TypeBool)
	case *pb.BoolToFloatFunction_BoolFunc:
		return ValidateBoolToBoolFunctionReferences(v.BoolFunc, state, input)
	case *pb.BoolToFloatFunction_IntFunc:
		return ValidateIntToBoolFunctionReferences(v.IntFunc, state, input)
	case *pb.BoolToFloatFunction_FloatFunc:
		return ValidateFloatToBoolFunctionReferences(v.FloatFunc, state, input)
	case *pb.BoolToFloatFunction_StringFunc:
		return ValidateStringToBoolFunctionReferences(v.StringFunc, state, input)
	case *pb.BoolToFloatFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.BoolToFloatFunction_BoolIntFunc:
		return ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.BoolToFloatFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.BoolToFloatFunction_BoolStringFunc:
		return ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.BoolToFloatFunction_IntBoolFunc:
		return ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.BoolToFloatFunction_IntIntFunc:
		return ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc, state, input)
	case *pb.BoolToFloatFunction_IntFloatFunc:
		return ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.BoolToFloatFunction_IntStringFunc:
		return ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc, state, input)
	case *pb.BoolToFloatFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.BoolToFloatFunction_FloatIntFunc:
		return ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.BoolToFloatFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.BoolToFloatFunction_FloatStringFunc:
		return ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.BoolToFloatFunction_StringBoolFunc:
		return ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.BoolToFloatFunction_StringIntFunc:
		return ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc, state, input)
	case *pb.BoolToFloatFunction_StringFloatFunc:
		return ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.BoolToFloatFunction_StringStringFunc:
		return ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc, state, input)
	case *pb.BoolToFloatFunction_If:
		return ValidateBoolValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from BoolToFloatFunction: %T", v)
	}

	return nil
}

func ValidateBoolToStringFunctionReferences(val *pb.BoolToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolToStringFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.BoolToStringFunction_Input:
		return ValidateReference(v.Input, input, TypeBool)
	case *pb.BoolToStringFunction_State:
		return ValidateReference(v.State, state, TypeBool)
	case *pb.BoolToStringFunction_BoolFunc:
		return ValidateBoolToBoolFunctionReferences(v.BoolFunc, state, input)
	case *pb.BoolToStringFunction_IntFunc:
		return ValidateIntToBoolFunctionReferences(v.IntFunc, state, input)
	case *pb.BoolToStringFunction_FloatFunc:
		return ValidateFloatToBoolFunctionReferences(v.FloatFunc, state, input)
	case *pb.BoolToStringFunction_StringFunc:
		return ValidateStringToBoolFunctionReferences(v.StringFunc, state, input)
	case *pb.BoolToStringFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.BoolToStringFunction_BoolIntFunc:
		return ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.BoolToStringFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.BoolToStringFunction_BoolStringFunc:
		return ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.BoolToStringFunction_IntBoolFunc:
		return ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.BoolToStringFunction_IntIntFunc:
		return ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc, state, input)
	case *pb.BoolToStringFunction_IntFloatFunc:
		return ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.BoolToStringFunction_IntStringFunc:
		return ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc, state, input)
	case *pb.BoolToStringFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.BoolToStringFunction_FloatIntFunc:
		return ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.BoolToStringFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.BoolToStringFunction_FloatStringFunc:
		return ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.BoolToStringFunction_StringBoolFunc:
		return ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.BoolToStringFunction_StringIntFunc:
		return ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc, state, input)
	case *pb.BoolToStringFunction_StringFloatFunc:
		return ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.BoolToStringFunction_StringStringFunc:
		return ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc, state, input)
	case *pb.BoolToStringFunction_If:
		return ValidateBoolValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from BoolToStringFunction: %T", v)
	}

	return nil
}

func ValidateIntToBoolFunctionReferences(val *pb.IntToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntToBoolFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.IntToBoolFunction_Input:
		return ValidateReference(v.Input, input, TypeInt)
	case *pb.IntToBoolFunction_State:
		return ValidateReference(v.State, state, TypeInt)
	case *pb.IntToBoolFunction_BoolFunc:
		return ValidateBoolToIntFunctionReferences(v.BoolFunc, state, input)
	case *pb.IntToBoolFunction_IntFunc:
		return ValidateIntToIntFunctionReferences(v.IntFunc, state, input)
	case *pb.IntToBoolFunction_FloatFunc:
		return ValidateFloatToIntFunctionReferences(v.FloatFunc, state, input)
	case *pb.IntToBoolFunction_StringFunc:
		return ValidateStringToIntFunctionReferences(v.StringFunc, state, input)
	case *pb.IntToBoolFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.IntToBoolFunction_BoolIntFunc:
		return ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.IntToBoolFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.IntToBoolFunction_BoolStringFunc:
		return ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.IntToBoolFunction_IntBoolFunc:
		return ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.IntToBoolFunction_IntIntFunc:
		return ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc, state, input)
	case *pb.IntToBoolFunction_IntFloatFunc:
		return ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.IntToBoolFunction_IntStringFunc:
		return ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc, state, input)
	case *pb.IntToBoolFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.IntToBoolFunction_FloatIntFunc:
		return ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.IntToBoolFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.IntToBoolFunction_FloatStringFunc:
		return ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.IntToBoolFunction_StringBoolFunc:
		return ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.IntToBoolFunction_StringIntFunc:
		return ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc, state, input)
	case *pb.IntToBoolFunction_StringFloatFunc:
		return ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.IntToBoolFunction_StringStringFunc:
		return ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc, state, input)
	case *pb.IntToBoolFunction_If:
		return ValidateIntValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from IntToBoolFunction: %T", v)
	}

	return nil
}

func ValidateIntToIntFunctionReferences(val *pb.IntToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntToIntFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.IntToIntFunction_Input:
		return ValidateReference(v.Input, input, TypeInt)
	case *pb.IntToIntFunction_State:
		return ValidateReference(v.State, state, TypeInt)
	case *pb.IntToIntFunction_BoolFunc:
		return ValidateBoolToIntFunctionReferences(v.BoolFunc, state, input)
	case *pb.IntToIntFunction_IntFunc:
		return ValidateIntToIntFunctionReferences(v.IntFunc, state, input)
	case *pb.IntToIntFunction_FloatFunc:
		return ValidateFloatToIntFunctionReferences(v.FloatFunc, state, input)
	case *pb.IntToIntFunction_StringFunc:
		return ValidateStringToIntFunctionReferences(v.StringFunc, state, input)
	case *pb.IntToIntFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.IntToIntFunction_BoolIntFunc:
		return ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.IntToIntFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.IntToIntFunction_BoolStringFunc:
		return ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.IntToIntFunction_IntBoolFunc:
		return ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.IntToIntFunction_IntIntFunc:
		return ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc, state, input)
	case *pb.IntToIntFunction_IntFloatFunc:
		return ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.IntToIntFunction_IntStringFunc:
		return ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc, state, input)
	case *pb.IntToIntFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.IntToIntFunction_FloatIntFunc:
		return ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.IntToIntFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.IntToIntFunction_FloatStringFunc:
		return ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.IntToIntFunction_StringBoolFunc:
		return ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.IntToIntFunction_StringIntFunc:
		return ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc, state, input)
	case *pb.IntToIntFunction_StringFloatFunc:
		return ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.IntToIntFunction_StringStringFunc:
		return ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc, state, input)
	case *pb.IntToIntFunction_If:
		return ValidateIntValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from IntToIntFunction: %T", v)
	}

	return nil
}

func ValidateIntToFloatFunctionReferences(val *pb.IntToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntToFloatFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.IntToFloatFunction_Input:
		return ValidateReference(v.Input, input, TypeInt)
	case *pb.IntToFloatFunction_State:
		return ValidateReference(v.State, state, TypeInt)
	case *pb.IntToFloatFunction_BoolFunc:
		return ValidateBoolToIntFunctionReferences(v.BoolFunc, state, input)
	case *pb.IntToFloatFunction_IntFunc:
		return ValidateIntToIntFunctionReferences(v.IntFunc, state, input)
	case *pb.IntToFloatFunction_FloatFunc:
		return ValidateFloatToIntFunctionReferences(v.FloatFunc, state, input)
	case *pb.IntToFloatFunction_StringFunc:
		return ValidateStringToIntFunctionReferences(v.StringFunc, state, input)
	case *pb.IntToFloatFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.IntToFloatFunction_BoolIntFunc:
		return ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.IntToFloatFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.IntToFloatFunction_BoolStringFunc:
		return ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.IntToFloatFunction_IntBoolFunc:
		return ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.IntToFloatFunction_IntIntFunc:
		return ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc, state, input)
	case *pb.IntToFloatFunction_IntFloatFunc:
		return ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.IntToFloatFunction_IntStringFunc:
		return ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc, state, input)
	case *pb.IntToFloatFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.IntToFloatFunction_FloatIntFunc:
		return ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.IntToFloatFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.IntToFloatFunction_FloatStringFunc:
		return ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.IntToFloatFunction_StringBoolFunc:
		return ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.IntToFloatFunction_StringIntFunc:
		return ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc, state, input)
	case *pb.IntToFloatFunction_StringFloatFunc:
		return ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.IntToFloatFunction_StringStringFunc:
		return ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc, state, input)
	case *pb.IntToFloatFunction_If:
		return ValidateIntValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from IntToFloatFunction: %T", v)
	}

	return nil
}

func ValidateIntToStringFunctionReferences(val *pb.IntToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntToStringFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.IntToStringFunction_Input:
		return ValidateReference(v.Input, input, TypeInt)
	case *pb.IntToStringFunction_State:
		return ValidateReference(v.State, state, TypeInt)
	case *pb.IntToStringFunction_BoolFunc:
		return ValidateBoolToIntFunctionReferences(v.BoolFunc, state, input)
	case *pb.IntToStringFunction_IntFunc:
		return ValidateIntToIntFunctionReferences(v.IntFunc, state, input)
	case *pb.IntToStringFunction_FloatFunc:
		return ValidateFloatToIntFunctionReferences(v.FloatFunc, state, input)
	case *pb.IntToStringFunction_StringFunc:
		return ValidateStringToIntFunctionReferences(v.StringFunc, state, input)
	case *pb.IntToStringFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.IntToStringFunction_BoolIntFunc:
		return ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.IntToStringFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.IntToStringFunction_BoolStringFunc:
		return ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.IntToStringFunction_IntBoolFunc:
		return ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.IntToStringFunction_IntIntFunc:
		return ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc, state, input)
	case *pb.IntToStringFunction_IntFloatFunc:
		return ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.IntToStringFunction_IntStringFunc:
		return ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc, state, input)
	case *pb.IntToStringFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.IntToStringFunction_FloatIntFunc:
		return ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.IntToStringFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.IntToStringFunction_FloatStringFunc:
		return ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.IntToStringFunction_StringBoolFunc:
		return ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.IntToStringFunction_StringIntFunc:
		return ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc, state, input)
	case *pb.IntToStringFunction_StringFloatFunc:
		return ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.IntToStringFunction_StringStringFunc:
		return ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc, state, input)
	case *pb.IntToStringFunction_If:
		return ValidateIntValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from IntToStringFunction: %T", v)
	}

	return nil
}

func ValidateFloatToBoolFunctionReferences(val *pb.FloatToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatToBoolFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.FloatToBoolFunction_Input:
		return ValidateReference(v.Input, input, TypeFloat)
	case *pb.FloatToBoolFunction_State:
		return ValidateReference(v.State, state, TypeFloat)
	case *pb.FloatToBoolFunction_BoolFunc:
		return ValidateBoolToFloatFunctionReferences(v.BoolFunc, state, input)
	case *pb.FloatToBoolFunction_IntFunc:
		return ValidateIntToFloatFunctionReferences(v.IntFunc, state, input)
	case *pb.FloatToBoolFunction_FloatFunc:
		return ValidateFloatToFloatFunctionReferences(v.FloatFunc, state, input)
	case *pb.FloatToBoolFunction_StringFunc:
		return ValidateStringToFloatFunctionReferences(v.StringFunc, state, input)
	case *pb.FloatToBoolFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.FloatToBoolFunction_BoolIntFunc:
		return ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.FloatToBoolFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.FloatToBoolFunction_BoolStringFunc:
		return ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.FloatToBoolFunction_IntBoolFunc:
		return ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.FloatToBoolFunction_IntIntFunc:
		return ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc, state, input)
	case *pb.FloatToBoolFunction_IntFloatFunc:
		return ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.FloatToBoolFunction_IntStringFunc:
		return ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc, state, input)
	case *pb.FloatToBoolFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.FloatToBoolFunction_FloatIntFunc:
		return ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.FloatToBoolFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.FloatToBoolFunction_FloatStringFunc:
		return ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.FloatToBoolFunction_StringBoolFunc:
		return ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.FloatToBoolFunction_StringIntFunc:
		return ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc, state, input)
	case *pb.FloatToBoolFunction_StringFloatFunc:
		return ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.FloatToBoolFunction_StringStringFunc:
		return ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc, state, input)
	case *pb.FloatToBoolFunction_If:
		return ValidateFloatValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from FloatToBoolFunction: %T", v)
	}

	return nil
}

func ValidateFloatToIntFunctionReferences(val *pb.FloatToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatToIntFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.FloatToIntFunction_Input:
		return ValidateReference(v.Input, input, TypeFloat)
	case *pb.FloatToIntFunction_State:
		return ValidateReference(v.State, state, TypeFloat)
	case *pb.FloatToIntFunction_BoolFunc:
		return ValidateBoolToFloatFunctionReferences(v.BoolFunc, state, input)
	case *pb.FloatToIntFunction_IntFunc:
		return ValidateIntToFloatFunctionReferences(v.IntFunc, state, input)
	case *pb.FloatToIntFunction_FloatFunc:
		return ValidateFloatToFloatFunctionReferences(v.FloatFunc, state, input)
	case *pb.FloatToIntFunction_StringFunc:
		return ValidateStringToFloatFunctionReferences(v.StringFunc, state, input)
	case *pb.FloatToIntFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.FloatToIntFunction_BoolIntFunc:
		return ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.FloatToIntFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.FloatToIntFunction_BoolStringFunc:
		return ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.FloatToIntFunction_IntBoolFunc:
		return ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.FloatToIntFunction_IntIntFunc:
		return ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc, state, input)
	case *pb.FloatToIntFunction_IntFloatFunc:
		return ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.FloatToIntFunction_IntStringFunc:
		return ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc, state, input)
	case *pb.FloatToIntFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.FloatToIntFunction_FloatIntFunc:
		return ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.FloatToIntFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.FloatToIntFunction_FloatStringFunc:
		return ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.FloatToIntFunction_StringBoolFunc:
		return ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.FloatToIntFunction_StringIntFunc:
		return ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc, state, input)
	case *pb.FloatToIntFunction_StringFloatFunc:
		return ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.FloatToIntFunction_StringStringFunc:
		return ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc, state, input)
	case *pb.FloatToIntFunction_If:
		return ValidateFloatValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from FloatToIntFunction: %T", v)
	}

	return nil
}

func ValidateFloatToFloatFunctionReferences(val *pb.FloatToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatToFloatFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.FloatToFloatFunction_Input:
		return ValidateReference(v.Input, input, TypeFloat)
	case *pb.FloatToFloatFunction_State:
		return ValidateReference(v.State, state, TypeFloat)
	case *pb.FloatToFloatFunction_BoolFunc:
		return ValidateBoolToFloatFunctionReferences(v.BoolFunc, state, input)
	case *pb.FloatToFloatFunction_IntFunc:
		return ValidateIntToFloatFunctionReferences(v.IntFunc, state, input)
	case *pb.FloatToFloatFunction_FloatFunc:
		return ValidateFloatToFloatFunctionReferences(v.FloatFunc, state, input)
	case *pb.FloatToFloatFunction_StringFunc:
		return ValidateStringToFloatFunctionReferences(v.StringFunc, state, input)
	case *pb.FloatToFloatFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.FloatToFloatFunction_BoolIntFunc:
		return ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.FloatToFloatFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.FloatToFloatFunction_BoolStringFunc:
		return ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.FloatToFloatFunction_IntBoolFunc:
		return ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.FloatToFloatFunction_IntIntFunc:
		return ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc, state, input)
	case *pb.FloatToFloatFunction_IntFloatFunc:
		return ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.FloatToFloatFunction_IntStringFunc:
		return ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc, state, input)
	case *pb.FloatToFloatFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.FloatToFloatFunction_FloatIntFunc:
		return ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.FloatToFloatFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.FloatToFloatFunction_FloatStringFunc:
		return ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.FloatToFloatFunction_StringBoolFunc:
		return ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.FloatToFloatFunction_StringIntFunc:
		return ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc, state, input)
	case *pb.FloatToFloatFunction_StringFloatFunc:
		return ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.FloatToFloatFunction_StringStringFunc:
		return ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc, state, input)
	case *pb.FloatToFloatFunction_If:
		return ValidateFloatValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from FloatToFloatFunction: %T", v)
	}

	return nil
}

func ValidateFloatToStringFunctionReferences(val *pb.FloatToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatToStringFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.FloatToStringFunction_Input:
		return ValidateReference(v.Input, input, TypeFloat)
	case *pb.FloatToStringFunction_State:
		return ValidateReference(v.State, state, TypeFloat)
	case *pb.FloatToStringFunction_BoolFunc:
		return ValidateBoolToFloatFunctionReferences(v.BoolFunc, state, input)
	case *pb.FloatToStringFunction_IntFunc:
		return ValidateIntToFloatFunctionReferences(v.IntFunc, state, input)
	case *pb.FloatToStringFunction_FloatFunc:
		return ValidateFloatToFloatFunctionReferences(v.FloatFunc, state, input)
	case *pb.FloatToStringFunction_StringFunc:
		return ValidateStringToFloatFunctionReferences(v.StringFunc, state, input)
	case *pb.FloatToStringFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.FloatToStringFunction_BoolIntFunc:
		return ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.FloatToStringFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.FloatToStringFunction_BoolStringFunc:
		return ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.FloatToStringFunction_IntBoolFunc:
		return ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.FloatToStringFunction_IntIntFunc:
		return ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc, state, input)
	case *pb.FloatToStringFunction_IntFloatFunc:
		return ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.FloatToStringFunction_IntStringFunc:
		return ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc, state, input)
	case *pb.FloatToStringFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.FloatToStringFunction_FloatIntFunc:
		return ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.FloatToStringFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.FloatToStringFunction_FloatStringFunc:
		return ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.FloatToStringFunction_StringBoolFunc:
		return ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.FloatToStringFunction_StringIntFunc:
		return ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc, state, input)
	case *pb.FloatToStringFunction_StringFloatFunc:
		return ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.FloatToStringFunction_StringStringFunc:
		return ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc, state, input)
	case *pb.FloatToStringFunction_If:
		return ValidateFloatValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from FloatToStringFunction: %T", v)
	}

	return nil
}

func ValidateStringToBoolFunctionReferences(val *pb.StringToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringToBoolFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.StringToBoolFunction_Input:
		return ValidateReference(v.Input, input, TypeString)
	case *pb.StringToBoolFunction_State:
		return ValidateReference(v.State, state, TypeString)
	case *pb.StringToBoolFunction_BoolFunc:
		return ValidateBoolToStringFunctionReferences(v.BoolFunc, state, input)
	case *pb.StringToBoolFunction_IntFunc:
		return ValidateIntToStringFunctionReferences(v.IntFunc, state, input)
	case *pb.StringToBoolFunction_FloatFunc:
		return ValidateFloatToStringFunctionReferences(v.FloatFunc, state, input)
	case *pb.StringToBoolFunction_StringFunc:
		return ValidateStringToStringFunctionReferences(v.StringFunc, state, input)
	case *pb.StringToBoolFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.StringToBoolFunction_BoolIntFunc:
		return ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.StringToBoolFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.StringToBoolFunction_BoolStringFunc:
		return ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.StringToBoolFunction_IntBoolFunc:
		return ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.StringToBoolFunction_IntIntFunc:
		return ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc, state, input)
	case *pb.StringToBoolFunction_IntFloatFunc:
		return ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.StringToBoolFunction_IntStringFunc:
		return ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc, state, input)
	case *pb.StringToBoolFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.StringToBoolFunction_FloatIntFunc:
		return ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.StringToBoolFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.StringToBoolFunction_FloatStringFunc:
		return ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.StringToBoolFunction_StringBoolFunc:
		return ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.StringToBoolFunction_StringIntFunc:
		return ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc, state, input)
	case *pb.StringToBoolFunction_StringFloatFunc:
		return ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.StringToBoolFunction_StringStringFunc:
		return ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc, state, input)
	case *pb.StringToBoolFunction_If:
		return ValidateStringValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from StringToBoolFunction: %T", v)
	}

	return nil
}

func ValidateStringToIntFunctionReferences(val *pb.StringToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringToIntFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.StringToIntFunction_Input:
		return ValidateReference(v.Input, input, TypeString)
	case *pb.StringToIntFunction_State:
		return ValidateReference(v.State, state, TypeString)
	case *pb.StringToIntFunction_BoolFunc:
		return ValidateBoolToStringFunctionReferences(v.BoolFunc, state, input)
	case *pb.StringToIntFunction_IntFunc:
		return ValidateIntToStringFunctionReferences(v.IntFunc, state, input)
	case *pb.StringToIntFunction_FloatFunc:
		return ValidateFloatToStringFunctionReferences(v.FloatFunc, state, input)
	case *pb.StringToIntFunction_StringFunc:
		return ValidateStringToStringFunctionReferences(v.StringFunc, state, input)
	case *pb.StringToIntFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.StringToIntFunction_BoolIntFunc:
		return ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.StringToIntFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.StringToIntFunction_BoolStringFunc:
		return ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.StringToIntFunction_IntBoolFunc:
		return ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.StringToIntFunction_IntIntFunc:
		return ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc, state, input)
	case *pb.StringToIntFunction_IntFloatFunc:
		return ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.StringToIntFunction_IntStringFunc:
		return ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc, state, input)
	case *pb.StringToIntFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.StringToIntFunction_FloatIntFunc:
		return ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.StringToIntFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.StringToIntFunction_FloatStringFunc:
		return ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.StringToIntFunction_StringBoolFunc:
		return ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.StringToIntFunction_StringIntFunc:
		return ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc, state, input)
	case *pb.StringToIntFunction_StringFloatFunc:
		return ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.StringToIntFunction_StringStringFunc:
		return ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc, state, input)
	case *pb.StringToIntFunction_If:
		return ValidateStringValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from StringToIntFunction: %T", v)
	}

	return nil
}

func ValidateStringToFloatFunctionReferences(val *pb.StringToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringToFloatFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.StringToFloatFunction_Input:
		return ValidateReference(v.Input, input, TypeString)
	case *pb.StringToFloatFunction_State:
		return ValidateReference(v.State, state, TypeString)
	case *pb.StringToFloatFunction_BoolFunc:
		return ValidateBoolToStringFunctionReferences(v.BoolFunc, state, input)
	case *pb.StringToFloatFunction_IntFunc:
		return ValidateIntToStringFunctionReferences(v.IntFunc, state, input)
	case *pb.StringToFloatFunction_FloatFunc:
		return ValidateFloatToStringFunctionReferences(v.FloatFunc, state, input)
	case *pb.StringToFloatFunction_StringFunc:
		return ValidateStringToStringFunctionReferences(v.StringFunc, state, input)
	case *pb.StringToFloatFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.StringToFloatFunction_BoolIntFunc:
		return ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.StringToFloatFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.StringToFloatFunction_BoolStringFunc:
		return ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.StringToFloatFunction_IntBoolFunc:
		return ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.StringToFloatFunction_IntIntFunc:
		return ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc, state, input)
	case *pb.StringToFloatFunction_IntFloatFunc:
		return ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.StringToFloatFunction_IntStringFunc:
		return ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc, state, input)
	case *pb.StringToFloatFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.StringToFloatFunction_FloatIntFunc:
		return ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.StringToFloatFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.StringToFloatFunction_FloatStringFunc:
		return ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.StringToFloatFunction_StringBoolFunc:
		return ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.StringToFloatFunction_StringIntFunc:
		return ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc, state, input)
	case *pb.StringToFloatFunction_StringFloatFunc:
		return ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.StringToFloatFunction_StringStringFunc:
		return ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc, state, input)
	case *pb.StringToFloatFunction_If:
		return ValidateStringValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from StringToFloatFunction: %T", v)
	}

	return nil
}

func ValidateStringToStringFunctionReferences(val *pb.StringToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringToStringFunction function missing")
	}

	switch v := val.Argument.(type) {
	case *pb.StringToStringFunction_Input:
		return ValidateReference(v.Input, input, TypeString)
	case *pb.StringToStringFunction_State:
		return ValidateReference(v.State, state, TypeString)
	case *pb.StringToStringFunction_BoolFunc:
		return ValidateBoolToStringFunctionReferences(v.BoolFunc, state, input)
	case *pb.StringToStringFunction_IntFunc:
		return ValidateIntToStringFunctionReferences(v.IntFunc, state, input)
	case *pb.StringToStringFunction_FloatFunc:
		return ValidateFloatToStringFunctionReferences(v.FloatFunc, state, input)
	case *pb.StringToStringFunction_StringFunc:
		return ValidateStringToStringFunctionReferences(v.StringFunc, state, input)
	case *pb.StringToStringFunction_BoolBoolFunc:
		return ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc, state, input)
	case *pb.StringToStringFunction_BoolIntFunc:
		return ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc, state, input)
	case *pb.StringToStringFunction_BoolFloatFunc:
		return ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc, state, input)
	case *pb.StringToStringFunction_BoolStringFunc:
		return ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc, state, input)
	case *pb.StringToStringFunction_IntBoolFunc:
		return ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc, state, input)
	case *pb.StringToStringFunction_IntIntFunc:
		return ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc, state, input)
	case *pb.StringToStringFunction_IntFloatFunc:
		return ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc, state, input)
	case *pb.StringToStringFunction_IntStringFunc:
		return ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc, state, input)
	case *pb.StringToStringFunction_FloatBoolFunc:
		return ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc, state, input)
	case *pb.StringToStringFunction_FloatIntFunc:
		return ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc, state, input)
	case *pb.StringToStringFunction_FloatFloatFunc:
		return ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc, state, input)
	case *pb.StringToStringFunction_FloatStringFunc:
		return ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc, state, input)
	case *pb.StringToStringFunction_StringBoolFunc:
		return ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc, state, input)
	case *pb.StringToStringFunction_StringIntFunc:
		return ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc, state, input)
	case *pb.StringToStringFunction_StringFloatFunc:
		return ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc, state, input)
	case *pb.StringToStringFunction_StringStringFunc:
		return ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc, state, input)
	case *pb.StringToStringFunction_If:
		return ValidateStringValueIf(v.If, state, input)
	default:
		return fmt.Errorf("unrecognized argument type from StringToStringFunction: %T", v)
	}

	return nil
}

// Binary Function Validation

func ValidateBoolAndBoolToBoolFunctionReferences(val *pb.BoolAndBoolToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndBoolToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndBoolToBoolFunction_Constant_1:
		return nil
	case *pb.BoolAndBoolToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndBoolToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndBoolToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToBoolFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndBoolToBoolFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndBoolToIntFunctionReferences(val *pb.BoolAndBoolToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndBoolToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndBoolToIntFunction_Constant_1:
		return nil
	case *pb.BoolAndBoolToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndBoolToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndBoolToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToIntFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndBoolToIntFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndBoolToFloatFunctionReferences(val *pb.BoolAndBoolToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndBoolToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndBoolToFloatFunction_Constant_1:
		return nil
	case *pb.BoolAndBoolToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndBoolToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndBoolToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToFloatFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndBoolToFloatFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndBoolToStringFunctionReferences(val *pb.BoolAndBoolToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndBoolToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndBoolToStringFunction_Constant_1:
		return nil
	case *pb.BoolAndBoolToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndBoolToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndBoolToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndBoolToStringFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndBoolToStringFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndIntToBoolFunctionReferences(val *pb.BoolAndIntToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndIntToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndIntToBoolFunction_Constant_1:
		return nil
	case *pb.BoolAndIntToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndIntToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndIntToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToBoolFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndIntToBoolFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndIntToIntFunctionReferences(val *pb.BoolAndIntToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndIntToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndIntToIntFunction_Constant_1:
		return nil
	case *pb.BoolAndIntToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndIntToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndIntToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToIntFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndIntToIntFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndIntToFloatFunctionReferences(val *pb.BoolAndIntToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndIntToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndIntToFloatFunction_Constant_1:
		return nil
	case *pb.BoolAndIntToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndIntToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndIntToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToFloatFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndIntToFloatFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndIntToStringFunctionReferences(val *pb.BoolAndIntToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndIntToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndIntToStringFunction_Constant_1:
		return nil
	case *pb.BoolAndIntToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndIntToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndIntToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndIntToStringFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndIntToStringFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndFloatToBoolFunctionReferences(val *pb.BoolAndFloatToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndFloatToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndFloatToBoolFunction_Constant_1:
		return nil
	case *pb.BoolAndFloatToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndFloatToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndFloatToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToBoolFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndFloatToBoolFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndFloatToIntFunctionReferences(val *pb.BoolAndFloatToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndFloatToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndFloatToIntFunction_Constant_1:
		return nil
	case *pb.BoolAndFloatToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndFloatToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndFloatToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToIntFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndFloatToIntFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndFloatToFloatFunctionReferences(val *pb.BoolAndFloatToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndFloatToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndFloatToFloatFunction_Constant_1:
		return nil
	case *pb.BoolAndFloatToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndFloatToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndFloatToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToFloatFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndFloatToFloatFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndFloatToStringFunctionReferences(val *pb.BoolAndFloatToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndFloatToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndFloatToStringFunction_Constant_1:
		return nil
	case *pb.BoolAndFloatToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndFloatToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndFloatToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndFloatToStringFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndFloatToStringFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndStringToBoolFunctionReferences(val *pb.BoolAndStringToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndStringToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndStringToBoolFunction_Constant_1:
		return nil
	case *pb.BoolAndStringToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndStringToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndStringToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToBoolFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndStringToBoolFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndStringToIntFunctionReferences(val *pb.BoolAndStringToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndStringToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndStringToIntFunction_Constant_1:
		return nil
	case *pb.BoolAndStringToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndStringToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndStringToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToIntFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndStringToIntFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndStringToFloatFunctionReferences(val *pb.BoolAndStringToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndStringToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndStringToFloatFunction_Constant_1:
		return nil
	case *pb.BoolAndStringToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndStringToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndStringToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToFloatFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndStringToFloatFunction: %T", v)
	}

	return nil
}

func ValidateBoolAndStringToStringFunctionReferences(val *pb.BoolAndStringToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("BoolAndStringToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.BoolAndStringToStringFunction_Constant_1:
		return nil
	case *pb.BoolAndStringToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeBool); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolFunc_1:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntFunc_1:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatFunc_1:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringFunc_1:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_If_1:
		if err := ValidateBoolValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from BoolAndStringToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.BoolAndStringToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.BoolAndStringToStringFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from BoolAndStringToStringFunction: %T", v)
	}

	return nil
}

func ValidateIntAndBoolToBoolFunctionReferences(val *pb.IntAndBoolToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndBoolToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndBoolToBoolFunction_Constant_1:
		return nil
	case *pb.IntAndBoolToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndBoolToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndBoolToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToBoolFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndBoolToBoolFunction: %T", v)
	}

	return nil
}

func ValidateIntAndBoolToIntFunctionReferences(val *pb.IntAndBoolToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndBoolToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndBoolToIntFunction_Constant_1:
		return nil
	case *pb.IntAndBoolToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndBoolToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndBoolToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToIntFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndBoolToIntFunction: %T", v)
	}

	return nil
}

func ValidateIntAndBoolToFloatFunctionReferences(val *pb.IntAndBoolToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndBoolToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndBoolToFloatFunction_Constant_1:
		return nil
	case *pb.IntAndBoolToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndBoolToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndBoolToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToFloatFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndBoolToFloatFunction: %T", v)
	}

	return nil
}

func ValidateIntAndBoolToStringFunctionReferences(val *pb.IntAndBoolToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndBoolToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndBoolToStringFunction_Constant_1:
		return nil
	case *pb.IntAndBoolToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndBoolToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndBoolToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndBoolToStringFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndBoolToStringFunction: %T", v)
	}

	return nil
}

func ValidateIntAndIntToBoolFunctionReferences(val *pb.IntAndIntToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndIntToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndIntToBoolFunction_Constant_1:
		return nil
	case *pb.IntAndIntToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndIntToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndIntToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToBoolFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndIntToBoolFunction: %T", v)
	}

	return nil
}

func ValidateIntAndIntToIntFunctionReferences(val *pb.IntAndIntToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndIntToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndIntToIntFunction_Constant_1:
		return nil
	case *pb.IntAndIntToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndIntToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndIntToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToIntFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndIntToIntFunction: %T", v)
	}

	return nil
}

func ValidateIntAndIntToFloatFunctionReferences(val *pb.IntAndIntToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndIntToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndIntToFloatFunction_Constant_1:
		return nil
	case *pb.IntAndIntToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndIntToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndIntToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToFloatFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndIntToFloatFunction: %T", v)
	}

	return nil
}

func ValidateIntAndIntToStringFunctionReferences(val *pb.IntAndIntToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndIntToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndIntToStringFunction_Constant_1:
		return nil
	case *pb.IntAndIntToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndIntToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndIntToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndIntToStringFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndIntToStringFunction: %T", v)
	}

	return nil
}

func ValidateIntAndFloatToBoolFunctionReferences(val *pb.IntAndFloatToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndFloatToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndFloatToBoolFunction_Constant_1:
		return nil
	case *pb.IntAndFloatToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndFloatToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndFloatToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToBoolFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndFloatToBoolFunction: %T", v)
	}

	return nil
}

func ValidateIntAndFloatToIntFunctionReferences(val *pb.IntAndFloatToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndFloatToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndFloatToIntFunction_Constant_1:
		return nil
	case *pb.IntAndFloatToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndFloatToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndFloatToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToIntFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndFloatToIntFunction: %T", v)
	}

	return nil
}

func ValidateIntAndFloatToFloatFunctionReferences(val *pb.IntAndFloatToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndFloatToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndFloatToFloatFunction_Constant_1:
		return nil
	case *pb.IntAndFloatToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndFloatToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndFloatToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToFloatFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndFloatToFloatFunction: %T", v)
	}

	return nil
}

func ValidateIntAndFloatToStringFunctionReferences(val *pb.IntAndFloatToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndFloatToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndFloatToStringFunction_Constant_1:
		return nil
	case *pb.IntAndFloatToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndFloatToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndFloatToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndFloatToStringFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndFloatToStringFunction: %T", v)
	}

	return nil
}

func ValidateIntAndStringToBoolFunctionReferences(val *pb.IntAndStringToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndStringToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndStringToBoolFunction_Constant_1:
		return nil
	case *pb.IntAndStringToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndStringToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndStringToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToBoolFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndStringToBoolFunction: %T", v)
	}

	return nil
}

func ValidateIntAndStringToIntFunctionReferences(val *pb.IntAndStringToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndStringToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndStringToIntFunction_Constant_1:
		return nil
	case *pb.IntAndStringToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndStringToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndStringToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToIntFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndStringToIntFunction: %T", v)
	}

	return nil
}

func ValidateIntAndStringToFloatFunctionReferences(val *pb.IntAndStringToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndStringToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndStringToFloatFunction_Constant_1:
		return nil
	case *pb.IntAndStringToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndStringToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndStringToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToFloatFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndStringToFloatFunction: %T", v)
	}

	return nil
}

func ValidateIntAndStringToStringFunctionReferences(val *pb.IntAndStringToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("IntAndStringToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.IntAndStringToStringFunction_Constant_1:
		return nil
	case *pb.IntAndStringToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeInt); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolFunc_1:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntFunc_1:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatFunc_1:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringFunc_1:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_If_1:
		if err := ValidateIntValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from IntAndStringToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.IntAndStringToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.IntAndStringToStringFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from IntAndStringToStringFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndBoolToBoolFunctionReferences(val *pb.FloatAndBoolToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndBoolToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndBoolToBoolFunction_Constant_1:
		return nil
	case *pb.FloatAndBoolToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndBoolToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndBoolToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToBoolFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndBoolToBoolFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndBoolToIntFunctionReferences(val *pb.FloatAndBoolToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndBoolToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndBoolToIntFunction_Constant_1:
		return nil
	case *pb.FloatAndBoolToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndBoolToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndBoolToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToIntFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndBoolToIntFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndBoolToFloatFunctionReferences(val *pb.FloatAndBoolToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndBoolToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndBoolToFloatFunction_Constant_1:
		return nil
	case *pb.FloatAndBoolToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndBoolToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndBoolToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToFloatFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndBoolToFloatFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndBoolToStringFunctionReferences(val *pb.FloatAndBoolToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndBoolToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndBoolToStringFunction_Constant_1:
		return nil
	case *pb.FloatAndBoolToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndBoolToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndBoolToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndBoolToStringFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndBoolToStringFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndIntToBoolFunctionReferences(val *pb.FloatAndIntToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndIntToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndIntToBoolFunction_Constant_1:
		return nil
	case *pb.FloatAndIntToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndIntToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndIntToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToBoolFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndIntToBoolFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndIntToIntFunctionReferences(val *pb.FloatAndIntToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndIntToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndIntToIntFunction_Constant_1:
		return nil
	case *pb.FloatAndIntToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndIntToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndIntToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToIntFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndIntToIntFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndIntToFloatFunctionReferences(val *pb.FloatAndIntToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndIntToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndIntToFloatFunction_Constant_1:
		return nil
	case *pb.FloatAndIntToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndIntToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndIntToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToFloatFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndIntToFloatFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndIntToStringFunctionReferences(val *pb.FloatAndIntToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndIntToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndIntToStringFunction_Constant_1:
		return nil
	case *pb.FloatAndIntToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndIntToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndIntToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndIntToStringFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndIntToStringFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndFloatToBoolFunctionReferences(val *pb.FloatAndFloatToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndFloatToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndFloatToBoolFunction_Constant_1:
		return nil
	case *pb.FloatAndFloatToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndFloatToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndFloatToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToBoolFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndFloatToBoolFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndFloatToIntFunctionReferences(val *pb.FloatAndFloatToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndFloatToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndFloatToIntFunction_Constant_1:
		return nil
	case *pb.FloatAndFloatToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndFloatToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndFloatToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToIntFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndFloatToIntFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndFloatToFloatFunctionReferences(val *pb.FloatAndFloatToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndFloatToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndFloatToFloatFunction_Constant_1:
		return nil
	case *pb.FloatAndFloatToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndFloatToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndFloatToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToFloatFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndFloatToFloatFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndFloatToStringFunctionReferences(val *pb.FloatAndFloatToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndFloatToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndFloatToStringFunction_Constant_1:
		return nil
	case *pb.FloatAndFloatToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndFloatToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndFloatToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndFloatToStringFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndFloatToStringFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndStringToBoolFunctionReferences(val *pb.FloatAndStringToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndStringToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndStringToBoolFunction_Constant_1:
		return nil
	case *pb.FloatAndStringToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndStringToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndStringToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToBoolFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndStringToBoolFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndStringToIntFunctionReferences(val *pb.FloatAndStringToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndStringToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndStringToIntFunction_Constant_1:
		return nil
	case *pb.FloatAndStringToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndStringToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndStringToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToIntFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndStringToIntFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndStringToFloatFunctionReferences(val *pb.FloatAndStringToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndStringToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndStringToFloatFunction_Constant_1:
		return nil
	case *pb.FloatAndStringToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndStringToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndStringToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToFloatFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndStringToFloatFunction: %T", v)
	}

	return nil
}

func ValidateFloatAndStringToStringFunctionReferences(val *pb.FloatAndStringToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("FloatAndStringToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.FloatAndStringToStringFunction_Constant_1:
		return nil
	case *pb.FloatAndStringToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeFloat); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolFunc_1:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntFunc_1:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatFunc_1:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringFunc_1:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_If_1:
		if err := ValidateFloatValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from FloatAndStringToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.FloatAndStringToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.FloatAndStringToStringFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from FloatAndStringToStringFunction: %T", v)
	}

	return nil
}

func ValidateStringAndBoolToBoolFunctionReferences(val *pb.StringAndBoolToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndBoolToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndBoolToBoolFunction_Constant_1:
		return nil
	case *pb.StringAndBoolToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndBoolToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndBoolToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToBoolFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndBoolToBoolFunction: %T", v)
	}

	return nil
}

func ValidateStringAndBoolToIntFunctionReferences(val *pb.StringAndBoolToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndBoolToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndBoolToIntFunction_Constant_1:
		return nil
	case *pb.StringAndBoolToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndBoolToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndBoolToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToIntFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndBoolToIntFunction: %T", v)
	}

	return nil
}

func ValidateStringAndBoolToFloatFunctionReferences(val *pb.StringAndBoolToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndBoolToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndBoolToFloatFunction_Constant_1:
		return nil
	case *pb.StringAndBoolToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndBoolToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndBoolToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToFloatFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndBoolToFloatFunction: %T", v)
	}

	return nil
}

func ValidateStringAndBoolToStringFunctionReferences(val *pb.StringAndBoolToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndBoolToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndBoolToStringFunction_Constant_1:
		return nil
	case *pb.StringAndBoolToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndBoolToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndBoolToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeBool); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeBool); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolFunc_2:
		if err := ValidateBoolToBoolFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntFunc_2:
		if err := ValidateIntToBoolFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatFunc_2:
		if err := ValidateFloatToBoolFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringFunc_2:
		if err := ValidateStringToBoolFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToBoolFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToBoolFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToBoolFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToBoolFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToBoolFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToBoolFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToBoolFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToBoolFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToBoolFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToBoolFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToBoolFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToBoolFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToBoolFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToBoolFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToBoolFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToBoolFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndBoolToStringFunction_If_2:
		if err := ValidateBoolValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndBoolToStringFunction: %T", v)
	}

	return nil
}

func ValidateStringAndIntToBoolFunctionReferences(val *pb.StringAndIntToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndIntToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndIntToBoolFunction_Constant_1:
		return nil
	case *pb.StringAndIntToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndIntToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndIntToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToBoolFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndIntToBoolFunction: %T", v)
	}

	return nil
}

func ValidateStringAndIntToIntFunctionReferences(val *pb.StringAndIntToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndIntToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndIntToIntFunction_Constant_1:
		return nil
	case *pb.StringAndIntToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndIntToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndIntToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToIntFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndIntToIntFunction: %T", v)
	}

	return nil
}

func ValidateStringAndIntToFloatFunctionReferences(val *pb.StringAndIntToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndIntToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndIntToFloatFunction_Constant_1:
		return nil
	case *pb.StringAndIntToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndIntToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndIntToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToFloatFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndIntToFloatFunction: %T", v)
	}

	return nil
}

func ValidateStringAndIntToStringFunctionReferences(val *pb.StringAndIntToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndIntToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndIntToStringFunction_Constant_1:
		return nil
	case *pb.StringAndIntToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndIntToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndIntToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeInt); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeInt); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolFunc_2:
		if err := ValidateBoolToIntFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntFunc_2:
		if err := ValidateIntToIntFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatFunc_2:
		if err := ValidateFloatToIntFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringFunc_2:
		if err := ValidateStringToIntFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToIntFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToIntFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToIntFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToIntFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToIntFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToIntFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToIntFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToIntFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToIntFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToIntFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToIntFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToIntFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToIntFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToIntFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToIntFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToIntFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndIntToStringFunction_If_2:
		if err := ValidateIntValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndIntToStringFunction: %T", v)
	}

	return nil
}

func ValidateStringAndFloatToBoolFunctionReferences(val *pb.StringAndFloatToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndFloatToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndFloatToBoolFunction_Constant_1:
		return nil
	case *pb.StringAndFloatToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndFloatToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndFloatToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToBoolFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndFloatToBoolFunction: %T", v)
	}

	return nil
}

func ValidateStringAndFloatToIntFunctionReferences(val *pb.StringAndFloatToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndFloatToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndFloatToIntFunction_Constant_1:
		return nil
	case *pb.StringAndFloatToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndFloatToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndFloatToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToIntFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndFloatToIntFunction: %T", v)
	}

	return nil
}

func ValidateStringAndFloatToFloatFunctionReferences(val *pb.StringAndFloatToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndFloatToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndFloatToFloatFunction_Constant_1:
		return nil
	case *pb.StringAndFloatToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndFloatToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndFloatToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToFloatFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndFloatToFloatFunction: %T", v)
	}

	return nil
}

func ValidateStringAndFloatToStringFunctionReferences(val *pb.StringAndFloatToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndFloatToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndFloatToStringFunction_Constant_1:
		return nil
	case *pb.StringAndFloatToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndFloatToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndFloatToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeFloat); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeFloat); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolFunc_2:
		if err := ValidateBoolToFloatFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntFunc_2:
		if err := ValidateIntToFloatFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatFunc_2:
		if err := ValidateFloatToFloatFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringFunc_2:
		if err := ValidateStringToFloatFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToFloatFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToFloatFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToFloatFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToFloatFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToFloatFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToFloatFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToFloatFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToFloatFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToFloatFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToFloatFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToFloatFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToFloatFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToFloatFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToFloatFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToFloatFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToFloatFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndFloatToStringFunction_If_2:
		if err := ValidateFloatValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndFloatToStringFunction: %T", v)
	}

	return nil
}

func ValidateStringAndStringToBoolFunctionReferences(val *pb.StringAndStringToBoolFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndStringToBoolFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndStringToBoolFunction_Constant_1:
		return nil
	case *pb.StringAndStringToBoolFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndStringToBoolFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndStringToBoolFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToBoolFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndStringToBoolFunction: %T", v)
	}

	return nil
}

func ValidateStringAndStringToIntFunctionReferences(val *pb.StringAndStringToIntFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndStringToIntFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndStringToIntFunction_Constant_1:
		return nil
	case *pb.StringAndStringToIntFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndStringToIntFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndStringToIntFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToIntFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndStringToIntFunction: %T", v)
	}

	return nil
}

func ValidateStringAndStringToFloatFunctionReferences(val *pb.StringAndStringToFloatFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndStringToFloatFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndStringToFloatFunction_Constant_1:
		return nil
	case *pb.StringAndStringToFloatFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndStringToFloatFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndStringToFloatFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToFloatFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndStringToFloatFunction: %T", v)
	}

	return nil
}

func ValidateStringAndStringToStringFunctionReferences(val *pb.StringAndStringToStringFunction, state, input *desc.MessageDescriptor) error {
	if val == nil {
		return fmt.Errorf("StringAndStringToStringFunction function missing")
	}

	switch v := val.Argument_1.(type) {
	case *pb.StringAndStringToStringFunction_Constant_1:
		return nil
	case *pb.StringAndStringToStringFunction_Input_1:
		if err := ValidateReference(v.Input_1, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_State_1:
		if err := ValidateReference(v.State_1, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolFunc_1:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntFunc_1:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatFunc_1:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringFunc_1:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolBoolFunc_1:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolIntFunc_1:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolFloatFunc_1:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolStringFunc_1:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntBoolFunc_1:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntIntFunc_1:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntFloatFunc_1:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntStringFunc_1:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatBoolFunc_1:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatIntFunc_1:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatFloatFunc_1:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatStringFunc_1:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringBoolFunc_1:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringIntFunc_1:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringFloatFunc_1:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringStringFunc_1:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_1, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_If_1:
		if err := ValidateStringValueIf(v.If_1, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_1 type from StringAndStringToStringFunction: %T", v)
	}

	switch v := val.Argument_2.(type) {
	case *pb.StringAndStringToStringFunction_Input_2:
		if err := ValidateReference(v.Input_2, input, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_State_2:
		if err := ValidateReference(v.State_2, state, TypeString); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolFunc_2:
		if err := ValidateBoolToStringFunctionReferences(v.BoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntFunc_2:
		if err := ValidateIntToStringFunctionReferences(v.IntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatFunc_2:
		if err := ValidateFloatToStringFunctionReferences(v.FloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringFunc_2:
		if err := ValidateStringToStringFunctionReferences(v.StringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolBoolFunc_2:
		if err := ValidateBoolAndBoolToStringFunctionReferences(v.BoolBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolIntFunc_2:
		if err := ValidateBoolAndIntToStringFunctionReferences(v.BoolIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolFloatFunc_2:
		if err := ValidateBoolAndFloatToStringFunctionReferences(v.BoolFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_BoolStringFunc_2:
		if err := ValidateBoolAndStringToStringFunctionReferences(v.BoolStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntBoolFunc_2:
		if err := ValidateIntAndBoolToStringFunctionReferences(v.IntBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntIntFunc_2:
		if err := ValidateIntAndIntToStringFunctionReferences(v.IntIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntFloatFunc_2:
		if err := ValidateIntAndFloatToStringFunctionReferences(v.IntFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_IntStringFunc_2:
		if err := ValidateIntAndStringToStringFunctionReferences(v.IntStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatBoolFunc_2:
		if err := ValidateFloatAndBoolToStringFunctionReferences(v.FloatBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatIntFunc_2:
		if err := ValidateFloatAndIntToStringFunctionReferences(v.FloatIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatFloatFunc_2:
		if err := ValidateFloatAndFloatToStringFunctionReferences(v.FloatFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_FloatStringFunc_2:
		if err := ValidateFloatAndStringToStringFunctionReferences(v.FloatStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringBoolFunc_2:
		if err := ValidateStringAndBoolToStringFunctionReferences(v.StringBoolFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringIntFunc_2:
		if err := ValidateStringAndIntToStringFunctionReferences(v.StringIntFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringFloatFunc_2:
		if err := ValidateStringAndFloatToStringFunctionReferences(v.StringFloatFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_StringStringFunc_2:
		if err := ValidateStringAndStringToStringFunctionReferences(v.StringStringFunc_2, state, input); err != nil {
			return err
		}
	case *pb.StringAndStringToStringFunction_If_2:
		if err := ValidateStringValueIf(v.If_2, state, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unrecognized argument_2 type from StringAndStringToStringFunction: %T", v)
	}

	return nil
}

// If-Else Validation

func ValidateBoolValueIf(v *pb.BoolValueIf, state, input *desc.MessageDescriptor) error {
	if v.Predicate == nil {
		return fmt.Errorf("predicate missing from 'if' statement")
	}
	if v.Then == nil {
		return fmt.Errorf("'then' missing from 'if' statement")
	}
	if v.Else == nil {
		return fmt.Errorf("'else' missing from 'if' statement")
	}

	if err := ValidateBoolValueReferences(v.Predicate, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement predicate: %w", err)
	}
	if err := ValidateBoolValueReferences(v.Then, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement 'then' branch: %w", err)
	}
	if err := ValidateBoolValueReferences(v.Else, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement 'else' branch: %w", err)
	}

	return nil
}

func ValidateIntValueIf(v *pb.IntValueIf, state, input *desc.MessageDescriptor) error {
	if v.Predicate == nil {
		return fmt.Errorf("predicate missing from 'if' statement")
	}
	if v.Then == nil {
		return fmt.Errorf("'then' missing from 'if' statement")
	}
	if v.Else == nil {
		return fmt.Errorf("'else' missing from 'if' statement")
	}

	if err := ValidateBoolValueReferences(v.Predicate, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement predicate: %w", err)
	}
	if err := ValidateIntValueReferences(v.Then, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement 'then' branch: %w", err)
	}
	if err := ValidateIntValueReferences(v.Else, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement 'else' branch: %w", err)
	}

	return nil
}

func ValidateFloatValueIf(v *pb.FloatValueIf, state, input *desc.MessageDescriptor) error {
	if v.Predicate == nil {
		return fmt.Errorf("predicate missing from 'if' statement")
	}
	if v.Then == nil {
		return fmt.Errorf("'then' missing from 'if' statement")
	}
	if v.Else == nil {
		return fmt.Errorf("'else' missing from 'if' statement")
	}

	if err := ValidateBoolValueReferences(v.Predicate, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement predicate: %w", err)
	}
	if err := ValidateFloatValueReferences(v.Then, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement 'then' branch: %w", err)
	}
	if err := ValidateFloatValueReferences(v.Else, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement 'else' branch: %w", err)
	}

	return nil
}

func ValidateStringValueIf(v *pb.StringValueIf, state, input *desc.MessageDescriptor) error {
	if v.Predicate == nil {
		return fmt.Errorf("predicate missing from 'if' statement")
	}
	if v.Then == nil {
		return fmt.Errorf("'then' missing from 'if' statement")
	}
	if v.Else == nil {
		return fmt.Errorf("'else' missing from 'if' statement")
	}

	if err := ValidateBoolValueReferences(v.Predicate, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement predicate: %w", err)
	}
	if err := ValidateStringValueReferences(v.Then, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement 'then' branch: %w", err)
	}
	if err := ValidateStringValueReferences(v.Else, state, input); err != nil {
		return fmt.Errorf("invalid 'if' statement 'else' branch: %w", err)
	}

	return nil
}
