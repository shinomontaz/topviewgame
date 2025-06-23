package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Level struct {
	Width  int
	Height int
	Tiles  []MapTile
}

func NewLevel() Level {
	l := Level{}
	tiles := l.createTiles()
	l.Tiles = tiles

	return l
}

type MapTile struct {
	PixelX  int
	PixelY  int
	Blocked bool
	Image   *ebiten.Image
}

func loadImage(name string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("assets/tiles/%s.png", name))

	return img, err
}

func getIndexFromXY(x, y int) int {
	gd := NewGameData()
	return y*gd.ScreenWidth + x
}

func (l *Level) createTiles() []MapTile {
	gd := NewGameData()
	tiles := make([]MapTile, gd.ScreenWidth*gd.ScreenHeight)

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

	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < gd.ScreenHeight; y++ {
			index := getIndexFromXY(x, y)
			if x == 0 || x == gd.ScreenWidth-1 || y == 0 || y == gd.ScreenHeight-1 {
				tiles[index] = MapTile{
					PixelX:  x * gd.TileWidth,
					PixelY:  y * gd.TileHeight,
					Blocked: true,
				}
			} else {
				tiles[index] = MapTile{
					PixelX:  x * gd.TileWidth,
					PixelY:  y * gd.TileHeight,
					Blocked: false,
				}
			}
		}
	}

	// Second pass: assign wall images based on neighbors using blob mask
	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < gd.ScreenHeight; y++ {
			index := getIndexFromXY(x, y)
			tile := &tiles[index]
			if !tile.Blocked {
				tile.Image = imageCache["floor"]
				continue
			}
			// Compute 8-bit mask
			mask := uint8(0)
			mask = computeMask(x, y, tiles, gd)

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

	return tiles
}
