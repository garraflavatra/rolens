package app

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DumpInfo struct {
	OutDir   string   `json:"outdir"`
	Filename string   `json:"filename"`
	HostKey  string   `json:"hostKey"`
	DbKey    string   `json:"dbKey"`
	CollKeys []string `json:"collKeys"`
}

func (a *App) PerformDump(jsonData string) bool {
	var info DumpInfo
	err := json.Unmarshal([]byte(jsonData), &info)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Could not unmarshal dump form: %s", err.Error())
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Malformed JSON",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	hosts, err := a.Hosts()
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Error getting hosts",
			Message: err.Error(),
			Type:    runtime.ErrorDialog,
		})
		return false
	}
	host := hosts[info.HostKey]

	if !a.Env.HasMongoDump {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:   "Additional tooling required",
			Message: "You need to install mongodump to perform a dump.",
			Type:    runtime.ErrorDialog,
		})
		return false
	}

	args := make([]string, 0)
	args = append(args, fmt.Sprintf(`--uri="%v"`, host.URI))
	fname := path.Join(info.OutDir, info.Filename)

	if info.DbKey != "" {
		args = append(args, fmt.Sprintf(`--db="%v"`, info.DbKey))

		if info.CollKeys != nil {
			args = append(args, fmt.Sprintf(`--collection="%v"`, info.CollKeys[0]))
		}
	}

	args = append(args, fmt.Sprintf(`--out="%v"`, fname))
	cmd := exec.Command("mongodump", args...)
	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	runtime.LogInfo(a.ctx, "Performing dump, executing command: mongodump "+strings.Join(args, " "))
	runtime.LogInfo(a.ctx, "mongodump stdout: "+stdout.String())
	runtime.LogInfo(a.ctx, "mongodump sterr: "+stderr.String())

	if err != nil {
		runtime.LogWarning(a.ctx, "Error while executing mongodump: "+err.Error())
		return false
	}

	println(fname)

	a.ui.Reveal(fname)
	return true
}
