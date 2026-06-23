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
	GetTransform(int, int) (ebiten.GeoM, float64)
	NextState() (int, bool)
	IsBusy() bool
}

type Player struct {
	state    stater
	states   map[int]stater
	frame    int
	time     float64
	lastMove float64
	dx       int
	dy       int
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
		{StateID: state.ATTACK, FrameCount: 6},
		{StateID: state.DEATH, FrameCount: 4},
		{StateID: state.WALK, FrameCount: 6},
	}

	animMap := animations.BuildMap(sheet, specs, 48, 48)
	pl.states = map[int]stater{
		state.STAND:  state.Stand(pl, animMap[state.STAND]),
		state.IDLE:   state.Idle(pl, animMap[state.IDLE]),
		state.DEATH:  state.Death(pl, animMap[state.DEATH]),
		state.ATTACK: state.Attack(pl, animMap[state.ATTACK]),
		state.WALK:   state.Walk(pl, animMap[state.WALK]),
	}

	pl.SetState(state.STAND)
	return pl
}

func (p *Player) GetImage(tileW, tileH int) Image {
	frame := p.state.GetFrame()
	geom, h := p.state.GetTransform(tileW, tileH)
	if p.dx == -1 {
		geom.Scale(-1, 1)
		geom.Translate(float64(frame.Bounds().Dx()), 0)
	}

	return Image{Image: frame, GeoM: geom, Height: h}
}

func (p *Player) SetState(newId int) {
	if p.state != nil && p.state.GetId() == newId {
		return
	}
	p.state = p.states[newId]
	p.state.Start()
}

func (p *Player) GetDirection() (int, int) {
	return p.dx, p.dy
}

func (p *Player) SetMoved(dx, dy int) {
	p.lastMove = 0
	p.dx = dx
	p.dy = dy
	p.SetState(state.WALK)
}

func (p *Player) SetAttacking(dx, dy int) {
	p.lastMove = 0
	p.dx = dx
	p.dy = dy
	p.SetState(state.ATTACK)
}

func (p *Player) Update(dt float64) {
	p.state.Update(dt)

	if nextId, ok := p.state.NextState(); ok && nextId != p.state.GetId() {
		p.SetState(nextId)
	}
}
