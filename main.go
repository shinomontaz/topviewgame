package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(w, h int) (int, int) {
	return 1280, 800
}

func main() {
	g := NewGame()

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("SGT Calculator")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
