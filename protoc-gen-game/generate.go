package main

import (
	"bytes"
	"go/format"
	"io"
	"io/ioutil"

	"github.com/jhump/protoreflect/desc"

	pb "angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
	"angelbeltran/game-engine/protoc-gen-game/types"
)

type (
	methodBundle struct {
		Method *desc.MethodDescriptor
		Input  types.Type
		Action *pb.Action
	}

	generationOptions struct {
		Package   string
		Service   *desc.ServiceDescriptor
		Methods   []methodBundle
		State     *desc.MessageDescriptor
		Response  *desc.MessageDescriptor
		StateType types.Type
	}
	run func() error
)

func generateAll(w io.Writer, opts generationOptions) error {
	out := bytes.NewBuffer([]byte{})

	if err := mainTemplate.Execute(out, mainTemplateParameters{
		Package:  opts.Package,
		Imports:  defaultMethodImports,
		Service:  opts.Service,
		Methods:  opts.Methods,
		State:    opts.State,
		Response: opts.Response,
	}); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(io.MultiReader(out))
	if err != nil {
		return err
	}

	b, err = format.Source(b)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}
