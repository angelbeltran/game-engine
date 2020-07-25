package main

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
)

func getResponseDescriptor(files []*desc.FileDescriptor) (*desc.MessageDescriptor, error) {
	responseMessages, err := getMessagesWithExtension(files, isActionServiceResponseExtensionFieldNumber)
	if err != nil {
		return nil, err
	}

	if len(responseMessages) > 1 {
		return nil, fmt.Errorf("multiple response messages found")
	}

	if len(responseMessages) == 0 {
		return nil, fmt.Errorf("no response message found")
	}

	return responseMessages[0], nil
}

func validateResponseMessage(state, response *desc.MessageDescriptor) error {
	var (
		sd *desc.FieldDescriptor
		ed *desc.FieldDescriptor
	)

	for _, f := range response.GetFields() {
		switch f.GetName() {
		case responseStateFieldName:
			sd = f
		case responseErrorFieldName:
			ed = f
		}
	}

	if sd == nil {
		return fmt.Errorf("no state field defined")
	}
	if ed == nil {
		return fmt.Errorf("no error field defined")
	}

	sdt := sd.GetMessageType()
	if sdt == nil || sdt.GetFullyQualifiedName() != state.GetFullyQualifiedName() {
		return fmt.Errorf("the '%s' field must be the state message", responseStateFieldName)
	}

	edt := ed.GetMessageType()
	if edt == nil || edt.GetFullyQualifiedName() != protoPackageName+"."+errorTypeName {
		return fmt.Errorf("the '%s' field must be the error message", responseErrorFieldName)
	}

	return nil
}
