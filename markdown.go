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

type markdownSelect struct {
	pattern string // CSS select pattern for cascadia
	printer printer
}

func (g markdownSelect) grep(path string, buf []byte) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}
	defer f.Close()

	// TODO Use the provided buffer?
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
	html_doc, err := html.Parse(bytes.NewReader(html_bytes))
	if err != nil {
		log.Fatalf("error parsing html '%s': %s\n", html_bytes, err)
	}

	// apply cascadia to the doc
	// TODO parse pattern at create, should be shared whole run.
	sel, err := cascadia.Parse(g.pattern)
	if err != nil {
		log.Fatalf("invalid pattern '%s': %s\n", g.pattern, err)
	}

	match := match{path: path}
	for _, p := range cascadia.QueryAll(html_doc, sel) {
		// TODO is there any way to get line and col #?
		match.add(0, 0, g.nodeStr(p), true)
	}
	g.printer.print(match)
}

func (g markdownSelect) nodeStr(node *html.Node) string {
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
