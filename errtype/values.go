package errtype

type emptyStringError struct {
	Label string
}

func NewEmptyStringError(label string) *emptyStringError {
	return &emptyStringError{Label: label}
}

func (e *emptyStringError) Error() string {
	return "Found empty string where non-empty string expected: `" + e.Label + "`"
}
