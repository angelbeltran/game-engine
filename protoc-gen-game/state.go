package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
)

func getStateDescriptor(files []*desc.FileDescriptor) (*desc.MessageDescriptor, error) {
	stateMessages, err := getMessagesWithExtension(files, isGameStateExtensionFieldNumber)
	if err != nil {
		return nil, err
	}

	if len(stateMessages) > 1 {
		return nil, fmt.Errorf("multiple state messages found")
	}

	if len(stateMessages) == 0 {
		return nil, fmt.Errorf("no state message found")
	}

	return stateMessages[0], nil
}

func getMessagesWithExtension(files []*desc.FileDescriptor, field int) ([]*desc.MessageDescriptor, error) {
	var stateMessages []*desc.MessageDescriptor

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
					stateMessages = append(stateMessages, md)
					break
				}
			}

			traversed[file.GetFullyQualifiedName()] = true
		}
	}

	return stateMessages, nil
}
