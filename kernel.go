package kernel

import (
	"github.com/driusan/kernel/acpi"
	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/descriptortables"
	"github.com/driusan/kernel/ide"
	"github.com/driusan/kernel/input/ps2"
	"github.com/driusan/kernel/interrupts"
	"github.com/driusan/kernel/memory"
	"github.com/driusan/kernel/pci"
	"github.com/driusan/kernel/terminal"
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
	//term := terminal.Terminal{}
	//terminal.Term = &term
	terminal.InitializeTerminal()

	// Initialize packages with package level variables
	pci.InitPkg()
	acpi.InitPkg()
	ide.InitPkg()
	ps2.InitPkg()

	ptr, err := acpi.FindRSDP()

	// if we don't declare this ahead of time gccgo complains about
	// goto skipping over its definition
	var rsdt *acpi.RSDT
	var drive ide.IDEDrive
	if err != nil {
		println(err.Error())
		goto errExit
	}

	println("Found ACPI Table at", ptr, " from OEM", string(ptr.OEMID[:]))
	rsdt, err = ptr.GetRSDT()
	if err != nil {
		println(err.Error())
		goto errExit
	}

	println("RSDT Signature:", string(rsdt.Signature[:]))
	// TODO: Initialize multiple CPUs based on the MADT table in ACPI.
	// There's not really much reason to do that until there's something
	// for the CPUs to do, though.
	// Should also probably try and enter long mode here.

	// Identify the by polling drive before interrupts are enabled.
	drive, err = ide.IdentifyDrive(ide.PrimaryDrive)
	if err != nil {
		println("Drive error:", err.Error())
	}

	ps2.EnableMouse()
	memory.InitializePaging()
	// Set up the GDT and interrupt handlers
	descriptortables.GDTInstall()
	descriptortables.IDTInstall()

	// Install handlers for Intel CPU exceptions
	interrupts.ISRSInstall()
	// and the PIC
	interrupts.IRQInstall()

	interrupts.InstallHandler(0, TimerHandler)
	interrupts.InstallHandler(1, ps2.KeyboardHandler)
	interrupts.InstallHandler(12, ps2.MouseHandler)
	interrupts.InstallHandler(14, ide.PrimaryDriveHandler)

	interrupts.Enable()

	// runs an STI instruction to enable interrupts

	// Now that everything is configured, print the memory.
	print(bi.MemLower, "kb of memory in lower memory.\n")
	print(bi.MemUpper, "kb of memory in upper memory.\n")
	print("Total ", (bi.MemLower+bi.MemUpper)/1024, "mb of memory.\n")

	print("PCI Devices on system: \n")
	pci.EnumerateDevices()

	err = ide.ReadLBA(drive, 1)
	err = ide.ReadLBA(drive, 0)
	if err != nil {
		println("Drive error:", err.Error())
	}

	// Just sit around waiting for an interrupt now that everything
	// is enabled.
	for {
		asm.HLT()
	}

	// If there's an error, this will return back to boot.s, which will
	// disable interrupts and HLT in a loop.
errExit:
}
