package app

import (
	"context"
	"errors"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

func (a *App) connectToHost(hostKey string) (*mongo.Client, context.Context, func(), error) {
	hosts, err := a.Hosts()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:  runtime.InfoDialog,
			Title: "Could not retrieve hosts",
		})
		return nil, nil, nil, errors.New("could not retrieve hosts")
	}

	h := hosts[hostKey]
	if len(h.URI) == 0 {
		runtime.LogInfo(a.ctx, "Invalid URI (len == 0) for host "+hostKey)
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Invalid uri",
			Message: "You haven't specified a valid uri for the selected host.",
		})
		return nil, nil, nil, errors.New("invalid uri")
	}

	client, err := mongo.NewClient(mongoOptions.Client().ApplyURI(h.URI))

	if err != nil {
		runtime.LogWarning(a.ctx, "Could not connect to host "+hostKey)
		runtime.LogWarning(a.ctx, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Could not connect to " + h.Name,
			Message: err.Error(),
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
		return nil
	}
	databases, err = client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		runtime.LogWarning(a.ctx, "Could not retrieve database names for host "+hostKey)
		runtime.LogWarning(a.ctx, err.Error())
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
