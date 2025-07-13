package main

import (
	"log"
	"topviewgame/animations"
	"topviewgame/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type stater interface {
	GetId() int
	Start()
	Update(dt float64)
	GetFrame() *ebiten.Image
	NextState() (int, bool)
}

type Player struct {
	state    stater
	states   map[int]stater
	frame    int
	time     float64
	level    *Level
	position *Position
	lastMove float64
}

func NewPlayer() *Player {
	pl := &Player{}
	sheet, _, err := ebitenutil.NewImageFromFile("assets/actors/GraveRobber_combined_spritesheet.png")
	if err != nil {
		log.Fatal(err)
	}

	specs := []animations.Config{
		{StateID: state.STAND, FrameCount: 1},
		{StateID: state.IDLE, FrameCount: 4},
	}

	animMap := animations.BuildMap(sheet, specs, 48, 48)
	pl.states = map[int]stater{
		state.STAND: state.Stand(pl, animMap[state.STAND]),
		state.IDLE:  state.Idle(pl, animMap[state.IDLE]),
	}

	pl.SetState(state.STAND)
	return pl
}

func (p *Player) GetImage() *ebiten.Image {
	return p.state.GetFrame()
}

func (p *Player) SetState(newId int) {
	if p.state != nil && p.state.GetId() == newId {
		return
	}
	p.state = p.states[newId]
	p.state.Start()
}

func (p *Player) SetMoved() {
	p.lastMove = 0
	p.SetState(state.STAND)
}

func (p *Player) Update(dt float64) {
	p.state.Update(dt)

	if nextId, ok := p.state.NextState(); ok && nextId != p.state.GetId() {
		p.SetState(nextId)
	}
}

type Position struct {
	X, Y int
}

type Renderable interface {
	GetImage() *ebiten.Image
}

type Movable struct {
}
