package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jhump/goprotoc/plugins"

	"angelbeltran/game-engine/protoc-gen-game/types"
)

//go:generate protoc -I=protos --go_out=$GOPATH/src game_engine.proto
//go:generate go mod vendor

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

	state := types.FromMessage(sd)

	// Find and validate the "action" options defined on rpc methods.

	var methods []methodBundle

	for _, method := range srv.GetMethods() {
		// Load the action option field.

		action, err := loadActionOptionMessage(method, actionExtensionFieldNumber)
		if err != nil {
			return err
		}
		if action == nil {
			methods = append(methods, methodBundle{Method: method})
			continue
		}

		// Parse the action option.

		if err := validateAction(state, action); err != nil {
			return err
		}

		methods = append(methods, methodBundle{
			Method: method,
			Action: action,
		})
	}

	w := resp.OutputFile("engine.go")

	if err := generateAll(w, generationOptions{
		Package:   pkgName,
		Service:   srv,
		Methods:   methods,
		State:     sd,
		StateType: state,
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
