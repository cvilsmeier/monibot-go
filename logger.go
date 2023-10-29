package monibot

import (
	"log"
)

// A Logger prints debug messages.
type Logger interface {

	// Debug prints a debug message.
	Debug(format string, args ...any)
}

// NewLogLogger creates a new Logger that logs to a log.Logger.
func NewLogLogger(out *log.Logger) Logger {
	return &logLogger{out}
}

type logLogger struct {
	out *log.Logger
}

var _ Logger = (*logLogger)(nil)

func (l *logLogger) Debug(format string, args ...any) {
	l.out.Printf(format+"\n", args...)
}

// NewDiscardLogger creates a new Logger that discards all output.
func NewDiscardLogger() Logger {
	return &discardLogger{}
}

type discardLogger struct {
}

var _ Logger = (*discardLogger)(nil)

func (l *discardLogger) Debug(format string, args ...any) {
}
