package main

import (
	"embed"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

const (
	windowWidth  = 1024
	windowHeight = 768
)

func main() {
	// Create an instance of the app structure
	app := NewApp()

	photosHandler := http.NewServeMux()
	photosHandler.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(photosDir))))

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "pi-beta",
		Width:  windowWidth,
		Height: windowHeight,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: photosHandler,
		},
		BackgroundColour: &colorBackground,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
