package app

import (
	"encoding/json"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"go.mongodb.org/mongo-driver/bson"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

type findResult struct {
	Total   int64       `json:"total"`
	Results interface{} `json:"results"`
}

func (a *App) PerformFind(hostKey, dbKey, collKey string, formJson string) findResult {
	var out findResult
	var form struct {
		Fields string `json:"fields"`
		Limit  int64  `json:"limit"`
		Query  string `json:"query"`
		Skip   int64  `json:"skip"`
		Sort   string `json:"sort"`
	}

	err := json.Unmarshal([]byte(formJson), &form)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Couldn't parse form",
			Message: err.Error(),
		})
		return out
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return out
	}

	defer close()
	var query bson.M
	var projection bson.M
	var sort bson.M

	err = json.Unmarshal([]byte(form.Query), &query)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Invalid query",
			Message: err.Error(),
		})
		return out
	}

	err = json.Unmarshal([]byte(form.Fields), &projection)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Invalid projection",
			Message: err.Error(),
		})
		return out
	}

	err = json.Unmarshal([]byte(form.Sort), &sort)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Invalid sort",
			Message: err.Error(),
		})
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
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while counting documents",
			Message: err.Error(),
		})
		return out
	}

	cur, err := client.Database(dbKey).Collection(collKey).Find(ctx, query, &opt)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while performing query",
			Message: err.Error(),
		})
		return out
	}

	defer cur.Close(ctx)
	var results []bson.M
	err = cur.All(ctx, &results)

	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while performing query",
			Message: err.Error(),
		})
		return out
	}

	out.Results = results
	out.Total = total
	return out
}
