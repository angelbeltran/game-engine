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
		stateField *desc.FieldDescriptor
		errorField *desc.FieldDescriptor
	)

	for _, f := range response.GetFields() {
		switch f.GetName() {
		case responseStateFieldName:
			stateField = f
		case responseErrorFieldName:
			errorField = f
		}
	}

	if stateField == nil {
		return fmt.Errorf("no state field defined")
	}
	if errorField == nil {
		return fmt.Errorf("no error field defined")
	}

	stateType := stateField.GetMessageType()
	if stateType == nil || stateType.GetFullyQualifiedName() != state.GetFullyQualifiedName() {
		return fmt.Errorf("the '%s' field must be the state message", responseStateFieldName)
	}

	errorType := errorField.GetMessageType()
	if errorType == nil || errorType.GetFullyQualifiedName() != protoPackageName+"."+errorTypeName {
		return fmt.Errorf("the '%s' field must be the error message", responseErrorFieldName)
	}

	return nil
}
