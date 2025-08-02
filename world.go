package main

import (
	"github.com/bytearena/ecs"
)

var (
	positionC    *ecs.Component
	renderableC  *ecs.Component
	playerC      *ecs.Component
	monsterC     *ecs.Component
	healthC      *ecs.Component
	meleeWeaponC *ecs.Component
	armorC       *ecs.Component
	nameC        *ecs.Component
)

func InitializeWorld(startingLevel Level) (*ecs.Manager, map[string]ecs.Tag) {
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	playerC = manager.NewComponent()
	positionC = manager.NewComponent()
	renderableC = manager.NewComponent()
	movableC := manager.NewComponent()
	monsterC = manager.NewComponent()
	healthC = manager.NewComponent()
	meleeWeaponC = manager.NewComponent()
	armorC = manager.NewComponent()
	nameC = manager.NewComponent()

	startingRoom := startingLevel.Rooms[0]
	playerX, playerY := startingRoom.Center()

	player := NewPlayer()

	manager.NewEntity().
		AddComponent(playerC, player).
		AddComponent(renderableC, player).
		AddComponent(movableC, Movable{}).
		AddComponent(positionC, &Position{X: playerX, Y: playerY}).
		AddComponent(healthC, &Health{Max: 30, Current: 30}).
		AddComponent(meleeWeaponC, &MeleeWeapon{
			Name:       "Fist",
			MinDamage:  1,
			MaxDamage:  3,
			ToHitBonus: 2,
		}).
		AddComponent(armorC, &Armor{
			Name:    "Burlap Sack",
			Defence: 1,
			Dodge:   1,
		}).
		AddComponent(nameC, &Name{Label: "Player"})

	for _, room := range startingLevel.Rooms {
		if room.X1 != startingRoom.X1 {
			var (
				monsterType MonsterType
				monsterName string
			)
			switch rnd.Intn(2) {
			case 0:
				monsterType = SKELETON
				monsterName = "Skeleton"
			case 1:
				monsterType = ZOMBIE
				monsterName = "Zombie"
			}
			monster := NewMonster(monsterType)
			mX, mY := room.Center()
			ent := manager.NewEntity().
				AddComponent(monsterC, monster).
				AddComponent(renderableC, monster).
				AddComponent(positionC, &Position{
					X: mX,
					Y: mY,
				}).
				AddComponent(nameC, &Name{Label: monsterName})

			if monsterType == SKELETON {
				ent.AddComponent(healthC, &Health{
					Max:     10,
					Current: 2,
				}).AddComponent(meleeWeaponC, &MeleeWeapon{
					Name:       "Short Sword",
					MinDamage:  2,
					MaxDamage:  6,
					ToHitBonus: 0,
				}).AddComponent(armorC, &Armor{
					Name:    "No armor",
					Defence: 0,
					Dodge:   5,
				})
			} else {
				ent.AddComponent(healthC, &Health{
					Max:     20,
					Current: 2,
				}).AddComponent(meleeWeaponC, &MeleeWeapon{
					Name:       "Khopesh",
					MinDamage:  1,
					MaxDamage:  4,
					ToHitBonus: 1,
				}).AddComponent(armorC, &Armor{
					Name:    "Rotten rags",
					Defence: 1,
					Dodge:   0,
				})
			}
		}
	}

	monsters := ecs.BuildTag(monsterC, positionC, healthC, meleeWeaponC, armorC, nameC)
	tags["monsters"] = monsters

	players := ecs.BuildTag(playerC, positionC, healthC, meleeWeaponC, armorC, nameC)
	tags["players"] = players

	renderables := ecs.BuildTag(renderableC, positionC)
	tags["renderables"] = renderables

	return manager, tags
}
