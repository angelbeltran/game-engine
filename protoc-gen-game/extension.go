package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"

	pb "github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
)

func getMessagesWithExtension(files []*desc.FileDescriptor, field int) ([]*desc.MessageDescriptor, error) {
	var matches []*desc.MessageDescriptor

	traversed := make(map[string]bool)

	for _, file := range files {
		for _, file := range append([]*desc.FileDescriptor{file}, file.GetDependencies()...) {
			if traversed[file.GetFullyQualifiedName()] {
				continue
			}

			for _, md := range file.GetMessageTypes() {
				opts := md.GetMessageOptions()
				if opts == nil {
					continue
				}

				extensions, err := proto.ExtensionDescs(opts)
				if err != nil {
					return nil, fmt.Errorf("failed to examine message extensions descriptions: %w", err)
				}

				var found bool
				for _, ext := range extensions {
					if ext != nil && ext.Field == int32(field) {
						found = true
						break
					}
				}

				if found {
					matches = append(matches, md)
					break
				}
			}

			traversed[file.GetFullyQualifiedName()] = true
		}
	}

	return matches, nil
}

func getMessagesAndExtensions(files []*desc.FileDescriptor, field int) (map[*desc.MessageDescriptor]interface{}, error) {
	matches := make(map[*desc.MessageDescriptor]interface{})

	traversed := make(map[string]bool)

	for _, file := range files {
		for _, file := range append([]*desc.FileDescriptor{file}, file.GetDependencies()...) {
			if traversed[file.GetFullyQualifiedName()] {
				continue
			}

			msgTypes := file.GetMessageTypes()
			for i := 0; i < len(msgTypes); i++ {
				md := msgTypes[i]

				ext, err := getMessageExtension(md, field)
				if err != nil {
					return nil, err
				}

				if ext != nil {
					matches[md] = ext
				}

				msgTypes = append(msgTypes, md.GetNestedMessageTypes()...)
			}

			traversed[file.GetFullyQualifiedName()] = true
		}
	}

	return matches, nil
}

func getMessageExtension(md *desc.MessageDescriptor, field int) (interface{}, error) {
	opts := md.GetMessageOptions()
	if opts == nil {
		return nil, nil
	}

	extensions, err := proto.ExtensionDescs(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to examine message extensions descriptions: %w", err)
	}

	var extDesc *proto.ExtensionDesc
	for _, ext := range extensions {
		if ext != nil && ext.Field == int32(field) {
			extDesc = ext
			break
		}
	}
	if extDesc == nil {
		return nil, nil
	}

	ext, err := proto.GetExtension(opts, extDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to examine message extension: %w", err)
	}

	return ext, nil
}

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
