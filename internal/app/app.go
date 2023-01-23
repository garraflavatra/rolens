package app

import (
	"context"
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
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
