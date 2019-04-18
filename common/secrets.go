package common

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/envutil"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/pkg/errors"
	"github.com/viktorbenei/bitrise-devops-microservice/services/notes/common/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Config ...
type Config map[string]string

// LoadSecret ...
func LoadSecret(secretKey string) (string, error) {
	if envVar := os.Getenv(secretKey); len(envVar) > 0 {
		return envVar, nil
	}

	// check in file
	secretConfigPath := envutil.GetenvWithDefault("SECRETS_CONFIG_FILE_PATH", "/var/secret/config.yaml")
	fileContBytes, err := fileutil.ReadBytesFromFile(secretConfigPath)
	if err != nil {
		logger.L.Error("Failed to open secret config file",
			zap.String("path", secretConfigPath),
		)
	} else {
		var config Config
		if err := yaml.Unmarshal(fileContBytes, &config); err != nil {
			logger.L.Error("Failed to parse secret config file",
				zap.String("error", fmt.Sprintf("%+v", err)),
			)
		}
		if val, ok := config[secretKey]; ok {
			return val, nil
		}
	}

	// not found
	return "", errors.Errorf("Secret (%s) not found", secretKey)
}
