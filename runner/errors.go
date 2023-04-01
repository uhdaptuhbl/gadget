package harness

import (
	"strings"
)

/*
 * NOTE: pointer receivers are used for these error types to ensure there is
 * one and only one way to use the errors and remove potential ambiguity
 * for any code doing error checking or that may care about types.
 */

type InvalidValueError struct {
	Label    string
	Value    string
	Expected string
}

func (e *InvalidValueError) Error() string {
	var msg strings.Builder
	msg.WriteString("invalid value for ")
	msg.WriteString(e.Label)
	msg.WriteString(": ")
	if e.Value == "" {
		msg.WriteString("-")
	} else {
		msg.WriteString(e.Value)
	}
	if e.Expected != "" {
		msg.WriteString("  {" + e.Expected + "}")
	}
	return msg.String()
}
