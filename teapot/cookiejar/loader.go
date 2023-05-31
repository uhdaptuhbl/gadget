package cookiejar

import (
	"net/http"
	"net/url"
)

// Loader
//
// TODO: should this be intended to load into provided jar or to keep the jar on itself?
// TODO: what about loading to a slice of jars instead of just one?
type Loader interface {
	SetJar(jar http.CookieJar) Loader
	ToJar(jar http.CookieJar, key ...*url.URL)
	Jar() http.CookieJar
	Load(hosts ...string) error
}
