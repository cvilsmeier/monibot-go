package monibot

import (
	"log"
)

// A Logger prints debug messages.
type Logger interface {

	// Debug prints a debug message.
	Debug(f string, a ...any)
}

// NewLogger creates a new Logger that logs to a log.Logger.
func NewLogger(out *log.Logger) Logger {
	return &logLogger{out}
}

type logLogger struct {
	out *log.Logger
}

var _ Logger = (*logLogger)(nil)

func (l *logLogger) Debug(f string, a ...any) {
	l.out.Printf(f+"\n", a...)
}

// NewDiscardLogger creates a new Logger that discards all output.
func NewDiscardLogger() Logger {
	return &discardLogger{}
}

type discardLogger struct {
}

var _ Logger = (*discardLogger)(nil)

func (l *discardLogger) Debug(f string, a ...any) {
}
