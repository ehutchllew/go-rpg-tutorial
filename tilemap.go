package main

import (
	"encoding/json"
	"os"
)

type TileMapLayerJSON struct {
	Data   []int  `json:"data"`
	Height int    `json:"height"`
	Name   string `json:"name"`
	Width  int    `json:"width"`
}

type TileMapJSON struct {
	Layers   []TileMapLayerJSON `json:"layers"`
	Tilesets []*TilesetJSON
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
