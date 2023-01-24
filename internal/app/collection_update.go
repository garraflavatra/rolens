package app

import (
	"encoding/json"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

func (a *App) UpdateItems(hostKey, dbKey, collKey string, formJson string) int64 {
	var form struct {
		Upsert     bool   `json:"upsert"`
		Many       bool   `json:"many"`
		Query      string `json:"query"`
		Parameters []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"parameters"`
	}

	err := json.Unmarshal([]byte(formJson), &form)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Couldn't parse form",
			Message: err.Error(),
		})
		return 0
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	defer close()
	var query bson.M
	update := bson.M{}

	err = bson.UnmarshalExtJSON([]byte(form.Query), true, &query)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Invalid query",
			Message: err.Error(),
		})
		return 0
	}

	for _, param := range form.Parameters {
		var unmarshalled bson.M
		err = json.Unmarshal([]byte(param.Value), &unmarshalled)
		if err == nil {
			update[param.Type] = unmarshalled
		} else {
			fmt.Println(err.Error())
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Title:   "Invalid query",
				Message: err.Error(),
			})
			return 0
		}
	}

	var result *mongo.UpdateResult
	options := mongoOptions.Update().SetUpsert(form.Upsert)

	if form.Many {
		result, err = client.Database(dbKey).Collection(collKey).UpdateMany(ctx, query, update, options)
	} else {
		result, err = client.Database(dbKey).Collection(collKey).UpdateOne(ctx, query, update, options)
	}

	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while updating items",
			Message: err.Error(),
		})
		return 0
	}

	return result.ModifiedCount
}
