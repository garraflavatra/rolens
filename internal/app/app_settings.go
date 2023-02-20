package app

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/ncruces/zenity"
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
		runtime.LogInfo(a.ctx, "Cannot open settings.json:")
		runtime.LogInfo(a.ctx, err.Error())
		return s
	}

	if len(jsonData) == 0 {
		return s
	} else {
		err = json.Unmarshal(jsonData, &s)

		if err != nil {
			runtime.LogWarning(a.ctx, "Cannot unmarshal settings.json:")
			runtime.LogWarning(a.ctx, err.Error())
			zenity.Info("Could not retrieve application settings, using defaults!", zenity.Title("Information"), zenity.WarningIcon)
		}
		return s
	}
}

func (a *App) UpdateSettings(jsonData string) Settings {
	s := a.Settings()
	err := json.Unmarshal([]byte(jsonData), &s)
	if err != nil {
		runtime.LogError(a.ctx, "Malformed JSON for settings file:")
		runtime.LogError(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Malformed JSON"), zenity.ErrorIcon)
		return s
	}

	newJson, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		runtime.LogError(a.ctx, "Could not marshal settings into JSON:")
		runtime.LogError(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Malformed JSON"), zenity.ErrorIcon)
		return s
	}

	filePath := path.Join(a.Env.DataDirectory, "settings.json")
	err = ioutil.WriteFile(filePath, newJson, os.ModePerm)
	if err != nil {
		runtime.LogError(a.ctx, "Could not update host list:")
		runtime.LogError(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Could not update host list"), zenity.ErrorIcon)
	}

	return s
}
