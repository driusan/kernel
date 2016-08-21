package filesystem

type Path string

func (p Path) HasPrefix(op Path) bool {
	return false
}

type File interface {
	Reader
	Writer
	Seeker
	Closer
	Name() string
	IsDirectory() bool
}

type Directory interface {
	Name() string
	Files() []File
}
type Filesystem interface {
	Root() Directory
}
