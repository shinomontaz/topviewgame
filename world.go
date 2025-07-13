package main

import (
	"github.com/bytearena/ecs"
)

var positionC *ecs.Component
var renderableC *ecs.Component
var playerC *ecs.Component

func InitializeWorld(startingLevel Level) (*ecs.Manager, map[string]ecs.Tag) {
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	playerC = manager.NewComponent()
	positionC = manager.NewComponent()
	renderableC = manager.NewComponent()
	movableC := manager.NewComponent()
	monsterC := manager.NewComponent()

	startingRoom := startingLevel.Rooms[0]
	playerX, playerY := startingRoom.Center()

	player := NewPlayer()

	manager.NewEntity().
		AddComponent(playerC, player).
		AddComponent(renderableC, player).
		AddComponent(movableC, Movable{}).
		AddComponent(positionC, &Position{X: playerX, Y: playerY})

	for _, room := range startingLevel.Rooms {
		if room.X1 != startingRoom.X1 {
			var monsterType MonsterType
			switch rnd.Intn(2) {
			case 0:
				monsterType = SKELETON
			case 1:
				monsterType = ZOMBIE
			}
			monster := NewMonster(monsterType)
			mX, mY := room.Center()
			manager.NewEntity().
				AddComponent(monsterC, monster).
				AddComponent(renderableC, monster).
				AddComponent(positionC, &Position{
					X: mX,
					Y: mY,
				})
		}
	}

	players := ecs.BuildTag(playerC, positionC)
	tags["players"] = players

	renderables := ecs.BuildTag(renderableC, positionC)
	tags["renderables"] = renderables

	return manager, tags
}
