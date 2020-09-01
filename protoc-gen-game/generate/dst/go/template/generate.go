// Generated by protoc-gen-game/generation. DO NOT EDIT.
package template

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/jhump/protoreflect/desc"

	pb "github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
)

//
// ----- Entrypoint -----
//

func GenerateService(w io.Writer, opts TemplateParams) error {

	// Parse template.

	tmpl, err := template.New("main").Funcs(template.FuncMap{
		"camelCase": names.CamelCase,
		"goTypeOfField": names.GoTypeOfField,
		"goTypeForMessage": names.GoTypeForMessage,
		"goTypeForMessageWithoutTargetPackageName": func(md *desc.MessageDescriptor) string {
			s := names.GoTypeForMessage(md).String()
			parts := strings.Split(s, ".")

			if len(parts) == 0 {
				return s
			}
			if parts[0] == opts.Package {
				return parts[1]
			}

			return s
		},
		"goTypeForMessageWithCamelCasePackage": func(md *desc.MessageDescriptor) string {
			var name string

			for _, part := range strings.Split(names.GoTypeForMessage(md).String(), ".") {
				name += names.CamelCase(part)
			}

			return name
		},
		"failNoFunctionName": failNoFunctionName,
		"failUndefinedEffect": failUndefinedEffect,
		"split": strings.Split,
		"errBadTypeName": func(typeName string) (interface{}, error) {
			return nil, fmt.Errorf("type name has more that one '.': %s", typeName)
		},
		"NewEffectParams": NewEffectParams,
		"NewMessageInitializerParams": NewMessageInitializerParams,
		"NewReferenceParams": NewReferenceParams,
		"NewValueParams": NewValueParams,
		"NewBoolValueParams": NewBoolValueParams,
		"NewIntValueParams": NewIntValueParams,
		"NewFloatValueParams": NewFloatValueParams,
		"NewStringValueParams": NewStringValueParams,
		"NewBoolToBoolFunctionParams": NewBoolToBoolFunctionParams,
		"NewBoolToIntFunctionParams": NewBoolToIntFunctionParams,
		"NewBoolToFloatFunctionParams": NewBoolToFloatFunctionParams,
		"NewBoolToStringFunctionParams": NewBoolToStringFunctionParams,
		"NewIntToBoolFunctionParams": NewIntToBoolFunctionParams,
		"NewIntToIntFunctionParams": NewIntToIntFunctionParams,
		"NewIntToFloatFunctionParams": NewIntToFloatFunctionParams,
		"NewIntToStringFunctionParams": NewIntToStringFunctionParams,
		"NewFloatToBoolFunctionParams": NewFloatToBoolFunctionParams,
		"NewFloatToIntFunctionParams": NewFloatToIntFunctionParams,
		"NewFloatToFloatFunctionParams": NewFloatToFloatFunctionParams,
		"NewFloatToStringFunctionParams": NewFloatToStringFunctionParams,
		"NewStringToBoolFunctionParams": NewStringToBoolFunctionParams,
		"NewStringToIntFunctionParams": NewStringToIntFunctionParams,
		"NewStringToFloatFunctionParams": NewStringToFloatFunctionParams,
		"NewStringToStringFunctionParams": NewStringToStringFunctionParams,
		"NewBoolAndBoolToBoolFunctionParams": NewBoolAndBoolToBoolFunctionParams,
		"NewBoolAndBoolToIntFunctionParams": NewBoolAndBoolToIntFunctionParams,
		"NewBoolAndBoolToFloatFunctionParams": NewBoolAndBoolToFloatFunctionParams,
		"NewBoolAndBoolToStringFunctionParams": NewBoolAndBoolToStringFunctionParams,
		"NewBoolAndIntToBoolFunctionParams": NewBoolAndIntToBoolFunctionParams,
		"NewBoolAndIntToIntFunctionParams": NewBoolAndIntToIntFunctionParams,
		"NewBoolAndIntToFloatFunctionParams": NewBoolAndIntToFloatFunctionParams,
		"NewBoolAndIntToStringFunctionParams": NewBoolAndIntToStringFunctionParams,
		"NewBoolAndFloatToBoolFunctionParams": NewBoolAndFloatToBoolFunctionParams,
		"NewBoolAndFloatToIntFunctionParams": NewBoolAndFloatToIntFunctionParams,
		"NewBoolAndFloatToFloatFunctionParams": NewBoolAndFloatToFloatFunctionParams,
		"NewBoolAndFloatToStringFunctionParams": NewBoolAndFloatToStringFunctionParams,
		"NewBoolAndStringToBoolFunctionParams": NewBoolAndStringToBoolFunctionParams,
		"NewBoolAndStringToIntFunctionParams": NewBoolAndStringToIntFunctionParams,
		"NewBoolAndStringToFloatFunctionParams": NewBoolAndStringToFloatFunctionParams,
		"NewBoolAndStringToStringFunctionParams": NewBoolAndStringToStringFunctionParams,
		"NewIntAndBoolToBoolFunctionParams": NewIntAndBoolToBoolFunctionParams,
		"NewIntAndBoolToIntFunctionParams": NewIntAndBoolToIntFunctionParams,
		"NewIntAndBoolToFloatFunctionParams": NewIntAndBoolToFloatFunctionParams,
		"NewIntAndBoolToStringFunctionParams": NewIntAndBoolToStringFunctionParams,
		"NewIntAndIntToBoolFunctionParams": NewIntAndIntToBoolFunctionParams,
		"NewIntAndIntToIntFunctionParams": NewIntAndIntToIntFunctionParams,
		"NewIntAndIntToFloatFunctionParams": NewIntAndIntToFloatFunctionParams,
		"NewIntAndIntToStringFunctionParams": NewIntAndIntToStringFunctionParams,
		"NewIntAndFloatToBoolFunctionParams": NewIntAndFloatToBoolFunctionParams,
		"NewIntAndFloatToIntFunctionParams": NewIntAndFloatToIntFunctionParams,
		"NewIntAndFloatToFloatFunctionParams": NewIntAndFloatToFloatFunctionParams,
		"NewIntAndFloatToStringFunctionParams": NewIntAndFloatToStringFunctionParams,
		"NewIntAndStringToBoolFunctionParams": NewIntAndStringToBoolFunctionParams,
		"NewIntAndStringToIntFunctionParams": NewIntAndStringToIntFunctionParams,
		"NewIntAndStringToFloatFunctionParams": NewIntAndStringToFloatFunctionParams,
		"NewIntAndStringToStringFunctionParams": NewIntAndStringToStringFunctionParams,
		"NewFloatAndBoolToBoolFunctionParams": NewFloatAndBoolToBoolFunctionParams,
		"NewFloatAndBoolToIntFunctionParams": NewFloatAndBoolToIntFunctionParams,
		"NewFloatAndBoolToFloatFunctionParams": NewFloatAndBoolToFloatFunctionParams,
		"NewFloatAndBoolToStringFunctionParams": NewFloatAndBoolToStringFunctionParams,
		"NewFloatAndIntToBoolFunctionParams": NewFloatAndIntToBoolFunctionParams,
		"NewFloatAndIntToIntFunctionParams": NewFloatAndIntToIntFunctionParams,
		"NewFloatAndIntToFloatFunctionParams": NewFloatAndIntToFloatFunctionParams,
		"NewFloatAndIntToStringFunctionParams": NewFloatAndIntToStringFunctionParams,
		"NewFloatAndFloatToBoolFunctionParams": NewFloatAndFloatToBoolFunctionParams,
		"NewFloatAndFloatToIntFunctionParams": NewFloatAndFloatToIntFunctionParams,
		"NewFloatAndFloatToFloatFunctionParams": NewFloatAndFloatToFloatFunctionParams,
		"NewFloatAndFloatToStringFunctionParams": NewFloatAndFloatToStringFunctionParams,
		"NewFloatAndStringToBoolFunctionParams": NewFloatAndStringToBoolFunctionParams,
		"NewFloatAndStringToIntFunctionParams": NewFloatAndStringToIntFunctionParams,
		"NewFloatAndStringToFloatFunctionParams": NewFloatAndStringToFloatFunctionParams,
		"NewFloatAndStringToStringFunctionParams": NewFloatAndStringToStringFunctionParams,
		"NewStringAndBoolToBoolFunctionParams": NewStringAndBoolToBoolFunctionParams,
		"NewStringAndBoolToIntFunctionParams": NewStringAndBoolToIntFunctionParams,
		"NewStringAndBoolToFloatFunctionParams": NewStringAndBoolToFloatFunctionParams,
		"NewStringAndBoolToStringFunctionParams": NewStringAndBoolToStringFunctionParams,
		"NewStringAndIntToBoolFunctionParams": NewStringAndIntToBoolFunctionParams,
		"NewStringAndIntToIntFunctionParams": NewStringAndIntToIntFunctionParams,
		"NewStringAndIntToFloatFunctionParams": NewStringAndIntToFloatFunctionParams,
		"NewStringAndIntToStringFunctionParams": NewStringAndIntToStringFunctionParams,
		"NewStringAndFloatToBoolFunctionParams": NewStringAndFloatToBoolFunctionParams,
		"NewStringAndFloatToIntFunctionParams": NewStringAndFloatToIntFunctionParams,
		"NewStringAndFloatToFloatFunctionParams": NewStringAndFloatToFloatFunctionParams,
		"NewStringAndFloatToStringFunctionParams": NewStringAndFloatToStringFunctionParams,
		"NewStringAndStringToBoolFunctionParams": NewStringAndStringToBoolFunctionParams,
		"NewStringAndStringToIntFunctionParams": NewStringAndStringToIntFunctionParams,
		"NewStringAndStringToFloatFunctionParams": NewStringAndStringToFloatFunctionParams,
		"NewStringAndStringToStringFunctionParams": NewStringAndStringToStringFunctionParams,
		"NewBoolsToBoolFunctionParams": NewBoolsToBoolFunctionParams,
		"NewBoolsToIntFunctionParams": NewBoolsToIntFunctionParams,
		"NewBoolsToFloatFunctionParams": NewBoolsToFloatFunctionParams,
		"NewBoolsToStringFunctionParams": NewBoolsToStringFunctionParams,
		"NewIntsToBoolFunctionParams": NewIntsToBoolFunctionParams,
		"NewIntsToIntFunctionParams": NewIntsToIntFunctionParams,
		"NewIntsToFloatFunctionParams": NewIntsToFloatFunctionParams,
		"NewIntsToStringFunctionParams": NewIntsToStringFunctionParams,
		"NewFloatsToBoolFunctionParams": NewFloatsToBoolFunctionParams,
		"NewFloatsToIntFunctionParams": NewFloatsToIntFunctionParams,
		"NewFloatsToFloatFunctionParams": NewFloatsToFloatFunctionParams,
		"NewFloatsToStringFunctionParams": NewFloatsToStringFunctionParams,
		"NewStringsToBoolFunctionParams": NewStringsToBoolFunctionParams,
		"NewStringsToIntFunctionParams": NewStringsToIntFunctionParams,
		"NewStringsToFloatFunctionParams": NewStringsToFloatFunctionParams,
		"NewStringsToStringFunctionParams": NewStringsToStringFunctionParams,
		"NewBoolValueIfParams": NewBoolValueIfParams,
		"NewIntValueIfParams": NewIntValueIfParams,
		"NewFloatValueIfParams": NewFloatValueIfParams,
		"NewStringValueIfParams": NewStringValueIfParams,
	}).Parse(serverTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	// Apply runtime parameters.

	out := bytes.NewBuffer([]byte{})

	if err := tmpl.Execute(out, opts); err != nil {
		return fmt.Errorf("failed to execute service template: %w", err)
	}

	// Format and write to file.

	b, err := ioutil.ReadAll(out)
	if err != nil {
		return fmt.Errorf("failed to read unformated templates: %w", err)
	}

	b, err = format.Source(b)
	if err != nil {
		return fmt.Errorf("failed to format generated templates: %w", err)
	}

	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("failed to write out formatted templates: %w", err)
	}

	return nil
}

// TemplateParams hold all the arguments needed for the main template.
type TemplateParams struct {
	Package             string
	Imports             []string
	Service             *desc.ServiceDescriptor
	Methods             []MethodInfo
	State               *desc.MessageDescriptor
	Response            *desc.MessageDescriptor
	StateVariable       string
	InputVariable       string
	EnumToFieldMappings map[*desc.MessageDescriptor]*desc.EnumDescriptor
	ResponseStateField  string
	ResponseErrorField  string
}

// MethodInfo is a method and action pair.
type MethodInfo struct {
	Method *desc.MethodDescriptor
	Action *pb.Action
}
