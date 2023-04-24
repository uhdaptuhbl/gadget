package settings

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/*
DefaultConfigDir determines where the program files should be loaded/stored.

NOTE: this could also be replaced by configdir but that means using
ApplicationSupport on OSX which is more cumbersome than just `~/.config`:
https://github.com/kirsle/configdir
https://specifications.freedesktop.org/basedir-spec/basedir-spec-0.8.html
*/
func DefaultConfigDir(appName string) string {
	var err error
	var dir = filepath.Join("/opt", appName)

	if runtime.GOOS == "darwin" {
		// NOTE: os.UserConfigDir() will generally give a location that
		// is in ApplicationSupport but we want to manually opt for ~/.config
		// NOTE: os.UserHomeDir() can technically error, but the
		// default is returned instead if it actually does.
		if dir, err = os.UserHomeDir(); err == nil {
			dir = filepath.Join(dir, ".config", appName)
		}
	}

	return dir
}

type NamespaceError struct {
	Namespace string
	Problem string
}

func (e *NamespaceError) Error() string {
	var msg strings.Builder

	if e.Namespace != "" {
		msg.WriteString("namespace: " + e.Namespace)
	}
	if e.Problem != "" {
		if e.Namespace != "" {
			msg.WriteString("; ")
		}
		msg.WriteString(e.Problem)
	}

	return msg.String()
}

func GetUserDirs(namespace string) (UserDirs, error) {
	var err error
	var dirs UserDirs

	if namespace == "" {
		return dirs, &NamespaceError{Problem: "empty namespace value"}
	}

	dirs.Namespace = namespace
	if dirs.home, err = os.UserHomeDir(); err != nil {
		return dirs, err
	}
	if dirs.cache, err = os.UserCacheDir(); err != nil {
		return dirs, err
	}
	if dirs.config, err = os.UserConfigDir(); err != nil {
		return dirs, err
	}

	return dirs, err
}

type UserDirs struct {
	Namespace string

	home string
	cache string
	config string
}

func (dirs UserDirs) Home() string {
	if dirs.Namespace != "" {
		return filepath.Join(dirs.home, dirs.Namespace)
	}
	return dirs.home
}

func (dirs UserDirs) Cache() string {
	if dirs.Namespace != "" {
		return filepath.Join(dirs.cache, dirs.Namespace)
	}
	return dirs.cache
}

func (dirs UserDirs) Config() string {
	if dirs.Namespace != "" {
		return filepath.Join(dirs.config, dirs.Namespace)
	}
	return dirs.config
}

/*
WriteDefaultConfigFile creates the path and writes a default config file to disk.
*/
// func WriteDefaultConfigFile(app string, name string, content []byte, dirAccess os.FileMode, fileAccess os.FileMode) (string, error) {
// 	// TODO: make this more generic where you can pass viper and struct and it will do just loading
// 	var err error

// 	if app == "" {
// 		return "", &InvalidValueError{Label: "app", Value: app, Expected: "non-empty string"}
// 	}
// 	if name == "" {
// 		return "", &InvalidValueError{Label: "name", Value: name, Expected: "non-empty string"}
// 	}

// 	var dir = DefaultConfigDir(app)
// 	var fpath = filepath.Clean(filepath.Join(dir, name))
// 	if err = os.MkdirAll(dir, dirAccess); err != nil {
// 		return fpath, err
// 	}

// 	var f *os.File
// 	var nbytes int
// 	// TODO: why is this line causing a linter error?
// 	if f, err = os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fileAccess); err != nil {
// 		return "", err
// 	}
// 	defer func() {
// 		if errClose := f.Close(); errClose != nil {
// 			err = errClose
// 		}
// 	}()
// 	if nbytes, err = f.Write(content); err != nil {
// 		return fpath, err
// 	}
// 	if nbytes < 1 {
// 		return fpath, fmt.Errorf("incomplete data written to file")
// 	}

// 	return fpath, err
// }
