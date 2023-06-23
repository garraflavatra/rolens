package app

import (
	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
)

type OpenDatabaseResult struct {
	Collections []string `json:"collections"`
	Stats       bson.M   `json:"stats"`
	StatsError  string   `json:"statsError"`
}

func (a *App) OpenDatabase(hostKey, dbKey string) (result OpenDatabaseResult) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return
	}
	defer close()

	command := bson.M{"dbStats": 1}
	err = client.Database(dbKey).RunCommand(ctx, command).Decode(&result.Stats)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not retrieve database stats for "+dbKey)
		runtime.LogWarning(a.ctx, err.Error())
		result.StatsError = err.Error()
	}

	result.Collections, err = client.Database(dbKey).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not retrieve collection list for db "+dbKey)
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while getting collections"), zenity.ErrorIcon)
	}

	return
}

func (a *App) DropDatabase(hostKey, dbKey string) bool {
	err := zenity.Question("Are you sure you want to drop "+dbKey+"?", zenity.Title("Confirm"), zenity.WarningIcon)
	if err == zenity.ErrCanceled {
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	err = client.Database(dbKey).Drop(ctx)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not drop db "+dbKey)
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while dropping database"), zenity.ErrorIcon)
		return false
	}

	return true
}
