package process

import (
	"github.com/driusan/kernel/filesystem"
)

type Namespace map[filesystem.Path]filesystem.Filesystem

// Given an absolute filesystem.Path, returns the filesystem that should serve
// that path, the filesystem.Path relative to that filesystem, and an error
// which is hopefully nil.
func (ns Namespace) Translate(file filesystem.Path) (filesystem.Filesystem, filesystem.Path, error) {
	var bestFit filesystem.Filesystem
	var bestPath filesystem.Path
	for path, fs := range ns {
		// TODO: Figure out why there are empty paths in the iteration
		if path != "" && file.HasPrefix(path) {
			if len(path) > len(bestPath) {
				bestFit = fs
				bestPath = path
			}
		}
	}
	if bestFit == nil {
		return nil, "", ProcessError("No namespaces map to path")
	}
	return bestFit, file[len(bestPath):], nil
}

// Returns a file handler as evaluated in this namespace.
func (ns Namespace) Open(file filesystem.Path) (filesystem.File, error) {
	if file[0] != '/' {
		return nil, ProcessError("Relative paths not supported")
	}
	fs, path, err := ns.Translate(file)
	if err != nil {
		return nil, err
	}
	return fs.Open(path)

}

type ProcessError string

func (pe ProcessError) Error() string {
	return string(pe)
}
