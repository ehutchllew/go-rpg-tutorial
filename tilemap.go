package main

import (
	"encoding/json"
	"os"
	"path"
)

type TileMapLayerJSON struct {
	Data   []int  `json:"data"`
	Height int    `json:"height"`
	Name   string `json:"name"`
	Width  int    `json:"width"`
}

type TileMapJSON struct {
	Layers   []TileMapLayerJSON `json:"layers"`
	Tilesets []map[string]any   `json:"tilesets"`
}

func (t *TileMapJSON) GenTilesets() ([]Tileset, error) {

	tilesets := make([]Tileset, 0)
	for _, tilesetData := range t.Tilesets {
		tilesetPath := path.Join("assets/maps/", tilesetData["source"].(string))
		tileset, err := NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}

		tilesets = append(tilesets, tileset)
	}

	return tilesets, nil
}

func NewTileMapJSON(filepath string) (*TileMapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tileMapJson TileMapJSON
	err = json.Unmarshal(contents, &tileMapJson)
	if err != nil {
		return nil, err
	}

	return &tileMapJson, nil
}
