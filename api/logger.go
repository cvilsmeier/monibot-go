package api

import "fmt"

// Logger prints log messages.
type Logger interface {

	// Infof prints a info messages.
	Infof(f string, a ...any)

	// Debugf prints a (verbose) debug messages.
	Debugf(f string, a ...any)
}

// NewLogger creates a new Logger. If verbose is true,
// debug output is enabled, otherwise it is discarded.
func NewLogger(verbose bool) Logger {
	return &loggerImpl{verbose}
}

type loggerImpl struct {
	verbose bool
}

func (l *loggerImpl) Infof(f string, a ...any) {
	fmt.Printf(f+"\n", a...)
}

func (l *loggerImpl) Debugf(f string, a ...any) {
	if l.verbose {
		fmt.Printf("DEBUG: "+f+"\n", a...)
	}
}
