package main

type GameMap struct {
	Dungeons []Dungeon
}

func NewGameMap() GameMap {
	l := NewLevel()
	d := Dungeon{
		Name:   "Dungeon 1",
		Levels: []Level{l},
	}
	gm := GameMap{
		Dungeons: []Dungeon{d},
	}

	return gm
}
