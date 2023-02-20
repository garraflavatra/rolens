package app

import (
	"github.com/ncruces/zenity"
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
		zenity.Info(err.Error(), zenity.Title("Error while getting collections"), zenity.ErrorIcon)
		return nil
	}

	defer close()
	return collections
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

	err = client.Database(dbKey).Drop(ctx)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not drop db "+dbKey)
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Error while dropping database"), zenity.ErrorIcon)
		return false
	}

	defer close()
	return true
}
