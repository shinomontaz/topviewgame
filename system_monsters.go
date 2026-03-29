package main

import (
	"topviewgame/state"

	"github.com/norendren/go-fov/fov"
)

func UpdateMonsters(g *Game) {
	l := g.Map.CurrentLevel
	playerPosition := Position{}
	var player *Player

	for _, plr := range g.World.QueryPlayers() {
		pos := g.World.GetPosition(plr)
		player = g.World.GetPlayer(plr).(*Player)
		playerPosition.X = pos.X
		playerPosition.Y = pos.Y
	}

	for _, result := range g.World.QueryMonsters() {
		pos := g.World.GetPosition(result)
		mon := g.World.GetMonster(result).(*Monster)
		if mon.IsDead() {
			continue
		}

		monsterSees := fov.New()
		monsterSees.Compute(l, pos.X, pos.Y, 7)
		if monsterSees.IsVisible(playerPosition.X, playerPosition.Y) {
			if pos.GetManhattanDistance(&playerPosition) == 1 {
				player.SetState(state.STAND)
				ProcessAttacks(g, pos, &playerPosition)
				dx := playerPosition.X - pos.X
				mon.SetAttacking(dx)
				if g.World.GetHealth(result).Current <= 0 {
					t := l.Tiles[l.GetIndexFromXY(pos.X, pos.Y)]
					t.Blocked = false
				}

				continue
			}
			var path []Position
			if mon.LastPlayerPos == playerPosition {
				path = mon.CachedPath
			} else {
				astar := AStar{}
				path = astar.GetPath(l, pos, &playerPosition)
			}
			mon.LastPlayerPos = playerPosition
			mon.CachedPath = path
			if len(path) > 1 {
				nextTile := l.Tiles[l.GetIndexFromXY(path[1].X, path[1].Y)]
				if !nextTile.Blocked {
					oldPos := *pos
					l.Tiles[l.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
					pos.X = path[1].X
					pos.Y = path[1].Y
					l.Tiles[l.GetIndexFromXY(pos.X, pos.Y)].Blocked = true

					g.gm.updateMonsterPosition(result.Entity, &oldPos, pos)
				}
			}
		}
	}

	g.Turn = PlayerTurn
}
