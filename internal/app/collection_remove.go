package app

import (
	"strings"

	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a *App) RemoveItems(hostKey, dbKey, collKey, jsonData string, many bool) int64 {
	var filter bson.M
	var err error
	jsonData = strings.TrimSpace(jsonData)

	if len(jsonData) == 0 {
		err := zenity.Question("Are you sure you want to drop all items in "+collKey+"?", zenity.Title("Confirm"), zenity.WarningIcon)
		if err == zenity.ErrCanceled {
			return 0
		}
	} else {
		err = bson.UnmarshalExtJSON([]byte(jsonData), true, &filter)
		if err != nil {
			runtime.LogError(a.ctx, "Could not parse remove query:")
			runtime.LogError(a.ctx, err.Error())
			zenity.Info(err.Error(), zenity.Title("Could not parse JSON"), zenity.ErrorIcon)
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
		runtime.LogWarning(a.ctx, "Encountered an error while performing remove:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Error while performing remove"), zenity.ErrorIcon)
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
		runtime.LogWarning(a.ctx, "Encountered an error while performing remove by id:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Error while performing remove"), zenity.ErrorIcon)

		return false
	}

	return err == nil
}
