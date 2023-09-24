package api

import "fmt"

type Logger interface {
	Infof(f string, a ...any)
	Debugf(f string, a ...any)
}

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
