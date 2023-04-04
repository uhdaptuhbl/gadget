package config

import (
	"os"
	"path/filepath"
	"runtime"
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
