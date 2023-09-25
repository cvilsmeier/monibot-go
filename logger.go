package monibot

import (
	"fmt"
	"io"
)

// A Logger prints debug messages.
type Logger interface {

	// Debug prints a debug message.
	Debug(f string, a ...any)
}

// NewLogger creates a new Logger. If w is nil then all output is discarded.
func NewLogger(w io.Writer) Logger {
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
