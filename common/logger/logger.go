package logger

import (
	"go.uber.org/zap"
)

// L is the global Logger
var L *zap.Logger

// Init ...
func Init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	L = logger
}

// Close ...
func Close() error {
	return L.Sync()
}
