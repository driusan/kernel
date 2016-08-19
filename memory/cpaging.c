#include <stdint.h>

// Go doesn't provide any obvious way to force the alignment on a 4096
// byte boundary, so these are defined in C and then the addressed
// retrieved in Go by the get* functions
uint32_t page_directory[1024] __attribute__((aligned(4096)));
uint32_t first_page_table[1024] __attribute__((aligned(4096)));

uint32_t* getPageDirectory() {
	return page_directory;
}

uint32_t* getInitialPageTable() {
	return first_page_table;
}

