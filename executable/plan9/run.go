package plan9

import (
	"io"
	"unsafe"

	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/memory"
)

type Plan9Error string

func (p Plan9Error) Error() string {
	return string(p)
}

func Run(h *ExecHeader, r io.Reader) error {
	textSize := h.Text.Uint32()
	textSegment := make([]byte, textSize)
	n, err := r.Read(textSegment)
	if n != int(textSize) {
		return Plan9Error("Could not read text segment in one shot. TODO: Make this more robust")
	}
	if err != nil {
		return err
	}

	dataSize := h.Data.Uint32()
	dataSegment := make([]byte, dataSize)
	n, err = r.Read(dataSegment)
	if n != int(dataSize) {
		return Plan9Error("Could not read data segment in one shot. TODO: Make this more robust")
	}
	if err != nil {
		return err
	}

	println("First byte", textSegment[0], "of", textSize, " at ", &textSegment[0])
	println("Entry point", h.EntryPoint.Uint32())
	println("Text segment (+10)")
	for i := 0; i < 10; i++ {
		print("i", i, ":", textSegment[i], " ")
		print("i", i, ":", textSegment[h.EntryPoint.Uint32()+uint32(i)], " ")
	}

	print(&textSegment[0])
	println("\nAddress By unsafe.Pointer: T0-10")
	start := uintptr(unsafe.Pointer(&textSegment[0]))
	for i := 0; i < 10; i++ {
		print("i", i, ":", *(*byte)(unsafe.Pointer(start + uintptr(i))), " ")

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
			textAddr,
			0,
			uint(textSize)},
		memory.MMapEntry{
			dAddr,
			0,
			uint(dataSize)},
	)
	if err != nil {
		return err
	}
	println("\nAddress 0-10")
	for i := 0; i < 10; i++ {
		print("i", i, ":", *(*byte)(unsafe.Pointer(uintptr(i))), " ")
		print("i", i, ":", *(*byte)(unsafe.Pointer(uintptr(h.EntryPoint.Uint32() + uint32(i)))), " ")

	}
	asm.JMP(unsafe.Pointer(uintptr(h.EntryPoint.Uint32())))
	return Plan9Error("Run Plan9 style a.out file not yet implemented")
}
