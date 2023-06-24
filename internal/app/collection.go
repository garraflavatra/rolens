package app

import (
	"fmt"

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
		runtime.LogWarningf(a.ctx, "Could not retrieve collection stats for %s: %s", collKey, err.Error())
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
		runtime.LogWarningf(a.ctx, "Could not rename collection %s: %s", collKey, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error renaming collection",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return true
}

func (a *App) TruncateCollection(hostKey, dbKey, collKey string) bool {
	choice, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to remove all items in " + collKey + "?",
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

	_, err = client.Database(dbKey).Collection(collKey).DeleteMany(ctx, bson.D{})
	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not truncate collection %s: %s", collKey, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error truncating collection",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return true
}

func (a *App) DropCollection(hostKey, dbKey, collKey string) bool {
	choice, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to drop " + collKey + "?",
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

	err = client.Database(dbKey).Collection(collKey).Drop(ctx)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not drop collection %s: %s", collKey, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error dropping collection",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return true
}
