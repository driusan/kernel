#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
 
uint8_t make_color(uint8_t fg, uint8_t bg) __asm__("boot.kernel.MakeColor");
uint16_t make_vgaentry(char c, uint8_t color) __asm__ ("boot.kernel.MakeVgaEntry");
 
 
size_t strlen(const char* str) {
	size_t len = 0;
	while (str[len])
		len++;
	return len;
}
 
typedef struct {
	uint16_t row;
	uint16_t column;
	uint8_t color;
	uint16_t* buffer;
} Terminal;

Terminal terminal;
void terminal_initialize(Terminal*) __asm__ ("boot.kernel.InitializeTerminal");
 

// This is easier to do in C than in Go, since t->buffer is a pointer, not an
// array. This is used as a helper from the Go side.
void setbuffer(Terminal *t, uint16_t idx, uint16_t val) {
	t->buffer[idx] = val;
} 

uint16_t getbuffer(Terminal *t, uint16_t idx) {
	return t->buffer[idx];
} 


/* TODO: Find out if the name gets mangled in a deterministic way */
void putentryat(Terminal *t, char c, uint8_t color, int16_t x, int16_t y) __asm__("boot.kernel.PutEntryAt.pN20_boot.kernel.Terminal");
void terminal_putchar(Terminal *t, char c) __asm__("boot.kernel.PutChar.pN20_boot.kernel.Terminal");

void putchar(char c) {
	terminal_putchar(&terminal, c);
}

void terminal_writestring(const char* data) {
	size_t datalen = strlen(data);
	for (size_t i = 0; i < datalen; i++)
		terminal_putchar(&terminal, data[i]);
}

void kernel_main() {
	/* Initialize terminal interface */
	terminal_initialize(&terminal);
/*	for(;;)
	terminal_writestring("Hello, kernel World!\n");
*/
}


