package main

import (
	"fmt"
	"strings"

	pb "angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
	"angelbeltran/game-engine/protoc-gen-game/types"
)

func validateAction(state types.Type, msg *pb.Action) error {
	if msg.Error == nil || (len(msg.Error.Code)+len(msg.Error.Msg) == 0) {
		return fmt.Errorf("no error defined on action")
	}
	if msg.Rule == nil {
		return fmt.Errorf("no rule defined on action")
	}
	if msg.Effect == nil {
		return fmt.Errorf("no effect defined on action")
	}
	if err := validateRule(state, msg.Rule); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}
	for _, effect := range msg.Effect {
		if err := validateEffect(state, effect); err != nil {
			return fmt.Errorf("invalid effect: %w", err)
		}
	}

	return nil
}

func validateEffect(state types.Type, effect *pb.Effect) error {
	if up := effect.GetUpdate(); up != nil {
		var srcType types.Type

		if up.Src == nil {
			return fmt.Errorf("missing destination on update effect")
		}

		switch src := up.Src.GetOperand().(type) {
		case *pb.Operand_Field:
			if src.Field == nil {
				return fmt.Errorf("empty effect source field on update effect")
			}

			var exists bool
			if srcType, exists = getField(state, src.Field.Name); !exists {
				return fmt.Errorf("update effect source field '%s' does not exist", strings.Join(src.Field.Name, "."))
			}

		case *pb.Operand_Value:
			if src.Value == nil {
				return fmt.Errorf("empty effect source value on update effect")
			}

			var err error
			if srcType, err = getValueType(src.Value); err != nil {
				return err
			}
		}

		dest := up.GetDest()
		if dest == nil {
			return fmt.Errorf("missing destination on update effect")
		}

		destType, exists := getField(state, dest.Name)
		if !exists {
			return fmt.Errorf("update effect destination field '%s' does not exist", strings.Join(dest.Name, "."))
		}

		if !srcType.IsSameType(destType) {
			return fmt.Errorf("mismatch operand types on update effect: %s and %s", srcType, destType)
		}

		return nil
	}

	return fmt.Errorf("unrecognized effect type: %T", effect.GetOperation())
}

func validateRule(state types.Type, rule *pb.Rule) error {
	if exp := rule.GetSingle(); exp != nil {
		if exp.Left == nil {
			return fmt.Errorf("missing left-hand operand")
		}
		if exp.Right == nil {
			return fmt.Errorf("missing right-hand operand")
		}
		if exp.Operator == 0 {
			return fmt.Errorf("missing operator")
		}

		lh, err := getOperandType(state, exp.Left)
		if err != nil {
			return fmt.Errorf("invalid left-hand operand: %w", err)
		}

		rh, err := getOperandType(state, exp.Right)
		if err != nil {
			return fmt.Errorf("invalid right-hand operand: %w", err)
		}

		if !lh.IsSameType(rh) {
			return fmt.Errorf("mismatch operand types: %s and %s", lh, rh)
		}

		if err := validateOperator(lh, exp.Operator); err != nil {
			return err
		}

		return nil
	}

	if and := rule.GetAnd(); and != nil {
		rules := and.GetRules()
		if len(rules) == 0 {
			return fmt.Errorf("empty 'and' condition")
		}

		for _, r := range rules {
			if err := validateRule(state, r); err != nil {
				return err
			}
		}

		return nil
	}

	if or := rule.GetOr(); or != nil {
		rules := or.GetRules()
		if len(rules) == 0 {
			return fmt.Errorf("empty 'or' condition")
		}

		for _, r := range rules {
			if err := validateRule(state, r); err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("rule with no conditions found")
}

func getOperandType(state types.Type, o *pb.Operand) (types.Type, error) {

	if v := o.GetValue(); v != nil {
		return getValueType(v)
	}

	if v := o.GetField(); v != nil {
		field, exists := getField(state, []string(v.Name))
		if !exists {
			return nil, fmt.Errorf("field '%s' does not exist", strings.Join(v.Name, "."))
		}

		return field, nil
	}

	return nil, fmt.Errorf("unexpected operand type: %T", o)
}

func getValueType(val *pb.Value) (types.Type, error) {
	switch val.Value.(type) {
	case *pb.Value_BoolValue:
		return types.Bool{}, nil
	case *pb.Value_IntegerValue:
		return types.Integer{}, nil
	case *pb.Value_FloatValue:
		return types.Float{}, nil
	case *pb.Value_StringValue:
		return types.String{}, nil
	}

	return nil, fmt.Errorf("unrecognized operand value type: %T", val.Value)
}

func getField(t types.Type, path []string) (field types.Type, found bool) {
	if len(path) == 0 {
		return nil, false
	}

	switch v := t.(type) {
	case types.OneOf:
		field, found = v[path[0]]
	case types.Structured:
		field, found = v[path[0]]
	case types.Map:
		field = v.Value
		found = true
	}

	if len(path) == 1 || !found {
		return field, found
	}

	return getField(field, path[1:])
}

func validateOperator(t types.Type, op pb.SingleRule_Operator) error {

	switch t.(type) {
	case types.Bool:
		switch op {
		case pb.SingleRule_EQ, pb.SingleRule_NEQ:
		default:
			return fmt.Errorf("operator %s is incompatible with boolean values", op)
		}
	case types.Integer:
		switch op {
		case pb.SingleRule_EQ, pb.SingleRule_NEQ,
			pb.SingleRule_LT, pb.SingleRule_LTE,
			pb.SingleRule_GT, pb.SingleRule_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with integer values", op)
		}
	case types.Float:
		switch op {
		case pb.SingleRule_EQ, pb.SingleRule_NEQ,
			pb.SingleRule_LT, pb.SingleRule_LTE,
			pb.SingleRule_GT, pb.SingleRule_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with float values", op)
		}
	case types.String:
		switch op {
		case pb.SingleRule_EQ, pb.SingleRule_NEQ,
			pb.SingleRule_LT, pb.SingleRule_LTE,
			pb.SingleRule_GT, pb.SingleRule_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with string values", op)
		}
	case types.Bytes:
		switch op {
		case pb.SingleRule_EQ, pb.SingleRule_NEQ,
			pb.SingleRule_LT, pb.SingleRule_LTE,
			pb.SingleRule_GT, pb.SingleRule_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with bytes values", op)
		}
	case types.Enum:
		switch op {
		case pb.SingleRule_EQ, pb.SingleRule_NEQ,
			pb.SingleRule_LT, pb.SingleRule_LTE,
			pb.SingleRule_GT, pb.SingleRule_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with enum values", op)
		}
	case types.OneOf:
		return fmt.Errorf("no direct operator support exists for oneofs at this time. please reference one of its fields")
	case types.List:
		return fmt.Errorf("no direct operator support exists for lists at this time")
	case types.Structured:
		return fmt.Errorf("no direct operator support exists for structured at this time. please reference one of its fields")
	case types.Map:
		return fmt.Errorf("no direct operator support exists for maps at this time. please reference one of its fields")
	default:
		return fmt.Errorf("unrecognized operand type: %T", t)
	}

	return nil
}
