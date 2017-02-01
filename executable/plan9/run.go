package plan9

import (
	"io"
	"unsafe"

	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/memory"
	"github.com/driusan/kernel/process"
)

type Plan9Error string

func (p Plan9Error) Error() string {
	return string(p)
}

// The currently running process, for syscalls to know what they're dealing with.
// This should be a function of the CPU, not the package, but there's no CPU
// type (or multiprocessing) yet.
var activeProc *process.Process

func Run(h *ExecHeader, r io.Reader, p *process.Process) error {
	if h.Magic.Uint32() != Magic386 {
		return Plan9Error("Invalid executable magic.")
	}

	// Use the Reader to load the data. This is unuseable,
	// because of the overhead from the slice changing the
	// exact memory layout that we care about.
	textSize := h.Text.Uint32()
	textSegment := make([]byte, textSize)
	n, err := r.Read(textSegment)
	if n != int(textSize) {
		return Plan9Error("Could not read text segment in one shot. TODO: Make this more robust")
	}
	if err != nil {
		return err
	}

	// Load the data segment. Same caveat.
	dataSize := h.Data.Uint32()
	dataSegment := make([]byte, dataSize)
	n, err = r.Read(dataSegment)
	if n != int(dataSize) {
		return Plan9Error("Could not read data segment in one shot. TODO: Make this more robust")
	}
	if err != nil {
		return err
	}

	textAddr, err := memory.GetPhysicalAddress(unsafe.Pointer(&textSegment[0]))
	if err != nil {
		return nil
	}
	dAddr, err := memory.GetPhysicalAddress(unsafe.Pointer(&dataSegment[0]))
	if err != nil {
		return nil
	}

	err = memory.LoadMap(
		memory.MMapEntry{
			uintptr(unsafe.Pointer(h)),
			0,
			32,
		},
		memory.MMapEntry{
			uintptr(textAddr),
			0,
			uint(textSize)},
		memory.MMapEntry{
			uintptr(dAddr),
			0,
			uint(dataSize)},
	)
	if err != nil {
		return err
	}

	activeProc = p
	asm.JMP(unsafe.Pointer(uintptr(h.EntryPoint.Uint32())))

	// This should never be reached. The program exits with an interrupt
	return Plan9Error("Run Plan9 style a.out file not yet implemented")
}
