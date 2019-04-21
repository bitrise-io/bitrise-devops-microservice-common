package common

import (
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/envutil"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/pkg/errors"
	"github.com/viktorbenei/bitrise-devops-microservice-common/common/logger"
	"go.uber.org/zap"
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
	if exist, err := pathutil.IsPathExists(secretConfigFilePath); err != nil {
		logger.L.Error("Failed to check if secret file exist at path",
			zap.String("path", secretConfigFilePath),
			zap.Error(err),
		)
	} else if exist {
		secretValue, err := fileutil.ReadStringFromFile(secretConfigFilePath)
		if err != nil {
			logger.L.Error("Failed to open secret config file",
				zap.String("path", secretConfigFilePath),
				zap.Error(err),
			)
		} else if len(secretValue) < 1 {
			logger.L.Error("Secret file found but was empty",
				zap.String("path", secretConfigFilePath),
			)
		} else {
			return secretValue, nil
		}
	}

	// Check if it's a file (with filename==secretKey) is a SUBDIRECTORY of secrets config dir.
	// Continue if not.
	filepathPattern := filepath.Join(secretConfigDirPath, "*", secretKey)
	files, err := filepath.Glob(filepathPattern)
	if err != nil {
		logger.L.Error("Failed to list files for pattern",
			zap.String("path", filepathPattern),
			zap.Error(err),
		)
	}
	if len(files) < 1 {
		logger.L.Error("No file found for pattern",
			zap.String("path", filepathPattern),
		)
	} else if len(files) > 1 {
		logger.L.Error("More than 1 file found for pattern",
			zap.String("path", filepathPattern),
		)
	} else {
		secretConfigFilePath := files[0]
		secretValue, err := fileutil.ReadStringFromFile(secretConfigFilePath)
		if err != nil {
			logger.L.Error("Failed to open secret config file",
				zap.String("path", secretConfigFilePath),
				zap.Error(err),
			)
		} else if len(secretValue) < 1 {
			logger.L.Error("Secret file found but was empty",
				zap.String("path", secretConfigFilePath),
			)
		} else {
			return secretValue, nil
		}
	}

	// not found
	return "", errors.Errorf("Secret (%s) not found", secretKey)
}
