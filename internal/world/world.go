package world

import (
	"topviewgame/internal/core"

	"github.com/bytearena/ecs"
)

// World provides a clean facade over the ECS system
// Encapsulates all ECS complexity and provides typed, convenient methods
type World struct {
	manager *ecs.Manager
	tags    map[string]ecs.Tag

	// Private component references - not exposed to client code
	positionC    *ecs.Component
	renderableC  *ecs.Component
	playerC      *ecs.Component
	monsterC     *ecs.Component
	healthC      *ecs.Component
	meleeWeaponC *ecs.Component
	armorC       *ecs.Component
	nameC        *ecs.Component
	cursorC      *ecs.Component
	userMessageC *ecs.Component
	stairsC      *ecs.Component
}

// NewWorld creates a new ECS world with all components initialized
func NewWorld() *World {
	manager := ecs.NewManager()

	w := &World{
		manager: manager,
		tags:    make(map[string]ecs.Tag),
	}

	// Initialize all components (private)
	w.playerC = manager.NewComponent()
	w.positionC = manager.NewComponent()
	w.renderableC = manager.NewComponent()
	_ = manager.NewComponent() // movableC - not used currently
	w.monsterC = manager.NewComponent()
	w.healthC = manager.NewComponent()
	w.meleeWeaponC = manager.NewComponent()
	w.armorC = manager.NewComponent()
	w.nameC = manager.NewComponent()
	w.userMessageC = manager.NewComponent()
	w.cursorC = manager.NewComponent()
	w.stairsC = manager.NewComponent()

	// Build tags (private)
	w.tags["monsters"] = ecs.BuildTag(w.monsterC, w.positionC, w.healthC, w.meleeWeaponC, w.armorC, w.nameC, w.userMessageC)
	w.tags["players"] = ecs.BuildTag(w.playerC, w.positionC, w.healthC, w.meleeWeaponC, w.armorC, w.nameC, w.userMessageC)
	w.tags["renderables"] = ecs.BuildTag(w.renderableC, w.positionC)
	w.tags["messengers"] = ecs.BuildTag(w.userMessageC)
	w.tags["cursors"] = ecs.BuildTag(w.cursorC)
	w.tags["stairs"] = ecs.BuildTag(w.stairsC, w.positionC, w.renderableC)

	return w
}

// Facade methods - clean API for client code

// QueryPlayers returns all player entities with their components
func (w *World) QueryPlayers() []*ecs.QueryResult {
	return w.manager.Query(w.tags["players"])
}

// QueryMonsters returns all monster entities with their components
func (w *World) QueryMonsters() []*ecs.QueryResult {
	return w.manager.Query(w.tags["monsters"])
}

// QueryRenderables returns all renderable entities
func (w *World) QueryRenderables() []*ecs.QueryResult {
	return w.manager.Query(w.tags["renderables"])
}

// QueryMessengers returns all entities with user messages
func (w *World) QueryMessengers() []*ecs.QueryResult {
	return w.manager.Query(w.tags["messengers"])
}

// QueryCursors returns all cursor entities
func (w *World) QueryCursors() []*ecs.QueryResult {
	return w.manager.Query(w.tags["cursors"])
}

func (w *World) QueryStairs() []*ecs.QueryResult {
	return w.manager.Query(w.tags["stairs"])
}

// GetPlayerCenter returns the center position of the player
func (w *World) GetPlayerCenter() core.Position {
	for _, entity := range w.QueryPlayers() {
		pos := entity.Components[w.positionC].(*core.Position)
		return *pos
	}
	return core.Position{}
}

// Component access methods - typed and safe
func (w *World) GetPosition(entity *ecs.QueryResult) *core.Position {
	return entity.Components[w.positionC].(*core.Position)
}

func (w *World) GetHealth(entity *ecs.QueryResult) *Health {
	return entity.Components[w.healthC].(*Health)
}

func (w *World) GetArmor(entity *ecs.QueryResult) *Armor {
	return entity.Components[w.armorC].(*Armor)
}

func (w *World) GetMeleeWeapon(entity *ecs.QueryResult) *MeleeWeapon {
	return entity.Components[w.meleeWeaponC].(*MeleeWeapon)
}

func (w *World) GetName(entity *ecs.QueryResult) *Name {
	return entity.Components[w.nameC].(*Name)
}

func (w *World) GetUserMessage(entity *ecs.QueryResult) *UserMessage {
	return entity.Components[w.userMessageC].(*UserMessage)
}

func (w *World) GetRenderable(entity *ecs.QueryResult) interface{} {
	return entity.Components[w.renderableC]
}

func (w *World) GetPlayer(entity *ecs.QueryResult) interface{} {
	return entity.Components[w.playerC]
}

func (w *World) GetMonster(entity *ecs.QueryResult) interface{} {
	return entity.Components[w.monsterC]
}

func (w *World) GetCursor(entity *ecs.QueryResult) interface{} {
	return entity.Components[w.cursorC]
}

func (w *World) GetStairs(entity *ecs.QueryResult) interface{} {
	return entity.Components[w.stairsC]
}

// Entity creation method
func (w *World) NewEntity() *ecs.Entity {
	return w.manager.NewEntity()
}

// Component access for entity creation (temporary - will be replaced with builder pattern)
func (w *World) PlayerComponent() *ecs.Component      { return w.playerC }
func (w *World) PositionComponent() *ecs.Component    { return w.positionC }
func (w *World) RenderableComponent() *ecs.Component  { return w.renderableC }
func (w *World) MonsterComponent() *ecs.Component     { return w.monsterC }
func (w *World) HealthComponent() *ecs.Component      { return w.healthC }
func (w *World) MeleeWeaponComponent() *ecs.Component { return w.meleeWeaponC }
func (w *World) ArmorComponent() *ecs.Component       { return w.armorC }
func (w *World) NameComponent() *ecs.Component        { return w.nameC }
func (w *World) UserMessageComponent() *ecs.Component { return w.userMessageC }
func (w *World) CursorComponent() *ecs.Component      { return w.cursorC }
func (w *World) StairsComponent() *ecs.Component      { return w.stairsC }

// Component types - moved from main package for better organization
type Health struct {
	Current, Max int
}

type MeleeWeapon struct {
	Name       string
	MinDamage  int
	MaxDamage  int
	ToHitBonus int
}

type Armor struct {
	Name    string
	Defence int
	Dodge   int
	Block   int
}

type Name struct {
	Label string
}

type UserMessage struct {
	AttackMessage    string
	DeadMessage      string
	GameStateMessage string
}
