package cookiejar

import (
	"net/http"

	"gadget/logging"
)

type ErrorChannel chan error
type ErrorHandler func(err error) bool
type Option func(jar *cookieContainer)

func Strict(jar *cookieContainer) {
	jar.strict = true
}

func WithJar(jar http.CookieJar) Option {
	return func(jar *cookieContainer) {
		jar.data = jar
	}
}

func GetErrors(ptr *chan<- error) Option {
	var ec = make(chan error, 1000)
	if ptr != nil {
		*ptr = ec
	}
	return func(jar *cookieContainer) {
		jar.errchan = ec
		jar.errhandler = func(err error) bool {
			ec <- err
			return false
		}
	}
}

func HandleErrors(handler ErrorHandler) Option {
	return func(jar *cookieContainer) {
		jar.errhandler = handler
	}
}

func Logger(log logging.Logger) Option {
	return func(jar *cookieContainer) {
		jar.log = log
	}
}
