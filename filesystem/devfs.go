package filesystem

// These should act the same as /dev/cons and /dev/consctl in Plan9.
var Cons DevCons
var ConsCtl File
var DevNull File
var DevFS Filesystem

func InitPkg() {
	Cons = DevCons{}
	DevNull = Null{}
	DevFS = devFS{}
}

type Null struct{}

func (f Null) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (f Null) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (f Null) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
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

// A consReader is something that has /dev/cons open.
// Bytes get writen to it.
// TODO: Handle raw vs line mode instead of assuming raw
type consReader struct {
	ByteWriter
	Raw bool
}

type DevCons struct {
	openers []consReader
}

func (f *DevCons) Open(callback ByteWriter) (consReader, error) {
	cr := consReader{ByteWriter: callback, Raw: true}
	f.openers = append(f.openers, cr)
	return cr, nil
}
func (f DevCons) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (f DevCons) Write(p []byte) (n int, err error) {
	print(string(p))
	return 0, nil
}

func (f DevCons) WriteByte(b byte) error {
	print(b)
	return nil
}
func (f DevCons) SendByte(b byte) error {
	for _, reader := range f.openers {
		reader.WriteByte(b)
	}
	return nil
}
func (f DevCons) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f DevCons) Close() error {
	return nil
}
func (f DevCons) Name() string {
	return "cons"
}

func (f DevCons) IsDirectory() bool {
	return false
}

type SimpleDirectory struct {
	name  string
	files []File
}

func (sd SimpleDirectory) Name() string {
	return sd.name
}

func (sd SimpleDirectory) Files() []File {
	return sd.files
}

type devFS struct{}

func (dfs devFS) Root() Directory {
	return SimpleDirectory{
		name:  "dev",
		files: []File{Cons, ConsCtl, DevNull},
	}
}

/*
type Directory interface {
	Name() string
	Files() []File
}
type Filesystem interface {
	Root() Directory
}
*/
