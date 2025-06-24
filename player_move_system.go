package main

import "github.com/hajimehoshi/ebiten/v2"

func TryMoveplayer(g *Game) {
	players := g.WorldTags["players"]

	x := 0
	y := 0

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		y = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		y = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		x = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		x = 1
	}

	for _, result := range g.World.Query(players) {
		pos := result.Components[position].(*Position)
		pos.X += x
		pos.Y += y
	}
}
