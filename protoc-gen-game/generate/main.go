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
	jsonSpec      = "spec.json"
	protoFileType = "proto"
	goFileType    = "go"
	sourceDir     = "src"
	outputDir     = "dst"
)

type (
	// NOTE: might be able to get rid of these altogether, use interface{}s, and let the json do it all!
	specDefinition struct {
		Types                []string
		ProtobufTypesToTypes map[string]string
		TypesToGoTypes       map[string]string
		TypesToProtobufTypes map[string]string
		Functions            functionDefinitions
	}

	functionDefinitions struct {
		Unary  map[string]map[string][]string
		Binary map[string]map[string]map[string][]string
		Nary   map[string]map[string][]string
	}
)

func main() {
	if err := generateProtos(jsonSpec, protoFileType); err != nil {
		log.Fatalf("failed to generate .proto files: %v", err)
	}

	if err := generateProtos(jsonSpec, goFileType); err != nil {
		log.Fatalf("failed to generate .go files: %v", err) // TODO: apply formatting?
	}
}

func generateProtos(specFile, fileTypeExtension string) error {
	spec, err := decodeSpecFile(specFile)
	if err != nil {
		return fmt.Errorf("failed to decode file %s: %w", specFile, err)
	}

	tmpl := template.New("base").Funcs(template.FuncMap{
		"toUpper":    strings.ToUpper,
		"capitalize": capitalize,
		"inc":        inc,
		"dec":        dec,
		"add":        add,
		"multiply":   multiply,
	})

	pattern := sourceDir + "/" + fileTypeExtension + "/*." + fileTypeExtension + ".tmpl"
	if tmpl, err = tmpl.ParseGlob(pattern); err != nil {
		return fmt.Errorf("failed to parse template files %s: %w", pattern, err)
	}

	for _, t := range tmpl.Templates() {
		ext := "proto"
		parts := strings.Split(strings.TrimSuffix(t.Name(), ".tmpl"), ".")
		if len(parts) > 1 {
			ext = parts[len(parts)-1]
		}

		target := outputDir + "/" + fileTypeExtension + "/" + strings.TrimSuffix(t.Name(), "."+fileTypeExtension+".tmpl") + "." + ext

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

func decodeSpecFile(specFile string) (spec specDefinition, err error) {
	fd, err := os.Open(specFile)
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
