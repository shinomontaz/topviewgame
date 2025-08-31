package main

import (
	"topviewgame/controller"

	"github.com/hajimehoshi/ebiten/v2"
)

func getDirection(g *Game) (int, int) {
	dx, dy := g.PlayerController.GetDirection()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseScreenTileX, mouseScreenTileY := g.PlayerController.GetMouseScreenTile()

		targetX := mouseScreenTileX + g.viewport.X1
		targetY := mouseScreenTileY + g.viewport.Y1
		level := g.Map.CurrentLevel
		targetPos := Position{targetX, targetY}

		if targetX >= 0 && targetX < level.gd.MapWidth && targetY >= 0 && targetY < level.gd.MapHeight {
			targetIdx := level.GetIndexFromXY(targetX, targetY)
			targetTile := level.Tiles[targetIdx]

			if !targetTile.IsRevealed {
				// calculate path from target to closest revealed tile
				// update targetTile to founded

				targetPos = level.ClosestVisibleOnLine(Position{targetX, targetY}, g.Center)
			}
			astar := AStar{}
			path := astar.GetPath(level, &g.Center, &targetPos)
			if len(path) == 0 { // cannot find path
				return dx, dy
			}

			dx = path[0].X - g.Center.X
			dy = path[0].Y - g.Center.Y
		}
	}

	return dx, dy
}

func ProcessPlayer(g *Game) {
	players := g.WorldTags["players"]
	level := g.Map.CurrentLevel
	action := g.PlayerController.GetAction()
	if action == controller.ActionPass {
		g.Turn = GetNextState(g.Turn)
		g.TurnCounter = 0

		return
	}

	dx, dy := getDirection(g)
	hasMoved := false

	for _, result := range g.World.Query(players) {
		pos := result.Components[positionC].(*Position)
		player := result.Components[playerC].(*Player)

		g.SetCenter(*pos)

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
				monsterPos := Position{X: newX, Y: newY}
				player.SetAttacking(dx)
				ProcessAttacks(g, pos, &monsterPos)
				hasMoved = true
			}
		}

		level.PlayerVisible.Compute(level, pos.X, pos.Y, 8)

		if hasMoved {
			g.Turn = GetNextState(g.Turn)
			g.TurnCounter = 0
		}
	}
}
