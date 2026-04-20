package main

import "github.com/hajimehoshi/ebiten/v2"

func DrawRenderables(g *Game, l Level, screen *ebiten.Image, viewport Rect) {
	tileWidth := g.gd.TileWidth
	tileHeight := g.gd.TileHeight
	offsetX := float64(viewport.X1 * tileWidth)
	offsetY := float64(viewport.Y1 * tileHeight)

	for _, result := range g.World.QueryRenderables() {
		pos := g.World.GetPosition(result)
		renderable := g.World.GetRenderable(result).(Renderable)
		img := renderable.GetImage()

		if l.PlayerVisible.IsVisible(pos.X, pos.Y) {
			index := l.GetIndexFromXY(pos.X, pos.Y)
			tile := l.Tiles[index]

			op := &ebiten.DrawImageOptions{}

			localX := float64(tile.PixelX) - offsetX
			localY := float64(tile.PixelY) - offsetY

			spriteOffsetX, spriteOffsetY := renderable.GetOffset(tileWidth, tileHeight)

			op.GeoM.Translate(localX-spriteOffsetX, localY-spriteOffsetY)
			screen.DrawImage(img, op)
		}
	}
}
