package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func ProcessUserLog(g *Game, screen *ebiten.Image) {
	if userLogImg == nil {
		userLogImg, _, err = ebitenutil.NewImageFromFile("assets/Back3.png")
		if err != nil {
			log.Fatal(err)
		}
	}

	if mplusNormalFont == nil {
		fontBytes, err := os.ReadFile("assets/fonts/ExpressionPro.ttf")
		if err != nil {
			log.Fatalf("failed to read font file: %v", err)
		}
		tt, err := opentype.Parse(fontBytes)

		const dpi = 72
		mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    24,
			DPI:     dpi,
			Hinting: font.HintingFull,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	gd := g.GetData()
	uiLocation := (gd.ScreenHeight) * gd.TileHeight
	var fontX = 16
	var fontY = uiLocation + 24
	op := &ebiten.DrawImageOptions{}
	screenW := screen.Bounds().Dx()
	imgW, imgH := userLogImg.Bounds().Dx(), userLogImg.Bounds().Dy()

	scaleX := float64(screenW) / float64(imgW)
	scaleY := float64(gd.UIHeight*gd.TileHeight) / float64(imgH)

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(0.), float64(uiLocation))
	screen.DrawImage(userLogImg, op)

	tmpMessages := make([]string, 0, 5)
	anyMessages := false
	for _, m := range g.World.Query(g.WorldTags["messengers"]) {
		messages := m.Components[userMessage].(*UserMessage)
		if messages.AttackMessage != "" {
			tmpMessages = append(tmpMessages, messages.AttackMessage)
			anyMessages = true
			messages.AttackMessage = ""
		}
	}
	for _, m := range g.World.Query(g.WorldTags["messengers"]) {
		messages := m.Components[userMessage].(*UserMessage)
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
		lastText = tmpMessages
	}
	for _, msg := range lastText {
		if msg != "" {
			text.Draw(screen, msg, mplusNormalFont, fontX, fontY, color.White)
			fontY += 16
		}
	}
}
