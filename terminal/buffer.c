#include <stdint.h>

// This is easier to do in C than in Go, since t->buffer is a pointer, not an
// array. This is used as a helper from the Go side.

typedef struct {
	uint16_t row;
	uint16_t column;
	uint8_t color;
	uint16_t* buffer;
} Terminal;


void setbuffer(Terminal *t, uint16_t idx, uint16_t val) {

	t->buffer[idx] = val;
} 

uint16_t getbuffer(Terminal *t, uint16_t idx) {
	return t->buffer[idx];
} 

