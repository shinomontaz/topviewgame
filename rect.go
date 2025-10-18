package main

import core "topviewgame/internal/core"

type Rect = core.Rect

func NewRect(x, y, w, h int) Rect {
	return core.NewRect(x, y, w, h)
}
