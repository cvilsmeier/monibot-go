package monibot

import (
	"log"
)

// A Logger prints debug messages.
// It is used by Api to print debug messages.
type Logger interface {

	// Debug prints a debug message.
	Debug(f string, a ...any)
}

// NewDefaultLogger creates a new Logger that logs to go log package.
func NewDefaultLogger() Logger {
	return &defaultLogger{}
}

type defaultLogger struct {
}

var _ Logger = (*defaultLogger)(nil)

func (x *defaultLogger) Debug(f string, a ...any) {
	log.Printf("DEBUG: "+f+"\n", a...)
}

// NewDiscardLogger creates a new Logger that discards all output.
func NewDiscardLogger() Logger {
	return &discardLogger{}
}

type discardLogger struct {
}

var _ Logger = (*discardLogger)(nil)

func (x *discardLogger) Debug(f string, a ...any) {
}
