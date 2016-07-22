/**
 * This file contains symbols that gccgo links to. They're mostly
 * stubs, but they need to be written in either C or asm because
 * GCCGO will prepend prefix.packagename to every symbol.
 */
#include <stddef.h>
#include <stdint.h>

// Used for __go_print_*
void terminal_writestring(const char* data);
void putchar(char c);

// Stuff that go uses that I don't understand. This should go in it's own file.
void __go_print_string(char *s) {
	terminal_writestring(s);
}

void __go_register_gc_roots(void) { }

char* itoa(uint64_t) __asm__ ("boot.kernel.Itoa");

// These keywords are *not* available, because there's no malloc or free
// defined, but they get linked to, so there needs to be a stub.
void __go_new(void) { }
void __go_append(void) { }
void __gccgo_personality_v0(void) {}

void __go_print_int64(int64_t i) {
	// Write it in hex. Just go through every byte and print the 0-F
	// representation. This is easier than converting to a decimal.
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
	//putchar("(");
	//putchar(itoa(i));
	//putchar(")");
}
void __go_print_uint64(uint64_t i) {
	 __go_print_int64((int64_t)i);
}

uintptr_t __go_type_hash_identity(const void *a, uintptr_t b) {
	--a; a++;
	// This is meaningless. It's just to make gcc stop complaining about
	return b;
}

typedef struct FuncVal FuncVal;
struct FuncVal {
	void (*fn)(void);
};

FuncVal __go_type_hash_identity_descriptor;
FuncVal __go_type_equal_identity_descriptor;
FuncVal __go_type_hash_error_descriptor;
FuncVal __go_type_equal_error_descriptor;

