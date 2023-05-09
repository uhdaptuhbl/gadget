package firefox

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	// "github.com/glebarez/sqlite" // Pure go SQLite driver
	"gorm.io/gorm"

	"gadget/teapot/cookiejar"
)

// ffxcookies implements the gadget/teapot/cookiejar.Loader interface.
type ffxcookies struct {
	core  *firefoxCore
	jar   http.CookieJar
	table string
}

func NewCookieLoader(jar http.CookieJar, options ...Option) *ffxcookies {
	var loader = &ffxcookies{
		core:  newFirefoxCore(options...),
		jar:   jar,
		table: "moz_cookies",
	}

	return loader
}

func (loader *ffxcookies) SetJar(jar http.CookieJar) cookiejar.Loader {
	loader.jar = jar
	return loader
}

func (loader *ffxcookies) ToJar(jar http.CookieJar, keys ...*url.URL) {
	for _, key := range keys {
		jar.SetCookies(key, loader.jar.Cookies(key))
	}
}

func (loader *ffxcookies) Jar() http.CookieJar {
	return loader.jar
}

func (loader *ffxcookies) Load(hosts ...string) error {
	var err error
	var result strings.Builder
	var db *gorm.DB
	var query *gorm.DB
	var ffxcookies []FirefoxSQLiteCookie
	var cookieCount int

	if err = loader.core.valid(); err != nil {
		return err
	}
	if loader.jar == nil {
		return new(cookiejar.NilCookieJarError)
	}

	var log = loader.core.log
	var dbpath = loader.core.cookieDBPath()
	log.Debugf("loading Firefox cookies: %s", dbpath)

	if db, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{}); err != nil {
		return err
	}
	query = db.Table(loader.table)
	if len(hosts) != 0 {
		if len(hosts) == 1 {
			query.Where("host = '" + hosts[0] + "'")
		} else {
			var stmts = make([]string, 0, len(hosts))
			for _, host := range hosts {
				host = strings.TrimSpace(host)
				if host == "" {
					continue
				}
				stmts = append(stmts, "host = '"+host+"'")
			}
			query.Where(strings.Join(stmts, " OR "))
		}
	}
	if query = query.Find(&ffxcookies); query.Error != nil {
		return query.Error
	}
	if len(ffxcookies) == 0 {
		return &cookiejar.NoCookiesFoundError{Hosts: hosts}
	} else {
		log.Debugf("%d cookies found", len(ffxcookies))
	}

	var wrappedCookies = make(map[string]wrapperSlice)
	var wrappedCookie cookieWrapper
	for ffxindex, ffxcookie := range ffxcookies {
		cookieCount++

		// NOTE: cookiejar.Jar will ignore any cookies you try to set without a Scheme!
		var loc = &url.URL{Host: ffxcookie.Host}
		if ffxcookie.Secure {
			loc.Scheme = "https"
		} else {
			loc.Scheme = "http"
		}

		wrappedCookie = cookieWrapper{
			host:      loc,
			ffxcookie: &ffxcookies[ffxindex],
			httpcookie: &http.Cookie{
				Name:     ffxcookie.Name,
				Value:    ffxcookie.Value,
				Path:     ffxcookie.Path,
				Domain:   ffxcookie.Host,
				Expires:  time.Unix(ffxcookie.Expiry, 0),
				Secure:   ffxcookie.Secure,
				HttpOnly: ffxcookie.HttpOnly,
				SameSite: http.SameSite(ffxcookie.SameSite),
			},
			valid: nil,
		}
		if ffxcookie.Host == "" {
			return &cookiejar.CookieError{Message: fmt.Sprintf("no host on cookie: %+v", *wrappedCookie.httpcookie)}
		}

		wrappedCookies[ffxcookie.Host] = append(wrappedCookies[ffxcookie.Host], wrappedCookie)
	}
	if len(wrappedCookies) == 0 || cookieCount == 0 {
		return new(cookiejar.NoCookiesFoundError)
	} else {
		log.Debugf("%d cookies processed", cookieCount)
	}

	for host, cookieList := range wrappedCookies {
		if len(cookieList) == 0 {
			return &cookiejar.NoCookiesFoundError{Hosts: []string{host}}
		}

		var loc *url.URL
		var cookies = make([]*http.Cookie, 0, len(cookieList))
		for _, cookie := range cookieList {
			var cookiestr string

			// log.Debugf(
			// 	"(cookie:%d) loc: `%s`  |  name: `%s`  |  valid: `%v`",
			// 	index,
			// 	cookie.host.String(),
			// 	cookie.httpcookie.Name,
			// 	cookie.valid,
			// )

			// log.Debugf(
			// 	"(cookie:%d) loc: `%s`  |  name: `%s`  |  valid: `%v`  |  cookie: `%s`  |  raw: `%+v`",
			// 	index,
			// 	cookie.host.String(),
			// 	cookie.httpcookie.Name,
			// 	cookie.valid,
			// 	cookie.httpcookie.String(),
			// 	cookie.httpcookie,
			// )

			if cookie.valid != nil {
				return fmt.Errorf("INVALID COOKIE: %+v", cookie)
			}

			if cookie.host == nil {
				return fmt.Errorf("cookie with nil URL: %+v", cookie)
			}
			if cookie.host.Host == "" {
				return fmt.Errorf("cookie URL without Host: %+v", map[string]string{
					"Scheme":      cookie.host.Scheme,
					"Opaque":      cookie.host.Opaque,
					"Host":        cookie.host.Host,
					"Path":        cookie.host.Path,
					"RawPath":     cookie.host.RawPath,
					"RawQuery":    cookie.host.RawQuery,
					"Fragment":    cookie.host.Fragment,
					"RawFragment": cookie.host.RawFragment,
				})
			}
			loc = cookie.host
			cookiestr = cookie.httpcookie.String()
			// if cookiestr == "" {
			// 	// log.Debugf("cookiestr: %+v", cookiestr)
			// 	return fmt.Errorf("cookie = %+v", cookie)
			// }
			// if cookie.httpcookie.Raw != "" {
			// 	cookiestr = strings.ReplaceAll(cookiestr, cookie.httpcookie.Name, cookie.httpcookie.Raw)
			// 	// cookie.httpcookie.Raw = ""
			// }

			if result.Len() != 0 {
				result.WriteString("; ")
			}
			result.WriteString(cookiestr)

			cookies = append(cookies, cookie.httpcookie)
		}
		loader.jar.SetCookies(loc, cookies)

		if len(loader.jar.Cookies(loc)) == 0 {
			// log.Debugf("%s", loc)
			return fmt.Errorf("empty cookie jar after setting `%s`: %+v", loc.Host, cookies)
		}
	}
	// log.Debugf("loaded cookiejar: " + dbpath)
	log.Debugf("cookiejar loaded")

	return err
}
