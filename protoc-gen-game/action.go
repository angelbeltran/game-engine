package main

import (
	"fmt"
	"strings"

	pb "angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
	"angelbeltran/game-engine/protoc-gen-game/types"
)

func validateAction(state, input types.Type, msg *pb.Action) error {
	if msg.Effect == nil {
		return fmt.Errorf("no effect defined on action")
	}
	if msg.Rule == nil {
		return fmt.Errorf("no rule defined on action")
	}
	if msg.Error == nil || (len(msg.Error.Code)+len(msg.Error.Msg) == 0) {
		return fmt.Errorf("no error defined on action")
	}

	for _, effect := range msg.Effect {
		if err := validateEffect(state, input, effect); err != nil {
			return fmt.Errorf("invalid effect: %w", err)
		}
	}

	if err := validateRule(state, input, msg.Rule); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	if err := validateResponse(state, msg.Response); err != nil {
		return fmt.Errorf("invalid response: %w", err)
	}

	return nil
}

func validateEffect(state, input types.Type, effect *pb.Effect) error {
	if up := effect.GetUpdate(); up != nil {
		var srcType types.Type

		if up.Src == nil {
			return fmt.Errorf("missing destination on update effect")
		}

		switch src := up.Src.GetOperand().(type) {
		case *pb.Operand_Value:
			if src.Value == nil {
				return fmt.Errorf("empty effect source value on update effect")
			}

			var err error
			if srcType, _, err = extractValue(src.Value); err != nil {
				return err
			}

		case *pb.Operand_Prop:
			if src.Prop == nil {
				return fmt.Errorf("empty effect source property on update effect")
			}

			var exists bool
			if srcType, exists = resolvePath(state, src.Prop.Path); !exists {
				return fmt.Errorf("update effect source property '%s' does not exist", strings.Join(src.Prop.Path, "."))
			}

		case *pb.Operand_Input:
			if src.Input == nil {
				return fmt.Errorf("empty effect source input on update effect")
			}

			var exists bool
			if srcType, exists = resolvePath(input, src.Input.Path); !exists {
				return fmt.Errorf("update effect source input '%s' does not exist", strings.Join(src.Input.Path, "."))
			}
		}

		dest := up.GetDest()
		if dest == nil {
			return fmt.Errorf("missing destination on update effect")
		}

		destType, exists := resolvePath(state, dest.Path)
		if !exists {
			return fmt.Errorf("update effect destination property '%s' does not exist", strings.Join(dest.Path, "."))
		}

		if !srcType.IsSameType(destType) {
			return fmt.Errorf("mismatch operand types on update effect: %s and %s", srcType, destType)
		}

		return nil
	}

	return fmt.Errorf("unrecognized effect type: %T", effect.GetOperation())
}

func validateRule(state, input types.Type, rule *pb.Rule) error {
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

		lh, err := resolveOperandType(state, input, exp.Left)
		if err != nil {
			return fmt.Errorf("invalid left-hand operand: %w", err)
		}

		rh, err := resolveOperandType(state, input, exp.Right)
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
			if err := validateRule(state, input, r); err != nil {
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
			if err := validateRule(state, input, r); err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("rule with no conditions found")
}

func validateResponse(state types.Type, res []*pb.Path) error {
	for _, p := range res {
		if p == nil {
			continue
		}

		if _, exists := resolvePath(state, p.Path); !exists {
			return fmt.Errorf("property '%s' does not exist", strings.Join(p.Path, "."))
		}
	}

	return nil
}

func resolveOperandType(state, input types.Type, o *pb.Operand) (types.Type, error) {

	if v := o.GetValue(); v != nil {
		t, _, err := extractValue(v)
		return t, err
	}

	if v := o.GetProp(); v != nil {
		property, exists := resolvePath(state, []string(v.Path))
		if !exists {
			return nil, fmt.Errorf("property '%s' does not exist", strings.Join(v.Path, "."))
		}

		return property, nil
	}

	if v := o.GetInput(); v != nil {
		input, exists := resolvePath(state, []string(v.Path))
		if !exists {
			return nil, fmt.Errorf("input '%s' does not exist", strings.Join(v.Path, "."))
		}

		return input, nil
	}

	return nil, fmt.Errorf("unexpected operand type: %T", o)
}

func resolvePath(t types.Type, path []string) (prop types.Type, found bool) {
	if len(path) == 0 {
		return nil, false
	}

	switch v := t.(type) {
	case types.OneOf:
		prop, found = v[path[0]]
	case types.Structured:
		prop, found = v[path[0]]
	case types.Map:
		prop = v.Value
		found = true
	}

	if len(path) == 1 || !found {
		return prop, found
	}

	return resolvePath(prop, path[1:])
}

func validateOperator(t types.Type, op pb.Rule_Single_Operator) error {

	switch t.(type) {
	case types.Bool:
		switch op {
		case pb.Rule_Single_EQ, pb.Rule_Single_NEQ:
		default:
			return fmt.Errorf("operator %s is incompatible with boolean values", op)
		}
	case types.Integer:
		switch op {
		case pb.Rule_Single_EQ, pb.Rule_Single_NEQ,
			pb.Rule_Single_LT, pb.Rule_Single_LTE,
			pb.Rule_Single_GT, pb.Rule_Single_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with integer values", op)
		}
	case types.Float:
		switch op {
		case pb.Rule_Single_EQ, pb.Rule_Single_NEQ,
			pb.Rule_Single_LT, pb.Rule_Single_LTE,
			pb.Rule_Single_GT, pb.Rule_Single_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with float values", op)
		}
	case types.String:
		switch op {
		case pb.Rule_Single_EQ, pb.Rule_Single_NEQ,
			pb.Rule_Single_LT, pb.Rule_Single_LTE,
			pb.Rule_Single_GT, pb.Rule_Single_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with string values", op)
		}
	case types.Bytes:
		switch op {
		case pb.Rule_Single_EQ, pb.Rule_Single_NEQ,
			pb.Rule_Single_LT, pb.Rule_Single_LTE,
			pb.Rule_Single_GT, pb.Rule_Single_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with bytes values", op)
		}
	case types.Enum:
		switch op {
		case pb.Rule_Single_EQ, pb.Rule_Single_NEQ,
			pb.Rule_Single_LT, pb.Rule_Single_LTE,
			pb.Rule_Single_GT, pb.Rule_Single_GTE:
		default:
			return fmt.Errorf("operator %s is incompatible with enum values", op)
		}
	case types.OneOf:
		return fmt.Errorf("no direct operator support exists for oneofs at this time. please reference one of its properties")
	case types.List:
		return fmt.Errorf("no direct operator support exists for lists at this time")
	case types.Structured:
		return fmt.Errorf("no direct operator support exists for structured at this time. please reference one of its properties")
	case types.Map:
		return fmt.Errorf("no direct operator support exists for maps at this time. please reference one of its properties")
	default:
		return fmt.Errorf("unrecognized operand type: %T", t)
	}

	return nil
}
