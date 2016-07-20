#include <stddef.h>
#include <stdint.h>

//typedef int uintptr_t;
void terminal_writestring(const char* data);
void terminal_putchar(char c);

// Stuff that go uses that I don't understand. This should go in it's own file.
void __go_print_string(char *s) {
	terminal_writestring(s);
}

void __go_print_int64(uint64_t i) {
	//0xFF FF FF FF FF FF FF FF
	if (i < 0) {
		terminal_putchar('-');
	}
	terminal_writestring("0x");
	for(char j = 15; j >= 0; j--) {
		uint64_t mask = 0xF << (j*4);
		char thebyte = (i & mask) >> (j*4);
		if (thebyte < 10) {
			terminal_putchar(thebyte+'0');
		} else {
			terminal_putchar(thebyte+('a'-10));
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
