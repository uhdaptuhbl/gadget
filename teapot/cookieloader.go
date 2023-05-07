package teapot

import (
	"net/http"
	"net/url"
)

// CookieLoader
type CookieLoader interface {
	SetJar(jar http.CookieJar) CookieLoader
	ToJar(jar http.CookieJar, key ...*url.URL)
	Jar() http.CookieJar
	Load(hosts ...string) error
}
