package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/taylorskalyo/goreader/epub"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	book            *epub.Rootfile
	currentPage     uint
	justPressed     bool
	width           uint
	height          uint
	mplusNormalFont font.Face
}

func (g *Game) Init() {
	fontData, err := os.ReadFile("./fonts/NotoSansJP-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	g.mplusNormalFont = mplusNormalFont
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !g.justPressed {
		x, _ := ebiten.CursorPosition()
		g.justPressed = true
		if x > int(g.width/2) {
			g.currentPage++
		} else {
			g.currentPage--
		}
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.justPressed = false
	}
	return nil
}

func (g *Game) Layout(width, height int) (int, int) {
	g.width = uint(width)
	g.height = uint(height)
	return width, height
}

func (g *Game) Draw(screen *ebiten.Image) {

	for _, item := range g.book.Manifest.Items {
		if item.MediaType == "image/jpeg" {
			fmt.Println(item.HREF)
			f, err := item.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			imgData, format, err := image.Decode(f)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(format)
			ebitenImg := ebiten.NewImageFromImage(imgData)
			screen.DrawImage(ebitenImg, &ebiten.DrawImageOptions{})
			break
		}
	}

	msg := "こんにちは、世界"
	text.Draw(screen, msg, g.mplusNormalFont, 50, 100, color.White)
}
