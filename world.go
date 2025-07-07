package main

import (
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var positionC *ecs.Component
var renderableC *ecs.Component

func InitializeWorld(startingLevel Level) (*ecs.Manager, map[string]ecs.Tag) {
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	playerC := manager.NewComponent()
	positionC = manager.NewComponent()
	renderableC = manager.NewComponent()
	movableC := manager.NewComponent()

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/actors/GraveRobber2.png")
	if err != nil {
		log.Fatal(err)
	}

	startingRoom := startingLevel.Rooms[0]
	playerX, playerY := startingRoom.Center()

	player := Player{
		Image: playerImg,
	}

	manager.NewEntity().
		AddComponent(playerC, player).
		AddComponent(renderableC, &player).
		AddComponent(movableC, Movable{}).
		AddComponent(positionC, &Position{X: playerX, Y: playerY})

	players := ecs.BuildTag(playerC, positionC)
	tags["players"] = players

	renderables := ecs.BuildTag(renderableC, positionC)
	tags["renderables"] = renderables

	return manager, tags
}
