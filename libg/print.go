package libg

import (
	_ "C"

	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/memory"
	"github.com/driusan/kernel/terminal"
)

func halt() {
	for {
		asm.CLI()
		asm.HLT()
	}

}
func GoPanic() {
	print("Kernel panic. TODO: add more interesting info here.")
	halt()
}

func GoPrintString(s string) {
	for i := 0; i < len(s); i++ {
		terminal.Term.PutChar(byte(s[i]))
	}
}

func GoPrintInt64(i int64) {
	terminal.PrintDec(i)
}

func GoPrintUint64(i uint64) {
	// TODO: Implement this properly.
	terminal.PrintDec(int64(i))
}

func GoPrintNewline() {
	terminal.Term.PutChar('\n')
}

func GoPrintSpace() {
	terminal.Term.PutChar(' ')
}

func GoPrintPointer(p uintptr) {
	terminal.PrintHex(uint64(p))
}

func GoRuntimePanicString(err string) {
	GoPrintString(err)
	halt()
}

func GoPrintBool(val bool) {
	if val {
		GoPrintString("true")
	} else {
		GoPrintString("false")
	}
}

func GoAlloc(size uint) uintptr {
	// This shouldn't really be here. For now it's just a stub
	// so that it compiles.
	ptr, err := memory.Malloc(size)

	if err != nil {
		print(err.Error())
		return 0
	}
	return ptr
}
func GoFree(uintptr) {
	// This shouldn't really be here. For now it's just a stub
	// so that it compiles.
}
