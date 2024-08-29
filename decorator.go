package mack

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	ColorReset = "\x1b[0m\x1b[K"

	SeparatorColon  = ":"
	SeparatorHyphen = "-"
)

type decorator interface {
	path(path string) string
	lineNumber(lineNum int) string
	columnNumber(columnNum int) string
	match(line string, matched bool) string
}

func newDecorator(pattern regex_pattern, option Option) decorator {
	if option.OutputOption.EnableColor {
		return newColor(pattern, option)
	} else {
		return plain{}
	}
}

type term_color struct {
	from   string
	to     string
	regexp *regexp.Regexp

	colorLineNumber string
	colorPath       string
	colorMatch      string
}

func newColor(pattern regex_pattern, option Option) term_color {
	color := term_color{
		colorLineNumber: ansiEscape(option.OutputOption.ColorCodeLineNumber),
		colorPath:       ansiEscape(option.OutputOption.ColorCodePath),
		colorMatch:      ansiEscape(option.OutputOption.ColorCodeMatch),
	}
	if pattern.regexp == nil {
		p := string(pattern.pattern)
		color.from = p
		color.to = color.colorMatch + p + ColorReset
	} else {
		color.to = color.colorMatch + "${1}" + ColorReset
		color.regexp = pattern.regexp
	}
	return color
}

func ansiEscape(code string) string {
	re := regexp.MustCompile("[^0-9;]")
	sanitized := re.ReplaceAllString(code, "")
	if sanitized == "" {
		sanitized = "0" // all attributes off
	}
	return "\x1b[" + sanitized + "m"
}

func (c term_color) path(path string) string {
	return c.colorPath + path + ColorReset
}

func (c term_color) lineNumber(lineNum int) string {
	return c.colorLineNumber + strconv.Itoa(lineNum) + ColorReset
}

func (c term_color) columnNumber(columnNum int) string {
	return strconv.Itoa(columnNum)
}

func (c term_color) match(line string, matched bool) string {
	if !matched {
		return line
	} else if c.regexp == nil {
		return strings.Replace(line, c.from, c.to, -1)
	} else {
		return c.regexp.ReplaceAllString(line, c.to)
	}
}

type plain struct {
}

func (p plain) path(path string) string {
	return path
}

func (p plain) lineNumber(lineNum int) string {
	return strconv.Itoa(lineNum)
}

func (p plain) columnNumber(columnNum int) string {
	return strconv.Itoa(columnNum)
}

func (p plain) match(line string, matched bool) string {
	return line
}
