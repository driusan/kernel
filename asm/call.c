// TODO: Port these to real assembly, instead of single line
// C functions.
#include <stdint.h>

void __attribute__ ((noinline)) call(uintptr_t addr)
{
	terminal_writestring("In C: Calling something at");
/*	__go_print_uint64(addr);
	terminal_writestring("\n");*/
    __asm__ __volatile__ ("ljmp *(%0)" : : "a" (addr));
}

