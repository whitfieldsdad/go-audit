package monitor

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/whitfieldsdad/go-audit/pkg/util"
)

type DirectoryEventSource struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"`
}

func NewDirectoryEventSource(path string, recursive bool) (*DirectoryEventSource, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}
	s := &DirectoryEventSource{Path: path, Recursive: recursive}
	return s, nil
}

// TODO
func (s DirectoryEventSource) Watch() error {
	watcher, err := util.CreateDirectoryWatcher(s.Path, s.Recursive)
	if err != nil {
		return err
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			path := event.Name
			log.Infof("Received file event: %s", path)
		case err, ok := <-watcher.Errors:
			if !ok {
				return err
			}
		}
	}
}
