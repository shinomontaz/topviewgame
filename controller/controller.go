package controller

import (
	"log"
	"topviewgame/event"

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

func (h Human) GetEvent() event.Event {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		return event.Event{Type: event.EventKey, Pos: [2]int{0, -1}}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		return event.Event{Type: event.EventKey, Pos: [2]int{0, 1}}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		return event.Event{Type: event.EventKey, Pos: [2]int{-1, 0}}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		return event.Event{Type: event.EventKey, Pos: [2]int{1, 0}}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return event.Event{Type: event.EventPass}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return event.Event{Type: event.EventClick, Pos: [2]int{x, y}}
	}

	return event.Event{Type: event.EventNone}
}

func (h Human) GetCursor() (int, int) {
	return ebiten.CursorPosition()
}
