package main

import (
	"encoding/json"
	"os"
)

type TileMapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type TileMapJSON struct {
	Layers []TileMapLayerJSON `json:"layers"`
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
