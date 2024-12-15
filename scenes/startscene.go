package scenes

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type StartScene struct {
	loaded bool
}

func NewStartScene() *StartScene {
	return &StartScene{}
}

func (s *StartScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255})
	ebitenutil.DebugPrint(screen, "Press enter to start")
}

func (s *StartScene) FirstLoad() {
}

func (s *StartScene) IsLoaded() bool {
	return s.loaded
}

func (s *StartScene) OnEnter() {
}

func (s *StartScene) OnExit() {
}

func (s *StartScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	return StartSceneId
}

var _ Scene = (*StartScene)(nil)
