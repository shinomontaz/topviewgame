package main

type GameMap struct {
	Dungeons     []Dungeon
	CurrentLevel Level
}

func NewGameMap() GameMap {
	l := NewLevel()
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
