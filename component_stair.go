package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Stairs struct {
	Img *ebiten.Image // cursor image

	NextLevel int
}

func NewStairs(w, h int) (*Stairs, error) {
	c := Stairs{}

	imgDefault, _, err := ebitenutil.NewImageFromFile("assets/tiles/stairs.png")
	if err != nil {
		return nil, err
	}
	c.Img = imgDefault

	return &c, nil
}

func (c *Stairs) GetImage(tW, tH int) Image {
	return Image{Image: c.Img, GeoM: ebiten.GeoM{}, Height: 0}
}

func (c *Stairs) Handle(g interface{ NextLevel(int) }) {
	g.NextLevel(c.NextLevel)
}
