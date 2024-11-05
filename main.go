package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Game struct {
	player      *Player
	enemies     []*Enemy
	potions     []*Potion
	tileMapJSON *TileMapJSON
	tileMapImg  *ebiten.Image
}

type Player struct {
	*Sprite
	Health uint
}

type Potion struct {
	*Sprite
	AmtHeal uint
}

type Sprite struct {
	Img *ebiten.Image
	X   float64
	Y   float64
}

func (g *Game) Update() error {
	// react to key presses
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Y -= 2
	}

	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.X += 1
			}
			if enemy.X > g.player.X {
				enemy.X -= 1
			}
			if enemy.Y < g.player.Y {
				enemy.Y += 1
			}
			if enemy.Y > g.player.Y {
				enemy.Y -= 1
			}
		}
	}

	for _, potion := range g.potions {
		if g.player.X > potion.X {
			g.player.Health += potion.AmtHeal
			fmt.Printf("Picked up potion! Health: (%d)\n", g.player.Health)
		}
	}

	return nil
}

/*
* NOTE: When drawing, assets/sprites get 'layered'.
* so, in order for something to appear on top it must
* be drawn after.
* Example: background drawn first, then trees, then
* player.
 */
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := ebiten.DrawImageOptions{}

	// loop over each layer
	for _, layer := range g.tileMapJSON.Layers {
		// loop over tiles in layer
		for imgIdx, imgId := range layer.Data {
			// get tile position of tile
			x := imgIdx % layer.Width
			y := imgIdx / layer.Width
			// convert tile position to pixel position
			x *= 16
			y *= 16

			// get the position on the TileSet image where the tile ID is
			srcX := (imgId - 1) % 22 // 22 hardcoded because tileset file shows last index on row as id 21 (0th based)
			srcY := (imgId - 1) / 22
			// convert the src tile position to src pixel position
			srcX *= 16
			srcY *= 16

			// draw tile at appropriate x,y position
			opts.GeoM.Translate(float64(x), float64(y))
			// draw the tile
			screen.DrawImage(
				// cropping out the tile we want from the spritesheet
				g.tileMapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)
			// reset the opts for the next tile
			opts.GeoM.Reset()
		}
	}

	drawSprite(screen, g.player.Sprite, &opts)

	for _, enemy := range g.enemies {
		drawSprite(screen, enemy.Sprite, &opts)
	}

	for _, potion := range g.potions {
		drawSprite(screen, potion.Sprite, &opts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func drawSprite(screen *ebiten.Image, sprite *Sprite, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(sprite.X, sprite.Y)

	screen.DrawImage(sprite.Img.SubImage(
		image.Rect(0, 0, 16, 16),
	).(*ebiten.Image),
		opts)

	opts.GeoM.Reset()
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("./assets/images/ninja.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("./assets/images/heart_potion.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("./assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	tileMapImg, _, err := ebitenutil.NewImageFromFile("./assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	tileMapJson, err := NewTileMapJSON("./assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&Game{
		player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				X:   50,
				Y:   50,
			},
			Health: 3,
		},
		enemies: []*Enemy{
			{
				Sprite: &Sprite{
					Img: skeletonImg,
					X:   100,
					Y:   100,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &Sprite{
					Img: skeletonImg,
					X:   150,
					Y:   150,
				},
				FollowsPlayer: false,
			},
		},
		potions: []*Potion{
			{
				Sprite: &Sprite{
					Img: potionImg,
					X:   210,
					Y:   100,
				},
				AmtHeal: 1,
			},
		},
		tileMapJSON: tileMapJson,
		tileMapImg:  tileMapImg,
	}); err != nil {
		log.Fatal(err)
	}
}
