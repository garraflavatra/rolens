package app

import (
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func menuCallbackEmit(a *App, eventName string, data ...interface{}) func(cd *menu.CallbackData) {
	return func(cd *menu.CallbackData) {
		wailsRuntime.EventsEmit(a.ctx, eventName, data...)
	}
}

func menuCallbackURL(a *App, url string) func(cd *menu.CallbackData) {
	return func(cd *menu.CallbackData) {
		wailsRuntime.BrowserOpenURL(a.ctx, url)
	}
}

func (a *App) Menu() *menu.Menu {
	appMenu := menu.NewMenu()

	aboutMenu := appMenu.AddSubmenu("Rolens")
	aboutMenu.AddText("About Rolens", nil, menuCallbackEmit(a, "OpenAboutModal"))
	aboutMenu.AddSeparator()
	aboutMenu.AddText("Preferences…", keys.CmdOrCtrl(","), menuCallbackEmit(a, "OpenPreferences"))
	aboutMenu.AddSeparator()
	aboutMenu.AddText("Open data directory…", nil, func(cd *menu.CallbackData) { a.ui.Reveal(a.Env.DataDirectory) })
	aboutMenu.AddText("Open log directory…", nil, func(cd *menu.CallbackData) { a.ui.Reveal(a.Env.LogDirectory) })
	aboutMenu.AddText("Purge logs", nil, func(cd *menu.CallbackData) { a.PurgeLogDirectory() })
	aboutMenu.AddSeparator()
	aboutMenu.AddText("Quit Rolens", keys.CmdOrCtrl("q"), func(cd *menu.CallbackData) { wailsRuntime.Quit(a.ctx) })

	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("New host…", keys.CmdOrCtrl("y"), menuCallbackEmit(a, "CreateHost"))
	fileMenu.AddText("New database", keys.CmdOrCtrl("y"), menuCallbackEmit(a, "CreateDatabase"))
	fileMenu.AddText("New collection…", keys.CmdOrCtrl("i"), menuCallbackEmit(a, "CreateCollection"))
	fileMenu.AddSeparator()
	fileMenu.AddText("Stats", keys.Combo("h", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuCallbackEmit(a, "OpenCollectionTab", "stats"))
	fileMenu.AddText("Find", keys.Combo("f", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuCallbackEmit(a, "OpenCollectionTab", "find"))
	fileMenu.AddText("Insert", keys.Combo("i", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuCallbackEmit(a, "OpenCollectionTab", "insert"))
	fileMenu.AddText("Update", keys.Combo("u", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuCallbackEmit(a, "OpenCollectionTab", "update"))
	fileMenu.AddText("Remove", keys.Combo("r", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuCallbackEmit(a, "OpenCollectionTab", "remove"))
	fileMenu.AddText("Indexes", keys.Combo("x", keys.CmdOrCtrlKey, keys.OptionOrAltKey), menuCallbackEmit(a, "OpenCollectionTab", "indexes"))

	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.EditMenu())
	}

	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("User guide", nil, menuCallbackURL(a, ""))
	helpMenu.AddSeparator()
	helpMenu.AddText("Report a problem", nil, menuCallbackURL(a, "https://github.com/garraflavatra/rolens/issues/new"))
	helpMenu.AddSeparator()
	helpMenu.AddText("Star Rolens on GitHub", nil, menuCallbackURL(a, "https://github.com/garraflavatra/rolens"))
	helpMenu.AddText("License", nil, menuCallbackURL(a, "https://github.com/garraflavatra/rolens/blob/main/LICENSE"))

	return appMenu
}
