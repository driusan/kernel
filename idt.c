/* Loads IDT entry. I haven't figured out how to do the 
 *__attribute((packed))__ in gccgo, so for now this is in C.
 *
 * Copied from Brendon's tutorial at 
 * http://www.osdever.net/bkerndev/Docs/idt.htm
 */
#include <stddef.h>
#include <stdint.h>
/* FIXME: This shouldn't be here. Copied from osdev Meaty Skeleton */
void* memset(void* bufptr, int value, size_t size)
{
	unsigned char* buf = (unsigned char*) bufptr;
	for ( size_t i = 0; i < size; i++ )
		buf[i] = (unsigned char) value;
	return bufptr;
}

extern void Halt(void);
/* Defines an IDT entry */
struct idt_entry
{
    unsigned short base_lo;
    unsigned short sel;        /* Our kernel segment goes here! */
    unsigned char always0;     /* This will ALWAYS be set to 0! */
    unsigned char flags;       /* Set using the above table! */
    unsigned short base_hi;
} __attribute__((packed));

struct idt_ptr
{
    uint16_t limit;
    uint32_t base;
} __attribute__((packed));

/* Declare an IDT of 256 entries. Although we will only use the
*  first 32 entries in this tutorial, the rest exists as a bit
*  of a trap. If any undefined IDT entry is hit, it normally
*  will cause an "Unhandled Interrupt" exception. Any descriptor
*  for which the 'presence' bit is cleared (0) will generate an
*  "Unhandled Interrupt" exception */
struct idt_entry idt[256];
struct idt_ptr idtp;

/* This exists in 'start.asm', and is used to load our IDT */
extern void idt_load();

/* Use this function to set an entry in the IDT. Alot simpler
*  than twiddling with the GDT ;) */
void idt_set_gate(unsigned char num, unsigned long base, unsigned short sel, unsigned char flags)
{
    idt[num].base_lo = base & 0xFFFF;
    idt[num].base_hi = (base >> 16) & 0xFFFF;

    idt[num].flags = flags;
    idt[num].sel = sel;
    idt[num].always0 = 0;

    /* We'll leave you to try and code this function: take the
    *  argument 'base' and split it up into a high and low 16-bits,
    *  storing them in idt[num].base_hi and base_lo. The rest of the
    *  fields that you must set in idt[num] are fairly self-
    *  explanatory when it comes to setup */
}

/* Installs the IDT */
void idt_install()
{
//	terminal_writestring("In IDT Install");
//	__go_print_int64((sizeof (struct idt_entry) * 256) -1);
//	terminal_writestring(" ");
//	__go_print_int64(sizeof (struct idt_entry));
    /* Sets the special IDT pointer up, just like in 'gdt.c' */
    idtp.limit = (sizeof (struct idt_entry) * 256) - 1;
    idtp.base = &idt;

    /* Clear out the entire IDT, initializing it to zeros */
    // never defined memset, so do this manually
    memset(idt, 0, sizeof(struct idt_entry)*256);
	for(uint16_t i = 0; i < 256; i++) {
		idt[i].base_lo = 0;
		idt[i].sel = 0;
		idt[i].always0 = 0;
		idt[i].flags = 0;
		idt[i].base_hi = 0;
	}

    /* Add any new ISRs to the IDT here using idt_set_gate */
	isrs_install();
    /* Points the processor's internal register to the new IDT */
    idt_load();
}
