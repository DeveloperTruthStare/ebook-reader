package ebook

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ikawaha/kagome/tokenizer"
	"github.com/luo-xuanxuan/gosho"
	"github.com/taylorskalyo/goreader/epub"
	"golang.org/x/net/html"
)

type ULLABook struct {
	Book   *epub.Rootfile
	Chunks []Chunk
}

type Chunk struct {
	IsImage bool
	Text    Sentence
}

type Sentence struct {
	Text   string
	Tokens []gosho.Word
}

func FindEBooks(searchDir string) []string {
	epubFiles := []string{}
	files, err := os.ReadDir(searchDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			epubFiles = append(epubFiles, FindEBooks(searchDir+file.Name()+"/")...)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			if filepath.Ext(file.Name()) == ".epub" {
				epubFiles = append(epubFiles, searchDir+file.Name())
			}
		}
	}
	return epubFiles
}

func IsWhitespaceOnly(text string) bool {
	for _, r := range text {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func NewChunk() *Chunk {
	return &Chunk{
		IsImage: false,
		Text: Sentence{
			Text:   "",
			Tokens: []gosho.Word{},
		},
	}
}

func ProcessHtml(n *html.Node, book *ULLABook) string {
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
	if n.Type == html.ElementNode && n.Data == "img" {
		book.Chunks[len(book.Chunks)-1].IsImage = true
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
			toReturn += ProcessHtml(c, book)
		}
		if n.Type == html.ElementNode && n.Data == "p" {

			if !book.Chunks[len(book.Chunks)-1].IsImage {
				content := ""
				inQuote := false
				for _, char := range toReturn {
					content += string(char)
					if char == '「' {
						// Add to the special stack
						inQuote = true
					} else if char == '」' {
						inQuote = false
					} else if char == '。' {
						if !inQuote {
							book.Chunks[len(book.Chunks)-1].Text.Text = content
							book.Chunks = append(book.Chunks, *NewChunk())
							content = ""
						}
					}
				}
				if content != "" {
					book.Chunks[len(book.Chunks)-1].Text.Text = content
					book.Chunks = append(book.Chunks, *NewChunk())
				}
			} else {
				book.Chunks[len(book.Chunks)-1].Text.Text = toReturn
				book.Chunks = append(book.Chunks, *NewChunk())
			}

			toReturn = ""
		}
	}

	return toReturn
}

func LoadEBook(epubFile string) (*ULLABook, error) {
	rc, err := epub.OpenReader(epubFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open .epub file: %w", err)
	}

	book := rc.Rootfiles[0]

	formattedBook := ULLABook{Book: book, Chunks: []Chunk{*NewChunk()}}

	for _, section := range book.Spine.Itemrefs {
		f, err := section.Open()
		if err != nil {
			return nil, fmt.Errorf("unable to open file %w", err)
		}
		defer f.Close()
		doc, err := html.Parse(f)
		if err != nil {
			return nil, fmt.Errorf("unable to parse file: %w", err)
		}

		ProcessHtml(doc, &formattedBook)
		t := tokenizer.New()
		for i, sentence := range formattedBook.Chunks {
			if len(sentence.Text.Tokens) == 0 {
				tokens := t.Analyze(formattedBook.Chunks[i].Text.Text, tokenizer.Normal)
				for _, token := range tokens {
					if token.Class != tokenizer.DUMMY {
						formattedBook.Chunks[i].Text.Tokens = append(formattedBook.Chunks[i].Text.Tokens, gosho.Word{Text: token.Surface})
					}
				}
			}
		}
	}
	return &formattedBook, nil
}
