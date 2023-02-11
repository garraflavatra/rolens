package app

import (
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *App) InsertItems(hostKey, dbKey, collKey, jsonData string) interface{} {
	var data []interface{}

	jsonData = strings.TrimSpace(jsonData)
	if strings.HasPrefix(jsonData, "{") {
		jsonData = "[" + jsonData + "]"
	}

	err := bson.UnmarshalExtJSON([]byte(jsonData), true, &data)
	if err != nil {
		runtime.LogError(a.ctx, "Could not parse insert JSON:")
		runtime.LogError(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Couldn't parse JSON",
			Message: err.Error(),
		})
		return nil
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return nil
	}

	defer close()
	res, err := client.Database(dbKey).Collection(collKey).InsertMany(ctx, data)
	if err != nil {
		runtime.LogWarning(a.ctx, "Encountered an error while performing insert:")
		runtime.LogWarning(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while performing query",
			Message: err.Error(),
		})
		return nil
	}

	return res.InsertedIDs
}
