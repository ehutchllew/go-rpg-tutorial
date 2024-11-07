package main

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
}

type UniformTileset struct {
	gid int
	img *ebiten.Image
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid
	// get the position on the TileSet image where the tile ID is
	srcX := (id - 1) % 22 // 22 hardcoded because tileset file shows last index on row as id 21 (0th based)
	srcY := (id - 1) / 22
	// convert the src tile position to src pixel position
	srcX *= 16
	srcY *= 16

	return u.img.SubImage(
		image.Rect(srcX, srcY, srcX+16, srcY+16),
	).(*ebiten.Image)
}

type TileJSON struct {
	Height int    `json:"imageheight"`
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
}

type DynamicTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynamicTileset struct {
	gid  int
	imgs []*ebiten.Image
}

func (d *DynamicTileset) Img(id int) *ebiten.Image {
	id -= d.gid

	return d.imgs[id]
}

func NewTileset(path string, gid int) (Tileset, error) {
	// temporary

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "TilesetBuilding") {
		// return dynamic tileset
		var dynamicTilesetJson DynamicTilesetJSON
		err = json.Unmarshal(content, &dynamicTilesetJson)
		if err != nil {
			return nil, err
		}

		dynamicTileset := DynamicTileset{
			gid:  gid,
			imgs: make([]*ebiten.Image, 0),
		}

		for _, tileJSON := range dynamicTilesetJson.Tiles {

			tileJSONPath := filepath.Clean(tileJSON.Path)
			tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = filepath.Join("assets/", tileJSONPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				return nil, err
			}

			dynamicTileset.imgs = append(dynamicTileset.imgs, img)
		}

		return &dynamicTileset, nil
	}

	// return inform tileset
	var uniformTilesetJson UniformTilesetJSON
	err = json.Unmarshal(content, &uniformTilesetJson)
	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{gid: gid}
	tileJSONPath := filepath.Clean(uniformTilesetJson.Path)
	tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = filepath.Join("assets/", tileJSONPath)
	img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
	if err != nil {
		return nil, err
	}

	uniformTileset.img = img

	return &uniformTileset, nil
}
