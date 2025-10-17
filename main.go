package main

import (
	_ "image/png"
	"log"
	"time"
	"topviewgame/controller"
	"topviewgame/event"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type controllable interface {
	GetEvent() event.Event
	GetCursor() (int, int)
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
	gm               *GameMap
	Center           Position
	viewport         Rect

	// Auto-movement system
	AutoMovePath   []Position
	AutoMoveIndex  int
	IsAutoMoving   bool
	AutoMoveTimer  float64
	AutoMoveDelay  float64
	PathVisualizer *PathVisualizer
}

func NewGame() *Game {
	gd := NewGameData()
	m := NewGameMap(gd)
	world, tags := InitializeWorld(m)

	pathVis, err := NewPathVisualizer()
	if err != nil {
		log.Println("Path visualization disabled:", err)
	}

	return &Game{
		PlayerController: controller.Human{TileWidth: gd.TileWidth, TileHeight: gd.TileHeight},
		Map:              m,
		World:            world,
		WorldTags:        tags,
		Turn:             PlayerTurn,
		TurnCounter:      0,
		last:             time.Now(),
		gd:               &gd,
		gm:               &m,
		AutoMoveDelay:    0.2, // 200ms delay between auto-movement steps
		PathVisualizer:   pathVis,
	}
}

func (g *Game) SetCenter(pos Position) {
	g.Center = pos
	playerX, playerY := g.Center.X, g.Center.Y
	g.viewport = Rect{
		X1: playerX - g.gd.ScreenWidth/2,
		Y1: playerY - g.gd.ScreenHeight/2,
		X2: playerX + g.gd.ScreenWidth/2,
		Y2: playerY + g.gd.ScreenHeight/2,
	}
}

func (g *Game) Viewport() Rect {
	return g.viewport
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
	UpdateCursor(g)
	g.UpdateAutoMoveTimer(g.dt) // Update auto-movement timer

	g.TurnCounter++
	if g.Turn == PlayerTurn && g.TurnCounter > 10 {
		UpdatePlayer(g)
	}
	if g.Turn == EnemyTurn {
		UpdateMonsters(g)
	}

	g.Turn = PlayerTurn

	return nil
}

func (g *Game) GetViewport() Rect {
	return g.viewport
}

// Auto-movement methods
func (g *Game) StartAutoMove(path []Position) {
	if len(path) > 1 {
		g.AutoMovePath = path
		g.AutoMoveIndex = 1 // Start from index 1 (0 is current position)
		g.IsAutoMoving = true
		g.AutoMoveTimer = 0 // Reset timer
	}
}

func (g *Game) StopAutoMove() {
	g.IsAutoMoving = false
	g.AutoMovePath = nil
	g.AutoMoveIndex = 0
	g.AutoMoveTimer = 0
}

func (g *Game) UpdateAutoMoveTimer(dt float64) {
	if g.IsAutoMoving {
		g.AutoMoveTimer += dt
	}
}

func (g *Game) CanAutoMoveStep() bool {
	return g.IsAutoMoving && g.AutoMoveTimer >= g.AutoMoveDelay
}

func (g *Game) GetNextAutoMoveStep() (int, int, bool) {
	if !g.IsAutoMoving || g.AutoMoveIndex >= len(g.AutoMovePath) {
		g.StopAutoMove()
		return 0, 0, false
	}

	nextPos := g.AutoMovePath[g.AutoMoveIndex]
	dx := nextPos.X - g.Center.X
	dy := nextPos.Y - g.Center.Y
	g.AutoMoveIndex++
	g.AutoMoveTimer = 0 // Reset timer for next step

	return dx, dy, true
}

func (g *Game) IsEnemyInSight() bool {
	level := g.Map.CurrentLevel
	for _, result := range g.World.Query(g.WorldTags["monsters"]) {
		pos := result.Components[positionC].(*Position)
		monster := result.Components[monsterC].(*Monster)

		// Skip dead monsters
		if monster.IsDead() {
			continue
		}

		// Check if monster is visible to player
		if level.PlayerVisible.IsVisible(pos.X, pos.Y) {
			return true
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	level := g.Map.CurrentLevel
	level.Draw(screen, g.viewport)

	DrawRenderables(g, level, screen, g.viewport)
	DrawPath(g, screen) // Draw path before cursor
	DrawUserLog(g, screen)
	DrawHUD(g, screen)
	DrawCursor(g, screen)
}

func main() {
	g := NewGame()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Topview game")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
