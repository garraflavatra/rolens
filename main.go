package main

import (
	"embed"

	"github.com/garraflavatra/mongodup/internal/app"
	"github.com/wailsapp/wails/v2"
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
		Title:            "Mongodup",
		Width:            1000,
		Height:           600,
		MinWidth:         1000,
		MinHeight:        600,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 139, A: 1},
		OnStartup:        app.Startup,
		Menu:             app.Menu(),

		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		Bind: []interface{}{
			app,
		},

		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Mongodup - MongoDB client",
				Message: "© 2023 Romein van Buren",
				Icon:    appIcon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
