package process

import (
	"github.com/driusan/kernel/filesystem"
	"github.com/driusan/kernel/filesystem/fat"
)

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

	// This is a horrible design
	if fat.Fat != nil {
		fs, err = filesystem.DevFS.Open("")
		if err == nil {
			ns["/dos"] = fat.Fat
			fsRoot.FilesMap["dos"] = fs
			if err := fat.Fat.Initialize(); err != nil {
				print(err.Error())
			}
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
