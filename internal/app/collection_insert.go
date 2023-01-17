package app

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) InsertItems(hostKey, dbKey, collKey, jsonData string) interface{} {
	var data []interface{}

	jsonData = strings.TrimSpace(jsonData)
	if strings.HasPrefix(jsonData, "{") {
		jsonData = "[" + jsonData + "]"
	}

	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Couldn't parse JSON",
			Message: err.Error(),
		})
		return nil
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	defer close()
	res, err := client.Database(dbKey).Collection(collKey).InsertMany(ctx, data)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while performing query",
			Message: err.Error(),
		})
		return nil
	}

	return res.InsertedIDs
}
