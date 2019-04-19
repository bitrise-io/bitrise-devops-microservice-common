package common

import (
	"os"

	"github.com/pkg/errors"
)

func requiredEnv(envKey string) (string, error) {
	// TODO: replace this with https://github.com/bitrise-io/go-utils/pull/85 once that's merged
	val := os.Getenv(envKey)
	if len(val) < 1 {
		return "", errors.Errorf("Required environment variable (%s) not provided", envKey)
	}
	return val, nil
}

// LoadConfig returns the config for the key.
func LoadConfig(key string) (string, error) {
	return requiredEnv(key)
}
