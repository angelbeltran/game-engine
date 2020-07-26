package main

import (
	"strings"

	"github.com/jhump/protoreflect/desc"

	pb "github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
	"github.com/angelbeltran/game-engine/protoc-gen-game/types"
)

type (
	serviceParams struct {
		Package   string
		Service   *desc.ServiceDescriptor
		Methods   []methodInfo
		State     *desc.MessageDescriptor
		Response  *desc.MessageDescriptor
		StateType types.Type
	}

	methodInfo struct {
		Method *desc.MethodDescriptor
		Input  types.Type
		Action *pb.Action
	}

	effectParams struct {
		State  *desc.MessageDescriptor
		Effect *pb.Effect
	}
	updateEffectParams struct {
		State  *desc.MessageDescriptor
		Update *pb.Effect_Update
	}

	initializeStatePropertyExpressionParams struct {
		Path       []string
		Identifier string
		State      *desc.MessageDescriptor
	}

	responseAppendExpressionTemplateParams struct {
		State *desc.MessageDescriptor
		Path  []string
	}
)

func (serviceParams) Imports() []string {
	return []string{
		"context",
		"fmt",
		"net",
		"sync",
		"google.golang.org/grpc",
		"github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb",
	}
}

func newEffectParams(state *desc.MessageDescriptor, effect *pb.Effect) effectParams {
	return effectParams{
		State:  state,
		Effect: effect,
	}
}

func newUpdateEffectParams(state *desc.MessageDescriptor, update *pb.Effect_Update) updateEffectParams {
	return updateEffectParams{
		State:  state,
		Update: update,
	}
}

func newInitializeStatePropertyExpressionParams(path []string, identifier string, state *desc.MessageDescriptor) initializeStatePropertyExpressionParams {
	return initializeStatePropertyExpressionParams{
		Path:       path,
		Identifier: identifier,
		State:      state,
	}
}

func newResponseAppendExpressionTemplateParams(state *desc.MessageDescriptor, path []string) responseAppendExpressionTemplateParams {
	return responseAppendExpressionTemplateParams{
		State: state,
		Path:  path,
	}
}

func joinCamelCase(a []string) string {
	b := make([]string, len(a))
	for i, s := range a {
		b[i] = goNames.CamelCase(s)
	}
	return strings.Join(b, ".")
}

func inc(i int) int {
	return i + 1
}

func dec(i int) int {
	return i - 1
}
