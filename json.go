package mack

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/antchfx/jsonquery"
	"github.com/itchyny/gojq"
	"golang.org/x/net/html"
)

type jsonQuery struct {
	pattern *gojq.Query
	printer printer
}

func NewJsonQuery(pattern pattern, printer printer) jsonQuery {
	// parse the pattern
	query, err := gojq.Parse(string(pattern.pattern))
	if err != nil {
		log.Fatalln(err)
	}
	return jsonQuery{
		pattern: query,
		printer: printer,
	}
}

func (g jsonQuery) grep(path string, buf []byte) {
	fileSpec, ok := GetLanguageSpec(path)
	if ok {
		switch fileSpec.Name {
		case "json":
			doc, err := loadJsonFile(path)
			if err != nil {
				log.Fatalf("Unable to load file: %s, %s\n", path, err)
			}
			list := g.pattern.Run(doc)
			g.printNode(path, list)
		default:
			log.Println("unknown file extention, skipping: ", filepath.Ext(path), " ", path)
		}
	} else {
		log.Fatalf("Unknow file type %s\n", path)
	}
}

func (g jsonQuery) printNode(path string, iter gojq.Iter) {
	match := match{path: path}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			log.Fatalln(err)
		}
		match.add(0, 0, v, true)
	}
	g.printer.print(match)
}

func loadJsonFile(path string) (map[string]interface{}, error) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}
	defer f.Close()

	var js map[string]interface{}

	// TODO Use the provided buffer in grep?
	file_buf, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("error reading file '%s': %s\n", newOutputOption().ColorCodePath, err)
		return js, err
	}

	json.Unmarshal(file_buf, &js)

	// parse the new html doc
	return js, nil
}

func RenderJson(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "") // fixme settings?
	if err != nil {
		return "", err
	}
	return string(val), nil
}

type json_printer struct {
	encoder          *json.Encoder
	w                io.Writer
	enableLineNumber bool
}

func NewJsonPrinter(w io.Writer, opts *OutputOption) json_printer {
	enc := json.NewEncoder(w)
	enc.SetIndent("", opts.JsonIndent)
	enc.SetEscapeHTML(false)
	return json_printer{encoder: enc, w: w, enableLineNumber: opts.EnableLineNumber}
}

func (f json_printer) renderNode(data any) interface{} {
	switch v := data.(type) {
	case nil:
		return "nil"
	case *jsonquery.Node:
		return v.Value()
	case *html.Node:
		return nodeStr(v)
	default:
		fmt.Printf("type unknown %T, blindly string-ifying: %v", v, v)
		return fmt.Sprintf("%v", v)
	}
}

func (f json_printer) print(match match) {
	for _, line := range match.lines {
		data := f.renderNode(line.text)
		if f.enableLineNumber {
			// wrap the match in a json struct that includes
			// match.path,
			data := map[string]interface{}{
				"path":    match.path,
				"value":   data,
				"matched": line.matched,
			}
			if line.num > 0 {
				data["line"] = line.num
			}
			if line.column > 0 {
				data["column"] = line.column
			}
		}

		// if enableLineNumber, print out non-matches as there is metadata to reflect the non-match
		// without that metadata, just print the data
		if line.matched || f.enableLineNumber {
			if err := f.encoder.Encode(data); err != nil {
				// invalid json, assume raw string and just print it, in quotes
				log.Fatalln("Error encoding data:", err, data)
				//fmt.Fprintln(f.w, data)
			}
		}
		// line didn't match
	}
}
