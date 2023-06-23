package app

import (
	"strings"

	"github.com/ncruces/zenity"
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
		zenity.Error(err.Error(), zenity.Title("Could not parse JSON"), zenity.ErrorIcon)
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
		zenity.Error(err.Error(), zenity.Title("Error while performing insert"), zenity.ErrorIcon)
		return nil
	}

	return res.InsertedIDs
}
