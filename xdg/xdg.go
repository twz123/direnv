// Package xdg is a minimal implementation of the XDG specification.
//
// https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html
package xdg

import (
	"path/filepath"
)

// DataDir returns the data folder for the application
func DataDir(env map[string]string, programName string) string {
	if dataHome := env["XDG_DATA_HOME"]; dataHome != "" {
		return filepath.Join(dataHome, programName)
	} else if home := env["HOME"]; home != "" {
		return filepath.Join(home, ".local", "share", programName)
	}
	// In theory we could also read /etc/passwd and look for the home based on
	// the process' UID
	return ""
}

// ConfigDir returns the config folder for the application
//
// The XDG_CONFIG_DIRS case is not being handled
func ConfigDir(env map[string]string, programName string) string {
	if configHome := env["XDG_CONFIG_HOME"]; configHome != "" {
		return filepath.Join(configHome, programName)
	} else if home := env["HOME"]; home != "" {
		return filepath.Join(home, ".config", programName)
	}
	// In theory we could also read /etc/passwd and look for the home based on
	// the process' UID
	return ""
}

// CacheDir returns the cache directory for the application
func CacheDir(env map[string]string, programName string) string {
	if cacheHome := env["XDG_CACHE_HOME"]; cacheHome != "" {
		return filepath.Join(cacheHome, programName)
	} else if home := env["HOME"]; home != "" {
		return filepath.Join(home, ".cache", programName)
	}
	// In theory we could also read /etc/passwd and look for the home based on
	// the process' UID
	return ""
}
