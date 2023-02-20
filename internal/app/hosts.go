package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/ncruces/zenity"
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
		runtime.LogInfo(a.ctx, "Could not open hosts.json")
		runtime.LogInfo(a.ctx, err.Error())
		return make(map[string]Host, 0), nil
	}

	if len(jsonData) == 0 {
		return make(map[string]Host, 0), nil
	} else {
		var hosts map[string]Host
		err = json.Unmarshal(jsonData, &hosts)

		if err != nil {
			runtime.LogInfo(a.ctx, "host.json file contains malformatted JSON data")
			runtime.LogInfo(a.ctx, err.Error())
			return nil, errors.New("host.json file contains malformatted JSON data")
		}
		return hosts, nil
	}
}

func (a *App) AddHost(jsonData string) error {
	hosts, err := a.Hosts()
	if err != nil {
		zenity.Info(err.Error(), zenity.Title("Error while retrieving hosts"), zenity.ErrorIcon)
		return errors.New("could not retrieve existing host list")
	}

	var newHost Host
	err = json.Unmarshal([]byte(jsonData), &newHost)
	if err != nil {
		runtime.LogError(a.ctx, "Add host: malformed form")
		runtime.LogError(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Could not parse JSON"), zenity.ErrorIcon)
		return errors.New("invalid JSON")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		runtime.LogError(a.ctx, "Add host: failed to generate a UUID")
		runtime.LogError(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Error while generating UUID"), zenity.ErrorIcon)
		return errors.New("could not generate a UUID")
	}

	hosts[id.String()] = newHost
	err = updateHostsFile(a, hosts)
	if err != nil {
		zenity.Info(err.Error(), zenity.Title("Error while updating host list"), zenity.ErrorIcon)
		return errors.New("could not update host list")
	}

	return nil
}

func (a *App) UpdateHost(hostKey string, jsonData string) error {
	hosts, err := a.Hosts()
	if err != nil {
		zenity.Info(err.Error(), zenity.Title("Error while getting hosts"), zenity.ErrorIcon)
		return errors.New("could not retrieve existing host list")
	}

	var host Host
	err = json.Unmarshal([]byte(jsonData), &host)
	if err != nil {
		runtime.LogError(a.ctx, "Could not parse update host JSON")
		runtime.LogError(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Could not parse JSON"), zenity.ErrorIcon)
		return errors.New("invalid JSON")
	}

	hosts[hostKey] = host
	err = updateHostsFile(a, hosts)
	if err != nil {
		zenity.Info(err.Error(), zenity.Title("Error while updating hosts"), zenity.ErrorIcon)
		return errors.New("could not update host list")
	}

	return nil
}

func (a *App) RemoveHost(key string) error {
	hosts, err := a.Hosts()
	if err != nil {
		zenity.Info(err.Error(), zenity.Title("Error while retrieving hosts"), zenity.ErrorIcon)
		return errors.New("could not retrieve existing host list")
	}

	err = zenity.Question("Are you sure you want to remove "+hosts[key].Name+"?", zenity.Title("Confirm"), zenity.WarningIcon)
	if err == zenity.ErrCanceled {
		return errors.New("operation aborted")
	}

	delete(hosts, key)
	err = updateHostsFile(a, hosts)

	if err != nil {
		zenity.Info(err.Error(), zenity.Title("Error while updating hosts"), zenity.ErrorIcon)
		return errors.New("could not update host list")
	}
	return nil
}
