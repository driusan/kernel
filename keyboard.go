package kernel

func KeyboardHandler(r *Registers) {
	scancode := inportb(0x60)

	if scancode&0x80 != 0 {
		println("Released key")
		// key has been released
	} else {
		println("Pressed key")
		// key has been pressed
	}
}
