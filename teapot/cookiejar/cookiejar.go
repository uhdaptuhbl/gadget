package cookiejar

import (
	"fmt"
	"net/http"
	stdjar "net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/publicsuffix"

	"gadget/logging"
)

// cookieContainer implements the http.CookieJar interface and wraps the std lib Jar.
//
//	type CookieJar interface {
//		SetCookies(u *url.URL, cookies []*Cookie)
//		Cookies(u *url.URL) []*Cookie
//	}
//
// The standard library cookie jar has a number of problems in how it deals
// with cookies. The biggest of those problems is that it will silently
// ignore any cookies that it deems "malformed"; spoiler alert - in the real
// world there are tons of cookies that don't conform to the proper specs, but
// they still need to be stored and be retrievable to successfully be able
// to interact with certain servers.
//
// The fact that they are silently ignored can lead to hours or days of trying
// to investigate and understand why requests to a 3rd party site are not
// working, only to find out that cookies aren't being set in the jar because
// the cookie name doesn't conform to specs... and even finding that out
// will likely require that one reads the actual std lib source code.
//
// The silence itself however, is, IMHO a problem in itself - the std lib
// http.CookieJar interface leaves no room for returning or handling errors,
// nor is there an option to disable such strict name checking when you know
// it will cause problems in your application.
//
// https://golangbyexample.com/set-cookie-http-golang/
// https://husni.dev/manage-http-cookie-in-go-with-cookie-jar/
type cookieContainer struct {
	strict     bool
	log        logging.Logger
	data       http.CookieJar
	nameMap    map[string]string
	nameLookup map[string]string

	// TODO: are these actually useful?
	errchan    chan error
	errhandler ErrorHandler
}

// NewcookieContainer constructs a new valid jar.
func New(options ...Option) *cookieContainer {
	var err error
	var jar = &cookieContainer{
		nameMap: make(map[string]string),
		nameLookup: make(map[string]string),
	}

	for _, option := range options {
		option(jar)
	}

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

// SetCookies func of the std lib interface.
func (jar *cookieContainer) SetCookies(uri *url.URL, cookies []*http.Cookie) {
	var name string

	// jar.log.Debugf("SetCookies()")

	// There is no error returned from the interface functions, so we have
	// to just return even though it would make more sense to error here.
	if uri == nil || uri.Scheme == "" || uri.Host == "" {
		var err = fmt.Errorf("empty input: %v %v", uri, cookies)
		jar.log.Error(err)
		if jar.errhandler(err) {
			return
		}
	}

	var cleaned []*http.Cookie
	if jar.strict {
		cleaned = cookies
	} else {
		cleaned = make([]*http.Cookie, 0, len(cookies))

		// The default implementation ignores any cookies with an improper name
		// but we want to store them anyway since requests may not work otherwise;
		// we replace the existing name with a stripped UUID we know is valid.
		for i, cookie := range cookies {
			// NOTE: http.Cookie.Valid() does not appear to catch all blacklisted characters!
			// var errValid = httpcookie.Valid()
			// if errValid == nil && strings.ContainsAny(ffxcookie.Name, "-_:;.\"'") {
			// 	errValid = &cookiejar.CookieError{"invalid Cookie.Name"}
			// }
			// if errValid == nil {
			// 	validCookies++
			// } else {
			// 	badCookies++

			// 	// we replace the cookie.Name with a stripped uuid we know is valid
			// 	// if strings.Contains(errValid.Error(), "invalid Cookie.Name") {
			// 	// 	var key = strings.ReplaceAll(uuid.New().String(), "-", "")
			// 	// 	httpcookie.Raw = httpcookie.Name
			// 	// 	httpcookie.Name = key
			// 	// 	errValid = nil
			// 	// }
			// 	// else {
			// 	// 	log.Debugf("cookie: %+v", httpcookie)
			// 	// 	return errValid
			// 	// }
			// }
			// if errValid != nil {
			// 	errValid = nil
			// 	return errValid
			// }
			if nameKey, ok := jar.nameLookup[cookie.Name]; ok {
				cookie.Name = nameKey
			} else {
				name = strings.ReplaceAll(uuid.New().String(), "-", "")
				jar.nameMap[name] = cookie.Name
				jar.nameLookup[cookie.Name] = name
				cookie.Name = name
			}
			cleaned = append(cleaned, cookies[i])
		}
	}

	if len(cleaned) == 0 {
		jar.log.Debug("THERE ARE NO CLEAN COOKIES TO SET!")
	}

	jar.data.SetCookies(uri, cleaned)
}

// Cookies func of the std lib interface.
func (jar *cookieContainer) Cookies(uri *url.URL) []*http.Cookie {
	// jar.log.Debugf("Cookies()")
	if true {
		// panic((&NoCookiesFoundError{Hosts: []string{uri.String()}}).Error())
	}

	// There is no error returned from the interface functions, so we have
	// to just return even though it would make more sense to error here.
	if uri == nil || uri.Scheme == "" || uri.Host == "" {
		jar.log.Errorf("URI IS EMPTY! %v", uri)
		return nil
	}

	var cookies = jar.data.Cookies(&url.URL{Scheme: uri.Scheme, Host: uri.Host})
	if len(cookies) == 0 {
		jar.log.Errorf("COOKIES ARE EMPTY! %+v %v", uri, cookies)
		jar.log.Errorf("%+v", jar)
		return nil
	}

	var cleaned = make([]*http.Cookie, 0, len(cookies))
	for i, cookie := range cookies {
		if _, ok := jar.nameMap[cookie.Name]; ok {
			cookies[i].Name = jar.nameMap[cookie.Name]
		}
		cleaned = append(cleaned, cookies[i])
	}

	return cleaned
}

func (jar *cookieContainer) defaultErrorHandler(err error) bool {
	if jar.errchan != nil {
		go func() {
			jar.errchan <- err
		}()
	}
	return true
}
