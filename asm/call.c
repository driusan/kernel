// TODO: Port these to real assembly, instead of single line
// C functions.
#include <stdint.h>

void __attribute__ ((noinline)) call(uintptr_t addr)
{
	terminal_writestring("In C: Calling something at");
/*	__go_print_uint64(addr);
	terminal_writestring("\n");*/
    __asm__ __volatile__ ("call *(%0)" : : "a" (addr));
}

void __attribute__ ((noinline)) jmp(uintptr_t addr)
{
    __asm__ __volatile__ ("jmp *%0" : : "a" (addr));
}


void __attribute__ ((noinline)) invlpg(uintptr_t addr) {
	__asm__ __volatile__ ("invlpg (%0)" : : "b"(addr) : "memory");
}
