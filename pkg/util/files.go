package util

import (
	"io/fs"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
)

func CreateDirectoryWatcher(root string, recursive bool) (*fsnotify.Watcher, error) {
	log.Infof("Initializing directory watcher for: %s", root)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	defer watcher.Close()

	total := 0
	if recursive {
		for path := range IterSubdirectories(root) {
			err = watcher.Add(path)
			if err != nil {
				return nil, err
			}
			total++
		}
	} else {
		log.Debugf("Watching directory %s", root)
		err = watcher.Add(root)
		if err != nil {
			return nil, err
		}
		total++
	}
	log.Infof("Watching %d directories for changes", total)
	return watcher, nil
}

func IterSubdirectories(path string) chan string {
	ch := make(chan string)
	total := 0
	startTime := time.Now()

	go func() {
		defer close(ch)

		log.Debugf("Searching for subdirectories of %s", path)

		fsys := os.DirFS(path)
		fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				ch <- p
				total++
			}
			return nil
		})
		elapsed := time.Since(startTime)
		log.Debugf("Found %d subdirectories of %s in %.2f seconds", total, path, elapsed.Seconds())
	}()
	return ch
}
