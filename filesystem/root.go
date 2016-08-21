package filesystem

type RootFS struct {
	stub bool
}

func (r RootFS) Open(p Path) (File, error) {
	switch string(p) {
	case "/", "":
		return SimpleDirectory{
			name:  "/",
			files: []File{File(DevFS)},
		}, nil
	default:
		return nil, FilesystemError("No such file or directory")
	}
}
func (r RootFS) Type() string {
	return "Filesystem root"
}
