package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	hudImg  *ebiten.Image   = nil
	hudErr  error           = nil
	hudFont font.Face       = nil
	tiles   []*ebiten.Image = nil
)

const fontSize = 24

func DrawHUD(g *Game, screen *ebiten.Image) {
	gd := g.GetData()
	uiY := gd.ScreenHeight * gd.TileHeight
	uiX := gd.ScreenWidth * gd.TileWidth / 2

	if hudImg == nil {
		hudImg, hudErr = createBorder(gd.ScreenWidth/2, gd.UIHeight, gd.TileWidth)
	}

	if hudFont == nil {
		const dpi = 72
		opts := &opentype.FaceOptions{
			Size:    fontSize,
			DPI:     dpi,
			Hinting: font.HintingFull,
		}
		hudFont, hudErr = loadFont("assets/fonts/ExpressionPro.ttf", opts)
		if hudErr != nil {
			log.Fatal(hudErr)
		}
	}

	fontX := uiX + fontSize
	fontY := uiY + fontSize + 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(uiX), float64(uiY))
	screen.DrawImage(hudImg, op)

	for _, p := range g.World.Query(g.WorldTags["players"]) {
		h := p.Components[healthC].(*Health)
		healthText := fmt.Sprintf("Health: %d/%d", h.Current, h.Max)
		text.Draw(screen, healthText, hudFont, fontX, fontY, color.White)
		fontY += fontSize + 1

		ac := p.Components[armorC].(*Armor)
		acText := fmt.Sprintf("Armor: %s (Dodge: %d, Def: %d, Block: %d)", ac.Name, ac.Dodge, ac.Defence, ac.Block)
		text.Draw(screen, acText, hudFont, fontX, fontY, color.White)
		fontY += fontSize + 1

		wp := p.Components[meleeWeaponC].(*MeleeWeapon)
		wpText := fmt.Sprintf("Weapon: %s (Dmg: %d-%d, ToHit: %d)", wp.Name, wp.MinDamage, wp.MaxDamage, wp.ToHitBonus)
		text.Draw(screen, wpText, hudFont, fontX, fontY, color.White)
	}
}
