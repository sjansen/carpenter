package logger

import (
	"io"
	"io/ioutil"

	isatty "github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Discard returns a logger that discards all log messages, most useful in tests
func Discard() *zap.SugaredLogger {
	return New(0, ioutil.Discard, nil)
}

// New returns a logger
//
// Valid levels are:
//   0 = errors only,
//   1 = include warnings,
//   2 = include informational messages,
//   3 = include debug messages.
//
// If debugLog is non-nil, all log messages will be written to it independent of
// the verbosity of the default log.
func New(verbosity int, defaultLog io.Writer, debugLog io.Writer) *zap.SugaredLogger {
	var level zapcore.Level
	switch {
	case verbosity >= 3:
		level = zapcore.DebugLevel
	case verbosity == 2:
		level = zapcore.InfoLevel
	case verbosity == 1:
		level = zapcore.WarnLevel
	default:
		level = zapcore.ErrorLevel
	}

	// default log
	encoder := zapcore.CapitalLevelEncoder
	if x, ok := defaultLog.(interface{ Fd() uintptr }); ok {
		if isatty.IsTerminal(x.Fd()) {
			encoder = zapcore.CapitalColorLevelEncoder
		}
	}
	cfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		MessageKey:     "msg",
		NameKey:        "logger",
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeLevel:    encoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(defaultLog),
		level,
	)

	// debug log
	if debugLog != nil {
		debugCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
				LevelKey:       "level",
				MessageKey:     "msg",
				NameKey:        "logger",
				TimeKey:        "time",
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
			}),
			zapcore.AddSync(debugLog),
			zapcore.DebugLevel,
		)
		core = zapcore.NewTee(core, debugCore)
	}

	log := zap.New(core).Sugar()
	if debugLog != nil {
		if x, ok := debugLog.(interface{ Name() string }); ok {
			log.Infof("logging to: %s", x.Name())
		}
	}
	return log
}
