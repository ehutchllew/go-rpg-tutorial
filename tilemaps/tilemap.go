package tilemaps

import (
	"encoding/json"
	"os"
	"path"

	"github.com/ev-the-dev/rpg-tutorial/tilesets"
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

func (t *TileMapJSON) GenTilesets() ([]tilesets.Tileset, error) {

	ts := make([]tilesets.Tileset, 0)
	for _, tilesetData := range t.Tilesets {
		tilesetPath := path.Join("assets/maps/", tilesetData["source"].(string))
		tileset, err := tilesets.NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}

		ts = append(ts, tileset)
	}

	return ts, nil
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
