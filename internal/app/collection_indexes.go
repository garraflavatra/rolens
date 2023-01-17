package app

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (a *App) DropIndex(hostKey, dbKey, collKey, indexName string) bool {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer close()

	_, err = client.Database(dbKey).Collection(collKey).Indexes().DropOne(ctx, indexName, &options.DropIndexesOptions{})
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while creating index cursor",
			Message: err.Error(),
		})
		return false
	}

	return true
}
