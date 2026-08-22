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

const (
	imagesRoutePrefix = "/images/"
	photosRoutePrefix = "/photos/"
)

func main() {
	// Create an instance of the app structure
	app := NewApp()

	photosHandler := http.NewServeMux()
	photosHandler.Handle(imagesRoutePrefix, http.StripPrefix(imagesRoutePrefix, http.FileServer(http.Dir(defaultPhotosDir))))
	photosHandler.HandleFunc(photosRoutePrefix, func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(photosRoutePrefix, http.FileServer(http.Dir(app.getPhotosDir()))).ServeHTTP(w, r)
	})

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
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
