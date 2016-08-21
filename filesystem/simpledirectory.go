package filesystem

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

func (sd SimpleDirectory) Close() error {
	return nil
}

func (sd SimpleDirectory) IsDirectory() bool {
	return true
}

func (sd SimpleDirectory) AsDirectory() (Directory, error) {
	return sd, nil
}

func (f SimpleDirectory) Read(p []byte) (n int, err error) {
	return 0, FilesystemError("File is a directory.")
}

func (f SimpleDirectory) Write(p []byte) (n int, err error) {
	return 0, FilesystemError("File is a directory.")
}

func (f SimpleDirectory) Seek(offset int64, whence int) (int64, error) {
	return 0, FilesystemError("File is a directory.")
}

func (f SimpleDirectory) ReadByte() (byte, error) {
	return 0, FilesystemError("File is a directory.")
}

func (f SimpleDirectory) WriteByte(b byte) error {
	return FilesystemError("File is a directory.")
}
func (f SimpleDirectory) WriteRune(r rune) error {
	return FilesystemError("File is a directory.")
}
