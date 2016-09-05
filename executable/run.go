// package executable handles the interpreting of executables, loading
// them from the filesystem and executing them.
//
// It's currently a stub.
package executable

import "io"

type ExeError string

func (e ExeError) Error() string {
	return string(e)
}

func Run(r io.Reader) error {
	return ExeError("Run not yet implemented")
}
