package filesystem

// There isn't enough of libg implemented to import "io", so we just redefine
// the same interfaces here. See standard Go docs for how to use these
// interfaces.

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Seeker interface {
	Seek(offset int64, whence int) (int64, error)
}

type Closer interface {
	Close() error
}

type ByteWriter interface {
	WriteByte(byte) error
}
