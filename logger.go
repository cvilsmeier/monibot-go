package monibot

import (
	"fmt"
	"io"
)

// A Logger prints debug messages.
// It is used by Api to print debug messages.
type Logger interface {

	// Debug prints a debug message.
	Debug(f string, a ...any)
}

// NewLogger creates a new Logger that writes to w.
func NewLogger(w io.Writer) Logger {
	if w == nil {
		w = io.Discard
	}
	return &loggerImpl{w}
}

type loggerImpl struct {
	w io.Writer
}

var _ Logger = (*loggerImpl)(nil)

func (x *loggerImpl) Debug(f string, a ...any) {
	if x.w != nil {
		fmt.Fprintf(x.w, "DEBUG: "+f+"\n", a...)
	}
}
