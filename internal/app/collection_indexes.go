package app

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *App) GetIndexes(hostKey, dbKey, collKey string) []bson.M {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer close()

	cur, err := client.Database(dbKey).Collection(collKey).Indexes().List(ctx)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while creating index cursor",
			Message: err.Error(),
		})
		return nil
	}

	var results []bson.M
	err = cur.All(ctx, &results)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while executing index cursor",
			Message: err.Error(),
		})
		return nil
	}

	return results
}
