package utils

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

var defaultLogger = NewLogger()

// APILogger the function sets logging format to the http server
func APILogger(message string, args ...interface{}) {
	sugar := defaultLogger.Sugar()
	sugar.Infof(message, args)
}

// NewLogger creates a new logger instance with the application specific settings
func NewLogger() *zap.Logger {
	// TODO: make it configurable from environment
	jsoncfg := []byte(`{
	  "level": "info",
	  "encoding": "json",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "severity",
		"levelEncoder": "lowercase",
		"timeKey": "ts",
		"timeEncoder": "millisecond"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(jsoncfg, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()
	return logger
}

// ContextLogger returns the logger stored in the context, if no logger is associated, return a default one
func ContextLogger(ctx context.Context) zap.Logger {
	logger, ok := ctx.Value(contextLogger).(zap.Logger)
	if !ok {
		logger = *defaultLogger
		logger.Info("Context logger not found, using the default logger", zap.Any("context", ctx))
	}
	return logger
}
