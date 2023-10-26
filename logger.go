package monibot

import (
	"log"
)

// A Logger prints debug messages.
type Logger interface {

	// Debug prints a debug message.
	Debug(f string, a ...any)
}

// NewDefaultLogger creates a new Logger that logs to a log.Logger.
func NewDefaultLogger(out *log.Logger) Logger {
	return &defaultLogger{out}
}

type defaultLogger struct {
	out *log.Logger
}

var _ Logger = (*defaultLogger)(nil)

func (l *defaultLogger) Debug(f string, a ...any) {
	l.out.Printf("DEBUG: "+f+"\n", a...)
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
