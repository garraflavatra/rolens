package app

import (
	"fmt"

	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
)

type OpenCollectionResult struct {
	Stats      bson.M `json:"stats"`
	StatsError string `json:"statsError"`
}

func (a *App) OpenCollection(hostKey, dbKey, collKey string) (result OpenCollectionResult) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return
	}
	defer close()

	command := bson.M{"collStats": collKey}
	err = client.Database(dbKey).RunCommand(ctx, command).Decode(&result.Stats)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not retrieve collection stats for "+collKey)
		runtime.LogWarning(a.ctx, err.Error())
		result.StatsError = err.Error()
	}

	return
}

func (a *App) RenameCollection(hostKey, dbKey, collKey, newCollKey string) bool {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	var result bson.M
	command := bson.D{
		bson.E{Key: "renameCollection", Value: fmt.Sprintf("%v.%v", dbKey, collKey)},
		bson.E{Key: "to", Value: fmt.Sprintf("%v.%v", dbKey, newCollKey)},
	}
	err = client.Database("admin").RunCommand(ctx, command).Decode(&result)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not rename collection "+collKey)
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while renaming collection"), zenity.ErrorIcon)
		return false
	}

	return true
}

func (a *App) TruncateCollection(hostKey, dbKey, collKey string) bool {
	err := zenity.Question("Are you sure you want to remove all items from "+collKey+"?", zenity.Title("Confirm"), zenity.WarningIcon)
	if err == zenity.ErrCanceled {
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	_, err = client.Database(dbKey).Collection(collKey).DeleteMany(ctx, bson.D{})
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not truncate collection "+collKey)
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while truncating collection"), zenity.ErrorIcon)
		return false
	}

	return true
}

func (a *App) DropCollection(hostKey, dbKey, collKey string) bool {
	err := zenity.Question("Are you sure you want to drop "+collKey+"?", zenity.Title("Confirm"), zenity.WarningIcon)
	if err == zenity.ErrCanceled {
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	err = client.Database(dbKey).Collection(collKey).Drop(ctx)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not drop collection "+collKey)
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while dropping collection"), zenity.ErrorIcon)
		return false
	}

	return true
}
