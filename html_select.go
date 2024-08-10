package mack

import (
	"fmt"
	"log"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/fatih/color"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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
	//nodes := make([]*html.Node, 0)
	//for _, p := range cascadia.QueryAll(doc, sel) {
	//nodes = append(nodes, p)
	// create lines for each match
	//fmt.Printf("%s #%d: %s\n", path, i, p)
	//	line := line{
	//		num     int
	//		column  int
	//		text
	//		matched: true
	//}
	//lines = append(lines, line)
	//}
	displayer.Display(cascadia.QueryAll(doc, sel))
	//g.printer.print(match{path: path, lines: lines})
}

var (
	maxPrintLevel int       = -1
	preformatted  bool      = false
	printColor    bool      = false
	escapeHTML    bool      = true
	indentString  string    = " "
	displayer     Displayer = TreeDisplayer{}

	// Colors
	tagColor     *color.Color = color.New(color.FgCyan)
	tokenColor                = color.New(color.FgCyan)
	attrKeyColor              = color.New(color.FgMagenta)
	quoteColor                = color.New(color.FgBlue)
	commentColor              = color.New(color.FgYellow)
)

type Displayer interface {
	Display([]*html.Node)
}

type TreeDisplayer struct {
}

func (t TreeDisplayer) Display(nodes []*html.Node) {
	for _, node := range nodes {
		t.printNode(node, 0)
	}
}

// The <pre> tag indicates that the text within it should always be formatted
// as is. See https://github.com/ericchiang/pup/issues/33
func (t TreeDisplayer) printPre(n *html.Node) {
	switch n.Type {
	case html.TextNode:
		s := n.Data
		if escapeHTML {
			// don't escape javascript
			if n.Parent == nil || n.Parent.DataAtom != atom.Script {
				s = html.EscapeString(s)
			}
		}
		fmt.Print(s)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			t.printPre(c)
		}
	case html.ElementNode:
		fmt.Printf("<%s", n.Data)
		for _, a := range n.Attr {
			val := a.Val
			if escapeHTML {
				val = html.EscapeString(val)
			}
			fmt.Printf(` %s="%s"`, a.Key, val)
		}
		fmt.Print(">")
		if !isVoidElement(n) {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				t.printPre(c)
			}
			fmt.Printf("</%s>", n.Data)
		}
	case html.CommentNode:
		data := n.Data
		if escapeHTML {
			data = html.EscapeString(data)
		}
		fmt.Printf("<!--%s-->\n", data)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			t.printPre(c)
		}
	case html.DoctypeNode, html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			t.printPre(c)
		}
	}
}

// Print a node and all of it's children to `maxlevel`.
func (t TreeDisplayer) printNode(n *html.Node, level int) {
	switch n.Type {
	case html.TextNode:
		s := n.Data
		if escapeHTML {
			// don't escape javascript
			if n.Parent == nil || n.Parent.DataAtom != atom.Script {
				s = html.EscapeString(s)
			}
		}
		s = strings.TrimSpace(s)
		if s != "" {
			t.printIndent(level)
			fmt.Println(s)
		}
	case html.ElementNode:
		t.printIndent(level)
		// TODO: allow pre with color
		if n.DataAtom == atom.Pre && !printColor && preformatted {
			t.printPre(n)
			fmt.Println()
			return
		}
		if printColor {
			tokenColor.Print("<")
			tagColor.Printf("%s", n.Data)
		} else {
			fmt.Printf("<%s", n.Data)
		}
		for _, a := range n.Attr {
			val := a.Val
			if escapeHTML {
				val = html.EscapeString(val)
			}
			if printColor {
				fmt.Print(" ")
				attrKeyColor.Printf("%s", a.Key)
				tokenColor.Print("=")
				quoteColor.Printf(`"%s"`, val)
			} else {
				fmt.Printf(` %s="%s"`, a.Key, val)
			}
		}
		if printColor {
			tokenColor.Println(">")
		} else {
			fmt.Println(">")
		}
		if !isVoidElement(n) {
			t.printChildren(n, level+1)
			t.printIndent(level)
			if printColor {
				tokenColor.Print("</")
				tagColor.Printf("%s", n.Data)
				tokenColor.Println(">")
			} else {
				fmt.Printf("</%s>\n", n.Data)
			}
		}
	case html.CommentNode:
		t.printIndent(level)
		data := n.Data
		if escapeHTML {
			data = html.EscapeString(data)
		}
		if printColor {
			commentColor.Printf("<!--%s-->\n", data)
		} else {
			fmt.Printf("<!--%s-->\n", data)
		}
		t.printChildren(n, level)
	case html.DoctypeNode, html.DocumentNode:
		t.printChildren(n, level)
	}
}

func (t TreeDisplayer) printChildren(n *html.Node, level int) {
	if maxPrintLevel > -1 {
		if level >= maxPrintLevel {
			t.printIndent(level)
			fmt.Println("...")
			return
		}
	}
	child := n.FirstChild
	for child != nil {
		t.printNode(child, level)
		child = child.NextSibling
	}
}

func (t TreeDisplayer) printIndent(level int) {
	for ; level > 0; level-- {
		fmt.Print(indentString)
	}
}

// Print the text of a node
type TextDisplayer struct{}

func (t TextDisplayer) Display(nodes []*html.Node) {
	for _, node := range nodes {
		if node.Type == html.TextNode {
			data := node.Data
			if escapeHTML {
				// don't escape javascript
				if node.Parent == nil || node.Parent.DataAtom != atom.Script {
					data = html.EscapeString(data)
				}
			}
			fmt.Println(data)
		}
		children := []*html.Node{}
		child := node.FirstChild
		for child != nil {
			children = append(children, child)
			child = child.NextSibling
		}
		t.Display(children)
	}
}

// Is this node a tag with no end tag such as <meta> or <br>?
// http://www.w3.org/TR/html-markup/syntax.html#syntax-elements
func isVoidElement(n *html.Node) bool {
	switch n.DataAtom {
	case atom.Area, atom.Base, atom.Br, atom.Col, atom.Command, atom.Embed,
		atom.Hr, atom.Img, atom.Input, atom.Keygen, atom.Link,
		atom.Meta, atom.Param, atom.Source, atom.Track, atom.Wbr:
		return true
	}
	return false
}
