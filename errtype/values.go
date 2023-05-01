package errtype

import (
	"fmt"
)

type emptyStringError struct {
	Label string
}

func EmptyStringError(label string, extra ...any) *emptyStringError {
	if len(extra) == 0 {
		return &emptyStringError{Label: label}
	}
	return &emptyStringError{Label: fmt.Sprintf(label, extra...)}
}

func (e *emptyStringError) Error() string {
	return "Found empty string where non-empty string expected: `" + e.Label + "`"
}

type notImplementedError struct {
	Item string
}

func NotImplementedError(item string, extra ...any) *notImplementedError {
	if len(extra) == 0 {
		return &notImplementedError{Item: item}
	}
	return &notImplementedError{Item: fmt.Sprintf(item, extra...)}
}

func (e *notImplementedError) Error() string {
	return "`" + e.Item + "` is not implemented"
}
