package main

import (
	_ "image/png"
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Map       GameMap
	World     *ecs.Manager
	WorldTags map[string]ecs.Tag
}

func NewGame() *Game {
	world, tags := InitializeWorld()
	return &Game{
		Map:       NewGameMap(),
		World:     world,
		WorldTags: tags,
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	gd := NewGameData()
	return gd.TileHeight * gd.ScreenWidth, gd.TileWidth * gd.ScreenHeight
}

func (g *Game) Update() error {
	TryMoveplayer(g)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	level := g.Map.Dungeons[0].Levels[0]
	level.Draw(screen)

	ProcessRenderables(g, level, screen)
}

func main() {
	g := NewGame()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Topview game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
