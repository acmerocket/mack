package mack

import (
	"github.com/jessevdk/go-flags"
)

// Top level options
type Option struct {
	Version        bool            `long:"version" description:"Show version"`
	OutputOption   *OutputOption   `group:"Output Options"`
	SearchOption   *SearchOption   `group:"Search Options"`
	FileTypeOption *FileTypeOption `group:"File Type Options"`
}

// Output options.
type OutputOption struct {
	Color               func()       `long:"color" description:"Print color codes in results (default: true)"`
	NoColor             func()       `long:"nocolor" description:"Don't print color codes in results (default: false)"`
	ForceColor          bool         // Force color. Not user option.
	EnableColor         bool         // Enable color. Not user option.
	ColorLineNumber     func(string) `long:"color-line-number" description:"Color codes for line numbers (default: 1;33)"`
	ColorPath           func(string) `long:"color-path" description:"Color codes for path names (default: 1;32)"`
	ColorMatch          func(string) `long:"color-match" description:"Color codes for result matches (default: 30;43)"`
	ColorCodeLineNumber string       // Color line numbers. Not user option.
	ColorCodePath       string       // Color path names. Not user option.
	ColorCodeMatch      string       // Color result matches. Not user option.
	Group               func()       `long:"group" description:"Print file name at header (default: true)"`
	NoGroup             func()       `long:"nogroup" description:"Don't print file name at header (default: false)"`
	ForceGroup          bool         // Force group. Not user option.
	EnableGroup         bool         // Enable group. Not user option.
	Null                bool         `short:"0" long:"null" description:"Separate filenames with null (for 'xargs -0') (default: false)"`
	Column              bool         `long:"column" description:"Print column (default: false)"`
	LineNumber          func()       `long:"numbers" description:"Print Line number. (default: true)"`
	NoLineNumber        func()       `short:"N" long:"nonumbers" description:"Omit Line number. (default: false)"`
	ForceLineNumber     bool         // Force line number. Not user option.
	EnableLineNumber    bool         // Enable line number. Not user option.
	After               int          `short:"A" long:"after" description:"Print lines after match"`
	Before              int          `short:"B" long:"before" description:"Print lines before match"`
	Context             int          `short:"C" long:"context" description:"Print lines before and after match"`
	FilesWithMatches    bool         `short:"l" long:"files-with-matches" description:"Only print filenames that contain matches"`
	Count               bool         `short:"c" long:"count" description:"Only print the number of matching lines for each input file."`
	OutputEncode        string       `short:"o" long:"output-encode" description:"Specify output encoding (none, jis, sjis, euc)"`
	OutputJson          bool         `long:"json" description:"Output results as JSON"`
	JsonIndent          string       `long:"indent" description:"Indent for JSON ouput"`
}

func newOutputOption() *OutputOption {
	opt := &OutputOption{}

	opt.Color = opt.SetEnableColor
	opt.NoColor = opt.SetDisableColor
	opt.EnableColor = true

	opt.Group = opt.SetEnableGroup
	opt.NoGroup = opt.SetDisableGroup
	opt.EnableGroup = true

	opt.LineNumber = opt.SetEnableLineNumber
	opt.NoLineNumber = opt.SetDisableLineNumber
	opt.EnableLineNumber = true

	opt.ColorLineNumber = opt.SetColorLineNumber
	opt.ColorPath = opt.SetColorPath
	opt.ColorMatch = opt.SetColorMatch
	opt.ColorCodeLineNumber = "1;33" // yellow with black background
	opt.ColorCodePath = "1;32"       // bold green
	opt.ColorCodeMatch = "30;43"     // black with yellow background

	return opt
}

func (o *OutputOption) SetEnableColor() {
	o.ForceColor = true
	o.EnableColor = true
}

func (o *OutputOption) SetDisableColor() {
	o.EnableColor = false
}

func (o *OutputOption) SetEnableLineNumber() {
	o.ForceLineNumber = true
	o.EnableLineNumber = true
}

func (o *OutputOption) SetDisableLineNumber() {
	o.EnableLineNumber = false
}

func (o *OutputOption) SetEnableGroup() {
	o.ForceGroup = true
	o.EnableGroup = true
}

func (o *OutputOption) SetDisableGroup() {
	o.EnableGroup = false
}

func (o *OutputOption) SetColorLineNumber(code string) {
	o.ColorCodeLineNumber = code
}

func (o *OutputOption) SetColorPath(code string) {
	o.ColorCodePath = code
}

func (o *OutputOption) SetColorMatch(code string) {
	o.ColorCodeMatch = code
}

// Search options.
type SearchOption struct {
	Regexp                 bool         `short:"e" description:"Parse PATTERN as a regular expression (default: false). Accepted syntax is the same as https://github.com/google/re2/wiki/Syntax except from \\C"`
	IgnoreCase             bool         `short:"i" long:"ignore-case" description:"Match case insensitively"`
	SmartCase              bool         `short:"S" long:"smart-case" description:"Match case insensitively unless PATTERN contains uppercase characters"`
	WordRegexp             bool         `short:"w" long:"word-regexp" description:"Only match whole words"`
	Ignore                 []string     `long:"ignore" description:"Ignore files/directories matching pattern"`
	VcsIgnore              []string     `long:"vcs-ignore" description:"VCS ignore files" default:".gitignore"`
	GlobalGitIgnore        bool         `long:"global-gitignore" description:"Use git's global gitignore file for ignore patterns"`
	HomePtIgnore           bool         `long:"home-ptignore" description:"Use $Home/.ptignore file for ignore patterns"`
	SkipVcsIgnore          bool         `short:"U" long:"skip-vcs-ignores" description:"Don't use VCS ignore file for ignore patterns"`
	FileNamesOnly          func()       `short:"f" description:"Only print the files selected, without searching."`
	FilesWithRegexp        func(string) `short:"g" description:"Print filenames matching PATTERN"`
	EnableFilesWithRegexp  bool         // Enable files with regexp. Not user option.
	PatternFilesWithRegexp string       // Pattern files with regexp. Not user option.
	FileSearchRegexp       string       `short:"G" long:"file-search-regexp" description:"PATTERN Limit search to filenames matching PATTERN"`
	Depth                  int          `long:"depth" default:"25" description:"Search up to NUM directories deep"`
	Follow                 bool         `long:"follow" description:"Follow symlinks"`
	Hidden                 bool         `long:"hidden" description:"Search hidden files and directories"`
	SearchStream           bool         // Input from pipe. Not user option.
	CssSelect              bool         `long:"css" description:"Parse PATTERN as a CSS selection against HTML and Markdown files"`
	XpathQuery             bool         `long:"xpath" description:"Parse PATTERN as an Xpath expression"`
	JsonQuery              bool         `long:"jq" description:"Parse PATTERN as a JSON query against JSON files"`
}

func (o *SearchOption) SetFilesWithRegexp(p string) {
	o.EnableFilesWithRegexp = true
	o.PatternFilesWithRegexp = p
}

func (o *SearchOption) SetFileNamesOnly() {
	o.EnableFilesWithRegexp = true
	o.PatternFilesWithRegexp = ""
}

func newSearchOption() *SearchOption {
	opt := &SearchOption{}
	opt.FilesWithRegexp = opt.SetFilesWithRegexp
	opt.FileNamesOnly = opt.SetFileNamesOnly
	return opt
}

type FileTypeOption struct {
	//-t X, --type=X  Include only X files, where X is a filetype, e.g. python, html, markdown, etc
	FileType []string `short:"t" long:"type" description:"Include only X files, where X is a filetype, e.g. python, html, markdown, etc"`

	//-T X, --type=noX Exclude X files, where X is a filetype.
	//ExcludeType []string `short:"T" long:"not-type" description:"Exclude X files, where X is a filetype."`
	// SKIPPING for now

	//-k, --known-types Include only files of types that ack recognizes.
	KnownTypes bool `short:"k" long:"known-types" description:"Include only files of types that mack recognizes."`

	//--help-types Display all known types, and how they're defined.
	ListTypes bool `long:"help-types" description:"Display all known types, and how they're defined."`

	// remove ExtSet map[string]bool // private, to store all the interesting file name regexes.
}

func newFileTypeOption() *FileTypeOption {
	opt := &FileTypeOption{}
	return opt
}

func newOptionParser(opts *Option) *flags.Parser {
	opts.OutputOption = newOutputOption()
	opts.SearchOption = newSearchOption()
	opts.FileTypeOption = newFileTypeOption()

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = "mack"
	parser.Usage = "[OPTIONS] PATTERN [PATH]"
	return parser
}
