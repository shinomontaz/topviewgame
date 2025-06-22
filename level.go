package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameData struct {
	ScreenWidth  int // in tiles
	ScreenHeight int

	TileWidth  int
	TileHeight int
}

func NewGameData() GameData {
	return GameData{
		ScreenWidth:  40,
		ScreenHeight: 25,
		TileWidth:    32,
		TileHeight:   32,
	}
}

type MapTile struct {
	PixelX  int
	PixelY  int
	Blocked bool
	Image   *ebiten.Image
}

func GetIndexFromXY(x, y int) int {
	gd := NewGameData()
	return y*gd.ScreenWidth + x
}

type WallDirection int

const (
	WALL_UP WallDirection = iota
	WALL_DOWN
	WALL_LEFT
	WALL_RIGHT
)

type ImageMap struct {
	Wall  map[int]*ebiten.Image
	Floor *ebiten.Image
}

func LoadImage(name string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("assets/tiles/%s.png", name))

	return img, err
}

func CreateTiles() []MapTile {
	gd := NewGameData()
	tiles := make([]MapTile, gd.ScreenWidth*gd.ScreenHeight)

	var wallImageCache = make(map[string]*ebiten.Image)
	var (
		err error
		img *ebiten.Image
	)

	img, err = LoadImage("floor")
	if err != nil {
		log.Fatalf("Failed to load image %s: %v", "floor", err)
	}
	wallImageCache["floor"] = img

	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < gd.ScreenHeight; y++ {
			index := GetIndexFromXY(x, y)
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
			index := GetIndexFromXY(x, y)
			tile := &tiles[index]
			if !tile.Blocked {
				tile.Image = wallImageCache["floor"]
				continue
			}
			// Compute 8-bit mask
			mask := uint8(0)
			mask = computeMask(x, y, tiles, gd)

			tileName := blobMaskToTile(mask)

			img, ok := wallImageCache[tileName]
			if !ok {
				img, err = LoadImage(tileName)
				if err != nil {
					log.Fatalf("Failed to load image %s: %v", tileName, err)
				}
				wallImageCache[tileName] = img
			}
			tile.Image = img
		}
	}

	return tiles
}
