package main

import (
	"context"
	"sync"
)

// App struct
type App struct {
	ctx            context.Context
	networkMonitor networkMonitor
	photosDir      string
	photosDirMutex sync.RWMutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{photosDir: defaultPhotosDir}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) getPhotosDir() string {
	a.photosDirMutex.RLock()
	defer a.photosDirMutex.RUnlock()
	return a.photosDir
}

func (a *App) setPhotosDir(directory string) {
	a.photosDirMutex.Lock()
	defer a.photosDirMutex.Unlock()
	a.photosDir = directory
}
