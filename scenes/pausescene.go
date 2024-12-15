package scenes

import (
	"image/color"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	loaded bool
}

func NewPauseScene() *PauseScene {
	return &PauseScene{}
}

func (s *PauseScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 255, 0, 255})
	ebitenutil.DebugPrint(screen, "Press enter to unpause.")
}

func (s *PauseScene) FirstLoad() {
}

func (s *PauseScene) IsLoaded() bool {
	return s.loaded
}

func (s *PauseScene) OnEnter() {
}

func (s *PauseScene) OnExit() {
}

func (s *PauseScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	return PauseSceneId
}

var _ Scene = (*PauseScene)(nil)
