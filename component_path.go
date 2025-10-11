package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PathVisualizer struct {
	pathImage *ebiten.Image
	width     int
	height    int
}

func NewPathVisualizer() (*PathVisualizer, error) {
	img, _, err := ebitenutil.NewImageFromFile("assets/cursor_path.png")
	if err != nil {
		return nil, err
	}

	return &PathVisualizer{
		pathImage: img,
		width:     img.Bounds().Dx(),
		height:    img.Bounds().Dy(),
	}, nil
}

func (pv *PathVisualizer) GetImage() *ebiten.Image {
	return pv.pathImage
}

func (pv *PathVisualizer) GetSize() (int, int) {
	return pv.width, pv.height
}
