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
		runtime.LogErrorf(a.ctx, "Could not parse insert JSON: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Malformed JSON",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
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
		runtime.LogWarningf(a.ctx, "Encountered an error while performing insert: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error performing insert",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return nil
	}

	return res.InsertedIDs
}
