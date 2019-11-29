package sys

import (
	"io"

	"go.uber.org/zap"
)

// IO encapsulates side effects
type IO struct {
	Log *zap.SugaredLogger

	Stdout io.Writer
	Stderr io.Writer
}

var _ Logger = &zap.SugaredLogger{}

// Logger provides methods for loosely-typed, structured logging
type Logger interface {
	Sync() error

	// Verbosity == 0
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	// Verbosity >= 1
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	// Verbosity >= 2
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	// Verbosity >= 3
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}
