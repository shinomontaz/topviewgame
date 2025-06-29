package main

type Rect struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

func NewRect(x, y, w, h int) Rect {
	return Rect{X1: x, Y1: y, X2: x + w, Y2: y + h}
}

func (r *Rect) Center() (int, int) {
	return (r.X1 + r.X2) / 2, (r.Y1 + r.Y2) / 2
}

func (r *Rect) Intersect(other Rect) bool {
	return r.X1 <= other.X2 && r.X2 >= other.X1 && r.Y1 <= other.Y2 && r.Y2 >= other.Y1
}
