package kernel

//extern initialize_paging
func InitializePaging()

//extern irq_install
func IrqInstall()

//extern enable_interrupts
func EnableInterrupts()

// Represents information passed along from multiboot compliant
// bootloader
type BootInfo struct {
	Flags,
	MemLower,
	MemUpper uint32
}

func KernelMain(bi *BootInfo) {
	// First init the video, so that we can print debug messages.
	InitializeTerminal()

	InitializePaging()

	// Set up the GDT and interrupt handlers
	GDTInstall()
	IDTInstall()

	// Install handlers for Intel CPU exceptions
	ISRSInstall()
	// and the PIC
	IrqInstall()
	IRQInstallGo()

	InstallHandler(1, KeyboardHandler)
	// runs an STI instruction to enable interrupts
	EnableInterrupts()

	// Now that everything is configured, print the memory.
	print(bi.MemLower, "kb of memory in lower memory.\n")
	print(bi.MemUpper, "kb of memory in upper memory.\n")
	print("Total ", (bi.MemLower+bi.MemUpper)/1024, "mb of memory.\n")
	println(bi.Flags)
	// Just sit around waiting for an interrupt now that everything
	// is enabled.
	for {
		//print("In loop")
		// TODO: After interrupts are working try putting a HLT
		// instruction here to avoid waiting CPU resources.
	}
}
