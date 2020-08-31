package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jhump/goprotoc/plugins"

	"github.com/angelbeltran/game-engine/protoc-gen-game/generate/dst/go/template"
)

//go:generate go generate ./generate
//go:generate protoc -I=./generate/dst/proto --go_out=$GOPATH/src action.proto effect.proto error.proto extensions.proto reference.proto type.proto values_and_functions.proto

//go:generate cp ./generate/dst/proto/action.proto protos/action.proto
//go:generate cp ./generate/dst/proto/effect.proto protos/effect.proto
//go:generate cp ./generate/dst/proto/error.proto protos/error.proto
//go:generate cp ./generate/dst/proto/extensions.proto protos/extensions.proto
//go:generate cp ./generate/dst/proto/reference.proto protos/reference.proto
//go:generate cp ./generate/dst/proto/type.proto protos/type.proto
//go:generate cp ./generate/dst/proto/values_and_functions.proto protos/values_and_functions.proto
//go:generate cp ./generate/dst/proto/bundle.proto protos/bundle.proto

//go:generate go mod vendor

// TODO: get rid of 'types' package once no longer needed.

func main() {
	output := os.Stdout
	os.Stdout = os.Stderr
	err := plugins.RunPlugin(os.Args[0], entrypoint, os.Stdin, output)
	if err != nil {
		os.Exit(1)
	}
}

func entrypoint(req *plugins.CodeGenRequest, resp *plugins.CodeGenResponse) error {
	// Find the 'actions' service.

	pkgName, err := parseArgs(req.Args)
	if err != nil {
		return err
	}

	srvs, err := getActionServices(req.Files)
	if err != nil {
		return err
	}
	if len(srvs) == 0 {
		return fmt.Errorf("no service with is_action_service option set")
	}
	if len(srvs) > 1 {
		return fmt.Errorf("multiple services with is_action_service option set")
	}

	srv := srvs[0]

	// Find and parse the "state" message

	sd, err := getStateDescriptor(req.Files)
	if err != nil {
		return err
	}

	// Find and parse the "response" message

	rd, err := getResponseDescriptor(req.Files)
	if err != nil {
		return err
	}

	if err := validateResponseMessage(sd, rd); err != nil {
		return fmt.Errorf("invalid response message: %w", err)
	}

	// Find and validate the "action" options defined on rpc methods.

	var methods []template.MethodInfo
	responseMessageName := rd.GetFullyQualifiedName()

	for _, method := range srv.GetMethods() {
		name := method.GetOutputType().GetFullyQualifiedName()
		if name != responseMessageName {
			return fmt.Errorf("all rpc methods must output the defined response message type, '%s'. '%s' has '%s' as its output type", responseMessageName, method.GetName(), name)
		}

		// Load the action option field.

		action, err := loadActionOptionMessage(method, actionExtensionFieldNumber)
		if err != nil {
			return err
		}
		if action == nil {
			methods = append(methods, template.MethodInfo{Method: method})
			continue
		}

		// Validate the action option.

		if err := validateAction(sd, method.GetInputType(), action); err != nil {
			return err
		}

		methods = append(methods, template.MethodInfo{
			Method: method,
			Action: action,
		})
	}

	// Find and validate all enum key options.

	msgs, err := getMessagesWithEnumKeys(req.Files)
	if err != nil {
		return fmt.Errorf("failed to parse enum keys options: %w", err)
	}

	if err = validateMessagesWithEnums(msgs); err != nil {
		return fmt.Errorf("invalid enum keys options: %w", err)
	}

	// Generate service files

	w := resp.OutputFile("engine.game.pb.go")

	if err := template.GenerateService(w, template.TemplateParams{
		Package: pkgName,
		Imports: []string{
			"context",
			"fmt",
			"net",
			"sync",
			"google.golang.org/grpc",
			"github.com/angelbeltran/game-engine/protoc-gen-game/game_engine_pb",
		},
		Service:             srv,
		Methods:             methods,
		State:               sd,
		Response:            rd,
		StateVariable:       stateVariable,
		InputVariable:       inputVariable,
		EnumToFieldMappings: msgs,
		ResponseStateField:  responseStateFieldNameCamelCase,
		ResponseErrorField:  responseErrorFieldNameCamelCase,
	}); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	return nil
}

func parseArgs(args []string) (pkgName string, err error) {
	packageParameterMessage := "go package required: --game_opt=package={goPackage}"

	if len(args) == 0 {
		return "", fmt.Errorf("no parameters passed. %s", packageParameterMessage)
	}

	if len(args) > 1 {
		return "", fmt.Errorf("multiple parameters passed. %s", packageParameterMessage)
	}

	parts := strings.Split(args[0], "=")

	if len(parts) != 2 {
		return "", fmt.Errorf("invalid parameter. %s", packageParameterMessage)
	}

	if strings.ToLower(parts[0]) != "package" {
		return "", fmt.Errorf("unrecogized parameter: %s. %s", parts[0], packageParameterMessage)
	}

	if parts[1] == "" {
		return "", fmt.Errorf("empty package name. %s", packageParameterMessage)
	}

	if pkgName = goPackageNamePattern.FindString(parts[1]); pkgName == "" {
		return "", fmt.Errorf("invalid go package: %s", parts[1])
	}

	return pkgName, nil
}

// https://golang.org/ref/spec#PackageName
var goPackageNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
