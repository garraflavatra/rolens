package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	goRuntime "runtime"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Settings struct {
	DefaultLimit    int64  `json:"defaultLimit"`
	DefaultSort     string `json:"defaultSort"`
	AutosubmitQuery bool   `json:"autosubmitQuery"`
}

func NewSettings() Settings {
	return Settings{
		DefaultLimit:    20,
		DefaultSort:     `{ "_id": 1 }`,
		AutosubmitQuery: true,
	}
}

func appDataDirectory() (string, error) {
	var err error
	homeDir, err := os.UserHomeDir()
	prefDir := ""

	switch goRuntime.GOOS {
	case "windows":
		prefDir = filepath.Join(homeDir, "/AppData/Local/Rolens")
	case "darwin":
		prefDir = filepath.Join(homeDir, "/Library/Application Support/Rolens")
	case "linux":
		prefDir = filepath.Join(homeDir, "/.config/Rolens")
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

func (a *App) Settings() Settings {
	s := NewSettings()
	filePath, err := appDataFilePath("settings.json")
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not retrieve application settings, using defaults!",
			Message: err.Error(),
		})
		return s
	}

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		// It's ok if the file cannot be opened, for example if it is not accessible.
		// Therefore no error is returned.
		fmt.Println(err.Error())
		return s
	}

	if len(jsonData) == 0 {
		return s
	} else {
		err = json.Unmarshal(jsonData, &s)

		if err != nil {
			fmt.Println(err.Error())
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Title:   "Could not retrieve application settings, using defaults!",
				Message: err.Error(),
			})
		}
		return s
	}
}

func (a *App) UpdateSettings(jsonData string) Settings {
	s := a.Settings()
	err := json.Unmarshal([]byte(jsonData), &s)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Malformed JSON",
			Message: err.Error(),
		})
		return s
	}

	filePath, err := appDataFilePath("settings.json")
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Could not update settings.json",
			Message: err.Error(),
		})
		return s
	}

	newJson, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Could not marshal settings into JSON",
			Message: err.Error(),
		})
		return s
	}

	err = ioutil.WriteFile(filePath, newJson, os.ModePerm)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Could not update host list",
			Message: err.Error(),
		})
	}

	return s
}
