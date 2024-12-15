package main

import (
	"github.com/ev-the-dev/rpg-tutorial/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	activeSceneId scenes.SceneId
	sceneMap      map[scenes.SceneId]scenes.Scene
}

func NewGame() *Game {
	activeSceneId := scenes.StartSceneId
	sceneMap := map[scenes.SceneId]scenes.Scene{
		scenes.GameSceneId:  scenes.NewGameScene(),
		scenes.PauseSceneId: scenes.NewPauseScene(),
		scenes.StartSceneId: scenes.NewStartScene(),
	}
	sceneMap[activeSceneId].FirstLoad()

	return &Game{
		activeSceneId,
		sceneMap,
	}
}

func (g *Game) Update() error {
	nextSceneId := g.sceneMap[g.activeSceneId].Update()

	if nextSceneId == scenes.ExitSceneId {
		g.sceneMap[g.activeSceneId].OnExit()
		return ebiten.Termination
	}
	// If true, game switched scenes
	if nextSceneId != g.activeSceneId {
		nextScene := g.sceneMap[nextSceneId]
		// check if scene loaded already
		if !nextScene.IsLoaded() {
			nextScene.FirstLoad()
		}

		nextScene.OnEnter()
		g.sceneMap[g.activeSceneId].OnExit()
	}

	g.activeSceneId = nextSceneId
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneMap[g.activeSceneId].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
