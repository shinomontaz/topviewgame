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
	size     int

	o owner
}

func Death(o owner, frames []*ebiten.Image) *DeathState {
	return &DeathState{
		id:     DEATH,
		frames: frames,
		size:   frames[0].Bounds().Dx(),
		o:      o,
	}
}

func (s *DeathState) IsBusy() bool {
	return !s.isPlayed
}

func (s *DeathState) GetId() int { return s.id }
func (s *DeathState) Start() {
}

func (s *DeathState) GetFrame() *ebiten.Image {
	return s.frames[s.index]
}

func (s *DeathState) GetTransform(tW, tH int) (ebiten.GeoM, float64) {
	cW := float64((s.size - tW) / 2)
	cH := float64(s.size - tH)

	g := ebiten.GeoM{}
	g.Translate(-cW, -cH)

	return g, 0
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
