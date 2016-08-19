package filesystem

type Path string

type Filesystem interface {
	Read()
}
