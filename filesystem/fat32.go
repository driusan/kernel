package filesystem

var Fat32 Filesystem

type Fat32FS struct {
	LBAStart uint64
	LBASize  uint64
}

/*
func (f Fat32FS) Read([]byte) (int64, error) {
	return 0, FilesystemError("Not yet implemented")
}
func (f Fat32FS) Write([]byte) (int64, error) {
	return 0, FilesystemError("Not yet implemented")
}

func (f Fat32FS) Seek(offset int64, whence int) (int64, error) {
	return 0, FilesystemError("Not yet implemented")
}

func (f Fat32FS) Close() error {
	return nil
}
*/
func (f Fat32FS) Type() string {
	return "FAT32"
}

func (f Fat32FS) AsDirectory() (Directory, error) {
	return nil, FilesystemError("Not yet implemented")
}
func (f Fat32FS) Open(name Path) (File, error) {
	return nil, FilesystemError("Not yet implemented")
}
