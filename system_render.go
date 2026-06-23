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
		img := renderable.GetImage(tileWidth, tileHeight)

		if l.PlayerVisible.IsVisible(pos.X, pos.Y) {
			index := l.GetIndexFromXY(pos.X, pos.Y)
			tile := l.Tiles[index]

			localX := float64(tile.PixelX) - offsetX
			localY := float64(tile.PixelY) - offsetY

			img.GeoM.Translate(localX, localY-img.Height)
			screen.DrawImage(img.Image, &ebiten.DrawImageOptions{GeoM: img.GeoM})
		}
	}
}
