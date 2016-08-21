package interrupts

import "github.com/driusan/kernel/asm"

//extern isrs_install
func isrsInstallC()

func ISRSInstall() {
	// for some reason the extern symbol doesn't get linked if this function
	// is empty, so just make a stub that calls the extern
	isrsInstallC()
}

// Stores the state of CPU registers
type Registers struct {
	gs, fs, es, ds                         uint32 /* pushed the segs last */
	edi, esi, ebp, esp, ebx, edx, ecx, eax uint32 /* pushed by 'pusha' */
	InterruptNo, err_code                  uint32 /* our 'push byte #' and ecodes do this */
	eip, cs, eflags, useresp, ss           uint32 /* pushed by the processor automatically */
}

func (r Registers) InterruptDescription() string {
	if r.InterruptNo >= 32 {
		return "Unknown interrupt type."
	}
	switch r.InterruptNo {
	case 0:
		return "Division By Zero"
	case 1:
		return "Debug"
	case 2:
		return "Non Maskable Interrupt"
	case 3:
		return "Breakpoint"
	case 4:
		return "Into Detected Overflow"
	case 5:
		return "Out of Bounds"
	case 6:
		return "Invalid Opcode"
	case 7:
		return "No Coprocessor"
	case 8:
		return "Double Fault"
	case 9:
		return "Coprocessor Segment Overrun"
	case 10:
		return "Bad TSS"
	case 11:
		return "Segment Not Present"
	case 12:
		return "Stack Fault"
	case 13:
		return "General Protection Fault"
	case 14:
		return "Page Fault"
	case 15:
		return "Unknown Interrupt"
	case 16:
		return "Coprocessor Fault"
	case 17:
		return "Alignment Check"
	case 18:
		return "Machine Check"
	default:
		return "Reserved"
	}
}

// All ISRs point to this function in boot.s. When it's called, interrupts
// have been disabled.
func CPUFaultHandler(r *Registers) {
	print("In system fault handler")
	if r.InterruptNo < 32 {
		print(r.InterruptDescription(), " Exception. System Halted!\n")
		for {
			asm.CLI()
			asm.HLT()
		}
	} else {
		print("Unknown CPU Fault, ", r.InterruptNo)
		panic("CPUFaultHandler called with invalid interrupt number")
	}
}
