package main

import (
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
func newLogger() *zap.Logger {
	verbosity := getEnv("VERBOSITY", "info")

	lvl, ok := getZapLevelFromString(verbosity)
	if !ok {
		panic("invalid verbosity level provided")
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(lvl),
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	cfg.Level.SetLevel(lvl)

	return logger
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
