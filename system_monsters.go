package main

import (
	"log"

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
		mon := result.Components[monsterC].(*Monster)

		monsterSees := fov.New()
		monsterSees.Compute(l, pos.X, pos.Y, 7)
		if monsterSees.IsVisible(playerPosition.X, playerPosition.Y) {
			log.Printf("%s sees you", mon.GetName())
		}
	}

	g.Turn = PlayerTurn
}
