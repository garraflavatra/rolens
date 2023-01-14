package app

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *App) OpenDatabase(hostKey, dbKey string) (collections []string) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	collections, err = client.Database(dbKey).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		fmt.Println(err.Error())
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
		Buttons:       []string{"yes", "no"},
		DefaultButton: "Yes",
		CancelButton:  "No",
	})
	if sure != "Yes" {
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	err = client.Database(dbKey).Drop(ctx)
	if err != nil {
		fmt.Println(err.Error())
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
