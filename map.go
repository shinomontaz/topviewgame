package main

import "github.com/bytearena/ecs"

type GameMap struct {
	Dungeons     []Dungeon
	CurrentLevel Level
	Gd           GameData

	monsterPositions map[Position]*ecs.Entity
}

func NewGameMap(gd GameData) GameMap {
	l := NewLevel(gd)
	d := Dungeon{
		Name:   "Dungeon 1",
		Levels: []Level{l},
	}
	gm := GameMap{
		Dungeons:         []Dungeon{d},
		CurrentLevel:     l,
		Gd:               gd,
		monsterPositions: make(map[Position]*ecs.Entity),
	}

	return gm
}

func (gm *GameMap) MonsterAt(p Position) bool {
	return gm.monsterPositions[p] != nil
}
