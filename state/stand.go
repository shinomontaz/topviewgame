package state

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type StandState struct {
	id     int
	frames []*ebiten.Image
	delay  float64
	index  int
	timer  float64
	o      owner
}

func Stand(o owner, frames []*ebiten.Image) *StandState {
	return &StandState{
		id:     STAND,
		frames: frames,
		o:      o,
		delay:  2,
	}
}

func (s *StandState) GetId() int { return s.id }
func (s *StandState) Start() {
	s.index = 0
	s.timer = 0
}
func (s *StandState) GetFrame() *ebiten.Image { return s.frames[s.index] }
func (s *StandState) Update(dt float64) {
	s.timer += dt
}

func (s *StandState) NextState() (int, bool) {
	if s.timer >= s.delay {
		return IDLE, true
	}

	return s.id, false
}
