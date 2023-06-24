package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Host struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

func updateHostsFile(a *App, newData map[string]Host) error {
	filePath := path.Join(a.Env.DataDirectory, "hosts.json")
	jsonData, err := json.MarshalIndent(newData, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, jsonData, os.ModePerm)
	return err
}

func (a *App) Hosts() (map[string]Host, error) {
	filePath := path.Join(a.Env.DataDirectory, "hosts.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		// It's ok if the file cannot be opened, for example if it is not accessible.
		// Therefore no error is returned.
		runtime.LogInfof(a.ctx, "Could not open hosts.json (%s), trying to create it.", err.Error())
		return make(map[string]Host, 0), nil
	}

	if len(jsonData) == 0 {
		return make(map[string]Host, 0), nil
	} else {
		var hosts map[string]Host
		err = json.Unmarshal(jsonData, &hosts)

		if err != nil {
			runtime.LogInfof(a.ctx, "host.json file contains malformatted JSON data: %s", err.Error())
			return nil, errors.New("host.json file contains malformatted JSON data")
		}
		return hosts, nil
	}
}

func (a *App) AddHost(jsonData string) string {
	hosts, err := a.Hosts()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error getting hosts",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return ""
	}

	var newHost Host
	err = json.Unmarshal([]byte(jsonData), &newHost)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Add host: malformed form: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Malformed JSON",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return ""
	}

	id, err := uuid.NewRandom()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Add host: failed to generate a UUID: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error generating UUID",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return ""
	}

	hosts[id.String()] = newHost
	err = updateHostsFile(a, hosts)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error updating host list",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return ""
	}

	return id.String()
}

func (a *App) UpdateHost(hostKey string, jsonData string) bool {
	hosts, err := a.Hosts()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error getting host list",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	var host Host
	err = json.Unmarshal([]byte(jsonData), &host)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Could not parse update host JSON: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Malformed JSON",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	hosts[hostKey] = host
	err = updateHostsFile(a, hosts)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error updating hosts",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return true
}

func (a *App) RemoveHost(key string) bool {
	hosts, err := a.Hosts()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error getting host list",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	choice, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to remove " + hosts[key].Name + "?",
		Buttons:       []string{"Yes", "Cancel"},
		DefaultButton: "Yes",
		CancelButton:  "Cancel",
	})
	if choice != "Yes" {
		return false
	}

	delete(hosts, key)
	err = updateHostsFile(a, hosts)

	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error updating host list",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return true
}
