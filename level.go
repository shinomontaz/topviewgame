package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type MapTile struct {
	PixelX  int
	PixelY  int
	Blocked bool
	Image   *ebiten.Image
}

type Level struct {
	Width  int
	Height int
	gd     GameData
	Tiles  []MapTile
}

func NewLevel() Level {
	l := Level{
		gd: NewGameData(),
	}
	l.build()
	l.adjust()

	return l
}

func (l *Level) GetIndexFromXY(x, y int) int {
	return l.gd.GetIndexFromXY(x, y)
}

func (l *Level) build() {
	l.Tiles = make([]MapTile, l.gd.ScreenWidth*l.gd.ScreenHeight)
	for x := 0; x < l.gd.ScreenWidth; x++ {
		for y := 0; y < l.gd.ScreenHeight; y++ {
			index := l.gd.GetIndexFromXY(x, y)
			if x == 0 || x == l.gd.ScreenWidth-1 || y == 0 || y == l.gd.ScreenHeight-1 {
				l.Tiles[index] = MapTile{
					PixelX:  x * l.gd.TileWidth,
					PixelY:  y * l.gd.TileHeight,
					Blocked: true,
				}
			} else {
				l.Tiles[index] = MapTile{
					PixelX:  x * l.gd.TileWidth,
					PixelY:  y * l.gd.TileHeight,
					Blocked: false,
				}
			}
		}
	}
}

func (l *Level) adjust() {
	var imageCache = make(map[string]*ebiten.Image)
	var (
		err error
		img *ebiten.Image
	)

	img, err = loadImage("floor")
	if err != nil {
		log.Fatalf("Failed to load image %s: %v", "floor", err)
	}
	imageCache["floor"] = img

	for x := 0; x < l.gd.ScreenWidth; x++ {
		for y := 0; y < l.gd.ScreenHeight; y++ {
			index := l.gd.GetIndexFromXY(x, y)
			tile := &l.Tiles[index]
			if !tile.Blocked {
				tile.Image = imageCache["floor"]
				continue
			}
			// Compute 8-bit mask
			mask := uint8(0)
			mask = computeMask(x, y, l.Tiles, l.gd)

			tileName := blobMaskToTile(mask)

			img, ok := imageCache[tileName]
			if !ok {
				img, err = loadImage(tileName)
				if err != nil {
					log.Fatalf("Failed to load image %s: %v", tileName, err)
				}
				imageCache[tileName] = img
			}
			tile.Image = img
		}
	}
}

func (l *Level) Draw(screen *ebiten.Image) {
	for x := 0; x < l.gd.ScreenWidth; x++ {
		for y := 0; y < l.gd.ScreenHeight; y++ {
			tile := l.Tiles[l.gd.GetIndexFromXY(x, y)]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			screen.DrawImage(tile.Image, op)
		}
	}
}
