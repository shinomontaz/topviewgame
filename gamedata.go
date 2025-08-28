package main

type GameData struct {
	ScreenWidth  int // in tiles
	ScreenHeight int

	MapWidth  int
	MapHeight int

	TileWidth  int
	TileHeight int

	UIHeight int
}

func NewGameData() GameData {
	return GameData{
		ScreenWidth:  40,
		ScreenHeight: 25,

		MapWidth:  80,
		MapHeight: 50,

		TileWidth:  32,
		TileHeight: 32,
		UIHeight:   5,
	}
}

func (gd GameData) GetIndexFromXY(x, y int) int {
	return y*gd.MapWidth + x
}
