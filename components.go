package main

import "github.com/hajimehoshi/ebiten/v2"

type Player struct{}

type Position struct {
	X, Y int
}

type Renderable struct {
	Image *ebiten.Image
}

type Movable struct {
}
