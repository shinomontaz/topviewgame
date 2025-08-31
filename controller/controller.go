package controller

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Human struct {
	TileWidth, TileHeight int
	cursorImg             *ebiten.Image
	err                   error
}

func (h Human) Draw(screen *ebiten.Image) {
	if h.cursorImg == nil {
		h.cursorImg, _, h.err = ebitenutil.NewImageFromFile("assets/cursor.png")
		if h.err != nil {
			log.Fatal(h.err)
		}
	}

	x, y := ebiten.CursorPosition()

	op := &ebiten.DrawImageOptions{}
	scaleX := float64(h.TileWidth) / float64(h.cursorImg.Bounds().Dx())
	scaleY := float64(h.TileHeight) / float64(h.cursorImg.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)

	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(h.cursorImg, op)
}

func (h Human) GetMouseScreenTile() (int, int) {
	x, y := ebiten.CursorPosition()
	tileX := x / h.TileWidth
	tileY := y / h.TileHeight
	return tileX, tileY
}

func (h Human) GetDirection() (int, int) {
	dx, dy := 0, 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx = 1
	}

	return dx, dy
}

const (
	ActionNone Action = iota
	ActionPass
	ActionAttack
	ActionPickup
)

type Action int

func (h Human) GetAction() Action {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return ActionPass
	}
	return ActionNone
}
