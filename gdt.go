package kernel

import "unsafe"

/* Loads a GDT entry. I haven't figured out how to do the
 *__attribute((packed))__ in gccgo, so for now this is in C.
 *
 * Copied from Brendon's tutorial at
 * http://www.osdever.net/bkerndev/Docs/gdt.htm
 * Mostly so that we can get on to interupts.
 */

/* Defines a GDT entry. We say packed, because it prevents the
*  compiler from doing things that it thinks is best: Prevent
*  compiler "optimization" by packing */
type GdtEntry [8]byte

// Limit is actually a 20 bit integer.
// TODO: This should return an error type if the parameter is invalid.
func (e *GdtEntry) SetLimit(limit uint32) {
	e[0] = byte((limit & 0x0FF00) >> 8)
	e[1] = byte((limit & 0x000FF))

	// The limit is packed in to the lower nibble of this byte, the
	// higher nibble is the flags, which need to be preserved.
	e[6] = byte((limit&0x0F0000)>>16) | (e[6] & 0xF0)
}

func (e *GdtEntry) SetBase(base uint32) {
	/* This was:
	   gdt[num].base_low = (base & 0xFFFF);
	   gdt[num].base_middle = (base >> 16) & 0xFF;
	   gdt[num].base_high = (base >> 24) & 0xFF;
	*/

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
func (e *GdtEntry) SetFlags(b byte) {
	e[6] = b&0xF0 | e[6]&0x0F
	// TODO: The flags should have consts defined to reference them.
}
func (e *GdtEntry) SetAccess(b byte) {
	e[5] = b
}

func (e *GdtEntry) SetGranularity(b byte) {
	e[6] = b// e[6] = b&0x0F | e[6]&0xF0
}

/*struct gdt_entry
{
    unsigned short limit_low; 2
    unsigned short base_low; 2
    unsigned char base_middle; 1
    unsigned char access; 1
    unsigned char granularity; 1
    unsigned char base_high; 1
} __attribute__((packed));
*/

type GdtPtr [6]byte

func (e *GdtPtr) SetSize(l uint16) {
	e[1] = byte((l & 0xFF00) >> 8)
	e[0] = byte(l & 0x00FF)
}

func (e *GdtPtr) SetBase(l uintptr) {
	e[5] = byte((l & 0xFF000000) >> 24)
	e[4] = byte((l & 0x00FF0000) >> 16)
	e[3] = byte((l & 0x0000FF00) >> 8)
	e[2] = byte((l & 0x000000FF))
}

/* Special pointer which includes the limit: The max bytes
*  taken up by the GDT, minus 1. Again, this NEEDS to be packed */
/*struct gdt_ptr
{
    unsigned short limit;
    unsigned int base;
} __attribute__((packed));
*/
/* Our GDT, with 3 entries, and finally our special GDT pointer */

var Gdt [3]GdtEntry

//struct gdt_entry gdt[3];
var Gp GdtPtr

//struct gdt_ptr gp;

/* This will be a function in start.asm. We use this to properly
*  reload the new segment registers */
//extern gdt_flush
func GdtFlush()

//extern halt
func Halt()

/* Setup a descriptor in the Global Descriptor Table */
func GdtSetGate(num int, base, limit uint32, access, gran byte) {
	gate := &Gdt[num]
	gate.SetBase(base)
	gate.SetLimit(limit)
	gate.SetAccess(access)
	gate.SetGranularity(gran)

	/*
	   // Setup the descriptor base address
	   gdt[num].base_low = (base & 0xFFFF);
	   gdt[num].base_middle = (base >> 16) & 0xFF;
	   gdt[num].base_high = (base >> 24) & 0xFF;

	   // Setup the descriptor limits
	   gdt[num].limit_low = (limit & 0xFFFF);
	   gdt[num].granularity = ((limit >> 16) & 0x0F);

	   // Finally, set up the granularity and access flags
	   gdt[num].granularity |= (gran & 0xF0);
	   gdt[num].access = access;
	*/
}

/* Should be called by main. This will setup the special GDT
*  pointer, set up the first 3 entries in our GDT, and then
*  finally call gdt_flush() in our assembler file in order
*  to tell the processor where the new GDT is and update the
*  new segment registers */
//export gdt_install_go
func GdtInstall() { //addr uintptr) {
	/* Setup the GDT pointer and limit */
	p := &Gp
	p.SetSize((8 /* sizeof gdt_entry */ * 3 /* num entries */) - 1)
	p.SetBase(uintptr(unsafe.Pointer(&Gdt)))

	/* Our NULL descriptor */
	GdtSetGate(0, 0, 0, 0, 0)

	/* The second entry is our Code Segment. The base address
	 *  is 0, the limit is 4GBytes, it uses 4KByte granularity,
	 *  uses 32-bit opcodes, and is a Code Segment descriptor.
	 *  Please check the table above in the tutorial in order
	 *  to see exactly what each value means */
	GdtSetGate(1, 0, 0xFFFFFFFF, 0x9A, 0xCF)

	/* The third entry is our Data Segment. It's EXACTLY the
	 *  same as our code segment, but the descriptor type in
	 *  this entry's access byte says it's a Data Segment */
	GdtSetGate(2, 0, 0xFFFFFFFF, 0x92, 0xCF)

	/* Flush out the old GDT and install the new changes! */
	GdtFlush()
}
