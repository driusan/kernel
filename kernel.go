// kernel represents the entry point to an operating system kernel written
// in (mostly) Go. It must be compiled and linked with gccgo and not the
// standard Go toolchain so that it can be linked in freestanding mode.
// Furthermore, it must be linked with gcc (the C) version to link in
// freestanding mode. gccgo and gcc (C) use the same (C) calling
// convention, so this isn't a big deal. However, it also means that the
// assembly portions of the code are written in GAS and not Plan9/Go style
// ASM so that gcc can compile them.
//
// All symbols expected by the Go runtime are not defined, so occasionally
// a language feature will try and use a symbol that can't be linked. When
// this happens putting the the gccgo frontend definition in libg/ and updating
// the makefile is (usually) enough.
//
// Most Go language features should be available. However, there's currently a
// bug where append() will go into an infinite loop if it needs to relocate the
// data. goroutines are also unlikely to have the symbols that they require
// defined, although I haven't tried.
//
// See the Makefile for details of how to compile/link.
package kernel

import (
	"github.com/driusan/kernel/acpi"
	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/descriptortables"
	"github.com/driusan/kernel/filesystem"
	"github.com/driusan/kernel/filesystem/fat"
	"github.com/driusan/kernel/ide"
	"github.com/driusan/kernel/input/ps2"
	"github.com/driusan/kernel/interrupts"
	"github.com/driusan/kernel/mbr"
	"github.com/driusan/kernel/memory"
	"github.com/driusan/kernel/pci"
	"github.com/driusan/kernel/shell"
	"github.com/driusan/kernel/terminal"
)

// Represents information passed along from multiboot compliant
// bootloader in the format specified by the multiboot spec.
type BootInfo struct {
	Flags,
	MemLower,
	MemUpper uint32
	BootDevice uint32
	Cmdline    uint32
	ModsCount  uint32
	ModsAddr   uint32
	ElfSec     ElfSectionHeaderTable
	MMapLength uint32
	MMapAddr   uint32
}

type ElfSectionHeaderTable struct {
	num   uint32
	size  uint32
	addr  uint32
	shndx uint32
}

// KernelMain represents the entry point of the kernel from a multiboot
// compliant bootloader. The bi parameter is a pointer to the boot information
// passed along from the bootloader.
func KernelMain(magic uint32, bi *BootInfo) {
	// First init the video, so that we can print debug messages.
	//term := terminal.Terminal{}
	//terminal.Term = &term
	terminal.InitializeTerminal()
	if magic != 0x2badb002 {
		println("Bad magic header. Not a multiboot bootloader")
		return
	}

	// Initialize packages with package level variables
	acpi.InitPkg()
	// if we don't declare these ahead of time gccgo complains about
	// goto skipping over their definition. They'll just be allocated
	// on the stack, so it's not a big deal to define them ahead of time.
	//var rsdt *acpi.RSDT
	var drive ide.IDEDrive
	var mbrdata ide.DriveSector
	var pts *mbr.Partitions
	var err error
	ptr, err := acpi.FindRSDP()
	if err != nil {
		println(err.Error())
		goto errExit
	}

	// the [:] hack causes a page fault until InitializePaging is done below,
	// so disable this for now.
	//println("Found ACPI Table at", ptr, " from OEM", string(ptr.OEMID[:]))
	_, err = ptr.GetRSDT()
	if err != nil {
		println(err.Error())
		goto errExit
	}

	// Not printing here for the same reason..
	// println("RSDT Signature:", string(rsdt.Signature[:]))
	// TODO: Initialize multiple CPUs based on the MADT table in ACPI.
	// There's not really much reason to do that until there's something
	// for the CPUs to do, though.
	// Should also probably try and enter long mode here.

	// Initialize paging does more than it claims. It reinitializes paging
	// from a higher level, and identity maps all the memory that the
	// bootloader told us about. It also initializes the structures used by
	// malloc and free. After this, we can allocate memory.
	memory.InitializePaging(uintptr(bi.MMapAddr), uintptr(bi.MMapLength))

	// Now that the heap is initialized, these packages "init" functions can
	// be run.
	filesystem.InitPkg()
	pci.InitPkg()
	ps2.InitPkg()
	ide.InitPkg()

	// Identify the by polling drive before interrupts are enabled.
	drive, err = ide.IdentifyDrive(ide.PrimaryDrive)
	if err != nil {
		println("Drive error:", err.Error())
	}
	// and enable the PS2 mouse, otherwise all packets will show delta x/y
	// of 0.
	ps2.EnableMouse()

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

	// runs an STI instruction to enable interrupts
	interrupts.Enable()

	// Now that everything is configured, print the memory.
	print(bi.MemLower, "kb of memory in lower memory.\n")
	print(bi.MemUpper, "kb of memory in upper memory.\n")
	print("Total ", (bi.MemLower+bi.MemUpper)/1024, "mb of memory.\n")
	//println("Flags", bi.Flags)
	//println("MMap Length", bi.MMapLength, " MMap Addr", bi.MMapAddr)
	//print(mmap.Memtype)
	print("PCI Devices on system: \n")
	pci.EnumerateDevices()

	// Read the partition table from the MBR.
	// TODO: Read it from the GPT if it exists.
	mbrdata, err = ide.ReadLBA(drive, 0)
	if err != nil {
		println("Drive error:", err.Error())
	}
	pts = mbr.ExtractPartitions(mbrdata.Data)

	for i, p := range pts {
		println("Partition", i, " active:", p.Active, " type", p.Type(), " LBA", p.LBAStart, " Size", p.LBASize)
		if p.Type() == "FAT32" {
			// This really needs to be redesigned, this is just a hack
			// to get something mounted without a mount command
			fat.Fat = &fat.FatFS{
				LBAStart: uint64(p.LBAStart),
				LBASize:  uint64(p.LBASize),
				Drive:    drive,
			}
		}
	}

	shell.Run()

	// Just sit around waiting for an interrupt now that everything
	// is enabled.
	for {
		asm.HLT()
	}

	// If there's an error, this will return back to boot.s, which will
	// disable interrupts and HLT in a loop.
errExit:
	println("Shutting down...")
}
