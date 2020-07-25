package main

import (
	"fmt"

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
