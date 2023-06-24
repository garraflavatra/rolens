package app

import (
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
		runtime.LogWarningf(a.ctx, "Could not retrieve database stats for %s: %s", dbKey, err.Error())
		result.StatsError = err.Error()
	}

	result.Collections, err = client.Database(dbKey).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not retrieve collection list for db %s: %s", dbKey, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error getting collection list",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
	}

	return
}

func (a *App) DropDatabase(hostKey, dbKey string) bool {
	choice, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to drop " + dbKey + "?",
		Buttons:       []string{"Yes", "Cancel"},
		DefaultButton: "Yes",
		CancelButton:  "Cancel",
	})
	if choice != "Yes" {
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	err = client.Database(dbKey).Drop(ctx)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not drop db %s: %s", dbKey, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error dropping database",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return true
}
