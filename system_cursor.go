package main

import "github.com/hajimehoshi/ebiten/v2"

func UpdateCursor(g *Game) {
	x, y := g.PlayerController.GetCursor()
	gd := g.GetData()
	viewport := g.Viewport()
	for _, res := range g.World.Query(g.WorldTags["cursors"]) {
		cursor := res.Components[cursorC].(*Cursor)
		cursor.Update(x, y, viewport, gd.TileWidth, gd.TileHeight)

		mx, my := cursor.MapPos()

		// Check map objects to set state
		if mx == g.Center.X && my == g.Center.Y {
			cursor.SetState(CursorSelect)
		} else {
			cursor.SetState(CursorDefault)
		}
	}
}

func DrawCursor(g *Game, screen *ebiten.Image) {
	for _, res := range g.World.Query(g.WorldTags["cursors"]) {
		cursor := res.Components[cursorC].(*Cursor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(cursor.scaleX, cursor.scaleY)
		sx, sy := cursor.ScreenPos()
		op.GeoM.Translate(float64(sx), float64(sy))
		screen.DrawImage(cursor.GetImage(), op)
	}
}
