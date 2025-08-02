package state

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type DeathState struct {
	id       int
	index    int
	timer    float64
	isPlayed bool
	frames   []*ebiten.Image
	o        owner
}

func Death(o owner, frames []*ebiten.Image) *DeathState {
	return &DeathState{
		id:     DEATH,
		frames: frames,
		o:      o,
	}
}

func (s *DeathState) GetId() int { return s.id }
func (s *DeathState) Start() {
}

func (s *DeathState) GetFrame() *ebiten.Image {
	return s.frames[s.index]
}

func (s *DeathState) Update(dt float64) {
	s.timer += dt
	if s.isPlayed {
		return
	}

	s.index = int(math.Floor(s.timer / 0.2))
	if s.index >= len(s.frames)-1 {
		s.index = len(s.frames) - 1
		s.isPlayed = true
	}
}

func (s *DeathState) NextState() (int, bool) {
	return s.id, false
}
