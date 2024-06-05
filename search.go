package mack

import (
	"io"
	"regexp"
)

type search struct {
	roots []string
	out   io.Writer
}

func (s search) start(pattern string) error {
	searchChan := make(chan string, 5000)
	done := make(chan struct{})

	if opts.SearchOption.WordRegexp {
		opts.SearchOption.Regexp = true
		pattern = "\\b" + pattern + "\\b"
	}

	if opts.SearchOption.SmartCase {
		if !regexp.MustCompile(`[[:upper:]]`).MatchString(pattern) {
			opts.SearchOption.IgnoreCase = true
		}
	}

	if opts.SearchOption.IgnoreCase {
		opts.SearchOption.Regexp = true
	}

	if opts.OutputOption.Context > 0 {
		opts.OutputOption.Before = opts.OutputOption.Context
		opts.OutputOption.After = opts.OutputOption.Context
	}

	var regFile *regexp.Regexp
	var err error
	if opts.SearchOption.FileSearchRegexp != "" {
		regFile, err = regexp.Compile(opts.SearchOption.FileSearchRegexp)
		if err != nil {
			return err
		}
	}
	if opts.SearchOption.EnableFilesWithRegexp {
		opts.OutputOption.FilesWithMatches = true
		if opts.SearchOption.PatternFilesWithRegexp != "" {
			regFile, err = regexp.Compile(opts.SearchOption.PatternFilesWithRegexp)
			if err != nil {
				return err
			}
		}
	}

	go find{
		out:  searchChan,
		opts: opts,
	}.start(s.roots, regFile)

	// original
	p, err := newPattern(pattern, opts)
	if err != nil {
		return err
	}

	go newGrep(
		p,
		searchChan,
		done,
		opts,
		newPrinter(p, s.out, opts),
	).start()

	<-done

	return nil
}
