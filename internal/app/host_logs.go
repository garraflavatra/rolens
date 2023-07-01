package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/bson"
)

type HostLogsResult struct {
	Total int32    `json:"total"`
	Logs  []string `json:"logs"`
	Error string   `json:"error"`
}

func (a *App) HostLogs(hostKey, filter string) (result HostLogsResult) {
	client, ctx, close, err := a.connectToHost(hostKey)
	if err != nil {
		result.Error = "Could not connect to host"
		return
	}
	defer close()

	var res bson.M
	err = client.Database("admin").RunCommand(ctx, bson.M{"getLog": filter}).Decode(&res)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Could not get %s logs for %s: %s", filter, hostKey, err.Error())
		result.Error = err.Error()
	}

	result.Total = res["totalLinesWritten"].(int32)
	result.Logs = make([]string, 0)

	for _, v := range res["log"].(bson.A) {
		result.Logs = append(result.Logs, v.(string))
	}

	return
}
