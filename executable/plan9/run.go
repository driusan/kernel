package plan9

import (
	"io"
	"unsafe"

	//	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/memory"
	"github.com/driusan/kernel/process"
)

type Plan9Error string

func (p Plan9Error) Error() string {
	return string(p)
}

func execAddr(start, stack memory.VirtualAddress)

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
	textSegment := make([]byte, textSize+32)

	n, err := r.Read(textSegment[32:])
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

	textAddr, err := memory.GetPhysicalAddress(memory.VirtualAddress(unsafe.Pointer(&textSegment[0])))
	if err != nil {
		return nil
	}
	dAddr, err := memory.GetPhysicalAddress(memory.VirtualAddress(unsafe.Pointer(&dataSegment[0])))
	if err != nil {
		return nil
	}

	// Reserve a 1 page stack
	stack := make([]byte, 4096)
	stackAddr, err := memory.GetPhysicalAddress(memory.VirtualAddress(unsafe.Pointer(&stack[0])))
	if err != nil {
		return nil
	}

	stackMmap := &memory.MMapEntry{
		stackAddr,
		0,
		memory.PageSize,
	}
	err = memory.LoadMap(
		&memory.MMapEntry{
			memory.PhysicalAddress(unsafe.Pointer(h)),
			0,
			32,
		},
		&memory.MMapEntry{
			textAddr,
			0,
			uint(textSize)},
		&memory.MMapEntry{
			dAddr,
			0,
			uint(dataSize),
		},
		stackMmap,
	)

	println("The new stack is at:", stackAddr)
	activeProc = p
	execAddr(memory.VirtualAddress(h.EntryPoint.Uint32()), stackMmap.VAddr)

	// This should never be reached. The program exits with an interrupt
	return Plan9Error("Run Plan9 style a.out file not yet implemented")
}
