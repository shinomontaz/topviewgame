package main

import (
	_ "image/png"
	"log"
	"time"
	"topviewgame/controller"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type controllable interface {
	GetDirection() (dx, dy int)
	GetAction() controller.Action
}

type Game struct {
	Map              GameMap
	World            *ecs.Manager
	WorldTags        map[string]ecs.Tag
	Turn             TurnState
	TurnCounter      int
	dt               float64
	last             time.Time
	PlayerController controllable
}

func NewGame() *Game {
	m := NewGameMap()
	world, tags := InitializeWorld(m.CurrentLevel)

	return &Game{
		PlayerController: controller.Human{},
		Map:              m,
		World:            world,
		WorldTags:        tags,
		Turn:             PlayerTurn,
		TurnCounter:      0,
		last:             time.Now(),
	}

}

func (g *Game) Layout(w, h int) (int, int) {
	gd := NewGameData()

	return gd.TileHeight * gd.ScreenWidth, gd.TileWidth * gd.ScreenHeight
}

func (g *Game) Update() error {
	g.dt = time.Since(g.last).Seconds()
	g.last = time.Now()

	UpdateAnimations(g.dt, g)

	g.TurnCounter++
	if g.Turn == PlayerTurn && g.TurnCounter > 10 {
		ProcessPlayer(g)
	}
	if g.Turn == EnemyTurn {
		UpdateMonsters(g)
	}

	g.Turn = PlayerTurn

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	level := g.Map.CurrentLevel
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
