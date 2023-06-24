package app

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Settings struct {
	DefaultLimit           int64  `json:"defaultLimit"`
	DefaultSort            string `json:"defaultSort"`
	AutosubmitQuery        bool   `json:"autosubmitQuery"`
	DefaultExportDirectory string `json:"defaultExportDirectory"`
}

func NewSettings(a *App) Settings {
	return Settings{
		DefaultLimit:           20,
		DefaultSort:            `{ "_id": 1 }`,
		AutosubmitQuery:        true,
		DefaultExportDirectory: path.Join(a.Env.HomeDirectory, "Downloads"),
	}
}

func (a *App) Settings() Settings {
	s := NewSettings(a)
	filePath := path.Join(a.Env.DataDirectory, "settings.json")

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		// It's ok if the file cannot be opened, for example if it is not accessible.
		// Therefore no error is returned.
		runtime.LogInfof(a.ctx, "Cannot open settings.json: %s", err.Error())
		return s
	}

	if len(jsonData) == 0 {
		return s
	} else {
		if err := json.Unmarshal(jsonData, &s); err != nil {
			runtime.LogWarningf(a.ctx, "Cannot unmarshal settings.json: %s", err.Error())
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Title:   "Settings malformed",
				Message: "Could not retrieve application settings: using defaults!",
				Type:    runtime.WarningDialog,
			})
		}
		return s
	}
}

func (a *App) UpdateSettings(jsonData string) Settings {
	s := a.Settings()
	err := json.Unmarshal([]byte(jsonData), &s)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Malformed JSON for settings file: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Settings malformed",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return s
	}

	newJson, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		runtime.LogErrorf(a.ctx, "Could not marshal settings into JSON: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "JSON is being awkward",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return s
	}

	filePath := path.Join(a.Env.DataDirectory, "settings.json")
	err = ioutil.WriteFile(filePath, newJson, os.ModePerm)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Could not update host list: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Could not update host list",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
	}

	return s
}
