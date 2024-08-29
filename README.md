# mack

mack (working title) - `ack` with markdown support

**`tl;dr find + grep + xpath + jq + css`**

`mack` provides a simple CLI for a three-stage parallel process to:
1. Select files from the file system, based on name or type
2. Query those files with regex, xpath, css or jq
3. Format those results: text, color-terminal, json


## Features

### File Selection
List files with `-f`. By default, `mack` searches recursively starting at the current directory.
```
mack -f
mack -f ./files
```

List files restricted by regex using `-g`:
```
mack -g .md
mack -g .html ./files
```

Select file by type with `-t`. `-l` to *list* files instead of applying a query:
```
mack -t markdown -l "##"
mack -t json -l "description" ./files
```

Limit file types based on all known file types with `-k` or `--known-types`.
```
mack -k "##"
mack --known-types "# "
```

To list all know file types:
```
mack --help-types
```

### Regex Matching
When not limited by `-t` or `-k`, `mack` will default to matching all files recursively. By default, `mack` uses a simple grep-like marching algorithm.
```
mack "##"
```

Using `-e` (same as the original `grep`):
```
mack -e "^# "
```

REGEX chect sheet: https://devhints.io/regexp, https://quickref.me/regex.html

### XPath
Search for `<table>`s in HTML and XML files.
```
mack --xpath //table
```

XPath cheat cheet: https://devhints.io/xpath, https://quickref.me/xpath

### CSS Selector
Match all the `<h2>` tags in HTML files:
```
mack -t html --css h2
```

CSS Selector cheat sheet: https://devhints.io/css, https://quickref.me/css3#css-selectors

### Json Query
Match `description` fields in JSON files:
```
mack -t json --jq .description
```

jq cheat sheet: https://www.devtoolsdaily.com/cheatsheets/jq/

### Speed

- It searches code about 3–5× faster than `ack`.
- It searches code as fast as `the_silver_searcher(ag)`.
- It ignores file patterns from your `.gitignore`.
- It ignores directories with names that start with `.`, eg `.config`. Use `--hidden` option, if you want to search.
- It searches `UTF-8`, `EUC-JP` and `Shift_JIS` files.
- It provides binaries for multi platform (macOS, Windows, Linux).

#### Benchmarks

```sh
cd ~/src/github.com/torvalds/linux
ack EXPORT_SYMBOL_GPL 30.18s user 2.32s system  99% cpu 32.613 total # ack
ag  EXPORT_SYMBOL_GPL  1.57s user 1.76s system 311% cpu  1.069 total # ag: It's faster than ack.
pt  EXPORT_SYMBOL_GPL  2.29s user 1.26s system 358% cpu  0.991 total # pt: It's faster than ag!!
```

## Usage
```
Usage:
  mack [OPTIONS] PATTERN [PATH]

Application Options:
      --version             Show version

Output Options:
      --color               Print color codes in results (default: true)
      --nocolor             Don't print color codes in results (default: false)
      --color-line-number=  Color codes for line numbers (default: 1;33)
      --color-path=         Color codes for path names (default: 1;32)
      --color-match=        Color codes for result matches (default: 30;43)
      --group               Print file name at header (default: true)
      --nogroup             Don't print file name at header (default: false)
  -0, --null                Separate filenames with null (for 'xargs -0') (default: false)
      --column              Print column (default: false)
      --numbers             Print Line number. (default: true)
  -N, --nonumbers           Omit Line number. (default: false)
  -A, --after=              Print lines after match
  -B, --before=             Print lines before match
  -C, --context=            Print lines before and after match
  -l, --files-with-matches  Only print filenames that contain matches
  -c, --count               Only print the number of matching lines for each input file.
  -o, --output-encode=      Specify output encoding (none, jis, sjis, euc)
      --json                Output results as JSON
      --indent=             Indent for JSON ouput

Search Options:
  -e                        Parse PATTERN as a regular expression (default: false). Accepted syntax is the same as
                            https://github.com/google/re2/wiki/Syntax except from \C
  -i, --ignore-case         Match case insensitively
  -S, --smart-case          Match case insensitively unless PATTERN contains uppercase characters
  -w, --word-regexp         Only match whole words
      --ignore=             Ignore files/directories matching pattern
      --vcs-ignore=         VCS ignore files (default: .gitignore)
      --global-gitignore    Use git's global gitignore file for ignore patterns
      --home-ptignore       Use $Home/.ptignore file for ignore patterns
  -U, --skip-vcs-ignores    Don't use VCS ignore file for ignore patterns
  -f                        Only print the files selected, without searching. The PATTERN must not be specified.
  -g=                       Print filenames matching PATTERN
  -G, --file-search-regexp= PATTERN Limit search to filenames matching PATTERN
      --depth=              Search up to NUM directories deep (default: 25)
      --follow              Follow symlinks
      --hidden              Search hidden files and directories
      --css                 Parse PATTERN as a CSS selection against HTML and Markdown files
      --xpath               Parse PATTERN as an Xpath expression
      --jq                  Parse PATTERN as a JSON query against JSON files

File Type Options:
  -t, --type=               Include only X files, where X is a filetype, e.g. python, html, markdown, etc
  -k, --known-types         Include only files of types that mack recognizes.
      --help-types          Display all known types, and how they're defined.

Help Options:
  -h, --help                Show this help message

ERROR Usage:
  mack [OPTIONS] PATTERN [PATH]

Application Options:
      --version             Show version

Output Options:
      --color               Print color codes in results (default: true)
      --nocolor             Don't print color codes in results (default: false)
      --color-line-number=  Color codes for line numbers (default: 1;33)
      --color-path=         Color codes for path names (default: 1;32)
      --color-match=        Color codes for result matches (default: 30;43)
      --group               Print file name at header (default: true)
      --nogroup             Don't print file name at header (default: false)
  -0, --null                Separate filenames with null (for 'xargs -0') (default: false)
      --column              Print column (default: false)
      --numbers             Print Line number. (default: true)
  -N, --nonumbers           Omit Line number. (default: false)
  -A, --after=              Print lines after match
  -B, --before=             Print lines before match
  -C, --context=            Print lines before and after match
  -l, --files-with-matches  Only print filenames that contain matches
  -c, --count               Only print the number of matching lines for each input file.
  -o, --output-encode=      Specify output encoding (none, jis, sjis, euc)
      --json                Output results as JSON
      --indent=             Indent for JSON ouput

Search Options:
  -e                        Parse PATTERN as a regular expression (default: false). Accepted syntax is the same as
                            https://github.com/google/re2/wiki/Syntax except from \C
  -i, --ignore-case         Match case insensitively
  -S, --smart-case          Match case insensitively unless PATTERN contains uppercase characters
  -w, --word-regexp         Only match whole words
      --ignore=             Ignore files/directories matching pattern
      --vcs-ignore=         VCS ignore files (default: .gitignore)
      --global-gitignore    Use git's global gitignore file for ignore patterns
      --home-ptignore       Use $Home/.ptignore file for ignore patterns
  -U, --skip-vcs-ignores    Don't use VCS ignore file for ignore patterns
  -f                        Only print the files selected, without searching. The PATTERN must not be specified.
  -g=                       Print filenames matching PATTERN
  -G, --file-search-regexp= PATTERN Limit search to filenames matching PATTERN
      --depth=              Search up to NUM directories deep (default: 25)
      --follow              Follow symlinks
      --hidden              Search hidden files and directories
      --css                 Parse PATTERN as a CSS selection against HTML and Markdown files
      --xpath               Parse PATTERN as an Xpath expression
      --jq                  Parse PATTERN as a JSON query against JSON files

File Type Options:
  -t, --type=               Include only X files, where X is a filetype, e.g. python, html, markdown, etc
  -k, --known-types         Include only files of types that mack recognizes.
      --help-types          Display all known types, and how they're defined.

Help Options:
  -h, --help                Show this help message
```

## Configuration

If you put configuration file on the following directories, pt use option in the file.

- $XDG\_CONFIG\_HOME/pt/config.toml
- $HOME/.ptconfig.toml
- .ptconfig.toml (current directory)

The file is TOML format like the following.

```toml
color = true
context = 3
ignore = ["dir1", "dir2"]
color-path = "1;34"
```

The options are same as command line options.

## Editor Integration

### Vim + Unite.vim

You can use pt with [Unite.vim](https://github.com/Shougo/unite.vim).

```vim
nnoremap <silent> ,g :<C-u>Unite grep:. -buffer-name=search-buffer<CR>
if executable('pt')
  let g:unite_source_grep_command = 'pt'
  let g:unite_source_grep_default_opts = '--nogroup --nocolor'
  let g:unite_source_grep_recursive_opt = ''
  let g:unite_source_grep_encoding = 'utf-8'
endif
```

### Emacs + pt.el

You can use pt with [pt.el](https://github.com/bling/pt.el), which can be installed from [MELPA](http://melpa.milkbox.net/).

## Installation

### Developer

```sh
$ go get -u github.com/acmerocket/mack/...
```

### User

Download from the following url.

- [https://github.com/acmerocket/mack/releases](https://github.com/acmerocket/mack/releases)

Or, you can use Homebrew (Only macOS).

```sh
$ brew install pt
```

`pt` is an alias for `mack` in Homebrew.

## Building

FIXME

## Contribution

1. Fork it
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

## License

[MIT](https://github.com/acmerocket/mack/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)
[acmerocket](https://github.com/acmerocket)
