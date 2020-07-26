package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
)

func generateService(w io.Writer, opts serviceParams) error {

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
