package firefox

import (
	"net/http"
	"net/url"
)

type EmptySqlitePathError struct {}

func (e *EmptySqlitePathError) Error() string {
	return "Empty path to sqlite Firefox cookie database"
}

// FirefoxSQLiteCookie
//
// https://firefox-source-docs.mozilla.org/devtools-user/storage_inspector/cookies/index.html
// http://fileformats.archiveteam.org/wiki/Firefox_cookie_database
type FirefoxSQLiteCookie struct {
	ID               int64  `gorm:"column:id:primaryKey"`
	OriginAttributes string `gorm:"column:originAttributes"`
	Name             string `gorm:"column:name"`
	Value            string `gorm:"column:value"`
	Host             string `gorm:"column:host"`
	Path             string `gorm:"column:path"`
	Expiry           int64  `gorm:"column:expiry"`
	LastAccessed     int64  `gorm:"column:lastAccessed"`
	CreationTime     int64  `gorm:"column:creationTime"`
	Secure           bool   `gorm:"column:isSecure"`
	HttpOnly         bool   `gorm:"column:isHttpOnly"`
	InBrowserElement int64  `gorm:"column:inBrowserElement"`
	SameSite         int64  `gorm:"column:sameSite"`
	RawSameSite      int64  `gorm:"column:rawSameSite"`
	SchemeMap        int64  `gorm:"column:schemeMap"`
}

type cookieWrapper struct {
	host       *url.URL
	ffxcookie  *FirefoxSQLiteCookie
	httpcookie *http.Cookie
	valid      error
}

type wrapperSlice []cookieWrapper
