package main

import "github.com/hajimehoshi/ebiten/v2"

func DrawPath(g *Game, screen *ebiten.Image) {
	// Draw path visualization if auto-moving
	if g.IsAutoMoving && len(g.AutoMovePath) > 0 && g.PathVisualizer != nil {
		pathWidth, pathHeight := g.PathVisualizer.GetSize()
		
		for i := g.AutoMoveIndex; i < len(g.AutoMovePath); i++ {
			pathPos := g.AutoMovePath[i]
			
			// Convert to screen coordinates
			screenX := (pathPos.X-g.viewport.X1)*g.gd.TileWidth + (g.gd.TileWidth-pathWidth)/2
			screenY := (pathPos.Y-g.viewport.Y1)*g.gd.TileHeight + (g.gd.TileHeight-pathHeight)/2

			// Only draw if on screen
			if screenX >= 0 && screenX < g.gd.ScreenWidth*g.gd.TileWidth &&
				screenY >= 0 && screenY < g.gd.ScreenHeight*g.gd.TileHeight {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(screenX), float64(screenY))
				screen.DrawImage(g.PathVisualizer.GetImage(), op)
			}
		}
	}
}
