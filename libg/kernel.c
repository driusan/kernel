/**
 * This file contains helpers to access things written in the Go kernel
 */ 
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

 
typedef struct {
	uint16_t row;
	uint16_t column;
	uint8_t color;
	uint16_t* buffer;
} Terminal;

extern Terminal *terminal __asm__("github_com_driusan_kernel_terminal.Term");
size_t strlen(const char* str) {
	size_t len = 0;
	while (str[len])
		len++;
	return len;
}

/* TODO: Find out if the name gets mangled in a deterministic way */
// void putentryat(Terminal *t, char c, uint8_t color, int16_t x, int16_t y) __asm__("boot.kernel.PutEntryAt.pN20_boot.kernel.Terminal");
// void terminal_putchar(Terminal *t, char c) __asm__("github_com_driusan_kernel.pN20_Terminal");
 

void putentryat(Terminal *t, char c, uint8_t color, int16_t x, int16_t y) __asm__("github_com_driusan_kernel_terminal.PutEntryAt.pN43_github_com_driusan_kernel_terminal.Terminal");;
void terminal_putchar(Terminal *t, char c) __asm__("github_com_driusan_kernel_terminal.PutChar.pN43_github_com_driusan_kernel_terminal.Terminal");

void putchar(char c) {
	terminal_putchar(terminal, c);
}

void terminal_writestring(const char* data) {
	size_t datalen = strlen(data);
	for (size_t i = 0; i < datalen; i++)
		terminal_putchar(terminal, data[i]);
}




