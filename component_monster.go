package main

import (
	"log"
	"topviewgame/animations"
	"topviewgame/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Monster struct {
	state    stater
	states   map[int]stater
	frame    int
	time     float64
	lastMove float64
}

type MonsterType int

const (
	SKELETON MonsterType = iota
	ZOMBIE

// GOBLIN
)

func NewMonster(t MonsterType) *Monster {
	pl := &Monster{}
	var (
		sheet *ebiten.Image
		err   error
	)
	switch t {
	case SKELETON:
		sheet, _, err = ebitenutil.NewImageFromFile("assets/actors/Mummy_combined_spritesheet.png")
	case ZOMBIE:
		sheet, _, err = ebitenutil.NewImageFromFile("assets/actors/Mummy2_combined_spritesheet.png")
	default:
		panic("unknown monster type")
	}
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

func (p *Monster) GetImage() *ebiten.Image {
	return p.state.GetFrame()
}

func (p *Monster) SetState(newId int) {
	if p.state != nil && p.state.GetId() == newId {
		return
	}
	p.state = p.states[newId]
	p.state.Start()
}

func (p *Monster) SetMoved() {
	p.lastMove = 0
	p.SetState(state.STAND)
}

func (p *Monster) Update(dt float64) {
	p.state.Update(dt)

	if nextId, ok := p.state.NextState(); ok && nextId != p.state.GetId() {
		p.SetState(nextId)
	}
}
