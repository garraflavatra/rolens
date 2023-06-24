package app

import (
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a *App) RemoveItems(hostKey, dbKey, collKey, jsonData string, many bool) int64 {
	var filter bson.M
	var err error
	jsonData = strings.TrimSpace(jsonData)

	if len(jsonData) == 0 {
		choice, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:         "Confirm",
			Message:       "Are you sure you want to drop all items in " + collKey + "?",
			Buttons:       []string{"Yes", "Cancel"},
			DefaultButton: "Yes",
			CancelButton:  "Cancel",
		})
		if choice != "Yes" {
			return 0
		}
	} else {
		err = bson.UnmarshalExtJSON([]byte(jsonData), true, &filter)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Could not parse remove query: %s", err.Error())
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Title:   "Malformed JSON",
				Message: err.Error(),
				Type:    runtime.ErrorDialog,
			})
			return 0
		}
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return 0
	}

	defer close()

	var res *mongo.DeleteResult

	if many {
		res, err = client.Database(dbKey).Collection(collKey).DeleteMany(ctx, filter)
	} else {
		res, err = client.Database(dbKey).Collection(collKey).DeleteOne(ctx, filter)
	}

	if err != nil {
		runtime.LogWarningf(a.ctx, "Encountered an error while performing remove: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error performing remove query",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return 0
	}

	return res.DeletedCount
}

func (a *App) RemoveItemById(hostKey, dbKey, collKey, itemId string) bool {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}

	defer close()

	filter := bson.M{"_id": itemId}
	err = client.Database(dbKey).Collection(collKey).FindOneAndDelete(ctx, filter).Err()

	if err != nil && err != mongo.ErrNoDocuments {
		runtime.LogWarningf(a.ctx, "Encountered an error while performing remove by id: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error performing remove query",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return err == nil
}
