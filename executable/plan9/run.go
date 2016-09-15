package plan9

import (
	"io"
	//"unsafe"
	//"github.com/driusan/kernel/asm"
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
		//println("Got", n, " want", dataSize)
		return Plan9Error("Could not read data segment in one shot. TODO: Make this more robust")
	}
	if err != nil {
		return err
	}

	/*
	println("First byte", textSegment[0], "of", textSize)
	println("Entry point", h.EntryPoint.Uint32(), " at", uintptr(unsafe.Pointer(&textSegment[h.EntryPoint.Uint32()])))
	for i := 0; i < 10; i++ {
		println("i", i, textSegment[h.EntryPoint.Uint32()+uint32(i)])
	}
*/

	//asm.CALL(uintptr(unsafe.Pointer(&textSegment[h.EntryPoint.Uint32()])))
	return Plan9Error("Run Plan9 style a.out file not yet implemented")
}
