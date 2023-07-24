package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

type Query struct {
	Fields string `json:"fields"`
	Limit  int64  `json:"limit"`
	Query  string `json:"query"`
	Skip   int64  `json:"skip"`
	Sort   string `json:"sort"`
}

type FindItemsResult struct {
	Total            int64    `json:"total"`
	Results          []string `json:"results"`
	ErrorTitle       string   `json:"errorTitle"`
	ErrorDescription string   `json:"errorDescription"`
}

type CountItemsResult struct {
	Total int64  `json:"total"`
	Error string `json:"error"`
}

func (a *App) FindItems(hostKey, dbKey, collKey, formJson string) (result FindItemsResult) {
	var form Query

	err := json.Unmarshal([]byte(formJson), &form)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Could not parse find form: %s", err.Error())
		result.ErrorTitle = "Could not parse form"
		result.ErrorDescription = err.Error()
		return
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return
	}
	defer close()

	var query bson.M
	var projection bson.M
	var sort bson.M

	err = bson.UnmarshalExtJSON([]byte(form.Query), true, &query)
	if err != nil {
		runtime.LogInfof(a.ctx, "Invalid find query: %s", err.Error())
		result.ErrorTitle = "Invalid query"
		result.ErrorDescription = err.Error()
		return
	}

	err = json.Unmarshal([]byte(form.Fields), &projection)
	if err != nil {
		runtime.LogInfof(a.ctx, "Invalid find projection: %s", err.Error())
		result.ErrorTitle = "Invalid projection"
		result.ErrorDescription = err.Error()
		return
	}

	err = json.Unmarshal([]byte(form.Sort), &sort)
	if err != nil {
		runtime.LogInfof(a.ctx, "Invalid find sort: %s", err.Error())
		result.ErrorTitle = "Invalid sort"
		result.ErrorDescription = err.Error()
		return
	}

	opt := mongoOptions.FindOptions{
		Limit:      &form.Limit,
		Projection: projection,
		Skip:       &form.Skip,
		Sort:       sort,
	}

	ctx2, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second))
	defer cancel()
	total, err := client.Database(dbKey).Collection(collKey).CountDocuments(ctx2, query, nil)
	if err == nil {
		result.Total = total
	} else {
		result.Total = -1
		runtime.LogWarningf(a.ctx, "Encountered an error while counting documents: %s", err.Error())
	}

	cur, err := client.Database(dbKey).Collection(collKey).Find(ctx, query, &opt)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Encountered an error while performing query: %s", err.Error())
		result.ErrorTitle = "Error while querying"
		result.ErrorDescription = err.Error()
		return
	}

	defer cur.Close(ctx)
	var results []bson.M
	err = cur.All(ctx, &results)

	if err != nil {
		runtime.LogWarningf(a.ctx, "Encountered an error while performing query: %s", err.Error())
		result.ErrorTitle = "Error while querying"
		result.ErrorDescription = err.Error()
		return
	}

	result.Results = make([]string, 0)

	for _, r := range results {
		marshalled, err := bson.MarshalExtJSON(r, true, true)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Failed to marshal find BSON: %s", err.Error())
			result.ErrorTitle = "Failed to marshal JSON"
			result.ErrorDescription = err.Error()
			return
		}
		result.Results = append(result.Results, string(marshalled))
	}

	return
}

func (a *App) CountItems(hostKey, dbKey, collKey, query string) (result CountItemsResult) {
	var q bson.M
	err := bson.UnmarshalExtJSON([]byte(query), true, &q)
	if err != nil {
		runtime.LogInfof(a.ctx, "Invalid count query: %s", err.Error())
		result.Error = "Invalid query"
		return
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return
	}
	defer close()

	ctx, _ = context.WithDeadline(ctx, time.Now().Add(2*time.Minute))
	total, err := client.Database(dbKey).Collection(collKey).CountDocuments(ctx, q, nil)
	if err == nil {
		result.Total = total
	} else {
		result.Total = -1
		result.Error = err.Error()
		runtime.LogWarningf(a.ctx, "Encountered an error while counting documents: %s", err.Error())
	}

	return
}

func (a *App) UpdateFoundDocument(hostKey, dbKey, collKey, idJson, newDocJson string) bool {
	var id bson.M
	if err := bson.UnmarshalExtJSON([]byte(idJson), true, &id); err != nil {
		runtime.LogWarningf(a.ctx, "Could not parse find/update query: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error parsing update query",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	var newDoc bson.M
	if err := bson.UnmarshalExtJSON([]byte(newDocJson), true, &newDoc); err != nil {
		runtime.LogWarningf(a.ctx, "Could not parse new find/update document: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error parsing document",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	if _, err := client.Database(dbKey).Collection(collKey).ReplaceOne(ctx, id, newDoc); err != nil {
		runtime.LogInfof(a.ctx, "Error while replacing document: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error replacing document",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	return true
}
