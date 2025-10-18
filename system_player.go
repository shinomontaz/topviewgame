package main

import (
	"topviewgame/event"
)

func getDirection(g *Game, ev event.Event) (int, int) {
	x, y := ev.Pos[0], ev.Pos[1]
	if ev.Type == event.EventKey {
		return x, y
	}

	// For mouse clicks, we don't return direction anymore
	// The auto-movement system handles pathfinding
	if ev.Type == event.EventClick {
		return 0, 0
	}

	return x, y
}

func UpdatePlayer(g *Game) {
	players := g.WorldTags["players"]
	level := g.Map.CurrentLevel
	ev := g.PlayerController.GetEvent()

	// Check for user input that should interrupt auto-movement
	if ev.Type != event.EventNone {
		g.StopAutoMove() // Stop auto-movement on any user input
	}

	// Handle manual input first
	if ev.Type == event.EventPass {
		g.Turn = GetNextState(g.Turn)
		g.TurnCounter = 0
		return
	}

	var dx, dy int
	var hasMoved bool

	// If we have manual input, process it
	if ev.Type != event.EventNone {
		dx, dy = getDirection(g, ev)

		// If it's a mouse click, start auto-movement
		if ev.Type == event.EventClick {
			targetX := ev.Pos[0]/level.gd.TileWidth + g.viewport.X1
			targetY := ev.Pos[1]/level.gd.TileHeight + g.viewport.Y1
			targetPos := Position{targetX, targetY}

			if targetX >= 0 && targetX < level.gd.MapWidth && targetY >= 0 && targetY < level.gd.MapHeight {
				targetIdx := level.GetIndexFromXY(targetX, targetY)
				targetTile := level.Tiles[targetIdx]

				if !targetTile.IsRevealed {
					targetPos = level.ClosestVisibleOnLine(Position{targetX, targetY}, g.Center)
				}

				astar := AStar{}
				path := astar.GetPath(level, &g.Center, &targetPos)
				if len(path) > 1 {
					g.StartAutoMove(path)
					// Don't move immediately, let auto-movement handle it next frame
					return
				}
			}
		}
	} else if g.IsAutoMoving {
		// Check if we should stop auto-movement due to enemy in sight
		if g.IsEnemyInSight() {
			g.StopAutoMove()
			return
		}

		// Check if enough time has passed for next step
		if !g.CanAutoMoveStep() {
			return // Wait for timer
		}

		// Get next step from auto-movement
		var hasStep bool
		dx, dy, hasStep = g.GetNextAutoMoveStep()
		if !hasStep {
			return
		}
	} else {
		// No input and no auto-movement
		return
	}

	for _, result := range g.World.Query(players) {
		pos := result.Components[positionC].(*Position)
		player := result.Components[playerC].(*Player)

		newX := pos.X + dx
		newY := pos.Y + dy
		index := level.GetIndexFromXY(newX, newY)

		if (dx != 0 || dy != 0) && index >= 0 && index < len(level.Tiles) {
			if !level.Tiles[index].Blocked {
				level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
				pos.X = newX
				pos.Y = newY
				player.SetMoved(dx)
				hasMoved = true
				level.Tiles[index].Blocked = true
			} else if level.Tiles[index].TileType != WALL {
				// Hit a monster - stop auto-movement and attack
				g.StopAutoMove()
				monsterPos := Position{X: newX, Y: newY}
				player.SetAttacking(dx)
				ProcessAttacks(g, pos, &monsterPos)
				hasMoved = true
			} else {
				// Hit a wall - stop auto-movement
				g.StopAutoMove()
			}
		}

		level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)

		if hasMoved {
			g.Turn = GetNextState(g.Turn)
			g.TurnCounter = 0
		}

		g.SetCenter(*pos)
	}
}
