package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Host struct {
	Name      string `json:"name"`
	URI       string `json:"uri"`
	Databases map[string]struct {
		Collections map[string]struct {
			ViewConfig struct {
				HideObjectIndicators bool `json:"hideObjectIndicators"`
				Columns              []struct {
					Key   string `json:"key"`
					Width int64  `json:"width"`
				} `json:"columns"`
			} `json:"viewConfig"`
		} `json:"collections"`
	} `json:"databases"`
}

func updateHostsFile(newData map[string]Host) error {
	filePath, err := appDataFilePath("hosts.json")
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(newData, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, jsonData, os.ModePerm)
	return err
}

func (a *App) Hosts() (map[string]Host, error) {
	filePath, err := appDataFilePath("hosts.json")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		// It's ok if the file cannot be opened, for example if it is not accessible.
		// Therefore no error is returned.
		fmt.Println(err.Error())
		return make(map[string]Host, 0), nil
	}

	if len(jsonData) == 0 {
		return make(map[string]Host, 0), nil
	} else {
		var hosts map[string]Host
		err = json.Unmarshal(jsonData, &hosts)

		if err != nil {
			fmt.Println(err.Error())
			return nil, errors.New("host.json file contains malformatted JSON data")
		}
		return hosts, nil
	}
}

func (a *App) AddHost(jsonData string) error {
	hosts, err := a.Hosts()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not retrieve hosts",
		})
		return errors.New("could not retrieve existing host list")
	}

	var newHost Host
	err = json.Unmarshal([]byte(jsonData), &newHost)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Malformed JSON",
		})
		return errors.New("invalid JSON")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Failed to generate a UUID",
		})
		return errors.New("could not generate a UUID")
	}

	hosts[id.String()] = newHost
	err = updateHostsFile(hosts)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not update host list",
		})
		return errors.New("could not update host list")
	}

	return nil
}

func (a *App) UpdateHost(hostKey string, jsonData string) error {
	hosts, err := a.Hosts()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not retrieve hosts",
		})
		return errors.New("could not retrieve existing host list")
	}

	var host Host
	err = json.Unmarshal([]byte(jsonData), &host)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Malformed JSON",
		})
		return errors.New("invalid JSON")
	}

	hosts[hostKey] = host
	err = updateHostsFile(hosts)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not update host list",
		})
		return errors.New("could not update host list")
	}

	return nil
}

func (a *App) RemoveHost(key string) error {
	hosts, err := a.Hosts()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not retrieve hosts",
		})
		return errors.New("could not retrieve existing host list")
	}

	sure, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to remove " + hosts[key].Name + "?",
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "Yes",
		CancelButton:  "No",
	})
	if sure != "Yes" {
		return errors.New("operation aborted")
	}

	delete(hosts, key)
	err = updateHostsFile(hosts)

	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not update host list",
		})
		return errors.New("could not update host list")
	}
	return nil
}
