package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/taylorskalyo/goreader/epub"
)

type Game struct {
	book        *epub.Rootfile
	currentPage uint
	justPressed bool
	width       uint
	height      uint
}

func (g *Game) Update() error {
	rf := g.book.Spine.Itemrefs[g.currentPage]
	fmt.Println(rf.HREF)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !g.justPressed {
		// Get the cursor position (x, y) in the window
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

}
