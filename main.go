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
	gd               *GameData
	Center           Position
}

func NewGame() *Game {
	gd := NewGameData()
	m := NewGameMap(gd)
	world, tags := InitializeWorld(m.CurrentLevel)

	return &Game{
		PlayerController: controller.Human{},
		Map:              m,
		World:            world,
		WorldTags:        tags,
		Turn:             PlayerTurn,
		TurnCounter:      0,
		last:             time.Now(),
		gd:               &gd,
		Center:           GetCenter(world, tags["players"]),
	}
}

func (g *Game) GetData() *GameData {
	return g.gd
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.gd.TileHeight * g.gd.ScreenWidth, g.gd.TileWidth * (g.gd.ScreenHeight + g.gd.UIHeight)
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
	g.Center = GetCenter(g.World, g.WorldTags["players"])
	playerX, playerY := g.Center.X, g.Center.Y

	// Calculate the viewport bounds
	viewport := Rect{
		X1: playerX - g.gd.ScreenWidth/2,
		Y1: playerY - g.gd.ScreenHeight/2,
		X2: playerX + g.gd.ScreenWidth/2,
		Y2: playerY + g.gd.ScreenHeight/2,
	}

	level := g.Map.CurrentLevel
	level.Draw(screen, viewport) //viewportLeft, viewportTop, viewportRight, viewportBottom)

	ProcessRenderables(g, level, screen, viewport) //Left, viewportTop, viewportRight, viewportBottom)
	ProcessUserLog(g, screen)
	ProcessHUD(g, screen)
}

func main() {
	g := NewGame()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Topview game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
