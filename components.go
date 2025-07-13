package main

import (
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
