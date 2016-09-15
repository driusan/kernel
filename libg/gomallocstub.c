// This is not smart and does not get GCed. It will just keep leaking
// memory until there's none left.
// It's just a stub so that things compile
#include <stdint.h>

#include "runtime.h"

void*
runtime_mallocgc(uintptr_t size, uintptr_t typ, uint32_t flag)
{
	if (size == 0) {
		return &runtime_zerobase;
	}
	return malloc(size);
}
