package kernel

import (
	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/interrupts"
)

func KeyboardHandler(r *interrupts.Registers) {
	scancode := asm.INB(0x60)

	if scancode&0x80 != 0 {
		println("Released key")
		// key has been released
	} else {
		println("Pressed key")
		// key has been pressed
	}
}
