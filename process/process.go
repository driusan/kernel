package process

import "github.com/driusan/kernel/filesystem"

type Process struct {
	Namespace
	Wd filesystem.Path
}

func (p *Process) Cwd(path filesystem.Path) error {
	if len(path) == 0 || path[0] != '/' {
		return ProcessError("Relative paths not yet supported.")
	}

	f, err := p.Open(path)
	if err != nil {
		return err
	}

	if !f.IsDirectory() {
		return ProcessError("Path is not a directory")
	}

	p.Wd = path
	return nil
}
