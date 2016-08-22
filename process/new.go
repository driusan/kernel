package process

import "github.com/driusan/kernel/filesystem"

func NewNamespace() Namespace {
	ns := make(Namespace)
	fsRoot := filesystem.RootFS{
		filesystem.SimpleDirectory{
			DirName:  "/",
			FilesMap: make(map[string]filesystem.File),
		},
	}
	ns["/"] = fsRoot

	fs, err := filesystem.DevFS.Open("")
	if err == nil {
		ns["/dev"] = filesystem.DevFS
		fsRoot.FilesMap["dev"] = fs
	}
	if filesystem.Fat != nil {
		fs, err = filesystem.DevFS.Open("")
		if err == nil {
			ns["/dos"] = filesystem.Fat
			fsRoot.FilesMap["dos"] = fs
		}
	}
	return ns
}

func New() Process {
	return Process{
		Namespace: NewNamespace(),
		Wd:        "/",
	}
}
