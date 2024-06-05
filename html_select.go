package mack

import (
	"fmt"
	"log"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type htmlSelect struct {
	pattern string
	printer printer
}

func (g htmlSelect) grep(path string, buf []byte) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}
	defer f.Close()

	doc, err := html.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	sel, err := cascadia.Parse(g.pattern)
	if err != nil {
		log.Fatalf("invalid pattern '%s': %s\n", g.pattern, err)
	}
	for i, p := range cascadia.QueryAll(doc, sel) {
		fmt.Printf(">> %s #%d: %s\n", path, i, text(p))
	}
}

func text(n *html.Node) string {
	result := ""
	if n.Type == html.TextNode {
		result = n.Data
	}

	// traverses the HTML of the webpage from the first child node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = result + " " + text(c)
	}
	return strings.TrimSpace(result)

}
