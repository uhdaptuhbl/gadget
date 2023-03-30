package logging

import (
	"fmt"
)

// TODO: these aliases are only here for backward compatibility and should
// eventually be fully deprecated and removed.
type ErrInvalidLogFormat = InvalidLogFormatError
type ErrInvalidLogLevel = InvalidLogLevelError
type ErrUnableToInitialize = InitializeError

// NOTE: pointer receivers are used for these error types to ensure there is
// one and only one way to use the errors and remove potential ambiguity
// for any code doing error checking or that may care about types.

type InvalidLogFormatError struct {
	Input string
}

func (e *InvalidLogFormatError) Error() string {
	return fmt.Sprintf("invalid log format '%s' expected one of: %s", e.Input, PrettyLogFormats())
}

type InvalidLogLevelError struct {
	Input string
}

func (e *InvalidLogLevelError) Error() string {
	return fmt.Sprintf("invalid log level '%s' expected one of: %s", e.Input, PrettyLogLevels())
}

type InvalidVerbosityError struct {
	Input string
}

func (e *InvalidVerbosityError) Error() string {
	return fmt.Sprintf("invalid log verbosity '%s' expected one of: %s", e.Input, PrettyLogVerbosities())
}

type InitializeError struct {
	err error
}

func (e *InitializeError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}
func (e *InitializeError) Unwrap() error {
	return e.err
}

type LoggingHandledError struct {
	err error
}

func (e *LoggingHandledError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}
func (e *LoggingHandledError) Unwrap() error {
	return e.err
}
