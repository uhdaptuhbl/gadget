package teapot

import (
	"crypto/tls"
	"net/http"

	"gadget/logging"
)

type Constructor interface {
	Logger(log logging.Logger) Constructor
	Config(cfg *Config) Constructor
	Transport(transport *http.Transport) Constructor
	TLS(tlsconfig *tls.Config) Constructor
	AddHeaders(headers http.Header) Constructor
	SetHeaders(headers http.Header) Constructor
	CookieJar(jar http.CookieJar) Constructor
	OnRequest(handlers ...RequestInterceptor) Constructor
	OnResponse(handlers ...ResponseInterceptor) Constructor
	Apply() Teapot
	Make() Teapot
	New() Teapot
}

type builder struct {
	tpt  *teapot  `mapstructure:"-" json:"-"`
	opts []Option `mapstructure:"-" json:"-"`
}

func constructorFromTeapot(tpt *teapot) Constructor {
	return &builder{tpt: tpt, opts: make([]Option, 0, 10)}
}

func Builder() Constructor {
	return constructorFromTeapot(nil)
}

func (bldr *builder) Logger(log logging.Logger) Constructor {
	bldr.opts = append(bldr.opts, UseLogger(log))
	return bldr
}

func (bldr *builder) Config(cfg *Config) Constructor {
	bldr.opts = append(bldr.opts, UseConfig(cfg))
	return bldr
}

func (bldr *builder) Transport(transport *http.Transport) Constructor {
	bldr.opts = append(bldr.opts, UseTransport(transport))
	return bldr
}

func (bldr *builder) TLS(tlsconfig *tls.Config) Constructor {
	bldr.opts = append(bldr.opts, UseTLS(tlsconfig))
	return bldr
}

func (bldr *builder) AddHeaders(h http.Header) Constructor {
	bldr.opts = append(bldr.opts, UseHeaders(h, false))
	return bldr
}

func (bldr *builder) SetHeaders(h http.Header) Constructor {
	bldr.opts = append(bldr.opts, UseHeaders(h, true))
	return bldr
}

func (bldr *builder) CookieJar(jar http.CookieJar) Constructor {
	bldr.opts = append(bldr.opts, UseCookieJar(jar))
	return bldr
}

// func (bldr *builder) JarLoader(loaders ...cookiejar.Loader) Constructor {
// 	bldr.opts = append(bldr.opts, UseJarLoaders(loaders))
// 	return bldr
// }

func (bldr *builder) OnRequest(handlers ...RequestInterceptor) Constructor {
	// NOTE: OnRequest will ADD handlers to what already exists!
	bldr.opts = append(bldr.opts, UseRequestHandlers(handlers))
	return bldr
}

func (bldr *builder) OnResponse(handlers ...ResponseInterceptor) Constructor {
	// NOTE: OnResponse will ADD handlers to what already exists!
	bldr.opts = append(bldr.opts, UseResponseHandlers(handlers))
	return bldr
}

// Apply will modify an existing teapot and return the pointer.
func (bldr *builder) Apply() Teapot {
	var tpt *teapot
	if bldr.tpt != nil {
		tpt = bldr.tpt
	} else {
		tpt = &teapot{}
	}
	for _, applyOption := range bldr.opts {
		applyOption(tpt)
	}
	return tpt
}

// Make will clone the teapot then modify and return the clone.
func (bldr *builder) Make() Teapot {
	var tpt *teapot
	if bldr.tpt != nil {
		tpt = bldr.tpt.clone()
	} else {
		tpt = &teapot{}
	}
	for _, applyOption := range bldr.opts {
		applyOption(tpt)
	}
	return tpt
}

// New will always create a brand new teapot every time it is called.
func (bldr *builder) New() Teapot {
	var tpt = &teapot{}
	for _, applyOption := range bldr.opts {
		applyOption(tpt)
	}
	return tpt
}
