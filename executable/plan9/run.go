package plan9

import (
	"io"
	"unsafe"

	//"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/memory"
)

type Plan9Error string

func (p Plan9Error) Error() string {
	return string(p)
}

func exec()

func Run(h *ExecHeader, r io.Reader) error {
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

	/*
		println("First byte", textSegment[0], "of", textSize, " at ", &textSegment[0])
		println("Entry point", h.EntryPoint.Uint32())
		println("Text segment (+10)")
		for i := 0; i < 10; i++ {
			//print("i", i, ":", textSegment[i], " ")
			//print("i", i, ":", textSegment[h.EntryPoint.Uint32()+uint32(i)], " ")
		}

		//print(&textSegment[0])
		//println("\nAddress By unsafe.Pointer: T0-10")
		start := uintptr(unsafe.Pointer(&textSegment[0]))
		for i := 0; i < 10; i++ {
			print("i", i, ":", *(*byte)(unsafe.Pointer(start + uintptr(i))), " ")

		}
	*/
	textAddr, err := memory.GetPhysicalAddress(unsafe.Pointer(&textSegment[0]))
	if err != nil {
		return nil
	}
	dAddr, err := memory.GetPhysicalAddress(unsafe.Pointer(&dataSegment[0]))
	if err != nil {
		return nil
	}

	// Now unsafely allocate the right amount without the
	// slice overhead, and memmove the data into it.
	/*
		unsafeTextAddr, err := memory.Malloc(uint(textSize))
		if err != nil {
			return err
		}
		unsafeDataAddr, err := memory.Malloc(uint(dataSize))
		if err != nil {
			return err
		}
		memory.Move(unsafe.Pointer(unsafeTextAddr), unsafe.Pointer(&textSegment[0]), int(textSize))
		memory.Move(unsafe.Pointer(unsafeDataAddr), unsafe.Pointer(&dataSegment[0]), int(dataSize))
	*/
	err = memory.LoadMap(
		memory.MMapEntry{
			uintptr(textAddr),
			0,
			uint(textSize)},
		memory.MMapEntry{0, 0, 1},
		memory.MMapEntry{
			uintptr(dAddr),
			0,
			uint(dataSize)},
	)
	if err != nil {
		return err
	}
	println("Text size", textSize, " Data Size", dataSize)
	println("Entry point? (Ignoring, using 0x20)", h.EntryPoint.Uint32())
	println("Byte at start of write: ", textSegment[0x168c])

	exec()
	//asm.JMP(unsafe.Pointer(uintptr(0x20)))//h.EntryPoint.Uint32())))

	// This should never be reached. The program exits with an interrupt
	return Plan9Error("Run Plan9 style a.out file not yet implemented")
}
