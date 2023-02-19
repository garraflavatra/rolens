package app

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *App) OpenCollection(hostKey, dbKey, collKey string) (result bson.M) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return nil
	}

	command := bson.M{"collStats": collKey}
	err = client.Database(dbKey).RunCommand(ctx, command).Decode(&result)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not retrieve collection stats for "+collKey)
		runtime.LogWarning(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not retrieve collection stats for " + collKey,
			Message: err.Error(),
		})
		return nil
	}

	defer close()
	return result
}

func (a *App) RenameCollection(hostKey, dbKey, collKey, newCollKey string) bool {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}

	var result bson.M
	command := bson.D{
		bson.E{Key: "renameCollection", Value: fmt.Sprintf("%v.%v", dbKey, collKey)},
		bson.E{Key: "to", Value: fmt.Sprintf("%v.%v", dbKey, newCollKey)},
	}
	err = client.Database("admin").RunCommand(ctx, command).Decode(&result)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not rename collection "+collKey)
		runtime.LogWarning(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not rename " + collKey,
			Message: err.Error(),
		})
		return false
	}

	defer close()
	return true
}

func (a *App) TruncateCollection(hostKey, dbKey, collKey string) bool {
	sure, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to remove all items from " + collKey + "?",
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "Yes",
		CancelButton:  "No",
		Type:          runtime.WarningDialog,
	})
	if sure != "Yes" {
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}

	_, err = client.Database(dbKey).Collection(collKey).DeleteMany(ctx, bson.D{})
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not truncate collection "+collKey)
		runtime.LogWarning(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not truncate " + collKey,
			Message: err.Error(),
		})
		return false
	}

	defer close()
	return true
}

func (a *App) DropCollection(hostKey, dbKey, collKey string) bool {
	sure, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to drop " + collKey + "?",
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "Yes",
		CancelButton:  "No",
		Type:          runtime.WarningDialog,
	})
	if sure != "Yes" {
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}

	err = client.Database(dbKey).Collection(collKey).Drop(ctx)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not drop collection "+collKey)
		runtime.LogWarning(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not drop " + collKey,
			Message: err.Error(),
		})
		return false
	}

	defer close()
	return true
}
