package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Map GameMap
}

func NewGame() *Game {
	return &Game{
		Map: NewGameMap(),
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	gd := NewGameData()
	return gd.TileHeight * gd.ScreenWidth, gd.TileWidth * gd.ScreenHeight
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	gd := NewGameData()
	level := g.Map.Dungeons[0].Levels[0]
	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < gd.ScreenHeight; y++ {
			tile := level.Tiles[getIndexFromXY(x, y)]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			screen.DrawImage(tile.Image, op)
		}
	}

}

func main() {
	g := NewGame()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("SGT Calculator")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
