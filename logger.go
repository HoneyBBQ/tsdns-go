package tsdns

import (
	"log"
	"os"
)

// Logger defines the interface for logging operations
type Logger interface {
	// Debug logs debug level messages
	Debug(format string, args ...interface{})
	// Info logs info level messages
	Info(format string, args ...interface{})
	// Warn logs warning level messages
	Warn(format string, args ...interface{})
	// Error logs error level messages
	Error(format string, args ...interface{})
	// Fatal logs fatal level messages and exits
	Fatal(format string, args ...interface{})
}

type stdLogger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

// NewStdLogger creates a new standard output logger
func newStdLogger() Logger {
	return &stdLogger{
		debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		info:  log.New(os.Stdout, "[INFO]  ", log.LstdFlags),
		warn:  log.New(os.Stdout, "[WARN]  ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		fatal: log.New(os.Stderr, "[FATAL] ", log.LstdFlags),
	}
}

func (l *stdLogger) Debug(format string, args ...interface{}) {
	l.debug.Printf(format, args...)
}

func (l *stdLogger) Info(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}

func (l *stdLogger) Warn(format string, args ...interface{}) {
	l.warn.Printf(format, args...)
}

func (l *stdLogger) Error(format string, args ...interface{}) {
	l.error.Printf(format, args...)
}

func (l *stdLogger) Fatal(format string, args ...interface{}) {
	l.fatal.Printf(format, args...)
	os.Exit(1)
}
