// This is not smart and does not get GCed. It will just keep leaking
// memory until there's none left.
// It's just a stub so that things compile
#include <stdint.h>

void*
runtime_mallocgc(uintptr_t size, uintptr_t typ, uint32_t flag)
{
	return malloc(size);
}
