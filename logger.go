package monibot

// A Logger prints debug messages.
type Logger interface {

	// Debug prints a debug message.
	Debug(format string, args ...any)
}

// A zeroLogger logs nothing (a.k.a. zero).
type zeroLogger struct{}

func (zeroLogger) Debug(format string, args ...any) {}
