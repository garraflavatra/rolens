package app

import (
	"encoding/json"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *App) Aggregate(hostKey, dbKey, collKey, pipelineJson, settingsJson string) {
	var settings *options.AggregateOptions
	if err := json.Unmarshal([]byte(settingsJson), &settings); err != nil {
		runtime.LogErrorf(a.ctx, "Could not parse aggregation settings: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Couldn't parse aggregation settings",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return
	}

	var pipeline mongo.Pipeline
	if err := bson.UnmarshalExtJSON([]byte(pipelineJson), true, &pipeline); err != nil {
		runtime.LogWarningf(a.ctx, "Could not parse aggregation pipeline: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Couldn't parse aggregation pipeline",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return
	}

	defer close()

	cursor, err := client.Database(dbKey).Collection(collKey).Aggregate(ctx, pipeline, settings)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not get aggregation cursor: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Couldn't get aggregation cursor",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		runtime.LogInfof(a.ctx, "Error while running aggregation pipeline: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error while running aggregation pipeline",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return
	}
}
