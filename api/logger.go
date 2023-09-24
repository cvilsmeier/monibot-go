package api

import (
	"fmt"
	"io"
)

// Logger prints log messages.
type Logger interface {

	// Infof prints a info messages.
	Infof(f string, a ...any)

	// Debugf prints a (verbose) debug messages.
	Debugf(f string, a ...any)
}

// NewLogger creates a new Logger. If verbose is true,
// debug output is enabled, otherwise it is discarded.
func NewLogger(w io.Writer, verbose bool) Logger {
	return &loggerImpl{w, verbose}
}

type loggerImpl struct {
	w       io.Writer
	verbose bool
}

var _ Logger = (*loggerImpl)(nil)

func (l *loggerImpl) Infof(f string, a ...any) {
	fmt.Fprintf(l.w, f+"\n", a...)
}

func (l *loggerImpl) Debugf(f string, a ...any) {
	if l.verbose {
		fmt.Fprintf(l.w, "DEBUG: "+f+"\n", a...)
	}
}
