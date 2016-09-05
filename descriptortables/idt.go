package descriptortables

import "unsafe"

type IDTPointer DescriptorTablePointer

// IdtEntrys have the structure
//	struct idt_entry
//	{
//		unsigned short base_lo;
//		unsigned short selector;
//		unsigned char always0;
//		unsigned char flags;
//		unsigned short base_hi;
//	} __attribute__((packed));
// but Go doesn't have any way to specify the packed attribute, so instead
// we use a byte array and have helper functions which manually put all of the
// bytes in the correct place by calling SetX receiver functions.
type IDTEntry [8]byte

// Sets the base address of this IDT Entry.
func (i *IDTEntry) SetBase(base uint32) {
	i[0] = byte(base & 0xFF)
	i[1] = byte((base & 0xFF00) >> 8)

	i[6] = byte((base & 0xFF0000) >> 16)
	i[7] = byte((base & 0xFF000000) >> 24)

}

// Set the flags for this IDT Entry
func (i *IDTEntry) SetFlags(flags byte) {
	i[5] = flags
}

// Set the selector for this IDT Entry
func (i *IDTEntry) SetSelector(flags uint16) {
	i[2] = byte(flags & 0xFF)
	i[3] = byte((flags & 0xFF00) >> 8)
}

// Declare an IDT of 256 entries. The first 32 handle Intel CPU exceptions, the
// next 32 handle interrupts from the PIC. The rest are currently unused, but
// can be used for software interrupts in the future.
var IDT [256]IDTEntry

// An IDT pointer is a DescriptorTablePointer of the same format as a GDTPtr.
var IDTPtr DescriptorTablePointer

// This is defined in asm. It calls LIDT to load our descriptor table.
//
// TODO: Move LIDT instruction to ASM package.
//
//extern idt_load
func IDTLoad()

// Sets an IDT gate.
func IDTSetGate(num byte, base uint32, selector uint16, flags byte) {
	gate := &IDT[num]
	gate.SetBase(base)
	gate.SetFlags(flags)
	gate.SetSelector(selector)
	gate[3] = 0
}

// Installs an empty IDT table to an area of memory to later get configured
// by the main Kernel.
func IDTInstall() {
	// Load the pointer in the format LIDT requires.
	p := &IDTPtr
	p.SetSize(8 /* sizeof IDTEntry */ *256 /* len(IDT) */ - 1)
	p.SetBase(uintptr(unsafe.Pointer(&IDT)))

	// Clear out the entire IDT, initializing it to zeros. I'm not sure if
	// this is required. The Go spec says things get initialized to zero,
	// but we haven't necessarily implemented enough of libg in order to be
	// sure that everything in the spec is usable in the kernel.
	for i := 0; i < 256; i++ {
		for j, _ := range IDT[i] {
			IDT[i][j] = 0
		}
	}

	// Load the IDT.
	IDTLoad()
}
