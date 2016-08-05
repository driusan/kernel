package interrupts

import "asm"

// Enables interrupts
func Enable() {
	asm.STI()
}

// Maps an interrupt to a handler for that interrupt
// TODO: This doesn't work because __go_new_map and
//	__go_map_index symbols aren't defined.
//var inthandlers map[int]func(*Registers)

var inthandlers [16]func(*Registers)

//extern irq_install
func irq_install()

func IRQInstall() {
	irq_install()
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
		asm.OUTB(0xa0, 0x20)
	} else {
		//print("IRQ0-14")
	}
	asm.OUTB(0x20, 0x20)
}
