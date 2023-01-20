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

func (a *App) Menu() *menu.Menu {
	appMenu := menu.NewMenu()

	aboutMenu := appMenu.AddSubmenu("About")
	aboutMenu.AddText("About…", nil, func(cd *menu.CallbackData) {
		wailsRuntime.EventsEmit(a.ctx, "OpenAboutModal")
	})
	aboutMenu.AddText("Prefrences…", keys.CmdOrCtrl(","), func(cd *menu.CallbackData) {
		wailsRuntime.EventsEmit(a.ctx, "OpenPrefrences")
	})
	aboutMenu.AddSeparator()
	aboutMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(cd *menu.CallbackData) {
		wailsRuntime.Quit(a.ctx)
	})

	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Hosts…", keys.CmdOrCtrl("k"), func(cd *menu.CallbackData) {
		wailsRuntime.EventsEmit(a.ctx, "OpenHostsModal")
	})

	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.EditMenu())
	}

	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("User guide", nil, func(cd *menu.CallbackData) {
		wailsRuntime.BrowserOpenURL(a.ctx, "")
	})

	return appMenu
}
