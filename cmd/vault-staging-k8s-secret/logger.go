package main

import (
	"errors"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// levelMap is a mapping of string based logging levels
	// to zapcore levels.
	levelMap = map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}
)

// getZapLevelFromString returns the appropriate zapcore level for the given
// string level, if it exists.
func getZapLevelFromString(lvlString string) (zapcore.Level, bool) {
	lvlString = strings.ToLower(lvlString)
	lvl, ok := levelMap[lvlString]
	return lvl, ok
}

// newLogger returns a logger based defaulted to INFO
// Configurable via an environment varaible VERBOSITY
func newLogger(service string) (*zap.Logger, error) {
	verbosity := getEnv("VERBOSITY", "info")

	lvl, ok := getZapLevelFromString(verbosity)
	if !ok {
		return nil, errors.New("invalid verbosity level provided")
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	cfg.OutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableStacktrace = true
	cfg.InitialFields = map[string]interface{}{
		"service": service,
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	cfg.Level.SetLevel(lvl)

	return logger, nil
}

// getEnv gets the value of an enviroinment variable named by the key.
// If the key is not found, a fallback value is used.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
