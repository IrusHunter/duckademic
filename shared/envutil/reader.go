package envutil

import (
	"fmt"
	"os"
	"strconv"
)

// GetIntFromENV gets a string environment variable by key and converts it to an integer.
//
// Returns an error if the string value is empty or cannot be parsed as an integer.
func GetIntFromENV(key string) (int, error) {
	resultStr := os.Getenv(key)
	if resultStr == "" {
		return 0, fmt.Errorf("environment variable %q is not set", key)
	}

	result, err := strconv.Atoi(resultStr)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q has invalid integer value %q: %w", key, resultStr, err)
	}

	return result, nil
}

// GetIntFromENV gets a string environment variable by key.
func GetStringFromENV(key string) string {
	return os.Getenv(key)
}
