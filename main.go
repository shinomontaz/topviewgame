package main

import (
	"fmt"
	_ "image/png"
	"log"
	"time"
	"topviewgame/controller"
	"topviewgame/event"
	"topviewgame/internal/world"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

type controllable interface {
	GetEvent() event.Event
	GetCursor() (int, int)
}

type Game struct {
	Map              GameMap
	World            *world.World
	Turn             TurnState
	TurnCounter      int
	dt               float64
	last             time.Time
	PlayerController controllable
	gd               *GameData
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
	pathVis, err := NewPathVisualizer()
	if err != nil {
		log.Println("Path visualization disabled:", err)
	}
	m := NewGameMap(gd)
	gameWorld := world.NewWorld()

	// Initialize entities using the old logic for now
	// This will be refactored in the next step
	InitializeWorldEntities(gameWorld, m)

	return &Game{
		PlayerController: controller.Human{TileWidth: gd.TileWidth, TileHeight: gd.TileHeight},
		Map:              m,
		World:            gameWorld,
		Turn:             PlayerTurn,
		TurnCounter:      0,
		last:             time.Now(),
		gd:               &gd,
		AutoMoveDelay:    0.2, // 200ms delay between auto-movement steps
		PathVisualizer:   pathVis,
	}
}

// InitializeWorldEntities initializes game entities using the new World
// This is a temporary bridge function that will be refactored
func InitializeWorldEntities(w *world.World, gm GameMap) {
	startingRoom := gm.CurrentLevel.Rooms[0]

	stairsComponent, err := NewStairs(gm.Gd.TileWidth, gm.Gd.TileHeight)
	if err != nil {
		log.Fatal(err)
	}

	stairsComponent.NextLevel = gm.LevelIndex + 1
	stairsPos := gm.CurrentLevel.StairsPos

	playerX, playerY := startingRoom.Center()
	player := NewPlayer()
	w.NewEntity().
		AddComponent(w.StairsComponent(), stairsComponent).
		AddComponent(w.RenderableComponent(), stairsComponent).
		AddComponent(w.PositionComponent(), &Position{X: stairsPos.X, Y: stairsPos.Y})

	w.NewEntity().
		AddComponent(w.PlayerComponent(), player).
		AddComponent(w.RenderableComponent(), player).
		AddComponent(w.PositionComponent(), &Position{X: playerX, Y: playerY}).
		AddComponent(w.HealthComponent(), &world.Health{Max: 30, Current: 30}).
		AddComponent(w.MeleeWeaponComponent(), &world.MeleeWeapon{
			Name:       "Fist",
			MinDamage:  1,
			MaxDamage:  3,
			ToHitBonus: 2,
		}).
		AddComponent(w.ArmorComponent(), &world.Armor{
			Name:    "Burlap Sack",
			Defence: 1,
			Dodge:   1,
		}).
		AddComponent(w.NameComponent(), &world.Name{Label: "Player"}).
		AddComponent(w.UserMessageComponent(), &world.UserMessage{
			AttackMessage:    "",
			DeadMessage:      "",
			GameStateMessage: "",
		})

	for _, room := range gm.CurrentLevel.Rooms {
		if room.X1 != startingRoom.X1 {
			var (
				monsterType MonsterType
				monsterName string
			)
			switch rnd.Intn(2) {
			case 0:
				monsterType = SKELETON
				monsterName = "Skeleton"
			case 1:
				monsterType = ZOMBIE
				monsterName = "Zombie"
			}
			monster := NewMonster(monsterType)
			mX, mY := room.Center()
			pos := Position{X: mX, Y: mY}
			ent := w.NewEntity().
				AddComponent(w.MonsterComponent(), monster).
				AddComponent(w.RenderableComponent(), monster).
				AddComponent(w.PositionComponent(), &pos).
				AddComponent(w.NameComponent(), &world.Name{Label: monsterName}).
				AddComponent(w.UserMessageComponent(), &world.UserMessage{
					AttackMessage:    "",
					DeadMessage:      "",
					GameStateMessage: "",
				})

			if monsterType == SKELETON {
				ent.AddComponent(w.HealthComponent(), &world.Health{
					Max:     10,
					Current: 2,
				}).AddComponent(w.MeleeWeaponComponent(), &world.MeleeWeapon{
					Name:       "Short Sword",
					MinDamage:  2,
					MaxDamage:  6,
					ToHitBonus: 0,
				}).AddComponent(w.ArmorComponent(), &world.Armor{
					Name:    "No armor",
					Defence: 0,
					Dodge:   5,
				})
			} else {
				ent.AddComponent(w.HealthComponent(), &world.Health{
					Max:     20,
					Current: 2,
				}).AddComponent(w.MeleeWeaponComponent(), &world.MeleeWeapon{
					Name:       "Khopesh",
					MinDamage:  1,
					MaxDamage:  4,
					ToHitBonus: 1,
				}).AddComponent(w.ArmorComponent(), &world.Armor{
					Name:    "Rotten rags",
					Defence: 1,
					Dodge:   0,
				})
			}
			gm.updateMonsterPosition(ent, nil, &pos)
		}
	}

	cursor, err := NewCursor(gm.Gd.TileWidth, gm.Gd.TileHeight)
	if err != nil {
		log.Fatal(err)
	}
	w.NewEntity().AddComponent(w.CursorComponent(), cursor)
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
	for _, result := range g.World.QueryMonsters() {
		pos := g.World.GetPosition(result)
		monster := g.World.GetMonster(result).(*Monster)

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

func (g *Game) NextLevel(level int) {
	fmt.Printf("Going to next level %d\n", level)

	g.Map.LevelIndex = level
	g.Map.CurrentLevel = NewLevel(*g.gd)
	g.Map.monsterPositions = make(map[Position]*ecs.Entity)

	gameWorld := world.NewWorld()

	InitializeWorldEntities(gameWorld, g.Map)

	g.StopAutoMove()

	g.World = gameWorld
	g.TurnCounter = 0
	g.last = time.Now()
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
