/**
 * This file contains symbols that gccgo links to. They're mostly
 * stubs, but they need to be written in either C or asm because
 * GCCGO will prepend prefix.packagename to every symbol.
 * Where possible, the version from the gccgo is used unmodified in the
 * libg directory. Where not possible, a stub is written here.
 */
#include <stddef.h>
#include <stdint.h>

#include "runtime.h"
#include "go-type.h"
#include "unwind.h"
#include "malloc.h"

extern void halt();
extern void* malloc(size_t);
extern void terminal_writestring(const char *);
extern uint8_t pagingInitialized();
void __go_register_gc_roots(struct root_list *roots __attribute__((unused))) { }

void *
__go_new (const struct __go_type_descriptor *td, uintptr_t size)
{
	// This shouldn't be called before paging is initialized. If it
	// is, just return an unused value and hope for the best
	if (!pagingInitialized()) {
		terminal_writestring("Warning: call to __go_new before memory is ready");
		return (void *)0xF000000;
	}
	return runtime_mallocgc (size,
		   (uintptr) td | TypeInfo_SingleObject,
		   td->__code & GO_NO_POINTERS ? FlagNoScan : 0);
}

// This is defined in Go because it's a C string string, and GoPrintString
// takes a Go style string.
void runtime_panicstring(const char* error) {
	terminal_writestring(error);
	// Halt only needs to be called once, but the for loop fixes a compiler
	// warning about __noreturn__ function returning.
	for (;;) halt();
}

/*
Lots of symbols aren't defined for this, so just register 
__gccgo_personality_v0 as void below and hope for the best.
*/
/*

_Unwind_Reason_Code
PERSONALITY_FUNCTION (int, _Unwind_Action, _Unwind_Exception_Class,
		      struct _Unwind_Exception *, struct _Unwind_Context *)
  __attribute__ ((no_split_stack, flatten)) {
}
*/

void __gccgo_personality_v0(void) {
	//__go_print_string("In __gccgo_personality_v0\n");

};// __attribute__ ((no_split_stack, flatten));


// Something calls this before paging is initialized, but I haven't investigated
// what yet. For now, this modifies the gccgo version to add a check and use
// an arbitrary address in that case, otherwise the video memory gets corrupted.
String
__go_byte_array_to_string (const void* p, intgo len)
{
	String ret;
        unsigned char *retdata;
	if (pagingInitialized()) {
		retdata = (unsigned char *)malloc((size_t )len); //(const unsigned char *)p;
	} else {
		terminal_writestring("Warning: call to __go_byte_array_to_string before memory is ready");
		retdata = (unsigned char *)0xFF00000;
	}

	__builtin_memcpy(retdata, p, len);
	ret.str = retdata;
	ret.len = len;
	return ret;
}
