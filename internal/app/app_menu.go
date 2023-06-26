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

func menuCallbackOpenURL(a *App, url string) func(cd *menu.CallbackData) {
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
	aboutMenu.AddText("Purge logs…", nil, func(cd *menu.CallbackData) { a.PurgeLogDirectory() })
	aboutMenu.AddSeparator()

	if runtime.GOOS == "darwin" {
		aboutMenu.AddText("Minimize", keys.CmdOrCtrl("M"), func(cd *menu.CallbackData) { wailsRuntime.WindowMinimise(a.ctx) })
		aboutMenu.AddText("Hide Rolens", keys.CmdOrCtrl("H"), func(cd *menu.CallbackData) { wailsRuntime.WindowHide(a.ctx) })
		aboutMenu.AddSeparator()

		appMenu.Append(menu.EditMenu())
	}

	aboutMenu.AddText("Quit Rolens", keys.CmdOrCtrl("q"), func(cd *menu.CallbackData) { wailsRuntime.Quit(a.ctx) })

	hostMenu := appMenu.AddSubmenu("Host")
	hostMenu.AddText("New…", nil, menuCallbackEmit(a, "ui.host.new"))
	hostMenu.AddText("Edit host…", nil, menuCallbackEmit(a, "ui.host.edit"))
	hostMenu.AddSeparator()
	hostMenu.AddText("Server status", nil, menuCallbackEmit(a, "ui.host.tab", "status"))
	hostMenu.AddText("System info", nil, menuCallbackEmit(a, "ui.host.tab", "systemInfo"))
	hostMenu.AddText("Shell", nil, menuCallbackEmit(a, "ui.host.tab", "shell"))
	hostMenu.AddSeparator()
	hostMenu.AddText("Remove host…", nil, menuCallbackEmit(a, "ui.host.remove"))

	dbMenu := appMenu.AddSubmenu("Database")
	dbMenu.AddText("New…", nil, menuCallbackEmit(a, "ui.db.new"))
	dbMenu.AddSeparator()
	dbMenu.AddText("Database statistics", nil, menuCallbackEmit(a, "ui.db.tab", "stats"))
	dbMenu.AddText("Shell", nil, menuCallbackEmit(a, "ui.db.tab", "shell"))
	dbMenu.AddSeparator()
	dbMenu.AddText("Dump…", nil, menuCallbackEmit(a, "ui.db.dump"))
	dbMenu.AddText("Drop…", nil, menuCallbackEmit(a, "ui.db.drop"))

	collMenu := appMenu.AddSubmenu("Collection")
	collMenu.AddText("New…", nil, menuCallbackEmit(a, "ui.coll.new"))
	collMenu.AddSeparator()
	collMenu.AddText("Collection statistics", nil, menuCallbackEmit(a, "ui.coll.tab", "stats"))
	collMenu.AddText("Find", nil, menuCallbackEmit(a, "ui.coll.tab", "find"))
	collMenu.AddText("Insert", nil, menuCallbackEmit(a, "ui.coll.tab", "insert"))
	collMenu.AddText("Update", nil, menuCallbackEmit(a, "ui.coll.tab", "update"))
	collMenu.AddText("Remove", nil, menuCallbackEmit(a, "ui.coll.tab", "remove"))
	collMenu.AddText("Indexes", nil, menuCallbackEmit(a, "ui.coll.tab", "indexes"))
	collMenu.AddText("Aggregate", nil, menuCallbackEmit(a, "ui.coll.tab", "aggregate"))
	collMenu.AddText("Shell", nil, menuCallbackEmit(a, "ui.coll.tab", "shell"))
	collMenu.AddSeparator()
	collMenu.AddText("Truncate…", nil, menuCallbackEmit(a, "ui.coll.truncate"))
	collMenu.AddText("Drop…", nil, menuCallbackEmit(a, "ui.coll.drop"))

	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("User guide", nil, menuCallbackOpenURL(a, "https://garraflavatra.github.io/rolens/user-guide/"))
	helpMenu.AddText("Website", nil, menuCallbackOpenURL(a, "https://garraflavatra.github.io/rolens/"))
	helpMenu.AddText("Discussion board", nil, menuCallbackOpenURL(a, "https://github.com/garraflavatra/rolens/discussions"))
	helpMenu.AddSeparator()
	helpMenu.AddText("Report a problem", nil, menuCallbackOpenURL(a, "https://github.com/garraflavatra/rolens/issues/new"))
	helpMenu.AddText("Ask a question", nil, menuCallbackOpenURL(a, "https://github.com/garraflavatra/rolens/discussions/new?category=questions"))
	helpMenu.AddSeparator()
	helpMenu.AddText("Star Rolens on GitHub", nil, menuCallbackOpenURL(a, "https://github.com/garraflavatra/rolens"))
	helpMenu.AddText("Changelog", nil, menuCallbackOpenURL(a, "https://garraflavatra.github.io/rolens/development/changelog/"))
	helpMenu.AddText("License", nil, menuCallbackOpenURL(a, "https://github.com/garraflavatra/rolens/blob/main/LICENSE"))

	return appMenu
}
