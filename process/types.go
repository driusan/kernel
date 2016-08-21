package process

import (
	"github.com/driusan/kernel/filesystem"
)

type Namespace map[filesystem.Path]filesystem.Filesystem

// Given a filesystem.Path, returns the filesystem that should serve that path,
// the filesystem.Path relative to that filesystem, and an error which is
// hopefully nil.
func (ns Namespace) Translate(file filesystem.Path) (filesystem.Filesystem, filesystem.Path, error) {
	var bestFit filesystem.Filesystem
	var bestPath filesystem.Path
	for path, fs := range ns {
		if file.HasPrefix(path) {
			bestFit = fs
			bestPath = path
		}
	}
	bestFit = nil
	if bestFit == nil {
		return nil, "", ProcessError("No such file or directory")
	}
	return bestFit, bestPath, nil
}

func (ns Namespace) Test() string {
	return "Test"
}

// Returns a file handler as evaluated in this namespace.
func (ns Namespace) Open(file filesystem.Path) (filesystem.File, error) {
	return nil, ProcessError("Could not open file")
}

type ProcessError string

func (pe ProcessError) Error() string {
	return string(pe)
}
