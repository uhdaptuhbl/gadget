package logging

import (
	"context"
	"strings"
)

type LogFormat string

func (lf LogFormat) String() string {
	return string(lf)
}

const LogFormatHuman = LogFormat("human")
const LogFormatJSON = LogFormat("json")

func LogFormats() []string {
	return []string{string(LogFormatHuman), string(LogFormatJSON)}
}
func PrettyLogFormats() string {
	return strings.Join(LogFormats(), ",")
}

type LogLevel string

func (ll LogLevel) String() string {
	return string(ll)
}

const LogLevelError = LogLevel("error")
const LogLevelWarn = LogLevel("warn")
const LogLevelInfo = LogLevel("info")
const LogLevelDebug = LogLevel("debug")

func LogLevels() []string {
	return []string{string(LogLevelError), string(LogLevelWarn), string(LogLevelInfo), string(LogLevelDebug)}
}
func PrettyLogLevels() string {
	return strings.Join(LogLevels(), ",")
}

type LogVerbosity string

func (lv LogVerbosity) String() string {
	return string(lv)
}

// LogVerbosityBare should be used to log only the message with no other info
const LogVerbosityBare = LogVerbosity("bare")

// LogVerbositySimple should be used to log the message with basic level and file:lineno
const LogVerbositySimple = LogVerbosity("simple")

// LogVerbosityVerbose should be used to log all information including all logger fields
const LogVerbosityVerbose = LogVerbosity("verbose")

func LogVerbosities() []string {
	return []string{string(LogVerbosityBare), string(LogVerbositySimple), string(LogVerbosityVerbose)}
}
func PrettyLogVerbosities() string {
	return strings.Join(LogVerbosities(), ",")
}

// CorrelationID for application logs
type CorrelationID int

// RequestIDKey for application logs
const RequestIDKey CorrelationID = iota

// Config - logging settings
type Config struct {
	Build     string       `mapstructure:"build" json:"build"`
	Version   string       `mapstructure:"version" json:"version"`
	Format    LogFormat    `mapstructure:"format" json:"format"`
	Level     LogLevel     `mapstructure:"level" json:"level"`
	Verbosity LogVerbosity `mapstructure:"verbosity" json:"verbosity"`

	// typically a local absolute file path but when using the zap logging there are some additional options. See: https://pkg.go.dev/go.uber.org/zap#Open
	OutputPaths []string `mapstructure:"outputpaths" json:"outputPaths"`
}

// Logger - Standard Team Cymru log interface
type Logger interface {
	// Configure - configure the logger based on a configuration struct
	Configure(Config) error

	// HandleError checks if err has already been logged, otherwise logs it, wraps, and returns it.
	HandleError(err error) error

	// Traced - returns an updated logger instance that includes tracing information (request id, spans etc)
	Traced(ctx context.Context) Logger

	// WithExtraFields - returns updated logger instance including the extra kev value pairs appended to all log lines.
	WithExtraFields(fields map[string]string) Logger

	// AddCallerSkip - adjust the callstack skip, useful when wrapping one logger with another
	AddCallerSkip(int) Logger

	// Error - write a log message at the error level
	Error(args ...interface{})

	// Errorf - write a string formatted log message at the error level
	Errorf(template string, args ...interface{})

	// Errorw - write a structured log message at the error level
	Errorw(msg string, keysAndValues ...interface{})

	// Warn - write a log message at the warn level
	Warn(args ...interface{})

	// Warnf - write a string formatted log message at the warn level
	Warnf(template string, args ...interface{})

	// Warnw - write a structured log message at the warn level
	Warnw(msg string, keysAndValues ...interface{})

	// Info - write a log message at the info level
	Info(args ...interface{})

	// Infof - write a string formatted log message at the info level
	Infof(template string, args ...interface{})

	// Infow - write a structured log message at the info level
	Infow(msg string, keysAndValues ...interface{})

	// Debug - write a log message at the debug level
	Debug(args ...interface{})

	// Debugf - write a string formatted log message at the debug level
	Debugf(template string, args ...interface{})

	// Debugw - write a structured log message at the debug level
	Debugw(msg string, keysAndValues ...interface{})
}
