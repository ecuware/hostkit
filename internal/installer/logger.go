package installer

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

// Logger handles installation logging
type Logger struct {
	writer io.Writer
	level  LogLevel
	prefix string
}

// LogLevel represents logging levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	SUCCESS
)

// NewLogger creates a new logger
func NewLogger() *Logger {
	return &Logger{
		writer: os.Stdout,
		level:  INFO,
		prefix: "[HostKit]",
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetWriter sets the output writer
func (l *Logger) SetWriter(w io.Writer) {
	l.writer = w
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("15:04:05")
	var prefix string

	switch level {
	case DEBUG:
		prefix = color.CyanString("[DEBUG]")
	case INFO:
		prefix = color.BlueString("[INFO]")
	case WARN:
		prefix = color.YellowString("[WARN]")
	case ERROR:
		prefix = color.RedString("[ERROR]")
	case SUCCESS:
		prefix = color.GreenString("[SUCCESS]")
	}

	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.writer, "%s %s %s %s\n", timestamp, l.prefix, prefix, message)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Success logs a success message
func (l *Logger) Success(format string, args ...interface{}) {
	l.log(SUCCESS, format, args...)
}

// Write implements io.Writer interface for streaming command output
func (l *Logger) Write(p []byte) (n int, err error) {
	return l.writer.Write(p)
}
