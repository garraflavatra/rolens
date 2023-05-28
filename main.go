package main

import (
	"context"
	"embed"

	"github.com/garraflavatra/rolens/internal"
	"github.com/garraflavatra/rolens/internal/app"
	uictrl "github.com/garraflavatra/rolens/internal/ui"
	"github.com/ncruces/zenity"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	//go:embed all:frontend/dist
	assets embed.FS

	//go:embed build/appicon.png
	appIcon []byte
)

func main() {
	app := app.NewApp()
	ui := uictrl.New()

	err := wails.Run(&options.App{
		Title: "Rolens",

		Width:     1000,
		Height:    600,
		MinWidth:  1000,
		MinHeight: 600,

		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 139, A: 1},
		Menu:             app.Menu(),
		Bind:             []interface{}{app, ui},
		AssetServer:      &assetserver.Options{Assets: assets},

		OnStartup: func(ctx context.Context) {
			defer func() {
				if p := recover(); p != nil {
					runtime.LogFatalf(ctx, "Application panicked: %v", p)
					zenity.Error("A fatal error occured.")
				}
			}()

			ui.Startup(ctx)
			app.Startup(ctx, ui)
		},
		OnShutdown: app.Shutdown,

		Logger:             internal.NewAppLogger(app.Env.LogDirectory, "rolens.log"),
		LogLevel:           logger.TRACE,
		LogLevelProduction: logger.INFO,

		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Rolens - Multiplatform MongoDB client",
				Message: "© 2023 Romein van Buren",
				Icon:    appIcon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
