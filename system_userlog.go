package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	userLogImg      *ebiten.Image = nil
	err             error         = nil
	mplusNormalFont font.Face     = nil
	lastText        []string      = make([]string, 0, 5)
)

func DrawUserLog(g *Game, screen *ebiten.Image) {
	if userLogImg == nil {
		gd := g.GetData()
		userLogImg, err = createBorder(gd.ScreenWidth/2, gd.UIHeight, gd.TileWidth)
		if err != nil {
			log.Println("failed to create user log border:", err)
			return
		}
	}

	if mplusNormalFont == nil {
		const dpi = 72
		opts := opentype.FaceOptions{
			Size:    24,
			DPI:     dpi,
			Hinting: font.HintingFull,
		}
		mplusNormalFont, err = loadFont("assets/fonts/ExpressionPro.ttf", &opts)
		if err != nil {
			log.Println("failed to load user log font:", err)
			return
		}
	}

	gd := g.GetData()
	uiLocation := (gd.ScreenHeight) * gd.TileHeight
	var fontX = 16
	var fontY = uiLocation + 24
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0.), float64(uiLocation))
	screen.DrawImage(userLogImg, op)

	tmpMessages := make([]string, 0, 5)
	anyMessages := false
	for _, m := range g.World.QueryMessengers() {
		messages := g.World.GetUserMessage(m)
		if messages.AttackMessage != "" {
			tmpMessages = append(tmpMessages, messages.AttackMessage)
			anyMessages = true
			messages.AttackMessage = ""
		}
	}
	for _, m := range g.World.QueryMessengers() {
		messages := g.World.GetUserMessage(m)
		if messages.DeadMessage != "" {
			tmpMessages = append(tmpMessages, messages.DeadMessage)
			anyMessages = true
			messages.DeadMessage = ""
		}
		if messages.GameStateMessage != "" {
			tmpMessages = append(tmpMessages, messages.GameStateMessage)
			anyMessages = true
		}

	}
	if anyMessages {
		for _, message := range tmpMessages {
			if message == "" {
				continue
			}
			lastText = append(lastText, message)
			if len(lastText) > 5 {
				lastText = lastText[len(lastText)-5:]
			}
		}
	}
	for _, msg := range lastText {
		if msg != "" {
			text.Draw(screen, msg, mplusNormalFont, fontX, fontY, color.White)
			fontY += 16
		}
	}
}
