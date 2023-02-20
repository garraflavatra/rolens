package app

import (
	"encoding/json"
	"fmt"

	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *App) Aggregate(hostKey, dbKey, collKey, pipelineJson, settingsJson string) {
	var settings *options.AggregateOptions
	if err := json.Unmarshal([]byte(settingsJson), &settings); err != nil {
		runtime.LogError(a.ctx, "Could not parse aggregation settings:")
		runtime.LogError(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Couldn't parse aggregation settings"), zenity.ErrorIcon)
		return
	}

	var pipeline mongo.Pipeline
	if err := bson.UnmarshalExtJSON([]byte(pipelineJson), true, &pipeline); err != nil {
		runtime.LogWarning(a.ctx, "Could not parse aggregation pipeline:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Couldn't parse aggregation pipeline"), zenity.ErrorIcon)
		return
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return
	}

	defer close()

	cursor, err := client.Database(dbKey).Collection(collKey).Aggregate(ctx, pipeline, settings)
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not get aggregation cursor:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Couldn't get aggregation cursor"), zenity.ErrorIcon)
		return
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		runtime.LogInfo(a.ctx, "Error while running aggregation pipeline:")
		runtime.LogInfo(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Error while running aggregation pipeline"), zenity.ErrorIcon)
		return
	}

	fmt.Println(results)
}
