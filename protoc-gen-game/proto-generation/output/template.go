package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/goprotoc/plugins"

	pb "github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
)

func GenerateService(w io.Writer, opts TemplateParams) error {

	// Parse template.

	tmpl, err := template.New("main").Funcs(template.FuncMap{
		"camelCase": (&plugins.GoNames{}).CamelCase,
		"goTypeOfField": (&plugins.GoNames{}).GoTypeOfField,
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
	}).Parse(Template)
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
	Package                         string
	Imports                         []string
	Service                         *desc.ServiceDescriptor
	Methods                         []MethodInfo
	State                           *desc.MessageDescriptor
	Response                        *desc.MessageDescriptor
	StateVariable                   string
	InputVariable                   string
	ResponseStateField              string
	ResponseErrorField              string
}

// MethodInfo is a method and action pair.
type MethodInfo struct {
	Method *desc.MethodDescriptor
	Action *pb.Action
}

const Template = `// Generated by protoc-gen-game. DO NOT EDIT.
package {{ .Package }}

import (
{{ range $_, $import := .Imports }}	{{ printf "%q" $import }}
	"github.com/angelbeltran/game-engine/protoc-gen-game/proto-generation/go_func"
{{ end }}
)

func NewServer(port uint) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, err
	}

	srv := grpc.NewServer()
	Register{{ .Service.GetName }}Server(srv, new(gameEngine))

	return srv, lis, nil
}

type gameEngine struct {
	Unimplemented{{ .Service.GetName }}Server
}

{{- $state := .State }}
{{- $responseType := .Response.GetName }}
{{- $responseStateField := .ResponseStateField }}
{{- $responseErrorField := .ResponseErrorField }}
{{- $stateVariable := .StateVariable }}
{{- $inputVariable := .InputVariable }}
{{- range $_, $bundle := .Methods }}
func (e *gameEngine) {{ $bundle.Method.GetName }}(ctx context.Context, in *{{ $bundle.Method.GetInputType.GetName }}) (*{{ $responseType }}, error) {
	{{- with $action := $bundle.Action -}}
		state.Lock()
		defer state.Unlock()

		// Enforce the rules

		allowed := {{ template "BoolValue" NewBoolValueParams $action.Rule $inputVariable $stateVariable }}
		if !allowed {
			return &{{ $responseType }}{
				{{ $responseErrorField }}: &game_engine_pb.Error{
					Code: {{ printf "%q" $action.Error.Code }},
					Msg: {{ printf "%q" $action.Error.Msg }},
				},
			}, nil
		}

		// Apply any effects

		{{- range $effect := $action.Effect }}
			{{ template "effect" NewEffectParams $effect $inputVariable $stateVariable }}
		{{- end }}

		// Construct the response
		res := NewResponse()

		{{- range $ref := $action.Response }}
			{{ template "Reference" NewReferenceParams (printf "res.%s" $responseStateField ) $ref }} = {{ template "Reference" NewReferenceParams $stateVariable $ref }}
		{{- end }}

		return &res, nil

	{{- else }}

		return &{{ $responseType }}{
			{{ $responseErrorField }}: &game_engine_pb.Error{
				Msg: "unimplemented",
			},
		}, nil

	{{- end }}
}
{{- end }}

var state = NewGameState()

type GameState struct {
	{{ .State.GetName }}
	sync.Mutex
}

func NewGameState() GameState {
	var s GameState
	
	{{ template "initialize" NewMessageInitializerParams (printf "s.%s" .State.GetName) .Package .State }}

	return s
}

func NewResponse() {{ $responseType }} {
	var res {{ $responseType }}

	{{ template "initialize" NewMessageInitializerParams "res" .Package .Response }}

	return res
}

{{- define "effect" }}
{{/* Expects EffectParams */}}
{{- if .Effect.GetUpdate }}
	{{- $up := .Effect.GetUpdate }}
	{{ template "Reference" NewReferenceParams .StateVariable $up.State }} = 
	{{- template "Value" NewValueParams $up.Value .InputVariable .StateVariable }}
{{- else }}
	{{ failUndefinedEffect }}
{{- end }}
{{- end }}

{{- define "initialize" }}
{{/* Expects MessageInitializerParams */}}
	{{- $ident := .Identifier }}
	{{- $package := .Package }}

	{{- range $field := .MessageDescriptor.GetFields }}
		{{- if or $field.GetMessageType $field.IsMap }}
			{{- $typeName := (goTypeOfField $field).String }}

			{{- $isPointer := eq (index $typeName 0) '*' }}
			{{- if $isPointer }}
				{{- $typeName = slice $typeName 1 }}
			{{- end }}

			{{- $parts := split $typeName "." }}
			{{- $n := len $parts }}
			{{- if or (eq $n 0) (gt $n 2 )}}
				{{- errBadTypeName $typeName }}
			{{- else if eq $n 2 }}
				{{- /*$typeName = index $parts 1 */}}

				{{- if eq (index $parts 0) $package }}
					{{- /* Target package. Trim package name. */}}
					{{- $typeName = index $parts 1 }}
				{{- else }}
				{{- end }}
			{{- end }}

			{{- if $field.GetMessageType }}
				{{- $lh := printf "%s.%s"  $ident (camelCase $field.GetName) }}

				{{- if $isPointer }}
		{{- $lh }} = new({{ $typeName }})
				{{- end }}

				{{- template "initialize" NewMessageInitializerParams $lh $package $field.GetMessageType }}

			{{- else }}
				{{- $lh := printf "%s.%s"  $ident (camelCase $field.GetName) }}

				{{- if $isPointer }}
		{{- $lh }} = make({{ $typeName }})
				{{- end }}

			{{- end }}

		{{- end }}
	{{- end }}
{{- end }}



{{/* Expression Evaluations */}}

{{- define "Reference" }}
{{- /* Expects a ReferenceParams */}}
{{- .Variable }}{{ range $_, $field := .Reference.Path }}.{{ camelCase $field }}{{ end }}
{{- end }}

{{- define "Value" }}
{{- if .Value.GetBool }}{{ template "BoolValue" NewBoolValueParams .Value.GetBool .InputVariable .StateVariable }}
{{- else if .Value.GetInt }}{{ template "IntValue" NewIntValueParams .Value.GetInt .InputVariable .StateVariable }}
{{- else if .Value.GetFloat }}{{ template "FloatValue" NewFloatValueParams .Value.GetFloat .InputVariable .StateVariable }}
{{- else if .Value.GetString }}{{ template "StringValue" NewStringValueParams .Value.GetString .InputVariable .StateVariable }}
{{- end }}
{{- end }}

{{- define "BoolValue" }}
{{- /* Expects a BoolValueParams */}}
{{- if .Value.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Value.GetInput }}
{{- else if .Value.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Value.GetState }}
{{- else if .Value.GetBoolFunc }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Value.GetBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntFunc }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Value.GetIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatFunc }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Value.GetFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringFunc }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Value.GetStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolBoolFunc }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Value.GetBoolBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolIntFunc }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Value.GetBoolIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolFloatFunc }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Value.GetBoolFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolStringFunc }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Value.GetBoolStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntBoolFunc }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Value.GetIntBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntIntFunc }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Value.GetIntIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntFloatFunc }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Value.GetIntFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntStringFunc }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Value.GetIntStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatBoolFunc }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Value.GetFloatBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatIntFunc }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Value.GetFloatIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatFloatFunc }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Value.GetFloatFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatStringFunc }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Value.GetFloatStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringBoolFunc }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Value.GetStringBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringIntFunc }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Value.GetStringIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringFloatFunc }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Value.GetStringFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringStringFunc }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Value.GetStringStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIf }}{{ template "BoolValueIf" NewBoolValueParams .Value.GetIf .InputVariable .StateVariable }}
{{- else }}{{ .Value.GetConstant }}
{{- end }}
{{- end }}

{{- define "IntValue" }}
{{- /* Expects a IntValueParams */}}
{{- if .Value.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Value.GetInput }}
{{- else if .Value.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Value.GetState }}
{{- else if .Value.GetBoolFunc }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Value.GetBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntFunc }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Value.GetIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatFunc }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Value.GetFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringFunc }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Value.GetStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolBoolFunc }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Value.GetBoolBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolIntFunc }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Value.GetBoolIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolFloatFunc }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Value.GetBoolFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolStringFunc }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Value.GetBoolStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntBoolFunc }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Value.GetIntBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntIntFunc }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Value.GetIntIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntFloatFunc }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Value.GetIntFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntStringFunc }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Value.GetIntStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatBoolFunc }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Value.GetFloatBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatIntFunc }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Value.GetFloatIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatFloatFunc }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Value.GetFloatFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatStringFunc }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Value.GetFloatStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringBoolFunc }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Value.GetStringBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringIntFunc }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Value.GetStringIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringFloatFunc }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Value.GetStringFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringStringFunc }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Value.GetStringStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIf }}{{ template "IntValueIf" NewIntValueParams .Value.GetIf .InputVariable .StateVariable }}
{{- else }}{{ .Value.GetConstant }}
{{- end }}
{{- end }}

{{- define "FloatValue" }}
{{- /* Expects a FloatValueParams */}}
{{- if .Value.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Value.GetInput }}
{{- else if .Value.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Value.GetState }}
{{- else if .Value.GetBoolFunc }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Value.GetBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntFunc }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Value.GetIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatFunc }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Value.GetFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringFunc }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Value.GetStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolBoolFunc }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Value.GetBoolBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolIntFunc }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Value.GetBoolIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolFloatFunc }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Value.GetBoolFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolStringFunc }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Value.GetBoolStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntBoolFunc }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Value.GetIntBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntIntFunc }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Value.GetIntIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntFloatFunc }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Value.GetIntFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntStringFunc }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Value.GetIntStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatBoolFunc }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Value.GetFloatBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatIntFunc }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Value.GetFloatIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatFloatFunc }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Value.GetFloatFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatStringFunc }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Value.GetFloatStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringBoolFunc }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Value.GetStringBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringIntFunc }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Value.GetStringIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringFloatFunc }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Value.GetStringFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringStringFunc }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Value.GetStringStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIf }}{{ template "FloatValueIf" NewFloatValueParams .Value.GetIf .InputVariable .StateVariable }}
{{- else }}{{ .Value.GetConstant }}
{{- end }}
{{- end }}

{{- define "StringValue" }}
{{- /* Expects a StringValueParams */}}
{{- if .Value.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Value.GetInput }}
{{- else if .Value.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Value.GetState }}
{{- else if .Value.GetBoolFunc }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Value.GetBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntFunc }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Value.GetIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatFunc }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Value.GetFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringFunc }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Value.GetStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolBoolFunc }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Value.GetBoolBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolIntFunc }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Value.GetBoolIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolFloatFunc }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Value.GetBoolFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetBoolStringFunc }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Value.GetBoolStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntBoolFunc }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Value.GetIntBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntIntFunc }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Value.GetIntIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntFloatFunc }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Value.GetIntFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIntStringFunc }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Value.GetIntStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatBoolFunc }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Value.GetFloatBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatIntFunc }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Value.GetFloatIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatFloatFunc }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Value.GetFloatFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetFloatStringFunc }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Value.GetFloatStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringBoolFunc }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Value.GetStringBoolFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringIntFunc }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Value.GetStringIntFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringFloatFunc }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Value.GetStringFloatFunc .InputVariable .StateVariable }}
{{- else if .Value.GetStringStringFunc }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Value.GetStringStringFunc .InputVariable .StateVariable }}
{{- else if .Value.GetIf }}{{ template "StringValueIf" NewStringValueParams .Value.GetIf .InputVariable .StateVariable }}
{{- else }}{{ .Value.GetConstant }}
{{- end }}
{{- end }}


{{- /* Unary Function Templates */}}


{{- define "BoolToBoolFunction" -}}
{{- /* Expects a BoolToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolToBoolFunction" }}{{ end }}go_func.BoolToBool_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "BoolValueIf" NewBoolValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "BoolToIntFunction" -}}
{{- /* Expects a BoolToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolToIntFunction" }}{{ end }}go_func.BoolToInt_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "BoolValueIf" NewBoolValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "BoolToFloatFunction" -}}
{{- /* Expects a BoolToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolToFloatFunction" }}{{ end }}go_func.BoolToFloat_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "BoolValueIf" NewBoolValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "BoolToStringFunction" -}}
{{- /* Expects a BoolToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolToStringFunction" }}{{ end }}go_func.BoolToString_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "BoolValueIf" NewBoolValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "IntToBoolFunction" -}}
{{- /* Expects a IntToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntToBoolFunction" }}{{ end }}go_func.IntToBool_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "IntValueIf" NewIntValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "IntToIntFunction" -}}
{{- /* Expects a IntToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntToIntFunction" }}{{ end }}go_func.IntToInt_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "IntValueIf" NewIntValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "IntToFloatFunction" -}}
{{- /* Expects a IntToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntToFloatFunction" }}{{ end }}go_func.IntToFloat_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "IntValueIf" NewIntValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "IntToStringFunction" -}}
{{- /* Expects a IntToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntToStringFunction" }}{{ end }}go_func.IntToString_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "IntValueIf" NewIntValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "FloatToBoolFunction" -}}
{{- /* Expects a FloatToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatToBoolFunction" }}{{ end }}go_func.FloatToBool_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "FloatValueIf" NewFloatValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "FloatToIntFunction" -}}
{{- /* Expects a FloatToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatToIntFunction" }}{{ end }}go_func.FloatToInt_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "FloatValueIf" NewFloatValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "FloatToFloatFunction" -}}
{{- /* Expects a FloatToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatToFloatFunction" }}{{ end }}go_func.FloatToFloat_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "FloatValueIf" NewFloatValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "FloatToStringFunction" -}}
{{- /* Expects a FloatToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatToStringFunction" }}{{ end }}go_func.FloatToString_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "FloatValueIf" NewFloatValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "StringToBoolFunction" -}}
{{- /* Expects a StringToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringToBoolFunction" }}{{ end }}go_func.StringToBool_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "StringValueIf" NewStringValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "StringToIntFunction" -}}
{{- /* Expects a StringToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringToIntFunction" }}{{ end }}go_func.StringToInt_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "StringValueIf" NewStringValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "StringToFloatFunction" -}}
{{- /* Expects a StringToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringToFloatFunction" }}{{ end }}go_func.StringToFloat_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "StringValueIf" NewStringValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}

{{- define "StringToStringFunction" -}}
{{- /* Expects a StringToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringToStringFunction" }}{{ end }}go_func.StringToString_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetInput }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput }}
	{{- else if .Function.GetState }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState }}
	{{- else if .Function.GetBoolFunc }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc .InputVariable .StateVariable }}

	{{- else if .If }}{{ template "StringValueIf" NewStringValueParams .If .InputVariable .StateVariable }}

	{{- end }},
))
{{- end }}



{{/* Binary Function Templates */}}


{{- define "BoolAndBoolToBoolFunction" -}}
{{- /* Expects a BoolAndBoolToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndBoolToBoolFunction" }}{{ end }}
go_func.BoolAndBoolToBool_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndBoolToIntFunction" -}}
{{- /* Expects a BoolAndBoolToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndBoolToIntFunction" }}{{ end }}
go_func.BoolAndBoolToInt_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndBoolToFloatFunction" -}}
{{- /* Expects a BoolAndBoolToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndBoolToFloatFunction" }}{{ end }}
go_func.BoolAndBoolToFloat_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndBoolToStringFunction" -}}
{{- /* Expects a BoolAndBoolToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndBoolToStringFunction" }}{{ end }}
go_func.BoolAndBoolToString_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "BoolAndIntToBoolFunction" -}}
{{- /* Expects a BoolAndIntToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndIntToBoolFunction" }}{{ end }}
go_func.BoolAndIntToBool_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndIntToIntFunction" -}}
{{- /* Expects a BoolAndIntToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndIntToIntFunction" }}{{ end }}
go_func.BoolAndIntToInt_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndIntToFloatFunction" -}}
{{- /* Expects a BoolAndIntToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndIntToFloatFunction" }}{{ end }}
go_func.BoolAndIntToFloat_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndIntToStringFunction" -}}
{{- /* Expects a BoolAndIntToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndIntToStringFunction" }}{{ end }}
go_func.BoolAndIntToString_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "BoolAndFloatToBoolFunction" -}}
{{- /* Expects a BoolAndFloatToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndFloatToBoolFunction" }}{{ end }}
go_func.BoolAndFloatToBool_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndFloatToIntFunction" -}}
{{- /* Expects a BoolAndFloatToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndFloatToIntFunction" }}{{ end }}
go_func.BoolAndFloatToInt_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndFloatToFloatFunction" -}}
{{- /* Expects a BoolAndFloatToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndFloatToFloatFunction" }}{{ end }}
go_func.BoolAndFloatToFloat_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndFloatToStringFunction" -}}
{{- /* Expects a BoolAndFloatToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndFloatToStringFunction" }}{{ end }}
go_func.BoolAndFloatToString_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "BoolAndStringToBoolFunction" -}}
{{- /* Expects a BoolAndStringToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndStringToBoolFunction" }}{{ end }}
go_func.BoolAndStringToBool_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndStringToIntFunction" -}}
{{- /* Expects a BoolAndStringToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndStringToIntFunction" }}{{ end }}
go_func.BoolAndStringToInt_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndStringToFloatFunction" -}}
{{- /* Expects a BoolAndStringToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndStringToFloatFunction" }}{{ end }}
go_func.BoolAndStringToFloat_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "BoolAndStringToStringFunction" -}}
{{- /* Expects a BoolAndStringToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "BoolAndStringToStringFunction" }}{{ end }}
go_func.BoolAndStringToString_{{- camelCase .Function.Name.String }}(bool(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "IntAndBoolToBoolFunction" -}}
{{- /* Expects a IntAndBoolToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndBoolToBoolFunction" }}{{ end }}
go_func.IntAndBoolToBool_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndBoolToIntFunction" -}}
{{- /* Expects a IntAndBoolToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndBoolToIntFunction" }}{{ end }}
go_func.IntAndBoolToInt_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndBoolToFloatFunction" -}}
{{- /* Expects a IntAndBoolToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndBoolToFloatFunction" }}{{ end }}
go_func.IntAndBoolToFloat_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndBoolToStringFunction" -}}
{{- /* Expects a IntAndBoolToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndBoolToStringFunction" }}{{ end }}
go_func.IntAndBoolToString_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "IntAndIntToBoolFunction" -}}
{{- /* Expects a IntAndIntToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndIntToBoolFunction" }}{{ end }}
go_func.IntAndIntToBool_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndIntToIntFunction" -}}
{{- /* Expects a IntAndIntToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndIntToIntFunction" }}{{ end }}
go_func.IntAndIntToInt_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndIntToFloatFunction" -}}
{{- /* Expects a IntAndIntToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndIntToFloatFunction" }}{{ end }}
go_func.IntAndIntToFloat_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndIntToStringFunction" -}}
{{- /* Expects a IntAndIntToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndIntToStringFunction" }}{{ end }}
go_func.IntAndIntToString_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "IntAndFloatToBoolFunction" -}}
{{- /* Expects a IntAndFloatToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndFloatToBoolFunction" }}{{ end }}
go_func.IntAndFloatToBool_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndFloatToIntFunction" -}}
{{- /* Expects a IntAndFloatToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndFloatToIntFunction" }}{{ end }}
go_func.IntAndFloatToInt_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndFloatToFloatFunction" -}}
{{- /* Expects a IntAndFloatToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndFloatToFloatFunction" }}{{ end }}
go_func.IntAndFloatToFloat_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndFloatToStringFunction" -}}
{{- /* Expects a IntAndFloatToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndFloatToStringFunction" }}{{ end }}
go_func.IntAndFloatToString_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "IntAndStringToBoolFunction" -}}
{{- /* Expects a IntAndStringToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndStringToBoolFunction" }}{{ end }}
go_func.IntAndStringToBool_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndStringToIntFunction" -}}
{{- /* Expects a IntAndStringToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndStringToIntFunction" }}{{ end }}
go_func.IntAndStringToInt_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndStringToFloatFunction" -}}
{{- /* Expects a IntAndStringToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndStringToFloatFunction" }}{{ end }}
go_func.IntAndStringToFloat_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "IntAndStringToStringFunction" -}}
{{- /* Expects a IntAndStringToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "IntAndStringToStringFunction" }}{{ end }}
go_func.IntAndStringToString_{{- camelCase .Function.Name.String }}(int(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "FloatAndBoolToBoolFunction" -}}
{{- /* Expects a FloatAndBoolToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndBoolToBoolFunction" }}{{ end }}
go_func.FloatAndBoolToBool_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndBoolToIntFunction" -}}
{{- /* Expects a FloatAndBoolToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndBoolToIntFunction" }}{{ end }}
go_func.FloatAndBoolToInt_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndBoolToFloatFunction" -}}
{{- /* Expects a FloatAndBoolToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndBoolToFloatFunction" }}{{ end }}
go_func.FloatAndBoolToFloat_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndBoolToStringFunction" -}}
{{- /* Expects a FloatAndBoolToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndBoolToStringFunction" }}{{ end }}
go_func.FloatAndBoolToString_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "FloatAndIntToBoolFunction" -}}
{{- /* Expects a FloatAndIntToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndIntToBoolFunction" }}{{ end }}
go_func.FloatAndIntToBool_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndIntToIntFunction" -}}
{{- /* Expects a FloatAndIntToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndIntToIntFunction" }}{{ end }}
go_func.FloatAndIntToInt_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndIntToFloatFunction" -}}
{{- /* Expects a FloatAndIntToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndIntToFloatFunction" }}{{ end }}
go_func.FloatAndIntToFloat_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndIntToStringFunction" -}}
{{- /* Expects a FloatAndIntToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndIntToStringFunction" }}{{ end }}
go_func.FloatAndIntToString_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "FloatAndFloatToBoolFunction" -}}
{{- /* Expects a FloatAndFloatToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndFloatToBoolFunction" }}{{ end }}
go_func.FloatAndFloatToBool_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndFloatToIntFunction" -}}
{{- /* Expects a FloatAndFloatToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndFloatToIntFunction" }}{{ end }}
go_func.FloatAndFloatToInt_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndFloatToFloatFunction" -}}
{{- /* Expects a FloatAndFloatToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndFloatToFloatFunction" }}{{ end }}
go_func.FloatAndFloatToFloat_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndFloatToStringFunction" -}}
{{- /* Expects a FloatAndFloatToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndFloatToStringFunction" }}{{ end }}
go_func.FloatAndFloatToString_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "FloatAndStringToBoolFunction" -}}
{{- /* Expects a FloatAndStringToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndStringToBoolFunction" }}{{ end }}
go_func.FloatAndStringToBool_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndStringToIntFunction" -}}
{{- /* Expects a FloatAndStringToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndStringToIntFunction" }}{{ end }}
go_func.FloatAndStringToInt_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndStringToFloatFunction" -}}
{{- /* Expects a FloatAndStringToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndStringToFloatFunction" }}{{ end }}
go_func.FloatAndStringToFloat_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "FloatAndStringToStringFunction" -}}
{{- /* Expects a FloatAndStringToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "FloatAndStringToStringFunction" }}{{ end }}
go_func.FloatAndStringToString_{{- camelCase .Function.Name.String }}(float64(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "StringAndBoolToBoolFunction" -}}
{{- /* Expects a StringAndBoolToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndBoolToBoolFunction" }}{{ end }}
go_func.StringAndBoolToBool_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndBoolToIntFunction" -}}
{{- /* Expects a StringAndBoolToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndBoolToIntFunction" }}{{ end }}
go_func.StringAndBoolToInt_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndBoolToFloatFunction" -}}
{{- /* Expects a StringAndBoolToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndBoolToFloatFunction" }}{{ end }}
go_func.StringAndBoolToFloat_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndBoolToStringFunction" -}}
{{- /* Expects a StringAndBoolToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndBoolToStringFunction" }}{{ end }}
go_func.StringAndBoolToString_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), bool(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToBoolFunction" NewBoolToBoolFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToBoolFunction" NewIntToBoolFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToBoolFunction" NewFloatToBoolFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToBoolFunction" NewStringToBoolFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToBoolFunction" NewBoolAndBoolToBoolFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToBoolFunction" NewBoolAndIntToBoolFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToBoolFunction" NewBoolAndFloatToBoolFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToBoolFunction" NewBoolAndStringToBoolFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToBoolFunction" NewIntAndBoolToBoolFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToBoolFunction" NewIntAndIntToBoolFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToBoolFunction" NewIntAndFloatToBoolFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToBoolFunction" NewIntAndStringToBoolFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToBoolFunction" NewFloatAndBoolToBoolFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToBoolFunction" NewFloatAndIntToBoolFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToBoolFunction" NewFloatAndFloatToBoolFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToBoolFunction" NewFloatAndStringToBoolFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToBoolFunction" NewStringAndBoolToBoolFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToBoolFunction" NewStringAndIntToBoolFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToBoolFunction" NewStringAndFloatToBoolFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToBoolFunction" NewStringAndStringToBoolFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "BoolValueIf" NewBoolValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "StringAndIntToBoolFunction" -}}
{{- /* Expects a StringAndIntToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndIntToBoolFunction" }}{{ end }}
go_func.StringAndIntToBool_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndIntToIntFunction" -}}
{{- /* Expects a StringAndIntToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndIntToIntFunction" }}{{ end }}
go_func.StringAndIntToInt_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndIntToFloatFunction" -}}
{{- /* Expects a StringAndIntToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndIntToFloatFunction" }}{{ end }}
go_func.StringAndIntToFloat_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndIntToStringFunction" -}}
{{- /* Expects a StringAndIntToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndIntToStringFunction" }}{{ end }}
go_func.StringAndIntToString_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), int(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToIntFunction" NewBoolToIntFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToIntFunction" NewIntToIntFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToIntFunction" NewFloatToIntFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToIntFunction" NewStringToIntFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToIntFunction" NewBoolAndBoolToIntFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToIntFunction" NewBoolAndIntToIntFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToIntFunction" NewBoolAndFloatToIntFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToIntFunction" NewBoolAndStringToIntFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToIntFunction" NewIntAndBoolToIntFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToIntFunction" NewIntAndIntToIntFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToIntFunction" NewIntAndFloatToIntFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToIntFunction" NewIntAndStringToIntFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToIntFunction" NewFloatAndBoolToIntFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToIntFunction" NewFloatAndIntToIntFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToIntFunction" NewFloatAndFloatToIntFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToIntFunction" NewFloatAndStringToIntFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToIntFunction" NewStringAndBoolToIntFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToIntFunction" NewStringAndIntToIntFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToIntFunction" NewStringAndFloatToIntFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToIntFunction" NewStringAndStringToIntFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "IntValueIf" NewIntValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "StringAndFloatToBoolFunction" -}}
{{- /* Expects a StringAndFloatToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndFloatToBoolFunction" }}{{ end }}
go_func.StringAndFloatToBool_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndFloatToIntFunction" -}}
{{- /* Expects a StringAndFloatToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndFloatToIntFunction" }}{{ end }}
go_func.StringAndFloatToInt_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndFloatToFloatFunction" -}}
{{- /* Expects a StringAndFloatToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndFloatToFloatFunction" }}{{ end }}
go_func.StringAndFloatToFloat_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndFloatToStringFunction" -}}
{{- /* Expects a StringAndFloatToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndFloatToStringFunction" }}{{ end }}
go_func.StringAndFloatToString_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), float64(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToFloatFunction" NewBoolToFloatFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToFloatFunction" NewIntToFloatFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToFloatFunction" NewFloatToFloatFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToFloatFunction" NewStringToFloatFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToFloatFunction" NewBoolAndBoolToFloatFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToFloatFunction" NewBoolAndIntToFloatFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToFloatFunction" NewBoolAndFloatToFloatFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToFloatFunction" NewBoolAndStringToFloatFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToFloatFunction" NewIntAndBoolToFloatFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToFloatFunction" NewIntAndIntToFloatFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToFloatFunction" NewIntAndFloatToFloatFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToFloatFunction" NewIntAndStringToFloatFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToFloatFunction" NewFloatAndBoolToFloatFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToFloatFunction" NewFloatAndIntToFloatFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToFloatFunction" NewFloatAndFloatToFloatFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToFloatFunction" NewFloatAndStringToFloatFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToFloatFunction" NewStringAndBoolToFloatFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToFloatFunction" NewStringAndIntToFloatFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToFloatFunction" NewStringAndFloatToFloatFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToFloatFunction" NewStringAndStringToFloatFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "FloatValueIf" NewFloatValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "StringAndStringToBoolFunction" -}}
{{- /* Expects a StringAndStringToBoolFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndStringToBoolFunction" }}{{ end }}
go_func.StringAndStringToBool_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndStringToIntFunction" -}}
{{- /* Expects a StringAndStringToIntFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndStringToIntFunction" }}{{ end }}
go_func.StringAndStringToInt_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndStringToFloatFunction" -}}
{{- /* Expects a StringAndStringToFloatFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndStringToFloatFunction" }}{{ end }}
go_func.StringAndStringToFloat_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}

{{- define "StringAndStringToStringFunction" -}}
{{- /* Expects a StringAndStringToStringFunctionParams */}}
{{- if not .Function.Name }}{{ failNoFunctionName "StringAndStringToStringFunction" }}{{ end }}
go_func.StringAndStringToString_{{- camelCase .Function.Name.String }}(string(
	{{- if .Function.GetConstant_1 }}{{ .Function.GetConstant_1 }}
	{{- else if .Function.GetInput_1 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_1 }}
	{{- else if .Function.GetState_1 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_1 }}
	{{- else if .Function.GetBoolFunc_1 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_1 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_1 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_1 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_1 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_1 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_1 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_1 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_1 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_1 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_1 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_1 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_1 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_1 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_1 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_1 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_1 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_1 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_1 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_1 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_1 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_1 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_1 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_1 .InputVariable .StateVariable }}

	{{- end }}), string(

	{{- if .Function.GetInput_2 }}{{ template "Reference" NewReferenceParams .InputVariable .Function.GetInput_2 }}
	{{- else if .Function.GetState_2 }}{{ template "Reference" NewReferenceParams .StateVariable .Function.GetState_2 }}
	{{- else if .Function.GetBoolFunc_2 }}{{ template "BoolToStringFunction" NewBoolToStringFunctionParams .Function.GetBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFunc_2 }}{{ template "IntToStringFunction" NewIntToStringFunctionParams .Function.GetIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFunc_2 }}{{ template "FloatToStringFunction" NewFloatToStringFunctionParams .Function.GetFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFunc_2 }}{{ template "StringToStringFunction" NewStringToStringFunctionParams .Function.GetStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolBoolFunc_2 }}{{ template "BoolAndBoolToStringFunction" NewBoolAndBoolToStringFunctionParams .Function.GetBoolBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolIntFunc_2 }}{{ template "BoolAndIntToStringFunction" NewBoolAndIntToStringFunctionParams .Function.GetBoolIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolFloatFunc_2 }}{{ template "BoolAndFloatToStringFunction" NewBoolAndFloatToStringFunctionParams .Function.GetBoolFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetBoolStringFunc_2 }}{{ template "BoolAndStringToStringFunction" NewBoolAndStringToStringFunctionParams .Function.GetBoolStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntBoolFunc_2 }}{{ template "IntAndBoolToStringFunction" NewIntAndBoolToStringFunctionParams .Function.GetIntBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntIntFunc_2 }}{{ template "IntAndIntToStringFunction" NewIntAndIntToStringFunctionParams .Function.GetIntIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntFloatFunc_2 }}{{ template "IntAndFloatToStringFunction" NewIntAndFloatToStringFunctionParams .Function.GetIntFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetIntStringFunc_2 }}{{ template "IntAndStringToStringFunction" NewIntAndStringToStringFunctionParams .Function.GetIntStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatBoolFunc_2 }}{{ template "FloatAndBoolToStringFunction" NewFloatAndBoolToStringFunctionParams .Function.GetFloatBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatIntFunc_2 }}{{ template "FloatAndIntToStringFunction" NewFloatAndIntToStringFunctionParams .Function.GetFloatIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatFloatFunc_2 }}{{ template "FloatAndFloatToStringFunction" NewFloatAndFloatToStringFunctionParams .Function.GetFloatFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetFloatStringFunc_2 }}{{ template "FloatAndStringToStringFunction" NewFloatAndStringToStringFunctionParams .Function.GetFloatStringFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringBoolFunc_2 }}{{ template "StringAndBoolToStringFunction" NewStringAndBoolToStringFunctionParams .Function.GetStringBoolFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringIntFunc_2 }}{{ template "StringAndIntToStringFunction" NewStringAndIntToStringFunctionParams .Function.GetStringIntFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringFloatFunc_2 }}{{ template "StringAndFloatToStringFunction" NewStringAndFloatToStringFunctionParams .Function.GetStringFloatFunc_2 .InputVariable .StateVariable }}
	{{- else if .Function.GetStringStringFunc_2 }}{{ template "StringAndStringToStringFunction" NewStringAndStringToStringFunctionParams .Function.GetStringStringFunc_2 .InputVariable .StateVariable }}

	{{- else if .Function.GetIf_2 }}{{ template "StringValueIf" NewStringValueParams .Function.GetIf_2 .InputVariable .StateVariable }}

	{{- end }}),
)
{{- end }}


{{- define "BoolValueIf" -}}
{{- /* Expects a *pb.BoolValueIf */}}
func() bool {
	if {{ template "BoolValue" NewBoolValueParams .If.Predicate .InputVariable .StateVariable }} {
		return {{ template "BoolValue" NewBoolValueParams .If.Then .InputVariable .StateVariable }}
	}

	return {{ template "BoolValue" NewBoolValueParams .If.Else .InputVariable .StateVariable }}
}()
{{- end }}

{{- define "IntValueIf" -}}
{{- /* Expects a *pb.IntValueIf */}}
func() int {
	if {{ template "BoolValue" NewBoolValueParams .If.Predicate .InputVariable .StateVariable }} {
		return {{ template "IntValue" NewIntValueParams .If.Then .InputVariable .StateVariable }}
	}

	return {{ template "IntValue" NewIntValueParams .If.Else .InputVariable .StateVariable }}
}()
{{- end }}

{{- define "FloatValueIf" -}}
{{- /* Expects a *pb.FloatValueIf */}}
func() float {
	if {{ template "BoolValue" NewBoolValueParams .If.Predicate .InputVariable .StateVariable }} {
		return {{ template "FloatValue" NewFloatValueParams .If.Then .InputVariable .StateVariable }}
	}

	return {{ template "FloatValue" NewFloatValueParams .If.Else .InputVariable .StateVariable }}
}()
{{- end }}

{{- define "StringValueIf" -}}
{{- /* Expects a *pb.StringValueIf */}}
func() string {
	if {{ template "BoolValue" NewBoolValueParams .If.Predicate .InputVariable .StateVariable }} {
		return {{ template "StringValue" NewStringValueParams .If.Then .InputVariable .StateVariable }}
	}

	return {{ template "StringValue" NewStringValueParams .If.Else .InputVariable .StateVariable }}
}()
{{- end }}

`

// References

type ReferenceParams struct {
	Reference *pb.Reference
	Variable string
}

func NewReferenceParams(variable string, reference *pb.Reference) ReferenceParams {
	return ReferenceParams{
		Variable: variable,
		Reference: reference,
	}
}

// Values

type ValueParams struct {
	Value *pb.Value
	InputVariable string
	StateVariable string
}

func NewValueParams(value *pb.Value, input, state string) ValueParams {
	return ValueParams{
		Value: value,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolValueParams struct {
	Value *pb.BoolValue
	InputVariable string
	StateVariable string
}

func NewBoolValueParams(value *pb.BoolValue, input, state string) BoolValueParams {
	return BoolValueParams{
		Value: value,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntValueParams struct {
	Value *pb.IntValue
	InputVariable string
	StateVariable string
}

func NewIntValueParams(value *pb.IntValue, input, state string) IntValueParams {
	return IntValueParams{
		Value: value,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatValueParams struct {
	Value *pb.FloatValue
	InputVariable string
	StateVariable string
}

func NewFloatValueParams(value *pb.FloatValue, input, state string) FloatValueParams {
	return FloatValueParams{
		Value: value,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringValueParams struct {
	Value *pb.StringValue
	InputVariable string
	StateVariable string
}

func NewStringValueParams(value *pb.StringValue, input, state string) StringValueParams {
	return StringValueParams{
		Value: value,
		InputVariable: input,
		StateVariable: state,
	}
}


// Unary Functions

type BoolToBoolFunctionParams struct {
	Function *pb.BoolToBoolFunction
	InputVariable string
	StateVariable string
}

func NewBoolToBoolFunctionParams(function *pb.BoolToBoolFunction, input, state string) BoolToBoolFunctionParams {
	return BoolToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolToIntFunctionParams struct {
	Function *pb.BoolToIntFunction
	InputVariable string
	StateVariable string
}

func NewBoolToIntFunctionParams(function *pb.BoolToIntFunction, input, state string) BoolToIntFunctionParams {
	return BoolToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolToFloatFunctionParams struct {
	Function *pb.BoolToFloatFunction
	InputVariable string
	StateVariable string
}

func NewBoolToFloatFunctionParams(function *pb.BoolToFloatFunction, input, state string) BoolToFloatFunctionParams {
	return BoolToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolToStringFunctionParams struct {
	Function *pb.BoolToStringFunction
	InputVariable string
	StateVariable string
}

func NewBoolToStringFunctionParams(function *pb.BoolToStringFunction, input, state string) BoolToStringFunctionParams {
	return BoolToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntToBoolFunctionParams struct {
	Function *pb.IntToBoolFunction
	InputVariable string
	StateVariable string
}

func NewIntToBoolFunctionParams(function *pb.IntToBoolFunction, input, state string) IntToBoolFunctionParams {
	return IntToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntToIntFunctionParams struct {
	Function *pb.IntToIntFunction
	InputVariable string
	StateVariable string
}

func NewIntToIntFunctionParams(function *pb.IntToIntFunction, input, state string) IntToIntFunctionParams {
	return IntToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntToFloatFunctionParams struct {
	Function *pb.IntToFloatFunction
	InputVariable string
	StateVariable string
}

func NewIntToFloatFunctionParams(function *pb.IntToFloatFunction, input, state string) IntToFloatFunctionParams {
	return IntToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntToStringFunctionParams struct {
	Function *pb.IntToStringFunction
	InputVariable string
	StateVariable string
}

func NewIntToStringFunctionParams(function *pb.IntToStringFunction, input, state string) IntToStringFunctionParams {
	return IntToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatToBoolFunctionParams struct {
	Function *pb.FloatToBoolFunction
	InputVariable string
	StateVariable string
}

func NewFloatToBoolFunctionParams(function *pb.FloatToBoolFunction, input, state string) FloatToBoolFunctionParams {
	return FloatToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatToIntFunctionParams struct {
	Function *pb.FloatToIntFunction
	InputVariable string
	StateVariable string
}

func NewFloatToIntFunctionParams(function *pb.FloatToIntFunction, input, state string) FloatToIntFunctionParams {
	return FloatToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatToFloatFunctionParams struct {
	Function *pb.FloatToFloatFunction
	InputVariable string
	StateVariable string
}

func NewFloatToFloatFunctionParams(function *pb.FloatToFloatFunction, input, state string) FloatToFloatFunctionParams {
	return FloatToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatToStringFunctionParams struct {
	Function *pb.FloatToStringFunction
	InputVariable string
	StateVariable string
}

func NewFloatToStringFunctionParams(function *pb.FloatToStringFunction, input, state string) FloatToStringFunctionParams {
	return FloatToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringToBoolFunctionParams struct {
	Function *pb.StringToBoolFunction
	InputVariable string
	StateVariable string
}

func NewStringToBoolFunctionParams(function *pb.StringToBoolFunction, input, state string) StringToBoolFunctionParams {
	return StringToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringToIntFunctionParams struct {
	Function *pb.StringToIntFunction
	InputVariable string
	StateVariable string
}

func NewStringToIntFunctionParams(function *pb.StringToIntFunction, input, state string) StringToIntFunctionParams {
	return StringToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringToFloatFunctionParams struct {
	Function *pb.StringToFloatFunction
	InputVariable string
	StateVariable string
}

func NewStringToFloatFunctionParams(function *pb.StringToFloatFunction, input, state string) StringToFloatFunctionParams {
	return StringToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringToStringFunctionParams struct {
	Function *pb.StringToStringFunction
	InputVariable string
	StateVariable string
}

func NewStringToStringFunctionParams(function *pb.StringToStringFunction, input, state string) StringToStringFunctionParams {
	return StringToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

// Binary Functions

type BoolAndBoolToBoolFunctionParams struct {
	Function *pb.BoolAndBoolToBoolFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndBoolToBoolFunctionParams(function *pb.BoolAndBoolToBoolFunction, input, state string) BoolAndBoolToBoolFunctionParams {
	return BoolAndBoolToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndBoolToIntFunctionParams struct {
	Function *pb.BoolAndBoolToIntFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndBoolToIntFunctionParams(function *pb.BoolAndBoolToIntFunction, input, state string) BoolAndBoolToIntFunctionParams {
	return BoolAndBoolToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndBoolToFloatFunctionParams struct {
	Function *pb.BoolAndBoolToFloatFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndBoolToFloatFunctionParams(function *pb.BoolAndBoolToFloatFunction, input, state string) BoolAndBoolToFloatFunctionParams {
	return BoolAndBoolToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndBoolToStringFunctionParams struct {
	Function *pb.BoolAndBoolToStringFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndBoolToStringFunctionParams(function *pb.BoolAndBoolToStringFunction, input, state string) BoolAndBoolToStringFunctionParams {
	return BoolAndBoolToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndIntToBoolFunctionParams struct {
	Function *pb.BoolAndIntToBoolFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndIntToBoolFunctionParams(function *pb.BoolAndIntToBoolFunction, input, state string) BoolAndIntToBoolFunctionParams {
	return BoolAndIntToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndIntToIntFunctionParams struct {
	Function *pb.BoolAndIntToIntFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndIntToIntFunctionParams(function *pb.BoolAndIntToIntFunction, input, state string) BoolAndIntToIntFunctionParams {
	return BoolAndIntToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndIntToFloatFunctionParams struct {
	Function *pb.BoolAndIntToFloatFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndIntToFloatFunctionParams(function *pb.BoolAndIntToFloatFunction, input, state string) BoolAndIntToFloatFunctionParams {
	return BoolAndIntToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndIntToStringFunctionParams struct {
	Function *pb.BoolAndIntToStringFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndIntToStringFunctionParams(function *pb.BoolAndIntToStringFunction, input, state string) BoolAndIntToStringFunctionParams {
	return BoolAndIntToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndFloatToBoolFunctionParams struct {
	Function *pb.BoolAndFloatToBoolFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndFloatToBoolFunctionParams(function *pb.BoolAndFloatToBoolFunction, input, state string) BoolAndFloatToBoolFunctionParams {
	return BoolAndFloatToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndFloatToIntFunctionParams struct {
	Function *pb.BoolAndFloatToIntFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndFloatToIntFunctionParams(function *pb.BoolAndFloatToIntFunction, input, state string) BoolAndFloatToIntFunctionParams {
	return BoolAndFloatToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndFloatToFloatFunctionParams struct {
	Function *pb.BoolAndFloatToFloatFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndFloatToFloatFunctionParams(function *pb.BoolAndFloatToFloatFunction, input, state string) BoolAndFloatToFloatFunctionParams {
	return BoolAndFloatToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndFloatToStringFunctionParams struct {
	Function *pb.BoolAndFloatToStringFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndFloatToStringFunctionParams(function *pb.BoolAndFloatToStringFunction, input, state string) BoolAndFloatToStringFunctionParams {
	return BoolAndFloatToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndStringToBoolFunctionParams struct {
	Function *pb.BoolAndStringToBoolFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndStringToBoolFunctionParams(function *pb.BoolAndStringToBoolFunction, input, state string) BoolAndStringToBoolFunctionParams {
	return BoolAndStringToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndStringToIntFunctionParams struct {
	Function *pb.BoolAndStringToIntFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndStringToIntFunctionParams(function *pb.BoolAndStringToIntFunction, input, state string) BoolAndStringToIntFunctionParams {
	return BoolAndStringToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndStringToFloatFunctionParams struct {
	Function *pb.BoolAndStringToFloatFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndStringToFloatFunctionParams(function *pb.BoolAndStringToFloatFunction, input, state string) BoolAndStringToFloatFunctionParams {
	return BoolAndStringToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type BoolAndStringToStringFunctionParams struct {
	Function *pb.BoolAndStringToStringFunction
	InputVariable string
	StateVariable string
}

func NewBoolAndStringToStringFunctionParams(function *pb.BoolAndStringToStringFunction, input, state string) BoolAndStringToStringFunctionParams {
	return BoolAndStringToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndBoolToBoolFunctionParams struct {
	Function *pb.IntAndBoolToBoolFunction
	InputVariable string
	StateVariable string
}

func NewIntAndBoolToBoolFunctionParams(function *pb.IntAndBoolToBoolFunction, input, state string) IntAndBoolToBoolFunctionParams {
	return IntAndBoolToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndBoolToIntFunctionParams struct {
	Function *pb.IntAndBoolToIntFunction
	InputVariable string
	StateVariable string
}

func NewIntAndBoolToIntFunctionParams(function *pb.IntAndBoolToIntFunction, input, state string) IntAndBoolToIntFunctionParams {
	return IntAndBoolToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndBoolToFloatFunctionParams struct {
	Function *pb.IntAndBoolToFloatFunction
	InputVariable string
	StateVariable string
}

func NewIntAndBoolToFloatFunctionParams(function *pb.IntAndBoolToFloatFunction, input, state string) IntAndBoolToFloatFunctionParams {
	return IntAndBoolToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndBoolToStringFunctionParams struct {
	Function *pb.IntAndBoolToStringFunction
	InputVariable string
	StateVariable string
}

func NewIntAndBoolToStringFunctionParams(function *pb.IntAndBoolToStringFunction, input, state string) IntAndBoolToStringFunctionParams {
	return IntAndBoolToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndIntToBoolFunctionParams struct {
	Function *pb.IntAndIntToBoolFunction
	InputVariable string
	StateVariable string
}

func NewIntAndIntToBoolFunctionParams(function *pb.IntAndIntToBoolFunction, input, state string) IntAndIntToBoolFunctionParams {
	return IntAndIntToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndIntToIntFunctionParams struct {
	Function *pb.IntAndIntToIntFunction
	InputVariable string
	StateVariable string
}

func NewIntAndIntToIntFunctionParams(function *pb.IntAndIntToIntFunction, input, state string) IntAndIntToIntFunctionParams {
	return IntAndIntToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndIntToFloatFunctionParams struct {
	Function *pb.IntAndIntToFloatFunction
	InputVariable string
	StateVariable string
}

func NewIntAndIntToFloatFunctionParams(function *pb.IntAndIntToFloatFunction, input, state string) IntAndIntToFloatFunctionParams {
	return IntAndIntToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndIntToStringFunctionParams struct {
	Function *pb.IntAndIntToStringFunction
	InputVariable string
	StateVariable string
}

func NewIntAndIntToStringFunctionParams(function *pb.IntAndIntToStringFunction, input, state string) IntAndIntToStringFunctionParams {
	return IntAndIntToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndFloatToBoolFunctionParams struct {
	Function *pb.IntAndFloatToBoolFunction
	InputVariable string
	StateVariable string
}

func NewIntAndFloatToBoolFunctionParams(function *pb.IntAndFloatToBoolFunction, input, state string) IntAndFloatToBoolFunctionParams {
	return IntAndFloatToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndFloatToIntFunctionParams struct {
	Function *pb.IntAndFloatToIntFunction
	InputVariable string
	StateVariable string
}

func NewIntAndFloatToIntFunctionParams(function *pb.IntAndFloatToIntFunction, input, state string) IntAndFloatToIntFunctionParams {
	return IntAndFloatToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndFloatToFloatFunctionParams struct {
	Function *pb.IntAndFloatToFloatFunction
	InputVariable string
	StateVariable string
}

func NewIntAndFloatToFloatFunctionParams(function *pb.IntAndFloatToFloatFunction, input, state string) IntAndFloatToFloatFunctionParams {
	return IntAndFloatToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndFloatToStringFunctionParams struct {
	Function *pb.IntAndFloatToStringFunction
	InputVariable string
	StateVariable string
}

func NewIntAndFloatToStringFunctionParams(function *pb.IntAndFloatToStringFunction, input, state string) IntAndFloatToStringFunctionParams {
	return IntAndFloatToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndStringToBoolFunctionParams struct {
	Function *pb.IntAndStringToBoolFunction
	InputVariable string
	StateVariable string
}

func NewIntAndStringToBoolFunctionParams(function *pb.IntAndStringToBoolFunction, input, state string) IntAndStringToBoolFunctionParams {
	return IntAndStringToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndStringToIntFunctionParams struct {
	Function *pb.IntAndStringToIntFunction
	InputVariable string
	StateVariable string
}

func NewIntAndStringToIntFunctionParams(function *pb.IntAndStringToIntFunction, input, state string) IntAndStringToIntFunctionParams {
	return IntAndStringToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndStringToFloatFunctionParams struct {
	Function *pb.IntAndStringToFloatFunction
	InputVariable string
	StateVariable string
}

func NewIntAndStringToFloatFunctionParams(function *pb.IntAndStringToFloatFunction, input, state string) IntAndStringToFloatFunctionParams {
	return IntAndStringToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type IntAndStringToStringFunctionParams struct {
	Function *pb.IntAndStringToStringFunction
	InputVariable string
	StateVariable string
}

func NewIntAndStringToStringFunctionParams(function *pb.IntAndStringToStringFunction, input, state string) IntAndStringToStringFunctionParams {
	return IntAndStringToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndBoolToBoolFunctionParams struct {
	Function *pb.FloatAndBoolToBoolFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndBoolToBoolFunctionParams(function *pb.FloatAndBoolToBoolFunction, input, state string) FloatAndBoolToBoolFunctionParams {
	return FloatAndBoolToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndBoolToIntFunctionParams struct {
	Function *pb.FloatAndBoolToIntFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndBoolToIntFunctionParams(function *pb.FloatAndBoolToIntFunction, input, state string) FloatAndBoolToIntFunctionParams {
	return FloatAndBoolToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndBoolToFloatFunctionParams struct {
	Function *pb.FloatAndBoolToFloatFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndBoolToFloatFunctionParams(function *pb.FloatAndBoolToFloatFunction, input, state string) FloatAndBoolToFloatFunctionParams {
	return FloatAndBoolToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndBoolToStringFunctionParams struct {
	Function *pb.FloatAndBoolToStringFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndBoolToStringFunctionParams(function *pb.FloatAndBoolToStringFunction, input, state string) FloatAndBoolToStringFunctionParams {
	return FloatAndBoolToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndIntToBoolFunctionParams struct {
	Function *pb.FloatAndIntToBoolFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndIntToBoolFunctionParams(function *pb.FloatAndIntToBoolFunction, input, state string) FloatAndIntToBoolFunctionParams {
	return FloatAndIntToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndIntToIntFunctionParams struct {
	Function *pb.FloatAndIntToIntFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndIntToIntFunctionParams(function *pb.FloatAndIntToIntFunction, input, state string) FloatAndIntToIntFunctionParams {
	return FloatAndIntToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndIntToFloatFunctionParams struct {
	Function *pb.FloatAndIntToFloatFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndIntToFloatFunctionParams(function *pb.FloatAndIntToFloatFunction, input, state string) FloatAndIntToFloatFunctionParams {
	return FloatAndIntToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndIntToStringFunctionParams struct {
	Function *pb.FloatAndIntToStringFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndIntToStringFunctionParams(function *pb.FloatAndIntToStringFunction, input, state string) FloatAndIntToStringFunctionParams {
	return FloatAndIntToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndFloatToBoolFunctionParams struct {
	Function *pb.FloatAndFloatToBoolFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndFloatToBoolFunctionParams(function *pb.FloatAndFloatToBoolFunction, input, state string) FloatAndFloatToBoolFunctionParams {
	return FloatAndFloatToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndFloatToIntFunctionParams struct {
	Function *pb.FloatAndFloatToIntFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndFloatToIntFunctionParams(function *pb.FloatAndFloatToIntFunction, input, state string) FloatAndFloatToIntFunctionParams {
	return FloatAndFloatToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndFloatToFloatFunctionParams struct {
	Function *pb.FloatAndFloatToFloatFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndFloatToFloatFunctionParams(function *pb.FloatAndFloatToFloatFunction, input, state string) FloatAndFloatToFloatFunctionParams {
	return FloatAndFloatToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndFloatToStringFunctionParams struct {
	Function *pb.FloatAndFloatToStringFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndFloatToStringFunctionParams(function *pb.FloatAndFloatToStringFunction, input, state string) FloatAndFloatToStringFunctionParams {
	return FloatAndFloatToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndStringToBoolFunctionParams struct {
	Function *pb.FloatAndStringToBoolFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndStringToBoolFunctionParams(function *pb.FloatAndStringToBoolFunction, input, state string) FloatAndStringToBoolFunctionParams {
	return FloatAndStringToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndStringToIntFunctionParams struct {
	Function *pb.FloatAndStringToIntFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndStringToIntFunctionParams(function *pb.FloatAndStringToIntFunction, input, state string) FloatAndStringToIntFunctionParams {
	return FloatAndStringToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndStringToFloatFunctionParams struct {
	Function *pb.FloatAndStringToFloatFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndStringToFloatFunctionParams(function *pb.FloatAndStringToFloatFunction, input, state string) FloatAndStringToFloatFunctionParams {
	return FloatAndStringToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type FloatAndStringToStringFunctionParams struct {
	Function *pb.FloatAndStringToStringFunction
	InputVariable string
	StateVariable string
}

func NewFloatAndStringToStringFunctionParams(function *pb.FloatAndStringToStringFunction, input, state string) FloatAndStringToStringFunctionParams {
	return FloatAndStringToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndBoolToBoolFunctionParams struct {
	Function *pb.StringAndBoolToBoolFunction
	InputVariable string
	StateVariable string
}

func NewStringAndBoolToBoolFunctionParams(function *pb.StringAndBoolToBoolFunction, input, state string) StringAndBoolToBoolFunctionParams {
	return StringAndBoolToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndBoolToIntFunctionParams struct {
	Function *pb.StringAndBoolToIntFunction
	InputVariable string
	StateVariable string
}

func NewStringAndBoolToIntFunctionParams(function *pb.StringAndBoolToIntFunction, input, state string) StringAndBoolToIntFunctionParams {
	return StringAndBoolToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndBoolToFloatFunctionParams struct {
	Function *pb.StringAndBoolToFloatFunction
	InputVariable string
	StateVariable string
}

func NewStringAndBoolToFloatFunctionParams(function *pb.StringAndBoolToFloatFunction, input, state string) StringAndBoolToFloatFunctionParams {
	return StringAndBoolToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndBoolToStringFunctionParams struct {
	Function *pb.StringAndBoolToStringFunction
	InputVariable string
	StateVariable string
}

func NewStringAndBoolToStringFunctionParams(function *pb.StringAndBoolToStringFunction, input, state string) StringAndBoolToStringFunctionParams {
	return StringAndBoolToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndIntToBoolFunctionParams struct {
	Function *pb.StringAndIntToBoolFunction
	InputVariable string
	StateVariable string
}

func NewStringAndIntToBoolFunctionParams(function *pb.StringAndIntToBoolFunction, input, state string) StringAndIntToBoolFunctionParams {
	return StringAndIntToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndIntToIntFunctionParams struct {
	Function *pb.StringAndIntToIntFunction
	InputVariable string
	StateVariable string
}

func NewStringAndIntToIntFunctionParams(function *pb.StringAndIntToIntFunction, input, state string) StringAndIntToIntFunctionParams {
	return StringAndIntToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndIntToFloatFunctionParams struct {
	Function *pb.StringAndIntToFloatFunction
	InputVariable string
	StateVariable string
}

func NewStringAndIntToFloatFunctionParams(function *pb.StringAndIntToFloatFunction, input, state string) StringAndIntToFloatFunctionParams {
	return StringAndIntToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndIntToStringFunctionParams struct {
	Function *pb.StringAndIntToStringFunction
	InputVariable string
	StateVariable string
}

func NewStringAndIntToStringFunctionParams(function *pb.StringAndIntToStringFunction, input, state string) StringAndIntToStringFunctionParams {
	return StringAndIntToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndFloatToBoolFunctionParams struct {
	Function *pb.StringAndFloatToBoolFunction
	InputVariable string
	StateVariable string
}

func NewStringAndFloatToBoolFunctionParams(function *pb.StringAndFloatToBoolFunction, input, state string) StringAndFloatToBoolFunctionParams {
	return StringAndFloatToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndFloatToIntFunctionParams struct {
	Function *pb.StringAndFloatToIntFunction
	InputVariable string
	StateVariable string
}

func NewStringAndFloatToIntFunctionParams(function *pb.StringAndFloatToIntFunction, input, state string) StringAndFloatToIntFunctionParams {
	return StringAndFloatToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndFloatToFloatFunctionParams struct {
	Function *pb.StringAndFloatToFloatFunction
	InputVariable string
	StateVariable string
}

func NewStringAndFloatToFloatFunctionParams(function *pb.StringAndFloatToFloatFunction, input, state string) StringAndFloatToFloatFunctionParams {
	return StringAndFloatToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndFloatToStringFunctionParams struct {
	Function *pb.StringAndFloatToStringFunction
	InputVariable string
	StateVariable string
}

func NewStringAndFloatToStringFunctionParams(function *pb.StringAndFloatToStringFunction, input, state string) StringAndFloatToStringFunctionParams {
	return StringAndFloatToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndStringToBoolFunctionParams struct {
	Function *pb.StringAndStringToBoolFunction
	InputVariable string
	StateVariable string
}

func NewStringAndStringToBoolFunctionParams(function *pb.StringAndStringToBoolFunction, input, state string) StringAndStringToBoolFunctionParams {
	return StringAndStringToBoolFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndStringToIntFunctionParams struct {
	Function *pb.StringAndStringToIntFunction
	InputVariable string
	StateVariable string
}

func NewStringAndStringToIntFunctionParams(function *pb.StringAndStringToIntFunction, input, state string) StringAndStringToIntFunctionParams {
	return StringAndStringToIntFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndStringToFloatFunctionParams struct {
	Function *pb.StringAndStringToFloatFunction
	InputVariable string
	StateVariable string
}

func NewStringAndStringToFloatFunctionParams(function *pb.StringAndStringToFloatFunction, input, state string) StringAndStringToFloatFunctionParams {
	return StringAndStringToFloatFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type StringAndStringToStringFunctionParams struct {
	Function *pb.StringAndStringToStringFunction
	InputVariable string
	StateVariable string
}

func NewStringAndStringToStringFunctionParams(function *pb.StringAndStringToStringFunction, input, state string) StringAndStringToStringFunctionParams {
	return StringAndStringToStringFunctionParams{
		Function: function,
		InputVariable: input,
		StateVariable: state,
	}
}

type EffectParams struct {
	Effect *pb.Effect
	InputVariable string
	StateVariable string
}

func NewEffectParams(effect *pb.Effect, input, state string) EffectParams {
	return EffectParams {
		Effect: effect,
		InputVariable: input,
		StateVariable: state,
	}
}

type ResponseParams struct {
	Reference *pb.Reference
	InputVariable string
	StateVariable string
}

func NewResponseParams(reference *pb.Reference, input, state string) ResponseParams {
	return ResponseParams {
		Reference: reference,
		InputVariable: input,
		StateVariable: state,
	}
}

type MessageInitializerParams struct {
	Identifier string
	Package string
	MessageDescriptor *desc.MessageDescriptor
}

func NewMessageInitializerParams(ident, pkg string, msg *desc.MessageDescriptor) MessageInitializerParams {
	return MessageInitializerParams {
		Identifier: ident,
		Package: pkg,
		MessageDescriptor: msg,
	}
}

func failNoFunctionName(funcType string) (interface{}, error) {
	return nil, fmt.Errorf("function name missing for function type %s", funcType)
}

func failUndefinedEffect() (interface{}, error) {
	return nil, fmt.Errorf("undefined effect")
}
