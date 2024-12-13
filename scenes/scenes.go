package scenes

import "github.com/hajimehoshi/ebiten/v2"

type SceneId uint

const (
	GameSceneId SceneId = iota
	StartSceneId
)

type Scene interface {
	Draw(screen *ebiten.Image)
	FirstLoad()
	OnEnter()
	OnExit()
	Update() SceneId
}
