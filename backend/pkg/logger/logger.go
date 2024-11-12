package logger

import (
	"strings"

	"github.com/luuisavelino/network-interface/pkg/envs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutputLogs()},
		Level:       zap.NewAtomicLevelAt(getLevelLogs()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			// LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "message",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			// EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	log, _ = logConfig.Build()
}

func Info(message string, tags ...zap.Field) {
	log.Info(message, tags...)
	log.Sync()
}

func Error(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Info(message, tags...)
	log.Sync()
}

func getOutputLogs() string {
	output := strings.ToLower(strings.TrimSpace(envs.Log.Output))
	if output == "" {
		return "stdout"
	}

	return output
}

func getLevelLogs() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(envs.Log.Level)) {
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}
