package main

import (
	"github.com/norendren/go-fov/fov"
)

func UpdateMonsters(g *Game) {
	l := g.Map.CurrentLevel
	playerPosition := Position{}

	for _, plr := range g.World.Query(g.WorldTags["players"]) {
		pos := plr.Components[positionC].(*Position)
		playerPosition.X = pos.X
		playerPosition.Y = pos.Y
	}

	for _, result := range g.World.Query(g.WorldTags["monsters"]) {
		pos := result.Components[positionC].(*Position)
		//		mon := result.Components[monsterC].(*Monster)

		monsterSees := fov.New()
		monsterSees.Compute(l, pos.X, pos.Y, 7)
		if monsterSees.IsVisible(playerPosition.X, playerPosition.Y) {
			astar := AStar{}
			path := astar.GetPath(l, pos, &playerPosition)
			if len(path) > 1 {
				nextTile := l.Tiles[l.GetIndexFromXY(path[1].X, path[1].Y)]
				if !nextTile.Blocked {
					l.Tiles[l.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
					pos.X = path[1].X
					pos.Y = path[1].Y
					l.Tiles[l.GetIndexFromXY(pos.X, pos.Y)].Blocked = true
				}
			}
		}
	}

	g.Turn = PlayerTurn
}
