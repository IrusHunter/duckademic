package jsonutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ReadFileTo(path string, v any) error {
	dat, err := os.ReadFile(filepath.Join("data", "services.json"))
	if err != nil {
		return fmt.Errorf("can't read file data/services.json: %w", err)
	}

	err = json.Unmarshal(dat, v)
	if err != nil {
		return fmt.Errorf("can't unmarshal data: %w", err)
	}

	return nil
}
