package monitor

import (
	"path/filepath"
)

type File struct {
	Path     string  `json:"path"`
	Filename string  `json:"filename"`
	Hashes   *Hashes `json:"hashes"`
}

func NewFile(path string) File {
	filename := filepath.Base(path)
	return File{
		Filename: filename,
		Path:     path,
	}
}

func GetFile(path string) (*File, error) {
	file := NewFile(path)
	hashes, err := GetFileHashes(file.Path)
	if err != nil {
		return nil, err
	}
	file.Hashes = hashes
	return &file, nil
}
