package common

import (
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/envutil"
	"github.com/bitrise-io/go-utils/fileutil"
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
func LoadSecret(secretKey, subPath string) (string, error) {
	if envVar := os.Getenv(secretKey); len(envVar) > 0 {
		return envVar, nil
	}

	// check in files
	secretConfigDirPath := envutil.GetenvWithDefault("SECRETS_CONFIG_DIR_PATH", "/var/secret")
	if len(subPath) > 0 {
		secretConfigDirPath = filepath.Join(secretConfigDirPath, subPath)
	}
	secretConfigFilePath := filepath.Join(secretConfigDirPath, secretKey)
	secretValue, err := fileutil.ReadStringFromFile(secretConfigFilePath)
	if err != nil {
		logger.L.Error("Failed to open secret config file",
			zap.String("path", secretConfigFilePath),
		)
	} else if len(secretValue) < 1 {
		logger.L.Error("Secret file found but was empty",
			zap.String("path", secretConfigFilePath),
		)
	} else {
		return secretValue, nil
	}

	// not found
	return "", errors.Errorf("Secret (%s) not found", secretKey)
}
