package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/garraflavatra/rolens/internal/ui"
	"github.com/ncruces/zenity"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sync/syncmap"
)

type EnvironmentInfo struct {
	Arch      string `json:"arch"`
	BuildType string `json:"buildType"`
	Platform  string `json:"platform"`
	Version   string `json:"version"`

	HasMongoExport bool `json:"hasMongoExport"`
	HasMongoDump   bool `json:"hasMongoDump"`

	HomeDirectory     string `json:"homeDirectory"`
	DataDirectory     string `json:"dataDirectory"`
	LogDirectory      string `json:"logDirectory"`
	DownloadDirectory string `json:"downloadDirectory"`
}

type App struct {
	Env   EnvironmentInfo
	State syncmap.Map
	ctx   context.Context
	ui    *ui.UI
}

func NewApp(version string) *App {
	a := &App{}
	a.Env.Version = strings.TrimSpace(version)

	_, err := exec.LookPath("mongodump")
	a.Env.HasMongoDump = err == nil

	_, err = exec.LookPath("mongoexport")
	a.Env.HasMongoExport = err == nil

	a.Env.HomeDirectory, err = os.UserHomeDir()
	if err != nil {
		panic(errors.New("encountered an error while getting home directory"))
	}

	switch runtime.GOOS {
	case "windows":
		a.Env.DataDirectory = filepath.Join(a.Env.HomeDirectory, "/AppData/Local/Rolens")
		a.Env.LogDirectory = filepath.Join(a.Env.HomeDirectory, "/AppData/Local/Rolens/Logs")
		a.Env.DownloadDirectory = filepath.Join(a.Env.HomeDirectory, "/Downloads")
	case "darwin":
		a.Env.DataDirectory = filepath.Join(a.Env.HomeDirectory, "/Library/Application Support/Rolens")
		a.Env.LogDirectory = filepath.Join(a.Env.HomeDirectory, "/Library/Logs/Rolens")
		a.Env.DownloadDirectory = filepath.Join(a.Env.HomeDirectory, "/Downloads")
	case "linux":
		a.Env.DataDirectory = filepath.Join(a.Env.HomeDirectory, "/.config/rolens")
		a.Env.LogDirectory = filepath.Join(a.Env.HomeDirectory, "/.config/rolens/logs")
		a.Env.DownloadDirectory = filepath.Join(a.Env.HomeDirectory, "/Downloads")
	default:
		panic(errors.New("unsupported platform"))
	}

	os.MkdirAll(a.Env.DataDirectory, os.ModePerm)
	os.MkdirAll(a.Env.LogDirectory, os.ModePerm)

	return a
}

func (a *App) Startup(ctx context.Context, ui *ui.UI) {
	a.ctx = ctx
	a.ui = ui
	wailsRuntime.LogInfo(a.ctx, "Runcycle: Startup")

	wailsEnv := wailsRuntime.Environment(a.ctx)
	a.Env.Arch = wailsEnv.Arch
	a.Env.BuildType = wailsEnv.BuildType
	a.Env.Platform = wailsEnv.Platform

	jsonEnv, err := json.MarshalIndent(a.Env, "", "    ")
	if err != nil {
		wailsRuntime.LogWarningf(a.ctx, "Could not marshal environment info: %s", err.Error())
	}
	err = os.WriteFile(path.Join(a.Env.LogDirectory, "environment.json"), jsonEnv, 0644)
	if err != nil {
		wailsRuntime.LogWarningf(a.ctx, "Could not save environment.json: %s", err.Error())
	}
}

func (a *App) Shutdown(ctx context.Context) {
	wailsRuntime.LogInfo(a.ctx, "Runcycle: Shutdown")
}

func (a *App) Environment() EnvironmentInfo {
	return a.Env
}

func (a *App) PurgeLogDirectory() {
	err := zenity.Question("Are you sure you want to remove all logfiles?", zenity.Title("Confirm"), zenity.WarningIcon)
	if err == zenity.ErrCanceled {
		return
	}

	err = os.RemoveAll(a.Env.LogDirectory)
	if err == nil {
		zenity.Info("Successfully purged log directory.", zenity.InfoIcon)
	} else {
		zenity.Error(err.Error(), zenity.Title("Encountered an error while purging log directory."), zenity.WarningIcon)
	}
}

func (a *App) ReportSharedStateVariable(key, value string) {
	a.State.Store(key, value)
	wailsRuntime.LogDebug(a.ctx, fmt.Sprintf("State: %s=\"%s\"", key, value))
}
