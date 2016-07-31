/**
 * This file contains symbols that gccgo links to. They're mostly
 * stubs, but they need to be written in either C or asm because
 * GCCGO will prepend prefix.packagename to every symbol.
 */
#include <stddef.h>
#include <stdint.h>

#include "runtime.h"
#include "go-type.h"
// Used for __go_print_*
extern void terminal_writestring(const char* data);
extern void putchar(char c);
extern void halt(void);

void __go_panic(void) {
	terminal_writestring("Kernel panic. TODO: Add more debug info here.");
	halt();
}
// Stuff that go uses that I don't understand. This should go in it's own file.
void __go_print_string(char *s) {
	terminal_writestring(s);
}

void __go_print_nl(void) {
	terminal_writestring("\n");
}

void __go_print_pointer(void *p) {
	terminal_writestring("Pointer: ");
	__go_print_uint64((uint64) p);
}
void __go_register_gc_roots(struct root_list *roots __attribute__((unused))) { }


void
runtime_panicstring(const char* error) { //__attribute__ ((noreturn)) {
	terminal_writestring(error);
	halt();
}
// This should be done in Go, but there's not enough of the go
// runtime implemented to do it properly yet.
void printdec(int64_t i) {
	// The highest int64_t is  18446744073709551616, a
	// 			12345678901234567890
	// 20 digit string. Since we don't have a malloc/free yet,
	// just use a 20 character array to store the string representation.
	// We need to do this, because the obvious algorithm counts it backwards
	// so we need to store an intermediary and then print the reverse.
	char c[21];
	
	if (i < 0) {
		putchar('-');
		i = -i;
	}
	int digit = 0;
	while(i) {
		c[digit++] = i % 10;
		i = i / 10;
	}

	while(digit--) {
		putchar(c[digit] + '0');
	}
	
	
}

// These keywords are *not* available, because there's no malloc or free
// defined, but they get linked to, so there needs to be a stub.
void __go_new(void) { }
void __go_append(void) { }
//void __gccgo_personality_v0(void) {}
void __go_print_int64(int64_t i) {
	// Write it in hex. Just go through every byte and print the 0-F
	// representation. This is easier than converting to a decimal.
	printdec(i);
	/*
	terminal_writestring("0x");

	for(char j = 15; j >= 0; j--) {
		uint64_t mask = 0xF << (j*4);
		char thebyte = (i & mask) >> (j*4);
		if (thebyte < 10) {
			putchar(thebyte+'0');
		} else {
			putchar(thebyte+('a'-10));
		}
	}
	putchar(' ');
	putchar('(');
	printdec(i);
	putchar(')');*/
}
void __go_print_uint64(uint64_t i) {
	 __go_print_int64((int64_t)i);
}

