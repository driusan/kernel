package filesystem

type RootFS struct {
	SimpleDirectory
}

func (r RootFS) Open(p Path) (File, error) {
	switch string(p) {
	case "/", "":
		return r, nil
	default:
		return nil, FilesystemError("No such file or directory")
	}
}
func (r RootFS) Type() string {
	return "Filesystem root"
}
