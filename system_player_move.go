package main

import "topviewgame/controller"

func TryMovePlayer(g *Game) {
	players := g.WorldTags["players"]
	level := g.Map.CurrentLevel

	dx, dy := g.PlayerController.GetDirection()
	action := g.PlayerController.GetAction()
	hasMoved := false

	for _, result := range g.World.Query(players) {
		pos := result.Components[positionC].(*Position)
		player := result.Components[playerC].(*Player)

		newX := pos.X + dx
		newY := pos.Y + dy
		index := level.GetIndexFromXY(newX, newY)

		if (dx != 0 || dy != 0) && index >= 0 && index < len(level.Tiles) && !level.Tiles[index].Blocked {
			level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
			pos.X = newX
			pos.Y = newY
			player.SetMoved(dx)
			hasMoved = true
			level.Tiles[index].Blocked = true
		}

		level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)

		if hasMoved || action == controller.ActionPass {
			g.Turn = GetNextState(g.Turn)
			g.TurnCounter = 0
		}
	}
}
