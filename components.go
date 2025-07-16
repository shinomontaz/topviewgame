package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Position struct {
	X, Y int
}

type Renderable interface {
	GetImage() *ebiten.Image
}

type Movable struct {
}

func (p *Position) GetManhattanDistance(other *Position) int {
	xDist := math.Abs(float64(p.X - other.X))
	yDist := math.Abs(float64(p.Y - other.Y))
	return int(xDist) + int(yDist)
}
