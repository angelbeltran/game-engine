package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:generate go run main.go

const (
	jsonSpec               = "spec.json"
	protoTemplateExtension = "proto.tmpl"
	goTemplateExtension    = "go.tmpl"
	outputDir              = "output"
)

type (
	// NOTE: might be able to get rid of these altogether, use interface{}s, and let the json do it all!
	specDefinition struct {
		Types                []string
		ProtobufTypesToTypes map[string]string
		TypesToProtobufTypes map[string]string
		Functions            functionDefinitions
	}

	functionDefinitions struct {
		Unary  map[string]map[string][]string
		Binary map[string]map[string]map[string][]string
	}
)

func main() {
	if err := generateProtos(jsonSpec, protoTemplateExtension); err != nil {
		log.Fatalf("failed to generate .proto files: %v", err)
	}

	if err := generateProtos(jsonSpec, goTemplateExtension); err != nil {
		log.Fatalf("failed to generate .go files: %v", err) // TODO: apply formatting?
	}
}

func generateProtos(filename, templateFileExtension string) error {
	spec, err := decodeSpecFile(filename)
	if err != nil {
		return fmt.Errorf("failed to decode file %s: %w", filename, err)
	}

	tmpl := template.New("base").Funcs(template.FuncMap{
		"toUpper":    strings.ToUpper,
		"capitalize": capitalize,
		"inc":        inc,
		"dec":        dec,
		"add":        add,
		"multiply":   multiply,
	})

	pattern := "*." + templateFileExtension
	if tmpl, err = tmpl.ParseGlob(pattern); err != nil {
		return fmt.Errorf("failed to parse template files %s: %w", pattern, err)
	}

	for _, t := range tmpl.Templates() {
		ext := "proto"
		parts := strings.Split(strings.TrimSuffix(t.Name(), ".tmpl"), ".")
		if len(parts) > 1 {
			ext = parts[len(parts)-1]
		}

		target := outputDir + "/" + strings.TrimSuffix(t.Name(), "."+templateFileExtension) + "." + ext

		fd, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer fd.Close()

		if err = t.Execute(fd, spec); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", t.Name(), err)
		}
	}

	return nil
}

func decodeSpecFile(filename string) (spec specDefinition, err error) {
	fd, err := os.Open(filename)
	if err != nil {
		return spec, err
	}
	defer fd.Close()

	err = json.NewDecoder(fd).Decode(&spec)

	return spec, err
}

// template utitlies

func inc(i int) int {
	return i + 1
}

func dec(i int) int {
	return i - 1
}

func add(n ...int) int {
	sum := 0
	for _, i := range n {
		sum += i
	}
	return sum
}

func multiply(n ...int) int {
	product := 1
	for _, i := range n {
		product *= i
	}
	return product
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
