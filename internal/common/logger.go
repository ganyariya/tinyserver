package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// LogLevelDebug represents debug level logging
	LogLevelDebug LogLevel = iota
	// LogLevelInfo represents info level logging
	LogLevelInfo
	// LogLevelWarn represents warning level logging
	LogLevelWarn
	// LogLevelError represents error level logging
	LogLevelError
)

// String returns the string representation of LogLevel
func (ll LogLevel) String() string {
	switch ll {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides a simple logging interface for TinyServer
type Logger struct {
	level  LogLevel
	output io.Writer
	logger *log.Logger
}

// NewLogger creates a new Logger instance
func NewLogger(level LogLevel, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}

	return &Logger{
		level:  level,
		output: output,
		logger: log.New(output, "", 0), // No default prefix or flags
	}
}

// NewDefaultLogger creates a logger with default settings (Info level, stdout)
func NewDefaultLogger() *Logger {
	return NewLogger(LogLevelInfo, os.Stdout)
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() LogLevel {
	return l.level
}

// shouldLog checks if a message should be logged based on the current level
func (l *Logger) shouldLog(level LogLevel) bool {
	return level >= l.level
}

// formatMessage formats a log message with timestamp and level
func (l *Logger) formatMessage(level LogLevel, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), message)
}

// log performs the actual logging
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}

	message := fmt.Sprintf(format, args...)
	formattedMessage := l.formatMessage(level, message)
	l.logger.Println(formattedMessage)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

// ErrorWithErr logs an error message with an error object
func (l *Logger) ErrorWithErr(err error, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(LogLevelError, "%s: %v", message, err)
}

// LogRequest logs an HTTP-like request
func (l *Logger) LogRequest(method, path, remoteAddr string) {
	l.Info("Request: %s %s from %s", method, path, remoteAddr)
}

// LogResponse logs an HTTP-like response
func (l *Logger) LogResponse(statusCode int, contentLength int64, duration time.Duration) {
	l.Info("Response: %d %d bytes in %v", statusCode, contentLength, duration)
}

// LogConnection logs a connection event
func (l *Logger) LogConnection(event, remoteAddr string) {
	l.Info("Connection %s: %s", event, remoteAddr)
}

// Global logger instance for package-level logging functions
var defaultLogger = NewDefaultLogger()

// SetDefaultLogger sets the default logger for package-level functions
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

// GetDefaultLogger returns the default logger
func GetDefaultLogger() *Logger {
	return defaultLogger
}

// Package-level logging functions that use the default logger

// Debug logs a debug message using the default logger
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info logs an info message using the default logger
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs a warning message using the default logger
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs an error message using the default logger
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// ErrorWithErr logs an error message with an error object using the default logger
func ErrorWithErr(err error, format string, args ...interface{}) {
	defaultLogger.ErrorWithErr(err, format, args...)
}

// LogRequest logs an HTTP-like request using the default logger
func LogRequest(method, path, remoteAddr string) {
	defaultLogger.LogRequest(method, path, remoteAddr)
}

// LogResponse logs an HTTP-like response using the default logger
func LogResponse(statusCode int, contentLength int64, duration time.Duration) {
	defaultLogger.LogResponse(statusCode, contentLength, duration)
}

// LogConnection logs a connection event using the default logger
func LogConnection(event, remoteAddr string) {
	defaultLogger.LogConnection(event, remoteAddr)
}

// FormatHTTPDate formats a time for HTTP Date header
func FormatHTTPDate() string {
	return time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}
