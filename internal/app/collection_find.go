package app

import (
	"encoding/json"

	"github.com/ncruces/zenity"
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

type QueryResult struct {
	Total   int64    `json:"total"`
	Results []string `json:"results"`
}

func (a *App) FindItems(hostKey, dbKey, collKey, formJson string) QueryResult {
	var out QueryResult
	var form Query

	err := json.Unmarshal([]byte(formJson), &form)
	if err != nil {
		runtime.LogError(a.ctx, "Could not parse find form:")
		runtime.LogError(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Could not parse form"), zenity.ErrorIcon)
		return out
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return out
	}

	defer close()
	var query bson.M
	var projection bson.M
	var sort bson.M

	err = bson.UnmarshalExtJSON([]byte(form.Query), true, &query)
	if err != nil {
		runtime.LogInfof(a.ctx, "Invalid find query: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Invalid query"), zenity.ErrorIcon)
		return out
	}

	err = json.Unmarshal([]byte(form.Fields), &projection)
	if err != nil {
		runtime.LogInfof(a.ctx, "Invalid find projection: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Invalid projection"), zenity.ErrorIcon)
		return out
	}

	err = json.Unmarshal([]byte(form.Sort), &sort)
	if err != nil {
		runtime.LogInfof(a.ctx, "Invalid find sort: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Invalid sort"), zenity.ErrorIcon)
		return out
	}

	opt := mongoOptions.FindOptions{
		Limit:      &form.Limit,
		Projection: projection,
		Skip:       &form.Skip,
		Sort:       sort,
	}

	total, err := client.Database(dbKey).Collection(collKey).CountDocuments(ctx, query, nil)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Encountered an error while counting documents: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while counting docs"), zenity.ErrorIcon)
		return out
	}

	cur, err := client.Database(dbKey).Collection(collKey).Find(ctx, query, &opt)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Encountered an error while performing query: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while querying"), zenity.ErrorIcon)
		return out
	}

	defer cur.Close(ctx)
	var results []bson.M
	err = cur.All(ctx, &results)

	if err != nil {
		runtime.LogWarningf(a.ctx, "Encountered an error while performing query: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while querying"), zenity.ErrorIcon)
		return out
	}

	out.Total = total
	out.Results = make([]string, 0)

	for _, r := range results {
		marshalled, err := bson.MarshalExtJSON(r, true, true)
		if err != nil {
			runtime.LogError(a.ctx, "Failed to marshal find BSON:")
			runtime.LogError(a.ctx, err.Error())
			zenity.Error(err.Error(), zenity.Title("Failed to marshal JSON"), zenity.ErrorIcon)
			return out
		}
		out.Results = append(out.Results, string(marshalled))
	}

	return out
}

func (a *App) UpdateFoundDocument(hostKey, dbKey, collKey, idJson, newDocJson string) bool {
	var id bson.M
	if err := bson.UnmarshalExtJSON([]byte(idJson), true, &id); err != nil {
		runtime.LogWarningf(a.ctx, "Could not parse find/update query: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Couldn't parse update query"), zenity.ErrorIcon)
		return false
	}

	var newDoc bson.M
	if err := bson.UnmarshalExtJSON([]byte(newDocJson), true, &newDoc); err != nil {
		runtime.LogWarningf(a.ctx, "Could not parse new find/update document: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Couldn't parse document"), zenity.ErrorIcon)
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return false
	}
	defer close()

	if _, err := client.Database(dbKey).Collection(collKey).ReplaceOne(ctx, id, newDoc); err != nil {
		runtime.LogInfof(a.ctx, "Error while replacing document: %s", err.Error())
		zenity.Error(err.Error(), zenity.Title("Unable to replace document"), zenity.ErrorIcon)
		return false
	}

	return true
}
