// package filesystem represents the abstractions shared across
// different filesystem implementations. The implementations of
// each filesystem type should usually be in a separate package
// for that filesystem.
package filesystem

import "io"

const EOF = FilesystemError("End of file")

// A Path is a special type of string denoting a file or path on a
// filesystem.
type Path string

// HasPrefix returns true iff op is a prefix of p.
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

// A File is something, generally on disk, which is read or written to.
// It's mostly a composition of various io interfaces.
type File interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
	io.ByteReader
	io.ByteWriter
	RuneWriter

	// IsDirectory() and AsDirectory() return true if this file is really
	// a directory. They should probably be deprecated in favour of
	// if d, ok := f.(Directory); ok { ... } now that enough of the Go
	// runtime is implemented to do so.
	IsDirectory() bool
	AsDirectory() (Directory, error)
}

// A Directory is a special type of file which contains other files. Its
// implementation is filesystem dependendent.
type Directory interface {
	File

	// Returns a map of files in the directory. Note that since it's a map,
	// the files are unordered when ranging through them, even if they're
	// ordered on disk.
	Files() map[string]File
}

// A filesystem represents an abstraction for accessing files (usually) on
// disk.
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
