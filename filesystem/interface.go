package filesystem

import "io"

type Path string

const EOF = FilesystemError("End of file")

func (p Path) HasPrefix(op Path) bool {
	if len(p) < len(op) {
		return false
	}

	for i, _ := range op {
		if p[i] != op[i] {
			return false
		}
	}
	return true
}

type File interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
	io.ByteReader
	io.ByteWriter
	RuneWriter
	//Name() string
	IsDirectory() bool
	AsDirectory() (Directory, error)
}

type Directory interface {
	File
	Files() map[string]File
}
type Filesystem interface {
	// Should initialize any internal data structures needed by the
	// filesystem being implemented before mounting
	Initialize() error

	// Opens a file relative to this filesystem. Open should generally not
	// cross filesystem boundaries.
	Open(Path) (File, error)

	// Returns a string identifying this filesystem handler
	Type() string
}

// RuneWriter is strangely missing from the io package, so we just define one
// here
type RuneWriter interface {
	WriteRune(r rune) error
}
