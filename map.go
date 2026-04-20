package main

import "github.com/bytearena/ecs"

type GameMap struct {
	LevelIndex   int
	CurrentLevel Level
	Gd           GameData

	monsterPositions map[Position]*ecs.Entity
}

func NewGameMap(gd GameData) GameMap {
	l := NewLevel(gd)
	gm := GameMap{
		LevelIndex:       1,
		CurrentLevel:     l,
		Gd:               gd,
		monsterPositions: make(map[Position]*ecs.Entity),
	}

	return gm
}

func (gm *GameMap) MonsterAt(p Position) bool {
	return gm.monsterPositions[p] != nil
}

func (gm *GameMap) updateMonsterPosition(entity *ecs.Entity, from, to *Position) {
	if gm.monsterPositions == nil {
		gm.monsterPositions = make(map[Position]*ecs.Entity)
	}
	if from != nil {
		delete(gm.monsterPositions, *from)
	}
	if to != nil {
		gm.monsterPositions[*to] = entity
	}
}
