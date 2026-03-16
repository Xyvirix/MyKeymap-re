package main

import (
	"log"

	"settings/internal/desktop"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func main() {
	app, err := desktop.NewDesktopApp()
	if err != nil {
		log.Fatal(err)
	}

	assets, err := app.Paths().FrontendFS()
	if err != nil {
		log.Fatal(err)
	}

	err = wails.Run(&options.App{
		Title:             "MyKeymap-re Desktop",
		Width:             1480,
		Height:            960,
		MinWidth:          1220,
		MinHeight:         760,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:     app.Startup,
		OnDomReady:    app.DomReady,
		OnShutdown:    app.Shutdown,
		OnBeforeClose: app.BeforeClose,
		Bind: []interface{}{
			app,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "mykeymap.desktop.singleton",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				app.ShowWindow()
			},
		},
		BackgroundColour: &options.RGBA{R: 243, G: 246, B: 250, A: 1},
		Windows: &windows.Options{
			DisableWindowIcon:    false,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
