package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/taylorskalyo/goreader/epub"
	"golang.org/x/net/html"
)

func printText(n *html.Node) {
	if n.Type == html.TextNode {
		fmt.Print(n.Data)
	}

	// Recursively traverse the child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		printText(c)
	}
}

func main() {
	filename := "./ebooks/支倉凍砂 - 狼と香辛料.epub"

	rc, err := epub.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	book := rc.Rootfiles[0]
	game := Game{book: book, currentPage: 0}

	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
