package process

import "github.com/driusan/kernel/filesystem"

type Namespace map[filesystem.Path]filesystem.Filesystem

func (ns Namespace) Test() string {
	return "Test"
}
