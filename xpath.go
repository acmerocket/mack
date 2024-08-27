package mack

import (
	"bufio"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xmlquery"

	"golang.org/x/net/html"
)

type xpathSelect struct {
	pattern string // xpath
	printer printer
}

func NewXpathQuery(pattern pattern, printer printer) xpathSelect {
	// parse the pattnern
	//sel, err := cascadia.Parse(string(pattern.pattern))
	//if err != nil {
	//	log.Fatalf("invalid pattern '%s': %s\n", string(pattern.pattern), err)
	//}
	return xpathSelect{
		pattern: string(pattern.pattern),
		printer: printer,
	}
}

func (g xpathSelect) grep(path string, buf []byte) {
	if hasExtension(path, "html") {
		doc, err := loadHtmlFile(path)
		if err != nil {
			log.Fatalf("Unable to load file: %s, %s\n", path, err)
		}
		list, err := htmlquery.QueryAll(doc, g.pattern)
		if err != nil {
			log.Fatalf("Error running query %s, on %s: %s\n", g.pattern, path, err)
		}
		g.printHtmlNode(path, list)
	} else if hasExtension(path, "markdown") {
		// parse the new html doc
		doc, err := loadMarkdownFile(path)
		if err != nil {
			log.Fatalf("error parsing file '%s': %s\n", path, err)
		}
		list, err := htmlquery.QueryAll(doc, g.pattern)
		if err != nil {
			log.Fatalf("Error running query %s, on %s: %s\n", g.pattern, path, err)
		}
		g.printHtmlNode(path, list)
	} else if hasExtension(path, "xml") {
		doc, err := loadXmlFile(path)
		if err != nil {
			log.Fatalf("error parsing file '%s': %s\n", path, err)
		}
		list, err := xmlquery.QueryAll(doc, g.pattern)
		if err != nil {
			log.Fatalf("Error running query %s, on %s: %s\n", g.pattern, path, err)
		}
		g.printXmlNode(path, list)
	} else {
		log.Fatal("unknown file type: ", path)
	}
}

func (g xpathSelect) printHtmlNode(path string, list []*html.Node) {
	match := match{path: path}
	for _, p := range list {
		// TODO is there any way to get line and col #?
		match.add(0, 0, nodeStr(p), true)
	}
	g.printer.print(match)
}

func (g xpathSelect) printXmlNode(path string, list []*xmlquery.Node) {
	match := match{path: path}
	for _, p := range list {
		// TODO is there any way to get line and col #?
		match.add(0, 0, p.InnerText(), true)
	}
	g.printer.print(match)
}

func loadHtmlFile(path string) (*html.Node, error) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}
	defer f.Close()

	return html.Parse(bufio.NewReader(f))
}

func loadXmlFile(path string) (*xmlquery.Node, error) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}
	defer f.Close()

	return xmlquery.Parse(bufio.NewReader(f))
}

func hasExtension(path, type_str string) bool {
	// based on ext (as file type)
	type_rec := known_languages[type_str]
	for _, ext := range type_rec.Exts {
		if strings.HasSuffix(path, "."+ext) {
			return true
		}

	}
	return false
}
