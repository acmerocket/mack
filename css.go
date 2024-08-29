package mack

import (
	"bytes"
	"io"
	"log"
	"path/filepath"

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

func NewCssSelect(pattern regex_pattern, printer printer) cssSelect {
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
	fileSpec, ok := GetLanguageSpec(path)
	if ok {
		switch fileSpec.Name {
		case "html":
			doc, err := loadHtmlFile(path)
			if err != nil {
				log.Fatalf("Unable to load file: %s, %s\n", path, err)
			}
			list := cascadia.QueryAll(doc, g.pattern)
			g.printHtmlNode(path, list)
		case "markdown":
		case "md":
			doc, err := loadMarkdownFile(path)
			if err != nil {
				log.Fatalf("Unable to load file: %s, %s\n", path, err)
			}
			list := cascadia.QueryAll(doc, g.pattern)
			g.printHtmlNode(path, list)
		default:
			log.Println("unknown file extention, skipping: ", filepath.Ext(path), " ", path)
		}
	} else {
		log.Fatalf("Unknow file type %s\n", path)
	}
}

func (g cssSelect) printHtmlNode(path string, list []*html.Node) {
	match := match{path: path}
	for _, p := range list {
		// TODO is there any way to get line and col #?
		match.add(0, 0, RenderHtml(p), true)
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

// what is area for "extracting text?" default css results?
func RenderHtml(node *html.Node) string {
	var buf bytes.Buffer
	err := html.Render(&buf, node) // TODO - abstract render(stream) insteam of print, render(s, result)
	if err != nil {
		log.Fatal("Error rendering ", node.Data, err)
	}

	return buf.String()
}
