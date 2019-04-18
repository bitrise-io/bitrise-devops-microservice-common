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

// LoadSecret ...
func LoadSecret(secretKey string) (string, error) {
	if envVar := os.Getenv(secretKey); len(envVar) > 0 {
		return envVar, nil
	}

	// check in files
	secretConfigDirPath := envutil.GetenvWithDefault("SECRETS_CONFIG_DIR_PATH", "/var/secret")
	secretConfigFilePath := filepath.Join(secretConfigDirPath, secretKey)
	secretValue, err := fileutil.ReadStringFromFile(secretConfigFilePath)
	if err != nil {
		logger.L.Error("Failed to open secret config file",
			zap.String("path", secretConfigFilePath),
		)
	} else {
		return secretValue, nil
	}

	// not found
	return "", errors.Errorf("Secret (%s) not found", secretKey)
}
