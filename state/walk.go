package state

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type WalkState struct {
	id       int
	index    int
	timer    float64
	duration float64
	frames   []*ebiten.Image
	o        owner
}

func Walk(o owner, frames []*ebiten.Image) *WalkState {
	return &WalkState{
		id:       WALK,
		frames:   frames,
		duration: float64(len(frames)) * 0.1,
		o:        o,
	}
}

func (s *WalkState) IsBusy() bool {
	return s.timer < s.duration
}

func (s *WalkState) GetId() int { return s.id }
func (s *WalkState) Start() {
	s.timer = 0
	s.index = 0
}

func (s *WalkState) GetFrame() *ebiten.Image {
	return s.frames[s.index]
}

func (s *WalkState) Update(dt float64) {
	s.timer += dt
	s.index = int(math.Floor(s.timer / 0.1))
	if s.index >= len(s.frames)-1 {
		s.index = 0
	}
}

func (s *WalkState) NextState() (int, bool) {
	if s.timer >= s.duration {
		return STAND, true
	}

	return s.id, false
}
