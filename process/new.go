package process

import "github.com/driusan/kernel/filesystem"

func NewNamespace() Namespace {
	ns := make(Namespace)
	ns["/dev"] = filesystem.DevFS
	return ns
}
