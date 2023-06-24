package app

import (
	"encoding/json"
	"math"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetIndexesResult struct {
	Indexes []bson.M `json:"indexes"`
	Error   string   `json:"error"`
}

func (a *App) GetIndexes(hostKey, dbKey, collKey string) (result GetIndexesResult) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return
	}
	defer close()

	cur, err := client.Database(dbKey).Collection(collKey).Indexes().List(ctx)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Encountered an error while creating index cursor: %s", err.Error())
		result.Error = err.Error()
		return
	}

	err = cur.All(ctx, &result.Indexes)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Encountered an error while executing index cursor: %s", err.Error())
		result.Error = err.Error()
	}

	return
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
		runtime.LogErrorf(a.ctx, "Could not parse index JSON: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Malformed JSON",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
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
		runtime.LogWarningf(a.ctx, "Encountered an error while creating index: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error creating index",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
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
		runtime.LogErrorf(a.ctx, "Encountered an error while creating index drop cursor: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error creating drop cursor",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return true
}
