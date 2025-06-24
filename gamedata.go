package main

type GameData struct {
	ScreenWidth  int // in tiles
	ScreenHeight int

	TileWidth  int
	TileHeight int
}

func NewGameData() GameData {
	return GameData{
		ScreenWidth:  40,
		ScreenHeight: 25,
		TileWidth:    32,
		TileHeight:   32,
	}
}

func (gd GameData) GetIndexFromXY(x, y int) int {
	return y*gd.ScreenWidth + x
}
