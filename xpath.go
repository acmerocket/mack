package mack

import (
	"bufio"
	"log"
	"path/filepath"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"

	"golang.org/x/net/html"
)

type xpathSelect struct {
	pattern string // xpath
	printer printer
}

func NewXpathQuery(pattern regex_pattern, printer printer) xpathSelect {
	return xpathSelect{
		pattern: string(pattern.pattern),
		printer: printer,
	}
}

func (g xpathSelect) grep(path string, buf []byte) {
	fileSpec, ok := GetLanguageSpec(path)
	if ok {
		switch fileSpec.Name {
		case "html":
			doc, err := loadHtmlFile(path)
			if err != nil {
				log.Fatalf("Unable to load file: %s, %s\n", path, err)
			}
			list, err := htmlquery.QueryAll(doc, g.pattern)
			//log.Println("query:", g.pattern, list, nodeStr(doc))
			if err != nil {
				log.Fatalf("Error running query %s, on %s: %s\n", g.pattern, path, err)
			}
			g.printHtmlNode(path, list)
		case "markdown":
		case "md":
			doc, err := loadMarkdownFile(path)
			if err != nil {
				log.Fatalf("Unable to load file: %s, %s\n", path, err)
			}
			list, err := htmlquery.QueryAll(doc, g.pattern)
			//log.Println("query:", g.pattern, list, nodeStr(doc))
			if err != nil {
				log.Fatalf("Error running query %s, on %s: %s\n", g.pattern, path, err)
			}
			g.printHtmlNode(path, list)
		case "xml":
			doc, err := loadXmlFile(path)
			if err != nil {
				log.Fatalf("error parsing file '%s': %s\n", path, err)
			}
			list, err := xmlquery.QueryAll(doc, g.pattern)
			if err != nil {
				log.Fatalf("Error running query %s, on %s: %s\n", g.pattern, path, err)
			}
			g.printXmlNode(path, list)
		case "json":
			doc, err := g.loadJsonFile(path)
			if err != nil {
				log.Fatalf("error parsing file '%s': %s\n", path, err)
			}
			list, err := jsonquery.QueryAll(doc, g.pattern)
			if err != nil {
				log.Fatalf("Error running query %s, on %s: %s\n", g.pattern, path, err)
			}
			g.printJsonNode(path, list)
		default:
			log.Println("unknown file extention, skipping: ", filepath.Ext(path), " ", path)
		}
	} else {
		log.Fatalf("Unknow file type %s\n", path)
	}
}

func (g xpathSelect) printHtmlNode(path string, list []*html.Node) {
	match := match{path: path}
	for _, p := range list {
		// TODO is there any way to get line and col #?
		match.add(0, 0, p, true)
	}
	g.printer.print(match)
}

func (g xpathSelect) printXmlNode(path string, list []*xmlquery.Node) {
	match := match{path: path}
	for _, p := range list {
		// TODO is there any way to get line and col #?
		match.add(0, 0, p, true)
	}
	g.printer.print(match)
}

func (g xpathSelect) printJsonNode(path string, list []*jsonquery.Node) {
	match := match{path: path}
	for _, p := range list {
		// TODO is there any way to get line and col #?
		// FIXME - json result?
		match.add(0, 0, p, true)
	}
	g.printer.print(match)
}

func loadHtmlFile(path string) (*html.Node, error) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("loadHtmlFile: %s\n", err)
	}
	defer f.Close()

	return html.Parse(bufio.NewReader(f))
}

func loadXmlFile(path string) (*xmlquery.Node, error) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("loadXmlFile: %s\n", err)
	}
	defer f.Close()

	return xmlquery.Parse(bufio.NewReader(f))
}

func (g xpathSelect) loadJsonFile(path string) (*jsonquery.Node, error) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("loadJsonFile: %s\n", err)
	}
	defer f.Close()

	return jsonquery.Parse(bufio.NewReader(f))
}

// func hasExtension(path, type_str string) bool {
// 	// based on ext (as file type)
// 	type_rec := known_languages[type_str]
// 	for _, ext := range type_rec.Exts {
// 		if strings.HasSuffix(path, "."+ext) {
// 			return true
// 		}

// 	}
// 	return false
// }
