package find

import (
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/ignore"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Finder struct {
	Out    chan *grep.Params
	Option *option.Option
}

func (self *Finder) Find(root, pattern string) {
	Walk(root, self.Option.Ignore, func(path string, info os.FileInfo, ig ignore.Ignore, onset bool, err error) (error, ignore.Ignore) {
		if info.IsDir() {
			ig.Patterns = append(ig.Patterns, ignore.IgnorePatterns(path, self.Option.VcsIgnores())...)
			// fmt.Printf("pattern -> %s = %s\n", path, ig.Patterns)
			for _, p := range ig.Patterns {
				files, _ := filepath.Glob(path + string(os.PathSeparator) + p)
				if files != nil {
					// fmt.Printf("matches -> %s = %s\n", path+string(os.PathSeparator)+p, files)
					ig.Matches = append(ig.Matches, files...)
				}
			}
			if !onset && isHidden(info.Name()) {
				return filepath.SkipDir, ig
			} else if contains(path, &ig.Matches) {
				// fmt.Printf("ignore  -> %s\n", path)
				return filepath.SkipDir, ig
			} else {
				return nil, ig
			}
		}

		if !onset && isHidden(info.Name()) {
			return nil, ig
		}
		if contains(path, &ig.Matches) {
			// fmt.Printf("ignore  -> %s\n", path)
			return nil, ig
		}
		fileType := file.IdentifyType(path)
		if fileType == file.BINARY {
			return nil, ig
		}
		self.Out <- &grep.Params{path, pattern, fileType}
		return nil, ig
	})
	close(self.Out)
}

type WalkFunc func(path string, info os.FileInfo, ig ignore.Ignore, onset bool, err error) (error, ignore.Ignore)

func Walk(root string, ignorePatterns []string, walkFn WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		walkError, _ := walkFn(root, nil, ignore.Ignore{}, true, err)
		return walkError
	}
	return walk(root, info, false, ignore.Ignore{Patterns: ignorePatterns}, true, walkFn)
}

func walkOnGoRoutine(path string, info os.FileInfo, notify chan int, went bool, parentIgnore ignore.Ignore, walkFn WalkFunc) {
	walk(path, info, went, parentIgnore, false, walkFn)
	notify <- 0
}

func walk(path string, info os.FileInfo, went bool, parentIgnore ignore.Ignore, onset bool, walkFn WalkFunc) error {
	err, ig := walkFn(path, info, parentIgnore, onset, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	list, err := ioutil.ReadDir(path)
	if err != nil {
		walkError, _ := walkFn(path, info, ig, onset, err)
		return walkError
	}

	notify := make(chan int, len(list))
	for _, fileInfo := range list {
		if !went {
			go walkOnGoRoutine(filepath.Join(path, fileInfo.Name()), fileInfo, notify, !went, ig, walkFn)

		} else {
			walk(filepath.Join(path, fileInfo.Name()), fileInfo, went, ig, false, walkFn)
		}
	}
	if !went {
		for i := 0; i < cap(notify); i++ {
			<-notify
		}
	}
	return nil
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".") && len(name) > 1
}

func contains(path string, patterns *[]string) bool {
	for _, p := range *patterns {
		if p == path {
			return true
		}
	}
	return false
}