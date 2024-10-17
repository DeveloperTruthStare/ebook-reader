package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/taylorskalyo/goreader/epub"
	"golang.org/x/net/html"
)

func IsWhitespaceOnly(text string) bool {
	for _, r := range text {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func printText(n *html.Node, book *ULLABook) string {
	toReturn := ""
	if n.Data == "head" {
		return ""
	}
	if n.Data == "body" {
		for _, d := range n.Attr {
			if d.Key == "class" {
				if d.Val == "p-caution" || d.Val == "p-colophon" {
					return ""
				}
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "p" {
		book.chunks = append(book.chunks, Chunk{isImage: false, text: ""})
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		book.chunks[len(book.chunks)-1].isImage = true
		for _, d := range n.Attr {
			if d.Key == "src" {
				toReturn += d.Val
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "ruby" {
		var kanjiText string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				kanjiText += c.Data
			} else if c.Type == html.ElementNode && c.Data == "rt" {
				continue
			}
		}
		if kanjiText != "" {
			toReturn += kanjiText
		}
	} else if n.Type == html.TextNode {
		text := n.Data

		if IsWhitespaceOnly(text) {
			return ""
		}

		toReturn += strings.TrimSpace(text)
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			toReturn += printText(c, book)
		}
		if n.Type == html.ElementNode && n.Data == "p" {

			book.chunks[len(book.chunks)-1].text = toReturn

			if len(strings.TrimSpace(toReturn)) <= 0 {
				book.chunks = book.chunks[:len(book.chunks)-1]
			} else {

			}

			toReturn = ""
		}
	}

	return toReturn
}

type ULLABook struct {
	book   *epub.Rootfile
	chunks []Chunk
}

type Chunk struct {
	isImage bool
	text    string
}

func AnalyzeBook(book *epub.Rootfile) (*ULLABook, error) {
	formattedBook := ULLABook{book: book}

	for _, section := range book.Spine.Itemrefs {
		f, err := section.Open()
		if err != nil {
			return nil, fmt.Errorf("Unable to open file %w", err)
		}
		defer f.Close()
		doc, err := html.Parse(f)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse file: %w", err)
		}

		printText(doc, &formattedBook)
	}
	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	for _, chunk := range formattedBook.chunks {
		if chunk.isImage {
			f.WriteString("IMAGE: ")
		} else {
			f.WriteString("TEXT: ")
		}
		f.WriteString(chunk.text + "\n")
	}
	f.Close()

	return &formattedBook, nil
}

func main() {
	loadedEpub := 1
	var epubFiles []string

	searchDir := "./ebooks/"
	directories, err := os.ReadDir(searchDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range directories {
		if dir.IsDir() {
			epubFiles = append(epubFiles, dir.Name())
		}
	}

	filename := searchDir + epubFiles[loadedEpub] + "/" + epubFiles[loadedEpub] + ".epub"

	rc, err := epub.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	book := rc.Rootfiles[0]

	bookULLAFormat, err := AnalyzeBook(book)
	if err != nil {
		log.Fatal(err)
	}

	game := Game{book: bookULLAFormat, currentPage: 0}
	game.Init()

	ebiten.SetWindowSize(540, 960)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
