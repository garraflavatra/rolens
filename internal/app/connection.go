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

type OpenConnectionResult struct {
	Databases       []string `json:"databases"`
	Status          bson.M   `json:"status"`
	StatusError     string   `json:"statusError"`
	SystemInfo      bson.M   `json:"systemInfo"`
	SystemInfoError string   `json:"systemInfoError"`
}

func (a *App) connectToHost(hostKey string) (*mongo.Client, context.Context, func(), error) {
	hosts, err := a.Hosts()
	if err != nil {
		runtime.LogInfof(a.ctx, "Error while getting hosts: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error getting hosts",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return nil, nil, nil, errors.New("could not retrieve hosts")
	}

	h := hosts[hostKey]
	if len(h.URI) == 0 {
		runtime.LogInfof(a.ctx, "Invalid URI (len == 0) for host %s", hostKey)
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Invalid host information",
			Message: "You haven't specified a valid uri for the selected host.",
			Type:    runtime.ErrorDialog,
		})
		return nil, nil, nil, errors.New("invalid uri")
	}

	client, err := mongo.NewClient(mongoOptions.Client().ApplyURI(h.URI))

	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not connect to host %s: %s", hostKey, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error while connecting to " + h.Name,
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
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

func (a *App) OpenConnection(hostKey string) (result OpenConnectionResult) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		return
	}
	defer close()

	result.Databases, err = client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not retrieve database names for host %s: %s", hostKey, err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error getting database list",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
	}

	command := bson.M{"serverStatus": 1}
	err = client.Database("admin").RunCommand(ctx, command).Decode(&result.Status)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not retrieve server status: %s", err.Error())
		result.StatusError = err.Error()
	}

	command = bson.M{"hostInfo": 1}
	err = client.Database("admin").RunCommand(ctx, command).Decode(&result.SystemInfo)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not retrieve system info: %s", err.Error())
		result.SystemInfoError = err.Error()
	}

	return
}
