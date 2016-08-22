package filesystem

var Fat Filesystem

type FatFS struct {
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
func (f FatFS) Type() string {
	return "FAT"
}

func (f FatFS) AsDirectory() (Directory, error) {
	return nil, FilesystemError("FAT filesystem not yet implemented")
}
func (f FatFS) Open(name Path) (File, error) {
	return nil, FilesystemError("FAT filesystem not yet implemented")
}
