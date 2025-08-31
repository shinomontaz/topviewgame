package state

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type AttackState struct {
	id       int
	index    int
	timer    float64
	frames   []*ebiten.Image
	finished bool
	o        owner
}

func Attack(o owner, frames []*ebiten.Image) *AttackState {
	return &AttackState{
		id:       ATTACK,
		frames:   frames,
		finished: false,
		o:        o,
	}
}

func (s *AttackState) GetId() int { return s.id }
func (s *AttackState) Start() {
	s.index = 0
	s.finished = false
	s.timer = 0
}

func (s *AttackState) GetFrame() *ebiten.Image {
	return s.frames[s.index]
}

func (s *AttackState) Update(dt float64) {
	s.timer += dt

	s.index = int(math.Floor(s.timer / 0.1))
	if s.index >= len(s.frames)-1 {
		s.index = 0
		s.finished = true
	}
}

func (s *AttackState) NextState() (int, bool) {
	if s.finished {
		return STAND, true
	}
	return s.id, false
}
