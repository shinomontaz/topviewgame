package state

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type IdleState struct {
	id       int
	frames   []*ebiten.Image
	index    int
	timer    float64
	o        owner
	isPlayed bool
}

func Idle(o owner, frames []*ebiten.Image) *IdleState {
	return &IdleState{
		id:     IDLE,
		frames: frames,
		o:      o,
	}
}

func (s *IdleState) GetId() int { return s.id }
func (s *IdleState) Start() {
	s.index = 0
	s.timer = 0
	s.isPlayed = false
}
func (s *IdleState) GetFrame() *ebiten.Image {
	return s.frames[s.index]
}

func (s *IdleState) Update(dt float64) {
	s.timer += dt

	s.index = int(math.Floor(s.timer / 0.2))
	if s.index >= len(s.frames) {
		s.index = 0
		s.isPlayed = true
	}
}

func (s *IdleState) NextState() (int, bool) {
	if s.isPlayed {
		return STAND, true
	}

	return s.id, false
}
