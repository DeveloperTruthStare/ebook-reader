package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"log"
	"os"
	"path/filepath"
	"ulla/ebook"
	wordrenderer "ulla/word_renderer"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	book            *ebook.ULLABook
	currentPage     uint
	justPressed     bool
	width           uint
	height          uint
	mplusNormalFont font.Face
}

func (g *Game) Init() {
	fontData, err := os.ReadFile("../fonts/NotoSansJP-Regular.ttf")
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
		_, y := ebiten.CursorPosition()
		g.justPressed = true
		if y > int(g.height-100) {
			g.currentPage++
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

func (g *Game) GetImage(href string) (*image.Image, error) {
	for _, img := range g.book.Book.Manifest.Items {
		if filepath.Base(href) == filepath.Base(img.HREF) {
			f, err := img.Open()
			if err != nil {
				return nil, fmt.Errorf("unable to open image: %w", err)
			}
			defer f.Close()
			imgData, _, err := image.Decode(f)
			if err != nil {
				return nil, fmt.Errorf("unable to decode image: %w", err)
			}
			return &imgData, nil
		}
	}
	return nil, nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.book.Chunks[g.currentPage].IsImage {
		img, err := g.GetImage(g.book.Chunks[g.currentPage].Text.Text)
		if err != nil {
			log.Fatal(err)
		}
		ebitenImg := ebiten.NewImageFromImage(*img)
		screen.DrawImage(ebitenImg, &ebiten.DrawImageOptions{})
	} else {
		// Draw the text to the screen
		line := 0
		margin := 50
		xOffset := 0
		yPadding := 1
		xPadding := 1
		word_renderer := wordrenderer.New()
		for _, word := range g.book.Chunks[g.currentPage].Text.Tokens {
			img, err := word_renderer.RenderWord(word.Text)
			if err != nil {
				log.Fatal(err)
			}
			if img.Bounds().Dx()+xOffset > int(g.width)-margin*2 {
				xOffset = 0
				line++
			}
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(margin+xOffset), float64(28+line*(img.Bounds().Dy()+yPadding)))
			screen.DrawImage(ebiten.NewImageFromImage(img), opts)
			xOffset += img.Bounds().Dx() + xPadding
		}
	}
	// Draw Continue Button
	continueButton := image.NewRGBA(image.Rect(0, 0, int(g.width)-50, 50))
	draw.Draw(continueButton, continueButton.Bounds(), &image.Uniform{color.RGBA{R: 0x24, G: 0x94, B: 0xd1, A: 0xff}}, image.Point{}, draw.Src)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(25, float64(g.height)-75)
	screen.DrawImage(ebiten.NewImageFromImage(continueButton), opts)
	ebiten.NewImageFromImage(continueButton)
}
