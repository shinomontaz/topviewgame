package main

func TryMovePlayer(g *Game) {
	players := g.WorldTags["players"]
	level := g.Map.CurrentLevel

	dx, dy := g.PlayerController.GetDirection()
	hasMoved := false

	for _, result := range g.World.Query(players) {
		pos := result.Components[positionC].(*Position)
		player := result.Components[playerC].(*Player)

		newX := pos.X + dx
		newY := pos.Y + dy
		index := level.GetIndexFromXY(newX, newY)

		if dx != 0 || dy != 0 && index >= 0 && index < len(level.Tiles) && !level.Tiles[index].Blocked {
			pos.X = newX
			pos.Y = newY
			player.SetMoved()
			hasMoved = true
		}

		level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)

		if hasMoved {
			g.Turn = GetNextState(g.Turn)
			g.TurnCounter = 0
		}
	}
}
