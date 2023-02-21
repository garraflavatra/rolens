package app

import (
	"encoding/json"
	"time"

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
	time.Sleep(2 * time.Second)
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
		runtime.LogInfo(a.ctx, "Invalid find query:")
		runtime.LogInfo(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Invalid query"), zenity.ErrorIcon)
		return out
	}

	err = json.Unmarshal([]byte(form.Fields), &projection)
	if err != nil {
		runtime.LogInfo(a.ctx, "Invalid find projection:")
		runtime.LogInfo(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Invalid projection"), zenity.ErrorIcon)
		return out
	}

	err = json.Unmarshal([]byte(form.Sort), &sort)
	if err != nil {
		runtime.LogInfo(a.ctx, "Invalid find sort:")
		runtime.LogInfo(a.ctx, err.Error())
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
		runtime.LogWarning(a.ctx, "Encountered an error while counting documents:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while counting docs"), zenity.ErrorIcon)
		return out
	}

	cur, err := client.Database(dbKey).Collection(collKey).Find(ctx, query, &opt)
	if err != nil {
		runtime.LogWarning(a.ctx, "Encountered an error while performing query:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while querying"), zenity.ErrorIcon)
		return out
	}

	defer cur.Close(ctx)
	var results []bson.M
	err = cur.All(ctx, &results)

	if err != nil {
		runtime.LogWarning(a.ctx, "Encountered an error while performing query:")
		runtime.LogWarning(a.ctx, err.Error())
		zenity.Error(err.Error(), zenity.Title("Error while querying"), zenity.ErrorIcon)
		return out
	}

	out.Results = make([]string, 0)
	out.Total = total
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
