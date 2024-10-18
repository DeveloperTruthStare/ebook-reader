package main

import (
	"log"
	"ulla/ebook"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Change this to the directory to search for .epub files
	searchDir := "../../ebooks/"
	loadedEpub := 1

	epubFiles := ebook.FindEBooks(searchDir)

	bookULLAFormat, err := ebook.LoadEBook(epubFiles[loadedEpub])
	if err != nil {
		log.Fatal(err)
	}

	game := Game{book: bookULLAFormat, currentPage: 0}
	game.Init()
	ebiten.SetWindowSize(540, 760)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
