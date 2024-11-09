package spritesheet

import "image"

type SpriteSheet struct {
	HeightInTiles int
	TileSize      int
	WidthInTiles  int
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := (index % s.WidthInTiles) * s.TileSize
	y := (index / s.WidthInTiles) * s.TileSize

	return image.Rect(x, y, x+s.TileSize, y+s.TileSize)
}

func NewSpriteSheet(w, h, t int) *SpriteSheet {
	return &SpriteSheet{
		HeightInTiles: h,
		TileSize:      t,
		WidthInTiles:  w,
	}
}
