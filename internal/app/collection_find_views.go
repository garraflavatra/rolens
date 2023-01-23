package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ViewType string

const (
	TableView ViewType = "table"
	ListView  ViewType = "list"
)

type ViewColumn struct {
	Key   string `json:"key"`
	Width int64  `json:"width"`
}

type View struct {
	Name                 string       `json:"name"`
	Host                 string       `json:"host"`
	Database             string       `json:"database"`
	Collection           string       `json:"collection"`
	Type                 ViewType     `json:"type"`
	HideObjectIndicators bool         `json:"hideObjectIndicators"`
	Columns              []ViewColumn `json:"columns"`
}

var BuiltInListView = View{
	Name: "List",
	Type: ListView,
}

type ViewStore map[string]View

func updateViewStore(newData ViewStore) error {
	filePath, err := appDataFilePath("views.json")
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

func (a *App) Views() (ViewStore, error) {
	views := make(ViewStore, 0)
	filePath, err := appDataFilePath("views.json")
	if err != nil {
		fmt.Println(err.Error())
		return views, err
	}

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		// It's ok if the file cannot be opened, for example if it is not accessible.
		// Therefore no error is returned.
		fmt.Println(err.Error())
		return views, nil
	}

	if len(jsonData) > 0 {
		err = json.Unmarshal(jsonData, &views)
		if err != nil {
			fmt.Println(err.Error())
			return nil, errors.New("views.json file contains malformatted JSON data")
		}
	}

	views["list"] = BuiltInListView
	return views, nil
}

// func (a *App) AddView(jsonData string) error {
// 	views, err := a.Views()
// 	if err != nil {
// 		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
// 			Type:  runtime.InfoDialog,
// 			Title: "Could not retrieve views",
// 		})
// 		return errors.New("could not retrieve existing view store")
// 	}

// 	var newView View
// 	err = json.Unmarshal([]byte(jsonData), &newView)
// 	if err != nil {
// 		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
// 			Type:  runtime.InfoDialog,
// 			Title: "Malformed JSON",
// 		})
// 		return errors.New("invalid JSON")
// 	}

// 	id, err := uuid.NewRandom()
// 	if err != nil {
// 		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
// 			Type:  runtime.InfoDialog,
// 			Title: "Failed to generate a UUID",
// 		})
// 		return errors.New("could not generate a UUID")
// 	}

// 	views[id.String()] = newView
// 	err = updateViewStore(views)
// 	if err != nil {
// 		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
// 			Type:  runtime.InfoDialog,
// 			Title: "Could not update view store",
// 		})
// 		return errors.New("could not update view store")
// 	}

// 	return nil
// }

func (a *App) UpdateViewStore(jsonData string) error {
	var viewStore ViewStore
	err := json.Unmarshal([]byte(jsonData), &viewStore)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Malformed JSON",
		})
		return errors.New("invalid JSON")
	}

	err = updateViewStore(viewStore)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not update view store",
		})
		return errors.New("could not update view store")
	}

	return nil
}

func (a *App) RemoveView(viewKey string) error {
	views, err := a.Views()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not retrieve views",
		})
		return errors.New("could not retrieve existing view store")
	}

	sure, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to remove " + views[viewKey].Name + "?",
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "Yes",
		CancelButton:  "No",
	})
	if sure != "Yes" {
		return errors.New("operation aborted")
	}

	delete(views, viewKey)
	err = updateViewStore(views)

	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not update view store",
		})
		return errors.New("could not update view store")
	}
	return nil
}
