package main

import (
	"encoding/json"
	"os"
)

const manifestPath = "static/rev-manifest.json"

// Preload the asset manifest on boot in production
var getAssetPath = func() func(string) string {
	if os.Getenv("CACHE_ASSET_MANIFEST") == "" {
		return nil
	}

	fn, err := assetPathFn()
	if err != nil {
		panic(err)
	}
	return fn
}()

func assetPathHelper() (func(string) string, error) {
	if getAssetPath != nil {
		return getAssetPath, nil
	}
	return assetPathFn()
}

func assetPathFn() (func(string) string, error) {
	manifest, err := loadManifest()
	if err != nil {
		return nil, err
	}

	return func(filename string) string {
		return staticPath + "/" + manifest[filename]
	}, nil
}

func loadManifest() (map[string]string, error) {
	var manifest map[string]string

	file, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(&manifest); err != nil {
		return nil, err
	}
	return manifest, nil
}
