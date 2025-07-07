package main

import "github.com/hajimehoshi/ebiten/v2"

type Stater interface {
	GetId() int
	Start()
	Update(dt float64)
}

type Player struct {
	Image  *ebiten.Image
	state  Stater
	states map[int]Stater
}

type Position struct {
	X, Y int
}

type Renderable interface {
	GetImage(dt float64) *ebiten.Image
}

type Movable struct {
}

func (p *Player) GetImage(dt float64) *ebiten.Image {
	return p.Image
}

func NewPlayer() *Player {
	pl := &Player{}

	sStand := state.New(state.STAND, p, p.anim)
	sIdle := state.New(state.IDLE, p, p.anim)

	p.states = map[int]Stater{
		state.STAND: sStand,
		state.IDLE:  sIdle,
	}

	p.SetState(state.STAND)
}

type State struct {
}
