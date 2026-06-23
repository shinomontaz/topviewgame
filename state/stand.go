package state

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type StandState struct {
	id     int
	frames []*ebiten.Image
	size   int
	delay  float64
	index  int
	timer  float64
	o      owner
}

func Stand(o owner, frames []*ebiten.Image) *StandState {
	return &StandState{
		id:     STAND,
		frames: frames,
		size:   frames[0].Bounds().Dx(),
		o:      o,
		delay:  2,
	}
}

func (s *StandState) IsBusy() bool {
	return false
}

func (s *StandState) GetId() int { return s.id }
func (s *StandState) Start() {
	s.index = 0
	s.timer = 0
}
func (s *StandState) GetFrame() *ebiten.Image { return s.frames[s.index] }

func (s *StandState) GetTransform(tW, tH int) (ebiten.GeoM, float64) {
	cW := float64((s.size - tW) / 2)
	cH := float64(s.size - tH)

	g := ebiten.GeoM{}
	g.Translate(-cW, -cH)

	return g, 0
}

func (s *StandState) Update(dt float64) {
	s.timer += dt
}

func (s *StandState) NextState() (int, bool) {
	if s.timer >= s.delay {
		return IDLE, true
	}

	return s.id, false
}
