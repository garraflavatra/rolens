package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

//go:embed all:frontend/dist
var assets embed.FS

type Host struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

var hosts = map[string]Host{
	"localhost": {Name: "Localhost", URI: "mongodb://localhost:27017"},
	"tig":       {Name: "cmdb.myinfra.nl"},
	"vbt":       {Name: "vbtverhuurmakelaars.nl"},
}

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Hosts() map[string]Host {
	return hosts
}

func (a *App) connectToHost(hostKey string) (*mongo.Client, context.Context, func(), error) {
	h := hosts[hostKey]
	if len(h.URI) == 0 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Invalid uri",
			Message: "You haven't specified a valid uri for the selected host.",
		})
		return nil, nil, nil, errors.New("invalid uri")
	}

	client, err := mongo.NewClient(mongoOptions.Client().ApplyURI(h.URI))

	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not connect",
			Message: "Failed to establish a connection with " + h.Name,
		})
		return nil, nil, nil, errors.New("could not establish a connection with " + h.Name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client.Connect(ctx)
	return client, ctx, func() {
		client.Disconnect(ctx)
		cancel()
	}, nil
}

func (a *App) OpenConnection(hostKey string) (databases []string) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	databases, err = client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not retrieve database list",
			Message: err.Error(),
		})
		return nil
	}
	defer close()
	return databases
}

func (a *App) OpenDatabase(hostKey, dbKey string) (collections []string) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	collections, err = client.Database(dbKey).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not retrieve collection list for " + dbKey,
			Message: err.Error(),
		})
		return nil
	}
	defer close()
	return collections
}

func (a *App) DropDatabase(hostKey, dbKey string) bool {
	sure, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         "Confirm",
		Message:       "Are you sure you want to drop " + dbKey + "?",
		Buttons:       []string{"yes", "no"},
		DefaultButton: "Yes",
		CancelButton:  "No",
	})
	if sure != "Yes" {
		return false
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	err = client.Database(dbKey).Drop(ctx)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not drop " + dbKey,
			Message: err.Error(),
		})
		return false
	}
	defer close()
	return true
}

func (a *App) OpenCollection(hostKey, dbKey, collKey string) (result bson.M) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	command := bson.M{"collStats": collKey}
	err = client.Database(dbKey).RunCommand(ctx, command).Decode(&result)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not retrieve collection list for " + dbKey,
			Message: err.Error(),
		})
		return nil
	}
	defer close()
	return result
}

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

func (a *App) PerformInsert(hostKey, dbKey, collKey, jsonData string) interface{} {
	var data []interface{}

	jsonData = strings.TrimSpace(jsonData)
	if strings.HasPrefix(jsonData, "{") {
		jsonData = "[" + jsonData + "]"
	}

	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Couldn't parse JSON",
			Message: err.Error(),
		})
		return nil
	}

	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	defer close()
	res, err := client.Database(dbKey).Collection(collKey).InsertMany(ctx, data)
	if err != nil {
		fmt.Println(err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Encountered an error while performing query",
			Message: err.Error(),
		})
		return nil
	}

	return res.InsertedIDs
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Mongodup",
		Width:  1000,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
