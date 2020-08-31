package main

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"

	pb "github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
	"github.com/angelbeltran/game-engine/protoc-gen-game/generate/dst/go/validation"
)

func validateAction(state, input *desc.MessageDescriptor, msg *pb.Action) error {
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

	if err := validation.ValidateBoolValueReferences(msg.Rule, state, input); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	if err := validateResponse(state, msg.Response); err != nil {
		return fmt.Errorf("invalid response: %w", err)
	}

	return nil
}

func validateEffect(state, input *desc.MessageDescriptor, effect *pb.Effect) error {
	if up := effect.GetUpdate(); up != nil {
		valueType, err := validation.ValidateValue(up.Value, state, input)
		if err != nil {
			return fmt.Errorf("invalid value: %w", err)
		}

		if up.State == nil {
			return fmt.Errorf("missing state on update effect")
		}

		if err := validation.ValidateReference(up.State, state, valueType); err != nil {
			return fmt.Errorf("invalid state reference: %w", err)
		}

		return nil
	}

	return fmt.Errorf("unrecognized effect type: %T", effect.GetOperation())
}

func validateResponse(state *desc.MessageDescriptor, refs []*pb.Reference) error {
	for _, ref := range refs {
		if ref == nil {
			return fmt.Errorf("empty state reference")
		}

		if err := validation.VerifyEndOfPath(ref.Path, state); err != nil {
			return fmt.Errorf("invalid state reference: %w", err)
		}
	}

	return nil
}
