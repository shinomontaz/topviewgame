package main

import (
	"fmt"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Heap[T any] struct {
	list []int
	obj  []T
}

func (h *Heap[T]) Push(key int, o T) {
	h.list = append(h.list, key)
	h.obj = append(h.obj, o)

	h.up()
}

func (h *Heap[T]) Pop() (int, T) {
	if len(h.list) == 0 {
		var zero T
		return -1, zero
	}

	h.list[0], h.list[len(h.list)-1] = h.list[len(h.list)-1], h.list[0]
	h.obj[0], h.obj[len(h.obj)-1] = h.obj[len(h.obj)-1], h.obj[0]

	x := h.list[len(h.list)-1]
	o := h.obj[len(h.obj)-1]

	h.list = h.list[:len(h.list)-1]
	h.obj = h.obj[:len(h.obj)-1]

	h.down()

	return x, o
}

func (h *Heap[T]) up() {
	i := len(h.list) - 1
	for i > 0 {
		a := (i - 1) / 2
		if h.list[a] <= h.list[i] {
			break
		}
		h.list[a], h.list[i] = h.list[i], h.list[a]
		h.obj[a], h.obj[i] = h.obj[i], h.obj[a]

		i = a
	}
}

func (h *Heap[T]) down() {
	i := 0
	for {
		a, b := 2*i+1, 2*i+2
		j := i
		if b < len(h.list) && h.list[b] < h.list[j] {
			j = b
		}
		if a < len(h.list) && h.list[a] < h.list[j] {
			j = a
		}
		if j == i {
			break
		}

		h.list[i], h.list[j] = h.list[j], h.list[i]
		h.obj[i], h.obj[j] = h.obj[j], h.obj[i]

		i = j
	}
}

func loadFont(path string, opts *opentype.FaceOptions) (font.Face, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file: %v", err)
	}
	tt, err := opentype.Parse(fontBytes)
	return opentype.NewFace(tt, opts)
}

func createBorder(tileW, tileH, tileSize int) (*ebiten.Image, error) {
	img := ebiten.NewImage(tileW*tileSize, tileH*tileSize)
	hudBorder, err := loadImage("HUD_border")
	if err != nil {
		return nil, err
	}
	hudCorner, err := loadImage("HUD_corner")
	if err != nil {
		return nil, err
	}
	for i := 0; i <= tileH; i++ {
		for j := 0; j <= tileW; j++ {
			if i == 0 && j == 0 {
				img.DrawImage(hudCorner, &ebiten.DrawImageOptions{})
			} else if (i == 0 || i == tileH) && (j != 0 && j != tileW) {
				op := &ebiten.DrawImageOptions{}
				if i == tileH {
					op.GeoM.Rotate(180 * math.Pi / 180)
				}
				op.GeoM.Translate(float64(j*tileSize), float64(i*tileSize))
				img.DrawImage(hudBorder, op)
			} else if (j == 0 || j == tileW) && (i != 0 && i != tileH) {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Rotate(90 * math.Pi / 180)
				op.GeoM.Translate(float64(j*tileSize), float64(i*tileSize))
				img.DrawImage(hudBorder, op)
			} else if i == tileH && j == tileW {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Rotate(180 * math.Pi / 180)
				op.GeoM.Translate(float64(j*tileSize), float64(i*tileSize))
				img.DrawImage(hudCorner, op)
			} else if i == 0 && j == tileW {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Rotate(90 * math.Pi / 180)
				op.GeoM.Translate(float64(j*tileSize), float64(i*tileSize))
				img.DrawImage(hudCorner, op)
			} else if i == tileH && j == 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(j*tileSize), float64(i*tileSize))
				op.GeoM.Rotate(-90 * math.Pi / 180)
				img.DrawImage(hudCorner, op)
			}
		}
	}

	return img, nil
}
