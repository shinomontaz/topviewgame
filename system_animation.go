package main

// UpdateAnimations ticks all entities with animated renderables
func UpdateAnimations(dt float64, g *Game) {
	for _, result := range g.World.Query(g.WorldTags["renderables"]) {
		renderable := result.Components[renderableC].(Renderable)

		// If the renderable supports time-based updates
		if updatable, ok := renderable.(interface{ Update(float64) }); ok {
			updatable.Update(dt)
		}
	}
}
