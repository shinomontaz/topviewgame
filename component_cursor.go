package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Cursor struct {
	sX, sY int // screen position
	mX, mY int // map-space coordinates
	//	Img            *ebiten.Image // cursor image
	scaleX, scaleY float64
	Width          int
	Height         int

	state  int
	images map[int]*ebiten.Image
}

func NewCursor(w, h int) (*Cursor, error) {
	c := Cursor{
		Width:  w,
		Height: h,
		images: make(map[int]*ebiten.Image),
	}

	imgDefault, _, err := ebitenutil.NewImageFromFile("assets/cursor.png")
	if err != nil {
		return nil, err
	}
	c.images[CursorDefault] = imgDefault
	c.scaleX = float64(w) / float64(imgDefault.Bounds().Dx())
	c.scaleY = float64(h) / float64(imgDefault.Bounds().Dy())

	imgAttack, _, err := ebitenutil.NewImageFromFile("assets/cursor_attack.png")
	if err != nil {
		return nil, err
	}
	c.images[CursorAttack] = imgAttack

	imgInfo, _, err := ebitenutil.NewImageFromFile("assets/cursor_info.png")
	if err != nil {
		return nil, err
	}
	c.images[CursorSelect] = imgInfo

	imgForbidden, _, err := ebitenutil.NewImageFromFile("assets/cursor_forbidden.png")
	if err != nil {
		return nil, err
	}
	c.images[CursorForbidden] = imgForbidden

	return &c, nil
}

func (c *Cursor) Update(x, y int, vp Rect, tileW, tileH int) {
	c.sX, c.sY = x, y
	c.mX = x/tileW + vp.X1
	c.mY = y/tileH + vp.Y1
}

func (c *Cursor) SetScreenPos(x, y int) {
	c.sX = x
	c.sY = y
}

func (c *Cursor) SetMapPos(x, y int) {
	c.mX = x
	c.mY = y
}

func (c *Cursor) ScreenPos() (int, int) { return c.sX, c.sY }
func (c *Cursor) MapPos() (int, int)    { return c.mX, c.mY }

func (c *Cursor) SetState(s int) { c.state = s }

func (c *Cursor) GetImage() *ebiten.Image {
	return c.images[c.state]
}

const (
	CursorDefault = iota
	CursorAttack
	CursorSelect
	CursorForbidden
)
