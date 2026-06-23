package state

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type WalkState struct {
	id        int
	index     int
	timer     float64
	duration  float64
	frames    []*ebiten.Image
	size      int
	tiltAngle float64

	o owner
}

func Walk(o owner, frames []*ebiten.Image) *WalkState {
	return &WalkState{
		id:        WALK,
		frames:    frames,
		size:      frames[0].Bounds().Dx(),
		duration:  float64(len(frames)) * 0.1,
		tiltAngle: 0.2,
		o:         o,
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

func (s *WalkState) GetTransform(tW, tH int) (ebiten.GeoM, float64) {

	progress := s.timer / s.duration
	dx, dy := s.o.GetDirection()                // -1, 0, +1
	offsetX := (1 - progress) * float64(-dx*tW) // от старой клетки к новой
	offsetY := (1 - progress) * float64(-dy*tH) // от старой клетки к новой

	height := 4 * progress * (1 - progress) * (float64(s.size) / 3.0) // парабола

	cW := float64((s.size - tW) / 2)
	cH := float64(s.size - tH)

	// здесь считаем параболу!
	// нам нужен вектор куда мы шли и положение куда мы пришли
	g := ebiten.GeoM{}
	rotation := s.tiltAngle * math.Sin(math.Pi*progress)
	g.Rotate(rotation)
	g.Translate(offsetX-cW, offsetY-cH)

	return g, height
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
