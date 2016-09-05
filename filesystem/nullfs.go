package filesystem

// DevNull represents /dev/null. Writes to it disappear, and reads
// return EOF.
var DevNull File

type Null struct {
	stub bool
}

func (f Null) Read(p []byte) (n int, err error) {
	return 0, EOF
}

func (f Null) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (f Null) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f Null) ReadByte() (byte, error) {
	return 0, EOF
}

func (f Null) WriteByte(b byte) error {
	return nil
}
func (f Null) WriteRune(r rune) error {
	return nil
}
func (f Null) Close() error {
	return nil
}

func (f Null) Name() string {
	return "null"
}

func (f Null) IsDirectory() bool {
	return false
}

func (f Null) AsDirectory() (Directory, error) {
	return nil, FilesystemError("Not a directory")
}
