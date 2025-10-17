package main

import "github.com/hajimehoshi/ebiten/v2"

func DrawRenderables(g *Game, l Level, screen *ebiten.Image, viewport Rect) {
	tileWidth := g.gd.TileWidth
	tileHeight := g.gd.TileHeight
	offsetX := float64(viewport.X1 * tileWidth)
	offsetY := float64(viewport.Y1 * tileHeight)

	for _, result := range g.World.Query(g.WorldTags["renderables"]) {
		pos := result.Components[positionC].(*Position)
		img := result.Components[renderableC].(Renderable).GetImage()

		if l.PlayerVisible.IsVisible(pos.X, pos.Y) {
			index := l.GetIndexFromXY(pos.X, pos.Y)
			tile := l.Tiles[index]

			op := &ebiten.DrawImageOptions{}

			localX := float64(tile.PixelX) - offsetX
			localY := float64(tile.PixelY) - offsetY

			spriteOffsetX := float64((48 - tileWidth) / 2)
			spriteOffsetY := float64(48 - tileHeight)

			op.GeoM.Translate(localX-spriteOffsetX, localY-spriteOffsetY)
			screen.DrawImage(img, op)
		}
	}
}
