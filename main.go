package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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

	for _, book := range epubFiles {
		UnzipEPUB(searchDir+book+"/"+book+".epub", searchDir+book+"/unzipped/")
	}

	filename := searchDir + epubFiles[1] + "/" + epubFiles[1] + ".epub"

	rc, err := epub.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	book := rc.Rootfiles[0]

	game := Game{book: book, currentPage: 0}
	game.Init()

	ebiten.SetWindowSize(4000, 4000)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func UnzipEPUB(epubPath, outputDir string) error {
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return fmt.Errorf("failed to open EPUB file: %w", err)
	}
	defer r.Close()

	for _, file := range r.File {
		outputPath := filepath.Join(outputDir, file.Name)

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(outputPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to open file in EPUB: %w", err)
		}

		zipFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer zipFile.Close()

		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer outputFile.Close()

		_, err = io.Copy(outputFile, zipFile)
		if err != nil {
			return fmt.Errorf("failed to copy file contents: %w", err)
		}

	}

	return nil
}

type Book struct {
	pages []Page
}

type Page struct {
	isImage bool
	text    string
}
