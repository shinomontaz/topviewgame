package animations

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	StateID    int
	FrameCount int
}

func BuildMap(sheet *ebiten.Image, specs []Config, w, h int) map[int][]*ebiten.Image {
	result := make(map[int][]*ebiten.Image)
	for _, s := range specs {
		result[s.StateID] = extractFramesFromRow(sheet, s.StateID, s.FrameCount, w, h)
	}

	return result
}

func extractFramesFromRow(sheet *ebiten.Image, row, frameCount, frameWidth, frameHeight int) []*ebiten.Image {
	frames := make([]*ebiten.Image, frameCount)
	offsetX := frameWidth
	for i := range frameCount {
		x := i*frameWidth + offsetX
		y := row * frameHeight
		sub := sheet.SubImage(image.Rect(x, y, x+frameWidth, y+frameHeight)).(*ebiten.Image)
		frames[i] = sub
	}

	return frames
}
