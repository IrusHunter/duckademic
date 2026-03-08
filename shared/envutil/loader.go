package envutil

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadENV loads the .env file if the DOCKER_ENV environment variable is not set.
//
// Returns an error if loading the .env file fails.
func LoadENV() error {
	if os.Getenv("DOCKER_ENV") != "" {
		return nil
	}

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}
	return nil
}
