#include <stddef.h>
#include <stdint.h>

//typedef int uintptr_t;
void terminal_writestring(const char* data);
void terminal_putchar(char c);

// Stuff that go uses that I don't understand. This should go in it's own file.
void __go_print_string(char *s) {
	terminal_writestring(s);
}

char* itoa(uint64_t) __asm__ ("boot.kernel.Itoa");
void __go_new(void) { }
void __go_append(void) { }
void __gccgo_personality_v0(void) {}
void __go_print_int64(uint64_t i) {
	// 2^64 = 18446744073709551615, a 20 digit number.
	//	12345678901234567890
	// so iterate over
	//terminal_writestring(itoa(i)); 
	/*
	if (i < 0) {
		terminal_putchar('-');
	}
	*/
	terminal_writestring("0x");

	// uint64_t can represent up to:
	//0xFF FF FF FF FF FF FF FF
	for(char j = 15; j >= 0; j--) {
		uint64_t mask = 0xF << (j*4);
		char thebyte = (i & mask) >> (j*4);
		if (thebyte < 10) {
			putchar(thebyte+'0');
		} else {
			putchar(thebyte+('a'-10));
		}
	}
}
uintptr_t __go_type_hash_identity(const void *a, uintptr_t b) { return 0;}

typedef struct FuncVal FuncVal;
struct FuncVal {
	void (*fn)(void);
};

FuncVal __go_type_hash_identity_descriptor;
FuncVal __go_type_equal_identity_descriptor;
FuncVal __go_type_hash_error_descriptor;
FuncVal __go_type_equal_error_descriptor;

