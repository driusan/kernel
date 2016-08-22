package process

import "github.com/driusan/kernel/filesystem"

func NewNamespace() Namespace {
	ns := make(Namespace)
	ns["/"] = filesystem.Root
	ns["/dev"] = filesystem.DevFS
	if filesystem.Fat32 != nil {
		ns["/dos"] = filesystem.Fat32
	}
	return ns
}

func New() Process {
	return Process{
		Namespace: NewNamespace(),
		Wd:        "/",
	}
}
