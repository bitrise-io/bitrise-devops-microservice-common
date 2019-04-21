package common

import (
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/envutil"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/pkg/errors"
)

// Config ...
type Config map[string]string

// LoadSecret returns the secret from Env Var or from a secret file.
// If subPath is specified it'll search for the secret file in that sub directory.
// The secret file's name has to match the secretKey, and should
// have no extension.
func LoadSecret(secretKey string) (string, error) {
	// First check if it's set as an env var.
	// Continue if not.
	if envVar := os.Getenv(secretKey); len(envVar) > 0 {
		return envVar, nil
	}

	// Check if it's a file (with filename==secretKey) in the secrets config dir.
	// Continue if not.
	secretConfigDirPath := envutil.GetenvWithDefault("SECRETS_CONFIG_DIR_PATH", "/var/secret")
	secretConfigFilePath := filepath.Join(secretConfigDirPath, secretKey)
	if exist, err := pathutil.IsPathExists(secretConfigFilePath); err == nil && exist {
		secretValue, err := fileutil.ReadStringFromFile(secretConfigFilePath)
		if err == nil && len(secretValue) > 0 {
			return secretValue, nil
		}
	}

	// Check if it's a file (with filename==secretKey) is a SUBDIRECTORY of secrets config dir.
	// Continue if not.
	filepathPattern := filepath.Join(secretConfigDirPath, "*", secretKey)
	files, err := filepath.Glob(filepathPattern)
	if err == nil && len(files) == 1 {
		secretConfigFilePath := files[0]
		secretValue, err := fileutil.ReadStringFromFile(secretConfigFilePath)
		if err == nil && len(secretValue) > 0 {
			return secretValue, nil
		}
	}

	// not found
	return "", errors.Errorf("Secret (%s) not found", secretKey)
}
