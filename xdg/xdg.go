// Package xdg is a minimal implementation of the XDG specification.
//
// https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html
package xdg

import "path/filepath"

// Env is the subset of environment access needed by xdg helpers.
type Env interface {
	Get(string) string
}

// DataDir returns the data folder for the application
func DataDir(env Env, programName string) string {
	if dataHome := env.Get("XDG_DATA_HOME"); dataHome != "" {
		return filepath.Join(dataHome, programName)
	} else if home := env.Get("HOME"); home != "" {
		return filepath.Join(home, ".local", "share", programName)
	}
	// In theory we could also read /etc/passwd and look for the home based on
	// the process' UID
	return ""
}

// ConfigDir returns the config folder for the application
//
// The XDG_CONFIG_DIRS case is not being handled
func ConfigDir(env Env, programName string) string {
	if configHome := env.Get("XDG_CONFIG_HOME"); configHome != "" {
		return filepath.Join(configHome, programName)
	} else if home := env.Get("HOME"); home != "" {
		return filepath.Join(home, ".config", programName)
	}
	// In theory we could also read /etc/passwd and look for the home based on
	// the process' UID
	return ""
}

// CacheDir returns the cache directory for the application
func CacheDir(env Env, programName string) string {
	if cacheHome := env.Get("XDG_CACHE_HOME"); cacheHome != "" {
		return filepath.Join(cacheHome, programName)
	} else if home := env.Get("HOME"); home != "" {
		return filepath.Join(home, ".cache", programName)
	}
	// In theory we could also read /etc/passwd and look for the home based on
	// the process' UID
	return ""
}
