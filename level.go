package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/norendren/go-fov/fov"
)

type TileType int

const (
	WALL TileType = iota
	FLOOR
)

type MapTile struct {
	PixelX     int
	PixelY     int
	TileType   TileType
	Blocked    bool
	IsRevealed bool
	Image      *ebiten.Image
}

type Level struct {
	Width         int
	Height        int
	gd            GameData
	Tiles         []*MapTile
	Rooms         []Rect
	PlayerVisible *fov.View
	shader        *ebiten.Shader
	OffScreen     *ebiten.Image
}

func NewLevel(gd GameData) Level {
	l := Level{
		gd:            gd,
		PlayerVisible: fov.New(),
	}
	l.build()
	l.adjust()

	bytes, err := os.ReadFile("assets/shaders/single_grayscale.kage")
	if err != nil {
		log.Fatal(err)
	}
	shader, err := ebiten.NewShader(bytes)
	if err != nil {
		log.Fatal(err)
	}
	l.shader = shader

	return l
}

func (l *Level) GetDimensions() (int, int) {
	return l.gd.ScreenWidth, l.gd.ScreenHeight
}

func (l Level) TileAt(x, y int) *MapTile {
	return l.Tiles[l.GetIndexFromXY(x, y)]
}

func (l *Level) GetIndexFromXY(x, y int) int {
	return l.gd.GetIndexFromXY(x, y)
}

func (l *Level) build() {
	MIN_SIZE := 6
	MAX_SIZE := 10
	MAX_ROOMS := 30

	l.Tiles = make([]*MapTile, l.gd.MapWidth*l.gd.MapHeight)
	for x := 0; x < l.gd.MapWidth; x++ {
		for y := 0; y < l.gd.MapHeight; y++ {
			index := l.gd.GetIndexFromXY(x, y)
			l.Tiles[index] = &MapTile{
				PixelX:   x * l.gd.TileWidth,
				PixelY:   y * l.gd.TileHeight,
				Blocked:  true,
				TileType: WALL,
			}
		}
	}

	// craete rooms
	for idx := 0; idx < MAX_ROOMS; idx++ {
		w := MIN_SIZE + rnd.Intn(MAX_SIZE-MIN_SIZE+1)
		h := MIN_SIZE + rnd.Intn(MAX_SIZE-MIN_SIZE+1)
		x := rnd.Intn(l.gd.MapWidth - w - 1)
		y := rnd.Intn(l.gd.MapHeight - h - 1)

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
			l.Tiles[index].TileType = FLOOR
		}
	}
}

func (l *Level) addHorizontalTunnel(x1, x2, y int) {
	for x := min(x1, x2); x < max(x1, x2)+1; x++ {
		index := l.GetIndexFromXY(x, y)
		if index == 0 || index >= l.gd.MapWidth*l.gd.MapHeight {
			continue
		}
		l.Tiles[index].Blocked = false
		l.Tiles[index].TileType = FLOOR
	}
}

func (l *Level) addVerticalTunnel(y1, y2, x int) {
	for y := min(y1, y2); y < max(y1, y2)+1; y++ {
		index := l.GetIndexFromXY(x, y)
		if index == 0 || index >= l.gd.MapWidth*l.gd.MapHeight {
			continue
		}
		l.Tiles[index].Blocked = false
		l.Tiles[index].TileType = FLOOR
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

	for x := 0; x < l.gd.MapWidth; x++ {
		for y := 0; y < l.gd.MapHeight; y++ {
			index := l.gd.GetIndexFromXY(x, y)
			tile := l.Tiles[index]
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

func (l Level) InBounds(x, y int) bool {
	if x < 0 || x > l.gd.MapWidth || y < 0 || y > l.gd.MapHeight {
		return false
	}

	return true
}

func (l Level) IsOpaque(x, y int) bool {
	return l.Tiles[l.GetIndexFromXY(x, y)].TileType == WALL
}

func (l *Level) Draw(screen *ebiten.Image, viewport Rect) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	vw := l.gd.ScreenWidth
	vh := l.gd.ScreenHeight
	tw := l.gd.TileWidth
	th := l.gd.TileHeight

	if l.OffScreen == nil || l.OffScreen.Bounds().Dx() != vw*tw || l.OffScreen.Bounds().Dy() != vh*th {
		l.OffScreen = ebiten.NewImage(vw*tw, vh*th)
	}
	l.OffScreen.Clear()

	visible := make([]float32, vw*vh)

	baseX := viewport.X1
	baseY := viewport.Y1

	x1 := max(0, viewport.X1)
	y1 := max(0, viewport.Y1)
	x2 := min(l.gd.MapWidth, viewport.X2)
	y2 := min(l.gd.MapHeight, viewport.Y2)

	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			dx := x - baseX
			dy := y - baseY

			if dx < 0 || dx >= vw || dy < 0 || dy >= vh {
				continue
			}

			idx := l.gd.GetIndexFromXY(x, y)
			tile := l.Tiles[idx]

			if l.PlayerVisible.IsVisible(x, y) {
				tile.IsRevealed = true
				visible[dy*vw+dx] = 1.0
			} else if tile.IsRevealed {
				visible[dy*vw+dx] = 0.5
			} else {
				visible[dy*vw+dx] = 0.0
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(dx*tw), float64(dy*th))
			l.OffScreen.DrawImage(tile.Image, op)
		}
	}

	shaderOpts := &ebiten.DrawRectShaderOptions{}
	shaderOpts.Images[0] = l.OffScreen
	shaderOpts.Uniforms = map[string]interface{}{
		"Visible":     visible,
		"ScreenWidth": vw,
		"TileWidth":   tw,
		"TileHeight":  th,
	}
	screen.DrawRectShader(l.OffScreen.Bounds().Dx(), l.OffScreen.Bounds().Dy(), l.shader, shaderOpts)
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
