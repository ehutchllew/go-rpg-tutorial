package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/ev-the-dev/rpg-tutorial/animations"
	"github.com/ev-the-dev/rpg-tutorial/components"
	"github.com/ev-the-dev/rpg-tutorial/constants"
	"github.com/ev-the-dev/rpg-tutorial/entities"
	"github.com/ev-the-dev/rpg-tutorial/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	camera            *Camera
	colliders         []image.Rectangle
	enemies           []*entities.Enemy
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	potions           []*entities.Potion
	tileMapImg        *ebiten.Image
	tileMapJSON       *TileMapJSON
	tilesets          []Tileset
}

func NewGame() *Game {
	playerImg, _, err := ebitenutil.NewImageFromFile("./assets/images/ninja.png")
	if err != nil {
		log.Fatalf("playerImg err: %v", err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, constants.Tilesize)

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

	return &Game{
		camera: NewCamera(0.0, 0.0),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
		player: &entities.Player{
			Animations: map[entities.PlayerState]*animations.Animation{
				entities.Up:    animations.NewAnimation(5, 13, 4, 20.0),
				entities.Down:  animations.NewAnimation(4, 12, 4, 20.0),
				entities.Left:  animations.NewAnimation(6, 14, 4, 20.0),
				entities.Right: animations.NewAnimation(7, 15, 4, 20.0),
			},
			CombatComp: components.NewBasicCombat(1, 3),
			Health:     3,
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   50,
				Y:   50,
			},
		},
		playerSpriteSheet: playerSpriteSheet,
		enemies: []*entities.Enemy{
			{
				CombatComp:    components.NewBasicCombat(1, 3),
				FollowsPlayer: true,
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   100,
					Y:   100,
				},
			},
			{
				CombatComp:    components.NewBasicCombat(1, 3),
				FollowsPlayer: false,
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   150,
					Y:   150,
				},
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
	}
}

func (g *Game) Update() error {

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
	checkCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy
	checkCollisionVertical(g.player.Sprite, g.colliders)

	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		activeAnim.Update()
	}

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
		checkCollisionHorizontal(enemy.Sprite, g.colliders)

		enemy.Y += enemy.Dy
		checkCollisionHorizontal(enemy.Sprite, g.colliders)
	}

	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	cX, cY := ebiten.CursorPosition()

	deadEnemies := make(map[int]struct{})
	for enemyIndex, enemy := range g.enemies {
		rect := image.Rect(
			int(enemy.X),
			int(enemy.Y),
			int(enemy.X)+constants.Tilesize,
			int(enemy.Y)+constants.Tilesize,
		)

		// is cursor within rect?
		if cX > rect.Min.X && cX <= rect.Max.X && cY > rect.Min.Y && cY <= rect.Max.Y {
			if clicked {
				fmt.Println("Damaging Enemy")
				enemy.CombatComp.Damage(g.player.CombatComp.AttackPower())

				if enemy.CombatComp.Health() <= 0 {
					fmt.Println("Enemy Eliminated")
					deadEnemies[enemyIndex] = struct{}{}
				}
			}
		}
	}
	if len(deadEnemies) > 0 {
		newEnemies := make([]*entities.Enemy, 0)
		for index, enemy := range g.enemies {
			if _, exists := deadEnemies[index]; !exists {
				newEnemies = append(newEnemies, enemy)
			}
		}
		g.enemies = newEnemies
	}

	// for _, potion := range g.potions {
	// 	if g.player.X > potion.X {
	// 		g.player.Health += potion.AmtHeal
	// 		fmt.Printf("Picked up potion! Health: (%d)\n", g.player.Health)
	// 	}
	// }

	screenWidth, screenHeight := ebiten.WindowSize()
	g.camera.FollowTarget(
		g.player.X+8, // +8 to center camera on middle of player sprite
		g.player.Y+8,
		float64(screenWidth),
		float64(screenHeight),
	)
	g.camera.Constrain(
		float64(g.tileMapJSON.Layers[0].Width*constants.Tilesize),
		float64(g.tileMapJSON.Layers[0].Height*constants.Tilesize),
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
			x *= constants.Tilesize
			y *= constants.Tilesize

			img := g.tilesets[layerIndex].Img(imgId)

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + constants.Tilesize))

			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(img, opts)

			opts.GeoM.Reset()

			// // get the position on the TileSet image where the tile ID is
			// srcX := (imgId - 1) % 22 // 22 hardcoded because tileset file shows last index on row as id 21 (0th based)
			// srcY := (imgId - 1) / 22
			// // convert the src tile position to src pixel position
			// srcX *= constants.Tilesize
			// srcY *= constants.Tilesize

			// // draw tile at appropriate x,y position
			// opts.GeoM.Translate(float64(x), float64(y))

			// opts.GeoM.Translate(g.camera.X, g.camera.Y)
			// // draw the tile
			// screen.DrawImage(
			// 	// cropping out the tile we want from the spritesheet
			// 	g.tileMapImg.SubImage(image.Rect(srcX, srcY, srcX+constants.Tilesize, srcY+constants.Tilesize)).(*ebiten.Image),
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

	playerFrame := 0
	activeAnim := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnim != nil {
		playerFrame = activeAnim.Frame()
	}

	screen.DrawImage(sprite.Img.SubImage(
		g.playerSpriteSheet.Rect(playerFrame),
	).(*ebiten.Image),
		opts)

	opts.GeoM.Reset()
}

func (g *Game) drawSprite(screen *ebiten.Image, sprite *entities.Sprite, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(sprite.X, sprite.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)

	screen.DrawImage(sprite.Img.SubImage(
		image.Rect(0, 0, constants.Tilesize, constants.Tilesize),
	).(*ebiten.Image),
		opts)

	opts.GeoM.Reset()
}

func checkCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X+constants.Tilesize), int(sprite.Y+constants.Tilesize))) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - constants.Tilesize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func checkCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X+constants.Tilesize), int(sprite.Y+constants.Tilesize))) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - constants.Tilesize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}
