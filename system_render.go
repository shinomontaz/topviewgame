package main

import "github.com/hajimehoshi/ebiten/v2"

func ProcessRenderables(g *Game, l Level, screen *ebiten.Image) {
	for _, result := range g.World.Query(g.WorldTags["renderables"]) {
		pos := result.Components[positionC].(*Position)
		img := result.Components[renderableC].(Renderable).GetImage()

		if l.PlayerVisible.IsVisible(pos.X, pos.Y) {
			index := l.GetIndexFromXY(pos.X, pos.Y)
			tile := l.Tiles[index]

			op := &ebiten.DrawImageOptions{}

			offsetX := float64((48 - 32) / 2)
			offsetY := float64(48 - 32) // align bottom of sprite with bottom of tile

			op.GeoM.Translate(float64(tile.PixelX)-offsetX, float64(tile.PixelY)-offsetY)
			screen.DrawImage(img, op)
		}

	}
}
