package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

const version = "0.1.0-dev"

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app, err := NewApp(version)
	if err != nil {
		log.Fatalf("failed to init app service: %v", err)
	}

	err = wails.Run(&options.App{
		Title:  "Go Wails Shot",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatalf("wails run failed: %v", err)
	}
}
