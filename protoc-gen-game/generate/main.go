package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:generate go run main.go

func main() {
	spec, err := decodeSpecFile(specFile)
	if err != nil {
		log.Fatalf("failed to decode file %s: %v", specFile, err)
	}

	if err := generateAll(spec); err != nil {
		log.Fatalf("failed to generate files: %v", err)
	}
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

func generateAll(spec specDefinition) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error while walking directory: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".tmpl") {
			return nil
		}

		parts := strings.Split(path, "/")
		base := parts[len(parts)-1]

		tmpl, err := template.New(base).Funcs(template.FuncMap{
			"toUpper":    strings.ToUpper,
			"capitalize": capitalize,
			"inc":        inc,
			"dec":        dec,
			"add":        add,
			"multiply":   multiply,
		}).ParseFiles(path)
		if err != nil {
			return fmt.Errorf("failed to parse template from file %s: %w", path, err)
		}

		target := outputDir + strings.TrimSuffix(strings.TrimPrefix(path, sourceDir), ".tmpl")

		parts = strings.Split(target, "/")
		dir := strings.Join(parts[:len(parts)-1], "/")
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		fd, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return fmt.Errorf("failed to create and truncate file %s: %w", target, err)
		}
		defer fd.Close()

		if err := tmpl.Execute(fd, spec); err != nil {
			return fmt.Errorf("failed to execute template from file %s: %w", path, err)
		}

		return nil
	})
}

// files and directories

const (
	specFile  = "spec.json"
	sourceDir = "src"
	outputDir = "dst"
)

// parameters available to templates during execution

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
