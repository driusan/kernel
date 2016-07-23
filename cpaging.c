/* This file contains the kernel paging code.
 * I haven't figured out how to force the alignemtn in Go,
 * so for now it's in C.
 */
#include <stdint.h>

extern void loadPageDirectory(uint32_t*);
extern void enablePaging();
const uint32_t PagePresent = 1;
const uint32_t PageReadWrite = 2;
const uint32_t PageUserspace = 4;
/* 
	PageWriteThrough = 8,
	PageCacheDisaled = 16,
	PageAccessed = 32,
	// What is bit 7? 
	PageIs4MB = 8,
	PageGlobal = 9,
	PageAddressMask = 0xfffff000
};
*/
// haven't figured out how to force the alignment in Go yet
// do this in 
uint32_t page_directory[1024] __attribute__((aligned(4096)));
uint32_t first_page_table[1024] __attribute__((aligned(4096)));

void initialize_paging() {
	uint16_t i;
	// Mark all pages as readwrite, but not present.
	for (i = 0; i < 1024; i++) {
		page_directory[i] = PageReadWrite;
	}
	// Create the first page table, and mark it as Present and ReadWrite
	// Mark the first page table
	for (i = 0; i < 1024; i++) {
		first_page_table[i] = (i*0x1000) | PagePresent | PageReadWrite;

		// Swap 2 arbitrary memory locations to make sure paging is working
		if (i*0x1000 == 0x60000) {
			first_page_table[i] = 0x70000 | 3;
		}
		if (i*0x1000 == 0x70000) {
			first_page_table[i] = 0x60000 | 3;
		}
	}
	page_directory[0] = ((uint32_t)first_page_table) | PagePresent | PageReadWrite;
	loadPageDirectory(page_directory);
	enablePaging();
}
