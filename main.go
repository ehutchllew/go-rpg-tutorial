package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/ev-the-dev/rpg-tutorial/animations"
	"github.com/ev-the-dev/rpg-tutorial/entities"
	"github.com/ev-the-dev/rpg-tutorial/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	camera                 *Camera
	colliders              []image.Rectangle
	enemies                []*entities.Enemy
	player                 *entities.Player
	playerRunningAnimation *animations.Animation
	playerSpriteSheet      *spritesheet.SpriteSheet
	potions                []*entities.Potion
	tileMapImg             *ebiten.Image
	tileMapJSON            *TileMapJSON
	tilesets               []Tileset
}

func (g *Game) Update() error {

	g.playerRunningAnimation.Update()

	g.player.Dx = 0.0
	g.player.Dy = 0.0

	// react to key presses
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Dx = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Dy = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Dy = -2
	}

	g.player.X += g.player.Dx
	CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy
	CheckCollisionVertical(g.player.Sprite, g.colliders)

	for _, enemy := range g.enemies {
		enemy.Dx = 0.0
		enemy.Dy = 0.0
		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.Dx += 1
			}
			if enemy.X > g.player.X {
				enemy.Dx -= 1
			}
			if enemy.Y < g.player.Y {
				enemy.Dy += 1
			}
			if enemy.Y > g.player.Y {
				enemy.Dy -= 1
			}
		}

		enemy.X += enemy.Dx
		CheckCollisionHorizontal(enemy.Sprite, g.colliders)

		enemy.Y += enemy.Dy
		CheckCollisionHorizontal(enemy.Sprite, g.colliders)
	}

	for _, potion := range g.potions {
		if g.player.X > potion.X {
			g.player.Health += potion.AmtHeal
			fmt.Printf("Picked up potion! Health: (%d)\n", g.player.Health)
		}
	}

	screenWidth, screenHeight := ebiten.WindowSize()
	g.camera.FollowTarget(
		g.player.X+8, // +8 to center camera on middle of player sprite
		g.player.Y+8,
		float64(screenWidth),
		float64(screenHeight),
	)
	g.camera.Constrain(
		float64(g.tileMapJSON.Layers[0].Width*16.0),
		float64(g.tileMapJSON.Layers[0].Height*16.0),
		float64(screenWidth),
		float64(screenHeight),
	)

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

	g.drawBackground(screen, &opts)

	g.drawPlayer(screen, g.player.Sprite, &opts)

	for _, enemy := range g.enemies {
		g.drawSprite(screen, enemy.Sprite, &opts)
	}

	for _, potion := range g.potions {
		g.drawSprite(screen, potion.Sprite, &opts)
	}

	for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.camera.X),
			float32(collider.Min.Y)+float32(g.camera.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			true,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func (g *Game) drawBackground(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// loop over each layer
	for layerIndex, layer := range g.tileMapJSON.Layers {
		// loop over tiles in layer
		for imgIdx, imgId := range layer.Data {

			if imgId == 0 {
				continue
			}

			// get tile position of tile
			x := imgIdx % layer.Width
			y := imgIdx / layer.Width
			// convert tile position to pixel position
			x *= 16
			y *= 16

			img := g.tilesets[layerIndex].Img(imgId)

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(img, opts)

			opts.GeoM.Reset()

			// // get the position on the TileSet image where the tile ID is
			// srcX := (imgId - 1) % 22 // 22 hardcoded because tileset file shows last index on row as id 21 (0th based)
			// srcY := (imgId - 1) / 22
			// // convert the src tile position to src pixel position
			// srcX *= 16
			// srcY *= 16

			// // draw tile at appropriate x,y position
			// opts.GeoM.Translate(float64(x), float64(y))

			// opts.GeoM.Translate(g.camera.X, g.camera.Y)
			// // draw the tile
			// screen.DrawImage(
			// 	// cropping out the tile we want from the spritesheet
			// 	g.tileMapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
			// 	opts,
			// )
			// // reset the opts for the next tile
			// opts.GeoM.Reset()
		}
	}
}

// Temp
func (g *Game) drawPlayer(screen *ebiten.Image, sprite *entities.Sprite, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(sprite.X, sprite.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)

	screen.DrawImage(sprite.Img.SubImage(
		g.playerSpriteSheet.Rect(g.playerRunningAnimation.Frame()),
	).(*ebiten.Image),
		opts)

	opts.GeoM.Reset()
}

func (g *Game) drawSprite(screen *ebiten.Image, sprite *entities.Sprite, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(sprite.X, sprite.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)

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
		log.Fatalf("playerImg err: %v", err)
	}

	playerRunningAnim := animations.NewAnimation(4, 12, 4, 20.0)

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)

	potionImg, _, err := ebitenutil.NewImageFromFile("./assets/images/heart_potion.png")
	if err != nil {
		log.Fatalf("potionImg err: %v", err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("./assets/images/skeleton.png")
	if err != nil {
		log.Fatalf("skeletonImg err: %v", err)
	}

	tileMapImg, _, err := ebitenutil.NewImageFromFile("./assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatalf("tileMapImg err: %v", err)
	}

	tileMapJson, err := NewTileMapJSON("./assets/maps/spawn.json")
	if err != nil {
		log.Fatalf("tileMapJson err: %v", err)
	}

	tilesets, err := tileMapJson.GenTilesets()
	if err != nil {
		log.Fatalf("tilesets err: %v", err)
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&Game{
		camera: NewCamera(0.0, 0.0),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   50,
				Y:   50,
			},
			Health: 3,
		},
		playerRunningAnimation: playerRunningAnim,
		playerSpriteSheet:      playerSpriteSheet,
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   100,
					Y:   100,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   150,
					Y:   150,
				},
				FollowsPlayer: false,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   210,
					Y:   100,
				},
				AmtHeal: 1,
			},
		},
		tileMapImg:  tileMapImg,
		tileMapJSON: tileMapJson,
		tilesets:    tilesets,
	}); err != nil {
		log.Fatal(err)
	}
}

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X+16), int(sprite.Y+16))) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - 16
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X+16), int(sprite.Y+16))) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - 16
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}
