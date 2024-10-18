package wordrenderer

/*
 * Given a sentence -> 「寿司が食べたい」
 * Break it into chunks -> 「寿司」、「が」、「食べたい」
 * Create an image that is the width, height and contains the text
 *
 */

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"unicode/utf8"

	"github.com/golang/freetype"
)

type WordRenderer struct {
	Margin     uint16
	XPadding   uint16
	YPadding   uint16
	RuneWidth  uint16
	RuneHeight uint16
	Font       string
}

func New() *WordRenderer {
	return &WordRenderer{
		Margin:     50,
		XPadding:   5,
		YPadding:   5,
		RuneWidth:  25,
		RuneHeight: 25,
		Font:       "../fonts/NotoSansJP-Regular.ttf",
	}
}

func (wr *WordRenderer) RenderWord(word string) (*image.RGBA, error) {
	runeCount := utf8.RuneCountInString(word)
	container := image.NewRGBA(image.Rect(0, 0, int(wr.RuneWidth)*int(runeCount)+int(wr.XPadding*2), int(wr.RuneHeight)+int(wr.YPadding*2)))
	// Set the background color of the container
	draw.Draw(container, container.Bounds(), &image.Uniform{color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}}, image.Point{}, draw.Src)
	fontBytes, err := ioutil.ReadFile(wr.Font)
	if err != nil {
		return nil, fmt.Errorf("error loading font: %w", err)
	}
	fontParsed, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing font: %w", err)
	}

	ftContext := freetype.NewContext()
	ftContext.SetDPI(72)
	ftContext.SetFont(fontParsed)
	ftContext.SetFontSize(24)
	ftContext.SetClip(container.Bounds())
	ftContext.SetDst(container)
	ftContext.SetSrc(image.White)

	pt := freetype.Pt(5, 25)
	_, err = ftContext.DrawString(word, pt)
	if err != nil {
		return nil, fmt.Errorf("error drawing text: %w", err)
	}

	return container, nil
}

type Word struct {
	Text     string
	Furigane string
	PoS      string
}
