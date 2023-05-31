package cookiejar

import (
	"fmt"
	"net/http"
	stdjar "net/http/cookiejar"

	"golang.org/x/net/publicsuffix"

	"gadget/logging"
)

type Constructor interface {
	Logger(log logging.Logger) Constructor
	Jar(jar http.CookieJar) Constructor
	Strict() Constructor
	HandleErrors(handler ErrorHandler) Constructor
	New() http.CookieJar
}

type builder struct {
	opts []Option
}

func Builder() Constructor {
	return &builder{opts: make([]Option, 0)}
}

func (bldr *builder) Logger(log logging.Logger) Constructor {
	bldr.opts = append(bldr.opts, WithLogger(log))
	return bldr
}

func (bldr *builder) Jar(jar http.CookieJar) Constructor {
	bldr.opts = append(bldr.opts, WithJar(jar))
	return bldr
}

func (bldr *builder) Strict() Constructor {
	bldr.opts = append(bldr.opts, Strict)
	return bldr
}

func (bldr *builder) HandleErrors(handler ErrorHandler) Constructor {
	bldr.opts = append(bldr.opts, HandleErrors(handler))
	return bldr
}

func (bldr *builder) New() http.CookieJar {
	// NOTE: creating an entirely new instance here allows `New()`
	// to be called multiple times, returning a separate instance each time.
	var jar = newCookieContainer()

	for _, option := range bldr.opts {
		option(jar)
	}

	// TODO: should this return *cookieContainer instead?
	return finalizeCookieContainer(jar)
}

type Option func(jar *cookieContainer)

func WithLogger(log logging.Logger) Option {
	return func(jar *cookieContainer) {
		jar.log = log
	}
}

func WithJar(jar http.CookieJar) Option {
	return func(jar *cookieContainer) {
		jar.data = jar
	}
}

func Strict(jar *cookieContainer) {
	jar.strict = true
}

func HandleErrors(handler ErrorHandler) Option {
	return func(jar *cookieContainer) {
		jar.errhandler = handler
	}
}

// func GetErrors(ptr *chan<- error) Option {
// 	var ec = make(chan error, 1000)
// 	if ptr != nil {
// 		*ptr = ec
// 	}
// 	return func(jar *cookieContainer) {
// 		jar.errchan = ec
// 		jar.errhandler = func(err error) bool {
// 			ec <- err
// 			return false
// 		}
// 	}
// }

func newCookieContainer() *cookieContainer {
	return &cookieContainer{
		nameMap:    make(map[string]string),
		nameLookup: make(map[string]string),
	}
}

func finalizeCookieContainer(jar *cookieContainer) *cookieContainer {
	var err error
	if jar.log == nil {
		jar.log = logging.NewNoopLogger()
	} else {
		jar.log.Debug("cookie jar using provided logger")
	}

	if jar.errhandler == nil {
		jar.errhandler = jar.defaultErrorHandler
	}

	if jar.data == nil {
		if jar.data, err = stdjar.New(&stdjar.Options{PublicSuffixList: publicsuffix.List}); err != nil {
			// NOTE: As of Go 1.16, cookiejar.New err is hardcoded nil:
			// https://cs.opensource.google/go/go/+/refs/tags/go1.20.3:src/net/http/cookiejar/jar.go;l=85
			panic(fmt.Sprintf("As of Go 1.16, cookiejar.New err value is SUPPOSED to be hardcoded as nil: %v", err))
		}
		jar.log.Debug("new std lib cookiejar created")
	}
	return jar
}
