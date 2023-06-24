package app

import (
	"encoding/json"

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
		runtime.LogErrorf(a.ctx, "Could not parse update form: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Malformed JSON",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return 0
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return 0
	}

	defer close()
	var query bson.M
	update := bson.M{}

	err = bson.UnmarshalExtJSON([]byte(form.Query), true, &query)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Invalid update query %v: %s", form.Query, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Invalid update query",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return 0
	}

	for _, param := range form.Parameters {
		var unmarshalled bson.M
		err = json.Unmarshal([]byte(param.Value), &unmarshalled)
		if err == nil {
			update[param.Type] = unmarshalled
		} else {
			runtime.LogWarningf(a.ctx, "Invalid update parameter value %v: %s", param.Value, err.Error())
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Title:   "Invalid update query",
				Message: err.Error(),
				Type:    runtime.ErrorDialog,
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
		runtime.LogWarningf(a.ctx, "Encountered an error while performing update: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error performing update query",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return 0
	}

	return result.ModifiedCount
}
