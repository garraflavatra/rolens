package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *App) OpenDatabase(hostKey, dbKey string) (collections []string) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return nil
	}

	collections, err = client.Database(dbKey).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not retrieve collection list for db "+dbKey)
		runtime.LogWarning(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not retrieve collection list for " + dbKey,
			Message: err.Error(),
		})
		return nil
	}

	defer close()
	return collections
}

func (a *App) DropDatabase(hostKey, dbKey string) bool {
	sure, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to drop " + dbKey + "?",
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

	err = client.Database(dbKey).Drop(ctx)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not drop db "+dbKey)
		runtime.LogWarning(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not drop " + dbKey,
			Message: err.Error(),
		})
		return false
	}

	defer close()
	return true
}
