package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

type Host struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

var hosts = map[string]Host{
	"localhost": {Name: "Localhost", URI: "mongodb://localhost:27017"},
	"tig":       {Name: "cmdb.myinfra.nl"},
	"vbt":       {Name: "vbtverhuurmakelaars.nl"},
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
