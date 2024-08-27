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
	pattern cascadia.Sel // CSS select pattern for cascadia
	printer printer
}

func NewMarkdownSelect(pattern pattern, printer printer) markdownSelect {
	// parse the pattnern
	sel, err := cascadia.Parse(string(pattern.pattern))
	if err != nil {
		log.Fatalf("invalid pattern '%s': %s\n", string(pattern.pattern), err)
	}
	return markdownSelect{
		pattern: sel,
		printer: printer,
	}
}

func (g markdownSelect) grep(path string, buf []byte) {
	// parse the new doc
	html_doc, err := loadMarkdownFile(path)
	if err != nil {
		log.Fatalf("error parsing file '%s': %s\n", path, err)
	}

	match := match{path: path}
	for _, p := range cascadia.QueryAll(html_doc, g.pattern) {
		// TODO is there any way to get line and col #?
		match.add(0, 0, nodeStr(p), true)
	}
	g.printer.print(match)
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
