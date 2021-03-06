// package executable handles the interpreting of executables, loading
// them from the filesystem and executing them.
//
// It's currently a stub.
package executable

import (
	"io"
	"unsafe"

	"github.com/driusan/kernel/executable/plan9"
	"github.com/driusan/kernel/process"
)

type ExeError string

func (e ExeError) Error() string {
	return string(e)
}

func Run(r io.Reader, p *process.Process) error {
	header := make([]byte, 32)
	n, err := r.Read(header)
	if n != 32 {
		return ExeError("Could not read program header.")
	}
	if err != nil {
		return err
	}
	// Making r a ReadSeeker is causing a kernel panic, disabling for now.
	//r.Seek(0, 0)
	return plan9.Run((*plan9.ExecHeader)(unsafe.Pointer(&header[0])), r, p)
}
