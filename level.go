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
	Rooms  []Rect
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
	MIN_SIZE := 6
	MAX_SIZE := 10
	MAX_ROOMS := 30

	l.Tiles = make([]MapTile, l.gd.ScreenWidth*l.gd.ScreenHeight)
	for x := 0; x < l.gd.ScreenWidth; x++ {
		for y := 0; y < l.gd.ScreenHeight; y++ {
			index := l.gd.GetIndexFromXY(x, y)
			l.Tiles[index] = MapTile{
				PixelX:  x * l.gd.TileWidth,
				PixelY:  y * l.gd.TileHeight,
				Blocked: true,
			}
		}
	}

	// craete rooms
	for idx := 0; idx < MAX_ROOMS; idx++ {
		w := MIN_SIZE + rnd.Intn(MAX_SIZE-MIN_SIZE+1)
		h := MIN_SIZE + rnd.Intn(MAX_SIZE-MIN_SIZE+1)
		x := rnd.Intn(l.gd.ScreenWidth - w - 1)
		y := rnd.Intn(l.gd.ScreenHeight - h - 1)

		newroom := NewRect(x, y, w, h)
		okToAdd := true

		for _, otherRoom := range l.Rooms {
			if newroom.Intersect(otherRoom) {
				okToAdd = false

				break
			}
		}

		if !okToAdd {
			continue
		}

		l.addRoom(newroom)
		l.Rooms = append(l.Rooms, newroom)
		if len(l.Rooms) > 1 {
			newX, newY := newroom.Center()
			prevX, prevY := l.Rooms[len(l.Rooms)-2].Center()

			dice := rnd.Intn(2)
			if dice == 1 {
				l.addHorizontalTunnel(prevX, newX, prevY)
				l.addVerticalTunnel(prevY, newY, newX)
			} else {
				l.addHorizontalTunnel(prevX, newX, newY)
				l.addVerticalTunnel(prevY, newY, prevX)
			}
		}
	}
}

func (l *Level) addRoom(room Rect) {
	for y := room.Y1 + 1; y < room.Y2; y++ {
		for x := room.X1 + 1; x < room.X2; x++ {
			index := l.GetIndexFromXY(x, y)
			l.Tiles[index].Blocked = false
		}
	}
}

func (l *Level) addHorizontalTunnel(x1, x2, y int) {
	for x := min(x1, x2); x < max(x1, x2)+1; x++ {
		index := l.GetIndexFromXY(x, y)
		if index == 0 || index >= l.gd.ScreenWidth*l.gd.ScreenHeight {
			continue
		}
		l.Tiles[index].Blocked = false
	}
}

func (l *Level) addVerticalTunnel(y1, y2, x int) {
	for y := min(y1, y2); y < max(y1, y2)+1; y++ {
		index := l.GetIndexFromXY(x, y)
		if index == 0 || index >= l.gd.ScreenWidth*l.gd.ScreenHeight {
			continue
		}
		l.Tiles[index].Blocked = false
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

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
