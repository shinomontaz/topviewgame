package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Position struct {
	X, Y int
}

func (p *Position) GetManhattanDistance(other *Position) int {
	xDist := math.Abs(float64(p.X - other.X))
	yDist := math.Abs(float64(p.Y - other.Y))
	return int(xDist) + int(yDist)
}

func (p *Position) IsEqual(other *Position) bool {
	return (p.X == other.X && p.Y == other.Y)
}

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
