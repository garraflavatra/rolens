package app

import (
	"encoding/json"
	"fmt"
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
		sure, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:         "Confirm",
			Message:       "Are you sure you want to drop all items in " + collKey + "?",
			Buttons:       []string{"Yes", "No"},
			DefaultButton: "Yes",
			CancelButton:  "No",
		})
		if sure != "Yes" {
			return 0
		}
	} else {
		err = json.Unmarshal([]byte(jsonData), &filter)
		if err != nil {
			fmt.Println(err.Error())
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Title:   "Couldn't parse JSON",
				Message: err.Error(),
			})
			return 0
		}
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while performing query",
			Message: err.Error(),
		})
		return 0
	}

	return res.DeletedCount
}