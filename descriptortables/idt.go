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

func (i *IDTEntry) SetBase(base uint32) {
	i[0] = byte(base & 0xFF)
	i[1] = byte((base & 0xFF00) >> 8)

	i[6] = byte((base & 0xFF0000) >> 16)
	i[7] = byte((base & 0xFF000000) >> 24)

}

func (i *IDTEntry) SetFlags(flags byte) {
	i[5] = flags
}

func (i *IDTEntry) SetSelector(flags uint16) {
	i[2] = byte(flags & 0xFF)
	i[3] = byte((flags & 0xFF00) >> 8)
}

/* Declare an IDT of 256 entries. Although we will only use the
*  first 32 entries in this tutorial, the rest exists as a bit
*  of a trap. If any undefined IDT entry is hit, it normally
*  will cause an "Unhandled Interrupt" exception. Any descriptor
*  for which the 'presence' bit is cleared (0) will generate an
*  "Unhandled Interrupt" exception */
var IDT [256]IDTEntry
var IDTPtr DescriptorTablePointer

/* This exists in 'start.asm', and is used to load our IDT */
//extern idt_load
func IDTLoad()

/* Use this function to set an entry in the IDT. Alot simpler
*  than twiddling with the GDT ;) */
func IDTSetGate(num byte, base uint32, selector uint16, flags byte) {
	gate := &IDT[num]
	gate.SetBase(base)
	gate.SetFlags(flags)
	gate.SetSelector(selector)
	gate[3] = 0
}

/* Installs the IDT */
func IDTInstall() {
	//	terminal_writestring("In IDT Install");
	//	__go_print_int64((sizeof (struct idt_entry) * 256) -1);
	//	terminal_writestring(" ");
	//	__go_print_int64(sizeof (struct idt_entry));
	/* Sets the special IDT pointer up, just like in 'gdt.c' */
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
	IDTLoad()

}
