package main

type GameMap struct {
	Dungeons     []Dungeon
	CurrentLevel Level
}

func NewGameMap(gd GameData) GameMap {
	l := NewLevel(gd)
	d := Dungeon{
		Name:   "Dungeon 1",
		Levels: []Level{l},
	}
	gm := GameMap{
		Dungeons:     []Dungeon{d},
		CurrentLevel: l,
	}

	return gm
}
