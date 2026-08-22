package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const noPhotosDir = ""
const choosePhotosFolderTitle = "Choose a folder containing photos"

var photoExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func (a *App) ListPhotos() ([]string, error) {
	photosDir := a.getPhotosDir()
	if photosDir == noPhotosDir {
		return []string{}, nil
	}

	entries, err := os.ReadDir(photosDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var photos []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if photoExtensions[strings.ToLower(filepath.Ext(entry.Name()))] {
			photos = append(photos, entry.Name())
		}
	}
	sort.Strings(photos)

	return photos, nil
}

func (a *App) ChoosePhotosFolder() (string, error) {
	selected, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: choosePhotosFolderTitle,
	})
	if err != nil || selected == "" {
		return "", err
	}

	a.setPhotosDir(selected)
	return selected, nil
}
