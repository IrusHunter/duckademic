package jsonutil

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadFileTo reads a JSON file from the given path and unmarshals its contents into v.
//
// It requires the file path to read from (path) and a pointer to the variable where the unmarshaled JSON will be stored (v).
//
// Returns an error if the file cannot be read or if the JSON is invalid.
func ReadFileTo(path string, v any) error {
	dat, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't read file %s: %w", path, err)
	}

	err = json.Unmarshal(dat, v)
	if err != nil {
		return fmt.Errorf("can't unmarshal data: %w", err)
	}

	return nil
}
