package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const photosDir = "images"
const logoFilename = "logo.png"

var photoExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func (a *App) ListPhotos() ([]string, error) {
	entries, err := os.ReadDir(photosDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var photos []string
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == logoFilename {
			continue
		}
		if photoExtensions[strings.ToLower(filepath.Ext(entry.Name()))] {
			photos = append(photos, entry.Name())
		}
	}
	sort.Strings(photos)

	return photos, nil
}
