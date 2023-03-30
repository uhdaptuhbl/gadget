package logging

import (
	"context"
)

// NopLogger - a logger that does nothing, good for unit tests
type NopLogger struct {
}

// NewNoOpLogger - create a new noop logger instance
func NewNoOpLogger() NopLogger {
	return NopLogger{}
}

func (n NopLogger) Configure(config Config) error {
	// noop
	return nil
}

func (n NopLogger) HandleError(err error) error {
	return err
}

func (n NopLogger) WrappedSource(source string) Logger {
	return n
}

func (n NopLogger) Traced(ctx context.Context) Logger {
	// noop
	return n
}

func (n NopLogger) WithExtraFields(fields map[string]string) Logger {
	// noop
	return n
}

func (n NopLogger) Error(args ...interface{}) {
	// noop
}

func (n NopLogger) Errorf(template string, args ...interface{}) {
	// noop
}

func (n NopLogger) Errorw(msg string, keysAndValues ...interface{}) {
	// noop
}

func (n NopLogger) Warn(args ...interface{}) {
	// noop
}

func (n NopLogger) Warnf(template string, args ...interface{}) {
	// noop
}

func (n NopLogger) Warnw(msg string, keysAndValues ...interface{}) {
	// noop
}

func (n NopLogger) Info(args ...interface{}) {
	// noop
}

func (n NopLogger) Infof(template string, args ...interface{}) {
	// noop
}

func (n NopLogger) Infow(msg string, keysAndValues ...interface{}) {
	// noop
}

func (n NopLogger) Debug(args ...interface{}) {
	// noop
}

func (n NopLogger) Debugf(template string, args ...interface{}) {
	// noop
}

func (n NopLogger) Debugw(msg string, keysAndValues ...interface{}) {
	// noop
}

func (n NopLogger) AddCallerSkip(skip int) Logger {
	// noop
	return n
}
