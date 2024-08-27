package mack

import (
	"bytes"
	"io"
	"log"

	"github.com/andybalholm/cascadia"
	"github.com/gomarkdown/markdown"
	ht "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/net/html"
)

type cssSelect struct {
	pattern cascadia.Sel // CSS select pattern for cascadia
	printer printer
}

func NewCssSelect(pattern pattern, printer printer) cssSelect {
	// parse the pattnern
	sel, err := cascadia.Parse(string(pattern.pattern))
	if err != nil {
		log.Fatalf("invalid pattern '%s': %s\n", string(pattern.pattern), err)
	}
	return cssSelect{
		pattern: sel,
		printer: printer,
	}
}

func (g cssSelect) grep(path string, buf []byte) {
	//var doc *html.Node = nil
	if hasExtension(path, "markdown") {
		// parse the new doc
		doc, err := loadMarkdownFile(path)
		if err != nil {
			log.Fatalf("error parsing file '%s': %s\n", path, err)
		}

		match := match{path: path}
		for _, p := range cascadia.QueryAll(doc, g.pattern) {
			// TODO is there any way to get line and col #?
			match.add(0, 0, nodeStr(p), true)
		}
		g.printer.print(match)
	} else if hasExtension(path, "html") {
		// parse the new doc
		doc, err := loadHtmlFile(path)
		if err != nil {
			log.Fatalf("error parsing file '%s': %s\n", path, err)
		}

		match := match{path: path}
		for _, p := range cascadia.QueryAll(doc, g.pattern) {
			// TODO is there any way to get line and col #?
			match.add(0, 0, nodeStr(p), true)
		}
		g.printer.print(match)
	} else {
		log.Println("Unknown file type: ", path)
	}
}

func loadMarkdownFile(path string) (*html.Node, error) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}
	defer f.Close()

	// TODO Use the provided buffer in grep?
	file_buf, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("error reading file '%s': %s\n", newOutputOption().ColorCodePath, err)
	}
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(file_buf)

	// create HTML renderer and render
	htmlFlags := ht.CommonFlags | ht.HrefTargetBlank
	opts := ht.RendererOptions{Flags: htmlFlags}
	renderer := ht.NewRenderer(opts)
	html_bytes := markdown.Render(doc, renderer)

	// parse the new html doc
	return html.Parse(bytes.NewReader(html_bytes))
}

// TODO - move to html package?
// what is area for "extracting text?" default css results?
func nodeStr(node *html.Node) string {
	var buf bytes.Buffer
	collectText(node, &buf)
	return buf.String()
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}
