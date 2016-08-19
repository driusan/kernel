package process

import "github.com/driusan/kernel/filesystem"

func NewNamespace() (ns Namespace) {
	ns = make(Namespace)
	ns["/dev/cons"] = filesystem.ConsoleFS
	return ns
}
