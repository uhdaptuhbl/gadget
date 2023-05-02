package closer

import (
	"io"
)

func CloseWithError(cf io.Closer, errp *error) {
	// TODO: check whether there is already an error at the destination and wrap it if so
	if err := cf.Close(); err != nil {
		*errp = err
	}
}

// func CloseWithLogger(cf io.Closer, errp *error) {
// 	if err := cf.Close(); err != nil {
// 		*errp = err
// 	}
// }

// func CloseExit(cf io.Closer, label string) {
// 	if err := cf.Close(); err != nil {
// 		//
// 	}
// }

// TODO: closer interface
