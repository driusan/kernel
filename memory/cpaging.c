#include <stdint.h>

uint32_t* getPageTable(int16_t i);

// Go doesn't provide any obvious way to force the alignment on a 4096
// byte boundary, so these are defined in C and then the addressed
// retrieved in Go by the get* functions
uint32_t page_directory[1024] __attribute__((aligned(4096)));
uint32_t page_tables[1024*1024] __attribute__((aligned(4096)));

uint32_t* getPageDirectory() {
	return page_directory;
}

uint32_t* getPageTable(int16_t i) {
	return &page_tables[i*1024];
}

uint32_t* getInitialPageTable() {
	return getPageTable(0);
}

