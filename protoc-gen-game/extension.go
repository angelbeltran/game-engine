package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"

	pb "angelbeltran/game-engine/protoc-gen-game/protos/game_engine_pb"
)

func createMessageFactoryWithActionExtension(files []*desc.FileDescriptor) (*dynamic.MessageFactory, error) {
	factory := dynamic.NewMessageFactoryWithDefaults()

	var gameEngineOptionsFile *desc.FileDescriptor

	for _, file := range files {
		for _, dep := range file.GetDependencies() {
			if dep.GetPackage() == protoPackageName {
				gameEngineOptionsFile = dep
				break
			}
		}

		if gameEngineOptionsFile != nil {
			break
		}
	}

	if gameEngineOptionsFile == nil {
		return factory, fmt.Errorf("no .proto files found within the %s package", protoPackageName)
	}

	actionRuleExt := gameEngineOptionsFile.FindExtension("google.protobuf.MethodOptions", actionExtensionFieldNumber)

	if err := factory.GetExtensionRegistry().AddExtension(actionRuleExt); err != nil {
		return factory, fmt.Errorf("failed to add new extension: %w", err)
	}

	return factory, nil
}

func loadActionOptionMessage(method *desc.MethodDescriptor, field int) (*pb.Action, error) {
	methodOpts := method.GetMethodOptions()
	if methodOpts == nil {
		return nil, nil
	}

	extDescs, err := proto.ExtensionDescs(methodOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to examine method extension descriptions: %w", err)
	}

	var extDesc *proto.ExtensionDesc
	for _, d := range extDescs {
		if d.Field == int32(field) {
			extDesc = d
			break
		}
	}
	if extDesc == nil {
		return nil, nil
	}

	ext, err := proto.GetExtension(methodOpts, extDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to examine method extension: %w", err)
	}

	m, ok := ext.(*pb.Action)
	if !ok {
		return nil, fmt.Errorf("unexpected action rule option type: %T", ext)
	}

	return m, nil
}
