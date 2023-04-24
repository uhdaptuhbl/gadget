package harness

import (
	"io"
)

const MinFilePermission = 0600
const MinDirPermission = 0750

// UNUSED allows unused variables to be included in Go programs: USE WITH CAUTION!
func UNUSED(x ...interface{}) {}

func CloseWithError(cf io.Closer, errp *error) {
	if err := cf.Close(); err != nil {
		*errp = err
	}
}
