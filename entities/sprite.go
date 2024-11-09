package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Dx  float64
	Dy  float64
	Img *ebiten.Image
	X   float64
	Y   float64
}
