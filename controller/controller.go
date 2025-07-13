package controller

import "github.com/hajimehoshi/ebiten/v2"

type Human struct{}

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
