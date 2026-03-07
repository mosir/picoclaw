package fileutil

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// HomeDir resolves the current user's home directory consistently across platforms.
// Environment variables are checked first so tests and sandboxed environments can
// override home resolution deterministically.
func HomeDir() (string, error) {
	if home := strings.TrimSpace(os.Getenv("HOME")); home != "" {
		return home, nil
	}

	if runtime.GOOS == "windows" {
		if profile := strings.TrimSpace(os.Getenv("USERPROFILE")); profile != "" {
			return profile, nil
		}
		drive := strings.TrimSpace(os.Getenv("HOMEDRIVE"))
		path := strings.TrimSpace(os.Getenv("HOMEPATH"))
		if drive != "" && path != "" {
			return drive + path, nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	home = strings.TrimSpace(home)
	if home == "" {
		return "", errors.New("home directory is empty")
	}
	return home, nil
}

// ExpandHome expands "~", "~/" and "~\\" prefixes to the current home directory.
// Paths like "~otheruser/..." are returned unchanged.
func ExpandHome(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}

	home, err := HomeDir()
	if err != nil || home == "" {
		return path
	}

	if len(path) == 1 {
		return home
	}

	if path[1] != '/' && path[1] != '\\' {
		return path
	}

	rest := strings.TrimLeft(path[1:], `/\`)
	if rest == "" {
		return home
	}
	return filepath.Join(home, rest)
}
