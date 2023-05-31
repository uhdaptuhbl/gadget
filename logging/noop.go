package logging

import (
	"context"
)

// NoopLogger - a logger that does nothing, good for unit tests
type NoopLogger struct {
}

// NewNoOpLogger - create a new noop logger instance
func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (n NoopLogger) Configure(config Config) error {
	// noop
	return nil
}

func (n NoopLogger) HandleError(err error) error {
	return err
}

func (n NoopLogger) WrappedSource(source string) Logger {
	return n
}

func (n NoopLogger) Traced(ctx context.Context) Logger {
	// noop
	return n
}

func (n NoopLogger) WithExtraFields(fields map[string]string) Logger {
	// noop
	return n
}

func (n NoopLogger) Error(args ...interface{}) {
	// noop
}

func (n NoopLogger) Errorf(template string, args ...interface{}) {
	// noop
}

func (n NoopLogger) Errorw(msg string, keysAndValues ...interface{}) {
	// noop
}

func (n NoopLogger) Warn(args ...interface{}) {
	// noop
}

func (n NoopLogger) Warnf(template string, args ...interface{}) {
	// noop
}

func (n NoopLogger) Warnw(msg string, keysAndValues ...interface{}) {
	// noop
}

func (n NoopLogger) Info(args ...interface{}) {
	// noop
}

func (n NoopLogger) Infof(template string, args ...interface{}) {
	// noop
}

func (n NoopLogger) Infow(msg string, keysAndValues ...interface{}) {
	// noop
}

func (n NoopLogger) Debug(args ...interface{}) {
	// noop
}

func (n NoopLogger) Debugf(template string, args ...interface{}) {
	// noop
}

func (n NoopLogger) Debugw(msg string, keysAndValues ...interface{}) {
	// noop
}

func (n NoopLogger) AddCallerSkip(skip int) Logger {
	// noop
	return n
}
