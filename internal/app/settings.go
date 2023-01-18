package app

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func appDataDirectory() (string, error) {
	var err error
	homeDir, err := os.UserHomeDir()
	prefDir := ""

	switch runtime.GOOS {
	case "windows":
		prefDir = filepath.Join(homeDir, "/AppData/Local/Mongodup")
	case "darwin":
		prefDir = filepath.Join(homeDir, "/Library/Application Support/Mongodup")
	case "linux":
		prefDir = filepath.Join(homeDir, "/.config/Mongodup")
	default:
		err = errors.New("unsupported platform")
	}

	_ = os.MkdirAll(prefDir, os.ModePerm)
	return prefDir, err
}

func appDataFilePath(filename string) (string, error) {
	dir, err := appDataDirectory()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, filename)
	return path, nil
}
