package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type EnvironmentInfo struct {
	Arch      string `json:"arch"`
	BuildType string `json:"buildType"`
	Platform  string `json:"platform"`

	HasMongoExport bool `json:"hasMongoExport"`
	HasMongoDump   bool `json:"hasMongoDump"`

	HomeDirectory string `json:"homeDirectory"`
	DataDirectory string `json:"dataDirectory"`
}

type App struct {
	ctx context.Context
	Env EnvironmentInfo
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	wailsEnv := wailsRuntime.Environment(a.ctx)

	a.Env.Arch = wailsEnv.Arch
	a.Env.BuildType = wailsEnv.BuildType
	a.Env.Platform = wailsEnv.Platform

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
	case "darwin":
		a.Env.DataDirectory = filepath.Join(a.Env.HomeDirectory, "/Library/Application Support/Rolens")
	case "linux":
		a.Env.DataDirectory = filepath.Join(a.Env.HomeDirectory, "/.config/Rolens")
	default:
		panic(errors.New("unsupported platform"))
	}

	_ = os.MkdirAll(a.Env.DataDirectory, os.ModePerm)
}

func (a *App) Environment() EnvironmentInfo {
	return a.Env
}

func menuEventEmitter(a *App, eventName string, data ...interface{}) func(cd *menu.CallbackData) {
	return func(cd *menu.CallbackData) {
		wailsRuntime.EventsEmit(a.ctx, eventName, data...)
	}
}

func (a *App) Menu() *menu.Menu {
	appMenu := menu.NewMenu()

	aboutMenu := appMenu.AddSubmenu("About")
	aboutMenu.AddText("About…", nil, menuEventEmitter(a, "OpenAboutModal"))
	aboutMenu.AddText("Prefrences…", keys.CmdOrCtrl(","), menuEventEmitter(a, "OpenPrefrences"))
	aboutMenu.AddSeparator()
	aboutMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(cd *menu.CallbackData) { wailsRuntime.Quit(a.ctx) })

	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.EditMenu())
	}

	hostMenu := appMenu.AddSubmenu("Host")
	hostMenu.AddText("New host", keys.CmdOrCtrl("y"), menuEventEmitter(a, "CreateHost"))

	databaseMenu := appMenu.AddSubmenu("Database")
	databaseMenu.AddText("New database", keys.CmdOrCtrl("u"), menuEventEmitter(a, "CreateDatabase"))

	collectionMenu := appMenu.AddSubmenu("Collection")
	collectionMenu.AddText("New collection", keys.CmdOrCtrl("i"), menuEventEmitter(a, "CreateCollection"))
	collectionMenu.AddSeparator()
	collectionMenu.AddText("Stats", keys.Combo("h", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuEventEmitter(a, "OpenCollectionTab", "stats"))
	collectionMenu.AddText("Find", keys.Combo("f", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuEventEmitter(a, "OpenCollectionTab", "find"))
	collectionMenu.AddText("Insert", keys.Combo("i", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuEventEmitter(a, "OpenCollectionTab", "insert"))
	collectionMenu.AddText("Update", keys.Combo("u", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuEventEmitter(a, "OpenCollectionTab", "update"))
	collectionMenu.AddText("Remove", keys.Combo("r", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuEventEmitter(a, "OpenCollectionTab", "remove"))
	collectionMenu.AddText("Indexes", keys.Combo("x", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuEventEmitter(a, "OpenCollectionTab", "indexes"))

	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("User guide", nil, func(cd *menu.CallbackData) { wailsRuntime.BrowserOpenURL(a.ctx, "") })

	return appMenu
}

func (a *App) OpenDirectory(id, title string) string {
	if title == "" {
		title = "Choose a directory"
	}

	dir, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title:                      title,
		CanCreateDirectories:       true,
		TreatPackagesAsDirectories: false,
	})

	if err != nil {
		fmt.Println(err.Error())
		wailsRuntime.MessageDialog(a.ctx, wailsRuntime.MessageDialogOptions{
			Type:    wailsRuntime.ErrorDialog,
			Title:   "Encountered an error while opening directory",
			Message: err.Error(),
		})
	}

	return dir
}
