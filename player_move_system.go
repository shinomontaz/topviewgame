package main

import "github.com/hajimehoshi/ebiten/v2"

func TryMovePlayer(g *Game) {
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

	level := g.Map.CurrentLevel
	for _, result := range g.World.Query(players) {
		pos := result.Components[position].(*Position)
		index := level.GetIndexFromXY(pos.X+x, pos.Y+y)

		tile := level.Tiles[index]
		if !tile.Blocked {
			pos.X += x
			pos.Y += y
			level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)
		}
	}

	if x != 0 || y != 0 {
		g.Turn = GetNextState(g.Turn)
		g.TurnCounter = 0
	}
}
