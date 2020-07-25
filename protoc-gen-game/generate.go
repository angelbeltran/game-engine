package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/jhump/goprotoc/plugins"
	"github.com/jhump/protoreflect/desc"

	pb "github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb"
	"github.com/angelbeltran/game-engine/protoc-gen-game/types"
)

func generateService(w io.Writer, opts generationOptions) error {

	// Load template functions.

	tmpl := template.New("service").Funcs(template.FuncMap{
		"printRule":                     printRule,
		"printEffect":                   printEffect,
		"printResponseAppendExpression": printResponseAppendExpression,
	})

	// Parse templates.

	tmpl, err := tmpl.Parse(serviceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	// Apply runtime parameters.

	out := bytes.NewBuffer([]byte{})

	if tmpl.Execute(out, opts); err != nil {
		return err
	}

	// Format and write to file.

	b, err := ioutil.ReadAll(out)
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

type (
	generationOptions struct {
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

func (generationOptions) Imports() []string {
	return []string{
		"context",
		"fmt",
		"net",
		"sync",
		"google.golang.org/grpc",
		"github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb",
	}
}

func (generationOptions) ResponseFieldName() string {
	return goNames.CamelCase(responseFieldName)
}

func (generationOptions) ResponseStateFieldName() string {
	return goNames.CamelCase(responseStateFieldName)
}

func (generationOptions) ResponseErrorFieldName() string {
	return goNames.CamelCase(responseErrorFieldName)
}

var goNames plugins.GoNames

func printRule(statePrefix, inputPrefix string, rule *pb.Rule) (string, error) {
	if s := rule.GetSingle(); s != nil {
		var op string

		switch v := s.GetOperator(); v {
		case pb.Rule_Single_EQ:
			op = "=="
		case pb.Rule_Single_NEQ:
			op = "!="
		case pb.Rule_Single_LT:
			op = "<"
		case pb.Rule_Single_LTE:
			op = "<="
		case pb.Rule_Single_GT:
			op = ">"
		case pb.Rule_Single_GTE:
			op = ">="
		default:
			return "", fmt.Errorf("unexpected operator: %s", v)
		}

		lh, err := printOperandWithNilChecks(statePrefix, inputPrefix, s.GetLeft())
		if err != nil {
			return "", err
		}

		rh, err := printOperandWithNilChecks(statePrefix, inputPrefix, s.GetRight())
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s %s %s", lh, op, rh), nil
	}

	if and := rule.GetAnd(); and != nil {
		printed := make([]string, len(and.Rules))

		for i, r := range and.Rules {
			str, err := printRule(statePrefix, inputPrefix, r)
			if err != nil {
				return "", err
			}

			printed[i] = "(" + str + ")"
		}

		return strings.Join(printed, " && "), nil
	}

	if or := rule.GetOr(); or != nil {
		printed := make([]string, len(or.Rules))

		for i, r := range or.Rules {
			str, err := printRule(statePrefix, inputPrefix, r)
			if err != nil {
				return "", err
			}

			printed[i] = "(" + str + ")"
		}

		return strings.Join(printed, " || "), nil
	}

	return "", fmt.Errorf("empty rule definition")
}

func printOperandWithNilChecks(statePrefix, inputPrefix string, op *pb.Operand) (string, error) {
	if v := op.GetValue(); v != nil {
		_, val, err := extractValue(v)
		if err != nil {
			return "", err
		}

		return fmt.Sprint(val), nil
	}

	if p := op.GetProp(); p != nil {
		allButLast := p.Path[:len(p.Path)-1]

		nilChecks := make([]string, len(allButLast))
		for i := range allButLast {
			nilChecks[i] = fmt.Sprintf("%s.%s != nil", statePrefix, joinCamelCase(p.Path[:i+1]))
		}

		return strings.Join(append(nilChecks, fmt.Sprintf("%s.%s", statePrefix, joinCamelCase(p.Path))), " && "), nil
	}

	if in := op.GetInput(); in != nil {
		allButLast := in.Path[:len(in.Path)-1]

		nilChecks := make([]string, len(allButLast))
		for i := range allButLast {
			nilChecks[i] = fmt.Sprintf("%s.%s != nil", inputPrefix, joinCamelCase(in.Path[:i+1]))
		}

		return strings.Join(append(nilChecks, fmt.Sprintf("%s.%s", inputPrefix, joinCamelCase(in.Path))), " && "), nil
	}

	return "", fmt.Errorf("undefined operand")
}

func printEffect(statePrefix, inputPrefix string, state *desc.MessageDescriptor, effect *pb.Effect) (string, error) {
	if up := effect.GetUpdate(); up != nil {
		return printUpdateEffect(statePrefix, inputPrefix, state, up)
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
		return "", fmt.Errorf("no property specified to be set in update")
	}

	initializations, err := printInitializeStatePropertyExpression(statePrefix, state, dst.Path)
	if err != nil {
		return "", fmt.Errorf("failed to print state property initialization expression: %w", err)
	}

	lh = fmt.Sprintf("%s.%s", statePrefix, joinCamelCase(dst.Path))

	// Build the right-hand side expression.

	src := up.GetSrc()
	if src == nil {
		return "", fmt.Errorf("no source specified to update with")
	}

	if v := src.GetValue(); v != nil {
		_, val, err := extractValue(v)
		if err != nil {
			return "", err
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
