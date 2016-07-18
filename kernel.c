#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

/* Ghetto sleep. */
void s(int x) {
	for(int i =0;i < x*100000; i++) {
		
	}
}

/* Hardware text mode color constants. */
enum vga_color {
	COLOR_BLACK = 0,
	COLOR_BLUE = 1,
	COLOR_GREEN = 2,
	COLOR_CYAN = 3,
	COLOR_RED = 4,
	COLOR_MAGENTA = 5,
	COLOR_BROWN = 6,
	COLOR_LIGHT_GREY = 7,
	COLOR_DARK_GREY = 8,
	COLOR_LIGHT_BLUE = 9,
	COLOR_LIGHT_GREEN = 10,
	COLOR_LIGHT_CYAN = 11,
	COLOR_LIGHT_RED = 12,
	COLOR_LIGHT_MAGENTA = 13,
	COLOR_LIGHT_BROWN = 14,
	COLOR_WHITE = 15,
};
 
uint8_t make_color(uint8_t fg, uint8_t bg) {
	return fg | bg << 4;
}
 
uint16_t make_vgaentry(char c, uint8_t color) {
	uint16_t c16 = c;
	uint16_t color16 = color;
	return c16 | color16 << 8;
}
 
size_t strlen(const char* str) {
	size_t len = 0;
	while (str[len])
		len++;
	return len;
}
 
static const size_t VGA_WIDTH = 80;
static const size_t VGA_HEIGHT = 25;
 
typedef struct {
	uint16_t row;
	uint16_t column;
	uint8_t color;
	uint16_t* buffer;
} Terminal;

Terminal terminal;
extern void terminal_initialize(Terminal*) __asm__ ("boot.kernel.InitializeTerminal");
 
void reset_terminal_buffer(Terminal* t) {
	t->buffer = (uint16_t*)(0xB8000);
	for (size_t y = 0; y < VGA_HEIGHT; y++) {
		for (size_t x = 0; x < VGA_WIDTH; x++) {
			const size_t index = y * VGA_WIDTH + x;
			t->buffer[index] = make_vgaentry(' ', t->color);
		}
	}

}
void terminal_setcolor(uint8_t color) {
	terminal.color = color;
}
 
void terminal_putentryat(char c, uint8_t color, size_t x, size_t y) {
	const size_t index = y * VGA_WIDTH + x;
	terminal.buffer[index] = make_vgaentry(c, color);
}
 
void terminal_putchar(char c) {
	if (c == '\n') {
		terminal.column = 0;
		if (terminal.row < VGA_HEIGHT-1) {
			terminal.row++;
		} else {
			/* scroll everything up 1 row */
			for(size_t y = 1; y < VGA_HEIGHT; y++) {
				for(size_t x = 0; x < VGA_WIDTH; x++) {
					const size_t idx = y * VGA_WIDTH + x;
					terminal.buffer[idx-VGA_WIDTH] = terminal.buffer[idx];
				}
			}
			/* clear the last row. */
			for(size_t x = 0; x < VGA_WIDTH; x++) {
				terminal_putentryat(' ', terminal.color, x, VGA_HEIGHT-1);
			}
		}
			
		return;
	}
	terminal_putentryat(c, terminal.color, terminal.column, terminal.row);
	if (++terminal.column == VGA_WIDTH) {
		terminal.column = 0;
		if (++terminal.row == VGA_HEIGHT) {
			terminal.row = 0;
		}
	}
}
 
void terminal_writestring(const char* data) {
	size_t datalen = strlen(data);
	for (size_t i = 0; i < datalen; i++)
		terminal_putchar(data[i]);
}

void kernel_main() {
	/* Initialize terminal interface */
	terminal_initialize(&terminal);
	for(;;)
	terminal_writestring("Hello, kernel World!\n");
}


// Stuff that go uses that I don't understand. This should go in it's own file.
void __go_print_string(char *s) {
	terminal_writestring(s);
}

uintptr_t __go_type_hash_identity(const void *a, uintptr_t b) { return 0;}

typedef struct FuncVal FuncVal;
struct FuncVal {
	void (*fn)(void);
};

FuncVal __go_type_hash_identity_descriptor;
FuncVal __go_type_equal_identity_descriptor;
