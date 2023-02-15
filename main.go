package main

import (
	"embed"
	"path"

	"github.com/garraflavatra/rolens/internal/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

var (
	//go:embed all:frontend/dist
	assets embed.FS

	//go:embed build/appicon.png
	appIcon []byte
)

func main() {
	app := app.NewApp()
	err := wails.Run(&options.App{
		Title: "Rolens",

		Width:     1000,
		Height:    600,
		MinWidth:  1000,
		MinHeight: 600,

		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 139, A: 1},
		Menu:             app.Menu(),
		Bind:             []interface{}{app},
		AssetServer:      &assetserver.Options{Assets: assets},

		OnStartup:  app.Startup,
		OnShutdown: app.Shutdown,

		Logger:             logger.NewFileLogger(path.Join(app.Env.LogDirectory, "rolens.log")),
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
