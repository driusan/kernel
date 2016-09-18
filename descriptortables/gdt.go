// package descriptortables handles the parsing and loading of IDT/GDT
// descriptor tables in Go.
//
// The main reason it's in a separate package is to avoid cyclical imports.
package descriptortables

import "unsafe"

// GDTEntry denotes a GDT entry. The packing is important, and
// Go doesn't have any way to define __attribute__((packed)), so
// it's defined as a byte array with Set() methods to manually place the
// bytes into the endian-correct packed location.
type GDTEntry [8]byte

// Sets the limit for this GDT entry. Limit is actually a 20 bit integer.
// TODO: This should return an error type if the parameter is invalid.
func (e *GDTEntry) SetLimit(limit uint32) {
	e[0] = byte((limit & 0x0FF00) >> 8)
	e[1] = byte((limit & 0x000FF))

	// The limit is packed in to the lower nibble of this byte, the
	// higher nibble is the flags, which need to be preserved.
	e[6] = byte((limit&0x0F0000)>>16) | (e[6] & 0xF0)
}

// Sets the base of this GDTEntry.
func (e *GDTEntry) SetBase(base uint32) {
	// The encoding required by x86 separates the high and low bytes
	// Base_low portion
	e[2] = byte((base & 0x000000FF))
	e[3] = byte((base & 0x0000FF00) >> 8)
	e[4] = byte((base & 0x00FF0000) >> 16)

	// The last byte of the base encoding is separated by the access byte,
	// limit and flags
	e[7] = byte((base >> 24) & 0xFF)
}

// This is actually a nibble, not a byte. It should probably be broken up into
// sub Set* helpers for each flag. For now, the bits should be set in the
// higher nibble of the byte in this parameter.
// This should also check the parameter and return an error when appropriate.
func (e *GDTEntry) SetFlags(b byte) {
	e[6] = b&0xF0 | e[6]&0x0F
	// TODO: The flags should have consts defined to reference them.
}

// Sets the access byte for this GDT entry.
func (e *GDTEntry) SetAccess(b byte) {
	e[5] = b
}

// Sets the granularity for this GDT entry.
func (e *GDTEntry) SetGranularity(b byte) {
	e[6] = b
}

// A DescriptorTablePointer is a pointer to either an IDT or a GDT encoded
// in a way that the LGDT or LIDT instructions are valid.
type DescriptorTablePointer [6]byte

// Set the size of the descriptor table
func (e *DescriptorTablePointer) SetSize(l uint16) {
	e[1] = byte((l & 0xFF00) >> 8)
	e[0] = byte(l & 0x00FF)
}

// Set the base address of the descriptor table.
func (e *DescriptorTablePointer) SetBase(l uintptr) {
	e[5] = byte((l & 0xFF000000) >> 24)
	e[4] = byte((l & 0x00FF0000) >> 16)
	e[3] = byte((l & 0x0000FF00) >> 8)
	e[2] = byte((l & 0x000000FF))
}

// Our GDT. The one currently loaded by our kernel only has 3 entries.
var Gdt [3]GDTEntry

// The pointer to the GDT used by our kernel.
var GDTPtr DescriptorTablePointer

//struct gdt_ptr gp;

// This is defined in assembly. It will load the GDT pointed to by GDTPtr.
//
//extern gdt_flush
func GDTFlush()

// Setup a descriptor in the Global Descriptor Table.
func GDTSetGate(num int, base, limit uint32, access, gran byte) {
	gate := &Gdt[num]
	gate.SetBase(base)
	gate.SetLimit(limit)
	gate.SetAccess(access)
	gate.SetGranularity(gran)
}

// GDTInstall is called by the kernel to setup and install the GDT.
func GDTInstall() {
	// Set up the GDTPtr
	p := &GDTPtr
	p.SetSize((8 /* sizeof gdt_entry */ * 3 /* num entries */) - 1)
	p.SetBase(uintptr(unsafe.Pointer(&Gdt)))

	// Setup the null descriptor
	GDTSetGate(0, 0, 0, 0, 0)

	// Set up the code segment to span the entire memory.
	// We're not as cautious as we should be with protecting things
	// with the GDT, but this allows us to setup the IDT.
	//
	// TODO: This should use fewer magic values. This has a 4GB limit,
	// 4KByte granularity, and 32-bit opcodes.
	GDTSetGate(1, 0, 0xFFFFFFFF, 0x9A, 0xCF)

	// Setup the data segment. It's exactly the same as the code segment,
	// but the access byte says it's data.
	GDTSetGate(2, 0, 0xFFFFFFFF, 0x92, 0xCF)

	// Install the new GDT.
	GDTFlush()
}
