package main

import "github.com/hajimehoshi/ebiten/v2"

func DrawRenderables(g *Game, l Level, screen *ebiten.Image, viewport Rect) {
	tileSize := 32
	offsetX := float64(viewport.X1 * tileSize)
	offsetY := float64(viewport.Y1 * tileSize)

	for _, result := range g.World.Query(g.WorldTags["renderables"]) {
		pos := result.Components[positionC].(*Position)
		img := result.Components[renderableC].(Renderable).GetImage()

		if l.PlayerVisible.IsVisible(pos.X, pos.Y) {
			index := l.GetIndexFromXY(pos.X, pos.Y)
			tile := l.Tiles[index]

			op := &ebiten.DrawImageOptions{}

			localX := float64(tile.PixelX) - offsetX
			localY := float64(tile.PixelY) - offsetY

			spriteOffsetX := float64((48 - tileSize) / 2)
			spriteOffsetY := float64(48 - tileSize) // align bottom of sprite with bottom of tile

			op.GeoM.Translate(localX-spriteOffsetX, localY-spriteOffsetY)
			screen.DrawImage(img, op)
		}
	}
}
