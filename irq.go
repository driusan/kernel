package kernel

//extern inportb
func inportb(port uint16) byte

//extern outportb
func outportb(port uint16, data byte)

// Maps an interrupt to a handler for that interrupt
// TODO: This doesn't work because __go_new_map and
//	__go_map_index symbols aren't defined.
//var inthandlers map[int]func(*Registers)

var inthandlers [16]func(*Registers)

func IRQInstallGo() {
	inthandlers = [16]func(*Registers){}
	//make(map[int]func(*Registers))
}

func InstallHandler(port int, handler func(*Registers)) {
	inthandlers[port] = handler
}

func IRQHandler(r *Registers) {

	irq := r.InterruptNo - 32

	if handler := inthandlers[irq]; handler != nil {
		handler(r)
	} else {
		println("No handler for IRQ", irq)
	}
	// Acknowledge the interrupt to the PIC so that the next
	// one will get sent.
	if r.InterruptNo >= 40 {
		//print("IRQ8-15")
		outportb(0xa0, 0x20)
	} else {
		//print("IRQ0-14")
	}
	outportb(0x20, 0x20)
}
