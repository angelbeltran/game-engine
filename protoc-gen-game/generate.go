package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"strings"

	"github.com/jhump/protoreflect/desc"

	pb "github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
	"github.com/angelbeltran/game-engine/protoc-gen-game/types"
)

func generateService(w io.Writer, opts serviceParameters) error {

	// Load templates.

	templates, err := parseDefinition(nil, definitions)
	if err != nil {
		return fmt.Errorf("failed to parse template definitions: %w", err)
	}

	// Apply runtime parameters.

	out := bytes.NewBuffer([]byte{})

	if err := templates["service"].Execute(out, opts); err != nil {
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

type (
	serviceParameters struct {
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
)

func (serviceParameters) Imports() []string {
	return []string{
		"context",
		"fmt",
		"net",
		"sync",
		"google.golang.org/grpc",
		"github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb",
	}
}

func (serviceParameters) ResponseFieldName() string {
	return goNames.CamelCase(responseFieldName)
}

func (serviceParameters) ResponseStateFieldName() string {
	return goNames.CamelCase(responseStateFieldName)
}

func (serviceParameters) ResponseErrorFieldName() string {
	return goNames.CamelCase(responseErrorFieldName)
}

func printEffect(statePrefix, inputPrefix string, state *desc.MessageDescriptor, effect *pb.Effect) (string, error) {
	if up := effect.GetUpdate(); up != nil {
		res, err := printUpdateEffect(statePrefix, inputPrefix, state, up)
		if err != nil {
			return "", fmt.Errorf("failed to print update effect: %w", err)
		}

		return res, nil
	}

	return "", fmt.Errorf("no effect specified")
}

func printUpdateEffect(statePrefix, inputPrefix string, state *desc.MessageDescriptor, up *pb.Effect_Update) (string, error) {
	var (
		nilCheck        string
		initializations []string
		lh              string
		rh              string
	)

	// Build the left-hand expression, to update the state.
	// But first, initialize the state as needed.

	dst := up.GetDest()
	if dst == nil {
		return "", fmt.Errorf("no property specified to be set")
	}

	initializations, err := printInitializeStatePropertyExpression(statePrefix, state, dst.Path)
	if err != nil {
		return "", fmt.Errorf("failed to print state property initialization expression: %w", err)
	}

	lh = fmt.Sprintf("%s.%s", statePrefix, joinCamelCase(dst.Path))

	// Build the right-hand side expression.

	src := up.GetSrc()
	if src == nil {
		return "", fmt.Errorf("no source specified")
	}

	if v := src.GetValue(); v != nil {
		_, val, err := extractValue(v)
		if err != nil {
			return "", fmt.Errorf("failed to extract value: %w", err)
		}

		rh = fmt.Sprint(val)
	} else {
		var p *pb.Path
		var prefix string

		if p = src.GetProp(); p != nil {
			prefix = statePrefix
		} else if p = src.GetInput(); p != nil {
			prefix = inputPrefix
		} else {
			return "", fmt.Errorf("no source property or value specified to update with")
		}

		nilChecks := make([]string, len(p.Path)-1)
		for i := range nilChecks {
			nilChecks[i] = fmt.Sprintf("(%s.%s != nil)", prefix, joinCamelCase(p.Path[:i+1]))
		}
		nilCheck = strings.Join(nilChecks, " && ")

		rh = fmt.Sprintf("%s.%s", prefix, joinCamelCase(p.Path))
	}

	res := lh + " = " + rh

	// TODO: instead of conditionally applying the update, fail if the nil checks don't pass.
	if nilCheck != "" {
		res = `if ` + nilCheck + ` { ` + res + ` }`
	}

	return strings.Join(initializations, "\n") + "\n\n" + res, nil
}

func printResponseAppendExpression(responsePrefix, statePrefix string, state *desc.MessageDescriptor, path []string) (string, error) {
	initializations, err := printInitializeStatePropertyExpression(responsePrefix, state, path)
	if err != nil {
		return "", fmt.Errorf("failed to print response property initialization expression: %w", err)
	}

	initializeParentFields := strings.Join(initializations, "\n")
	resolvedProperty := joinCamelCase(path)
	setField := fmt.Sprintf("%s.%s = %s.%s", responsePrefix, resolvedProperty, statePrefix, resolvedProperty)

	return initializeParentFields + "\n" + setField, nil
}

func printInitializeStatePropertyExpression(identifier string, state *desc.MessageDescriptor, path []string) ([]string, error) {
	allButLast := path[:len(path)-1]
	initializations := make([]string, len(allButLast))

	msgType := state

	for i, name := range allButLast {
		field := msgType.FindFieldByName(name)
		if field == nil {
			return nil, fmt.Errorf("no field under the path %s found in the state message description", strings.Join(path[:i+1], "."))
		}

		if msgType = field.GetMessageType(); msgType == nil {
			return nil, fmt.Errorf("failed to look up message type for field %s in the state message description", strings.Join(path[:i+1], ","))
		}

		typeName := goNames.GoTypeOfField(field).String()

		if !strings.HasPrefix(typeName, "*") {
			// no initialization needed since the
			// field is a not a pointer.
			continue
		}

		typeName = typeName[1:]

		parts := strings.Split(typeName, ".")

		switch len(parts) {
		case 1:
		case 2:
			// remove the package name prefix.
			typeName = parts[1]
		default:
			return nil, fmt.Errorf("type name has more that one '.': %s", typeName)
		}

		fullResponseFieldName := fmt.Sprintf("%s.%s", identifier, joinCamelCase(allButLast[:i+1]))

		initializations = append(initializations, `
			if `+fullResponseFieldName+` == nil {
				`+fullResponseFieldName+` = &`+typeName+`{}
			}
		`)
	}

	return initializations, nil
}

func joinCamelCase(a []string) string {
	b := make([]string, len(a))
	for i, s := range a {
		b[i] = goNames.CamelCase(s)
	}
	return strings.Join(b, ".")
}
