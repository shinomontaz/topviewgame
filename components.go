package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	core "topviewgame/internal/core"
)

type Position = core.Position

type Renderable interface {
	GetImage() *ebiten.Image
}

type Movable struct {
}

type Name struct {
	Label string
}

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

type UserMessage struct {
	AttackMessage    string
	DeadMessage      string
	GameStateMessage string
}
