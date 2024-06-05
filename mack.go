package mack

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/monochromegane/conflag"
	"github.com/monochromegane/go-home"
	"github.com/monochromegane/terminal"
)

const version = "0.1.0"

const (
	ExitCodeOK = iota
	ExitCodeError
)

var opts Option

type MarkdownAck struct {
	Out, Err io.Writer
}

func (p MarkdownAck) Run(args []string) int {
	InitFileTypes()
	parser := newOptionParser(&opts)

	conflag.LongHyphen = true
	conflag.BoolValue = false
	for _, c := range [...]string{
		filepath.Join(xdgConfigHomeDir(), "pt", "config.toml"),
		filepath.Join(home.Dir(), ".ptconfig.toml"),
		".ptconfig.toml",
	} {
		if args, err := conflag.ArgsFrom(c); err == nil {
			parser.ParseArgs(args)
		}
	}

	args, err := parser.ParseArgs(args)
	if err != nil {
		fmt.Fprintf(p.Err, "ERROR %s\n", err)

		if ferr, ok := err.(*flags.Error); ok && ferr.Type == flags.ErrHelp {
			return ExitCodeOK
		}
		return ExitCodeError
	}

	if opts.Version {
		fmt.Printf("mack version %s\n", version)
		return ExitCodeOK
	}

	if opts.FileTypeOption.ListTypes {
		for _, element := range known_languages {
			// place-holder
			fmt.Printf("  %s: %s\n", element.Name, strings.Join(element.Exts, ", "))
		}
		return ExitCodeOK
	}

	// override file types based on specialized query language
	if opts.SearchOption.HtmlQuery {
		opts.FileTypeOption.FileType = []string{"html"}
	}
	if opts.SearchOption.MarkdownQuery {
		opts.FileTypeOption.FileType = []string{"md"}
	}
	if opts.SearchOption.JsonQuery {
		opts.FileTypeOption.FileType = []string{"json"}
	}

	if len(opts.FileTypeOption.FileType) > 0 {
		// got filetype option, gather all extentions
		extentions := uniq_exts_from_file_types(opts.FileTypeOption.FileType)
		regex_str := regex_from_file_exts(extentions)

		if opts.SearchOption.FileNamesOnly {
			opts.SearchOption.EnableFilesWithRegexp = true
			opts.SearchOption.PatternFilesWithRegexp = regex_str
		} else {
			opts.SearchOption.FileSearchRegexp = regex_str
		}
	}

	if opts.FileTypeOption.KnownTypes {
		// build a map of *all* extensions
		regex_str := regex_from_file_exts(known_extentions)

		if opts.SearchOption.FileNamesOnly {
			opts.SearchOption.EnableFilesWithRegexp = true
			opts.SearchOption.PatternFilesWithRegexp = regex_str
		} else {
			opts.SearchOption.FileSearchRegexp = regex_str
		}
	}

	if len(args) == 0 && !(opts.SearchOption.EnableFilesWithRegexp) {
		fmt.Printf("No regular expression found.\n")
		parser.WriteHelp(p.Err)
		return ExitCodeError
	}

	if !terminal.IsTerminal(os.Stdout) {
		if !opts.OutputOption.ForceColor {
			opts.OutputOption.EnableColor = false
		}
		if !opts.OutputOption.ForceGroup {
			opts.OutputOption.EnableGroup = false
		}
	}

	if p.givenStdin() && p.noRootPathIn(args) {
		opts.SearchOption.SearchStream = true
	}

	if opts.SearchOption.EnableFilesWithRegexp {
		args = append([]string{""}, args...)
	}

	if opts.OutputOption.Count {
		opts.OutputOption.Before = 0
		opts.OutputOption.After = 0
		opts.OutputOption.Context = 0
	}

	// html match selected
	if opts.SearchOption.HtmlQuery {
		// interpret PATTERN as CSS selector
		// execute against html files

	}

	// normal path
	search := search{
		roots: p.rootsFrom(args),
		out:   p.Out,
	}
	if err = search.start(p.patternFrom(args)); err != nil {
		fmt.Fprintf(p.Err, "%s\n", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func (p MarkdownAck) patternFrom(args []string) string {
	return args[0]
}

func (p MarkdownAck) rootsFrom(args []string) []string {
	if len(args) > 1 {
		return args[1:]
	} else {
		return []string{"."}
	}
}

func (p MarkdownAck) givenStdin() bool {
	fi, err := os.Stdin.Stat()
	if runtime.GOOS == "windows" {
		if err == nil {
			return true
		}
	} else {
		if err != nil {
			return false
		}

		mode := fi.Mode()
		if (mode&os.ModeNamedPipe != 0) || mode.IsRegular() {
			return true
		}
	}
	return false
}

func (p MarkdownAck) noRootPathIn(args []string) bool {
	return len(args) == 1
}

func xdgConfigHomeDir() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(home.Dir(), ".config")
	}
	return xdgConfigHome
}
