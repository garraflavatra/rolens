package app

import (
	"encoding/json"
	"math"

	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *App) GetIndexes(hostKey, dbKey, collKey string) []bson.M {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return nil
	}
	defer close()

	cur, err := client.Database(dbKey).Collection(collKey).Indexes().List(ctx)
	if err != nil {
		runtime.LogWarning(a.ctx, "Encountered an error while creating index cursor:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Error while creating cursor"), zenity.ErrorIcon)
		return nil
	}

	var results []bson.M
	err = cur.All(ctx, &results)
	if err != nil {
		runtime.LogWarning(a.ctx, "Encountered an error while executing index cursor:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Error while running cursor"), zenity.ErrorIcon)
		return nil
	}

	return results
}

func (a *App) CreateIndex(hostKey, dbKey, collKey, jsonData string) string {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return ""
	}
	defer close()

	type modelItem struct {
		Key  string      `json:"key"`
		Sort interface{} `json:"sort"`
	}

	var form struct {
		Name           string      `json:"name"`
		Background     bool        `json:"background"`
		Unique         bool        `json:"unique"`
		DropDuplicates bool        `json:"dropDuplicates"`
		Sparse         bool        `json:"sparse"`
		Model          []modelItem `json:"model"`
	}

	err = json.Unmarshal([]byte(jsonData), &form)
	if err != nil {
		runtime.LogError(a.ctx, "Could not parse index JSON:")
		runtime.LogError(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Could not parse JSON"), zenity.ErrorIcon)
		return ""
	}

	if form.Model == nil {
		form.Model = make([]modelItem, 0)
	}

	var keys bson.D
	for _, v := range form.Model {
		asFloat, canBeFloat := v.Sort.(float64)
		if canBeFloat {
			v.Sort = int8(math.Floor(asFloat))
		}

		keys = append(keys, bson.E{
			Key:   v.Key,
			Value: v.Sort,
		})
	}

	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(form.Unique).SetSparse(form.Sparse),
	}

	name, err := client.Database(dbKey).Collection(collKey).Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		runtime.LogWarning(a.ctx, "Encountered an error while creating index:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Error while creating index"), zenity.ErrorIcon)
		return ""
	}

	return name
}

func (a *App) DropIndex(hostKey, dbKey, collKey, indexName string) bool {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	_, err = client.Database(dbKey).Collection(collKey).Indexes().DropOne(ctx, indexName, &options.DropIndexesOptions{})
	if err != nil {
		runtime.LogError(a.ctx, "Encountered an error while creating index drop cursor:")
		runtime.LogError(a.ctx, err.Error())
		zenity.Info(err.Error(), zenity.Title("Error while creating drop cursor"), zenity.ErrorIcon)
		return false
	}

	return true
}
