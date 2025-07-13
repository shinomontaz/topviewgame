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

	startingRoom := startingLevel.Rooms[0]
	playerX, playerY := startingRoom.Center()

	player := NewPlayer()
	player.level = &startingLevel
	player.position = &Position{X: playerX, Y: playerY}

	manager.NewEntity().
		AddComponent(playerC, player).
		AddComponent(renderableC, player).
		AddComponent(movableC, Movable{}).
		AddComponent(positionC, &Position{X: playerX, Y: playerY})

	players := ecs.BuildTag(playerC, positionC)
	tags["players"] = players

	renderables := ecs.BuildTag(renderableC, positionC)
	tags["renderables"] = renderables

	return manager, tags
}
