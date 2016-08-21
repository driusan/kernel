package filesystem

import "github.com/driusan/kernel/asm"
import "github.com/driusan/kernel/terminal"

// These should act the same as /dev/cons and /dev/consctl in Plan9.
var Cons DevCons

//var ConsCtl File
var DevFS devFS
var Root Filesystem

type FilesystemError string

func (f FilesystemError) Error() string {
	return string(f)
}

func InitPkg() {
	Cons = DevCons{}
	DevNull = Null{}
	DevFS = devFS{}
	Root = RootFS{}
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

func (f DevCons) WriteRune(r rune) error {
	terminal.PrintRune(r)
	return nil
}

// This whole ReadByte/SendByte thing needs to be re-thought out.
// It won't work with multiple readers.
var lastbyte byte

func (f DevCons) ReadByte() (b byte, err error) {
	for {
		if lastbyte != 0 {
			b = lastbyte
			lastbyte = 0
			return
		}
		asm.HLT()
	}
}
func (f DevCons) SendByte(b byte) error {
	lastbyte = b
	if f.openers != nil {
		for _, reader := range f.openers {
			reader.WriteByte(b)
		}
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

func (f DevCons) AsDirectory() (Directory, error) {
	return nil, FilesystemError("File is not a directory")
}

type devFS struct {
	// TODO: debug why struct{} results in an index out of range error
	// This is just a placeholder
	stub bool
}

/*
func (dfs devFS) Root() Directory {
	return SimpleDirectory{
		name:  "dev",
		files: nil, //[]File{Cons, ConsCtl, DevNull},
	}
}
*/
func (dfs devFS) Open(p Path) (File, error) {
	switch string(p) {
	case "", "/":
		return dfs, nil
	case "cons", "/cons":
		return Cons, nil
	case "null", "/null":
		return DevNull, nil
	}
	return nil, FilesystemError("No such file or directory")
}

func (dfs devFS) Type() string {
	return "DevFS"
}

func (dfs devFS) Name() string {
	return "dev"
}

func (dfs devFS) Files() []File {
	return []File{Cons, DevNull}
}

func (dfs devFS) Close() error {
	return nil
}

func (dfs devFS) IsDirectory() bool {
	return true
}

func (dfs devFS) AsDirectory() (Directory, error) {
	return dfs, nil
}

func (dfs devFS) Read(p []byte) (n int, err error) {
	return 0, FilesystemError("File is a directory.")
}

func (dfs devFS) Write(p []byte) (n int, err error) {
	return 0, FilesystemError("File is a directory.")
}

func (dfs devFS) Seek(offset int64, whence int) (int64, error) {
	return 0, FilesystemError("File is a directory.")
}

func (dfs devFS) ReadByte() (byte, error) {
	return 0, FilesystemError("File is a directory.")
}

func (dfs devFS) WriteByte(b byte) error {
	return FilesystemError("File is a directory.")
}
func (dfs devFS) WriteRune(r rune) error {
	return FilesystemError("File is a directory.")
}
