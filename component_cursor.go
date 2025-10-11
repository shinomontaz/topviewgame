package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Cursor struct {
	sX, sY int // screen position
	mX, mY int // map-space coordinates
	tX, tY int // tile position
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


	imgTileSelect, _, err := ebitenutil.NewImageFromFile("assets/tile_selection.png")
	if err != nil {
		return nil, err
	}
	c.images[Selection] = imgTileSelect

	return &c, nil
}

func (c *Cursor) Update(x, y int, vp Rect, tileW, tileH int) {
	c.sX, c.sY = x, y
	c.tX = x / tileW
	c.tY = y / tileH
	c.mX = c.tX + vp.X1
	c.mY = c.tY + vp.Y1
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

func (c *Cursor) GetSelection() *ebiten.Image {
	return c.images[Selection]
}


const (
	CursorDefault = iota
	CursorAttack
	CursorSelect
	CursorForbidden
	Selection
)
