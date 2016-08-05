package kernel

import (
	"asm"
	"descriptortables"
	"interrupts"
	"memory"
	"pci"
)

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

	memory.InitializePaging()

	// Set up the GDT and interrupt handlers
	descriptortables.GDTInstall()
	descriptortables.IDTInstall()

	// Install handlers for Intel CPU exceptions
	interrupts.ISRSInstall()
	// and the PIC
	interrupts.IRQInstall()

	interrupts.InstallHandler(0, TimerHandler)
	interrupts.InstallHandler(1, KeyboardHandler)
	// runs an STI instruction to enable interrupts
	interrupts.Enable()

	// Now that everything is configured, print the memory.
	print(bi.MemLower, "kb of memory in lower memory.\n")
	print(bi.MemUpper, "kb of memory in upper memory.\n")
	print("Total ", (bi.MemLower+bi.MemUpper)/1024, "mb of memory.\n")

	pci.EnumerateDevices()

	// Just sit around waiting for an interrupt now that everything
	// is enabled.
	for {
		asm.HLT()
	}
}
