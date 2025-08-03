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
	lastMove float64
	dir      int
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
		state.WALK:   state.Attack(pl, animMap[state.WALK]),
	}

	pl.SetState(state.STAND)
	return pl
}

func (p *Player) GetImage() *ebiten.Image {
	frame := p.state.GetFrame()
	if p.dir == -1 {
		mirroredFrame := ebiten.NewImage(frame.Bounds().Dx(), frame.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(frame.Bounds().Dx()), 0)
		mirroredFrame.DrawImage(frame, op)
		return mirroredFrame
	}

	return frame
}

func (p *Player) SetState(newId int) {
	if p.state != nil && p.state.GetId() == newId {
		return
	}
	p.state = p.states[newId]
	p.state.Start()
}

func (p *Player) SetMoved(dir int) {
	p.lastMove = 0
	p.dir = dir
	p.SetState(state.STAND)
}

func (p *Player) SetAttacking(dir int) {
	p.lastMove = 0
	p.dir = dir
	p.SetState(state.ATTACK)
}

func (p *Player) Update(dt float64) {
	p.state.Update(dt)

	if nextId, ok := p.state.NextState(); ok && nextId != p.state.GetId() {
		p.SetState(nextId)
	}
}
