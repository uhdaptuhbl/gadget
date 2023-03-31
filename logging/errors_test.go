package logging

import (
	"errors"
	"testing"
)

func TestErrInvalidLogFormatReturns_Error(t *testing.T) {
	expected := "invalid log format invalid expected one of: human,json"
	e := ErrInvalidLogFormat{Input: "invalid"}
	eStr := e.Error()
	if eStr != expected {
		t.Errorf("did not get expected error '%s' != '%s'", eStr, expected)
	}
}

func TestErrInvalidLogFormatWorksWith_As(t *testing.T) {
	var err = &ErrInvalidLogFormat{Input: "invalid"}
	if !errors.As(err, new(*ErrInvalidLogFormat)) {
		t.Error("expected error type was not found")
	}
}

func TestErrInvalidLogLevelReturns_Error(t *testing.T) {
	expected := "invalid log level invalid expected one of: error,warn,info,debug"
	e := ErrInvalidLogLevel{Input: "invalid"}
	eStr := e.Error()
	if eStr != expected {
		t.Errorf("did not get expected error '%s' != '%s'", eStr, expected)
	}
}

func TestErrInvalidLogLevelWorksWith_As(t *testing.T) {
	var err = &ErrInvalidLogLevel{Input: "invalid"}
	if !errors.As(err, new(*ErrInvalidLogLevel)) {
		t.Error("expected error type was not found")
	}
}

func TestErrUnableToInitialize_Error(t *testing.T) {
	expected := "test"
	e := ErrUnableToInitialize{err: errors.New("test")}
	eStr := e.Error()
	if eStr != expected {
		t.Errorf("did not get expected error '%s' != '%s'", eStr, expected)
	}
}

func TestErrUnableToInitialize_As(t *testing.T) {
	var err = &ErrUnableToInitialize{err: errors.New("test")}
	if !errors.As(err, new(*ErrUnableToInitialize)) {
		t.Error("expected error type was not found")
	}
}
