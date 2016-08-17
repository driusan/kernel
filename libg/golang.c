/**
 * This file contains symbols that gccgo links to. They're mostly
 * stubs, but they need to be written in either C or asm because
 * GCCGO will prepend prefix.packagename to every symbol.
 */
#include <stddef.h>
#include <stdint.h>

#include "runtime.h"
#include "go-type.h"
#include "unwind.h"

void __go_register_gc_roots(struct root_list *roots __attribute__((unused))) { }

// This is defined in Go because it's a C string string, and GoPrintString
// takes a Go style string.
void runtime_panicstring(const char* error) {
	terminal_writestring(error);
	// Halt only needs to be called once, but the for loop fixes a compiler
	// warning about __noreturn__ function returning.
	for (;;) halt();
}

// These keywords are *not* available, because there's no malloc or free
// defined, but they get linked to, so there needs to be a stub.
void __go_new(void) { }
//void __go_append(void) { }

/*
Lots of symbols aren't defined for this, so just register 
__gccgo_personality_v0 as void below and hope for the best.
*/
/*
_Unwind_Reason_Code
PERSONALITY_FUNCTION (int, _Unwind_Action, _Unwind_Exception_Class,
		      struct _Unwind_Exception *, struct _Unwind_Context *)
  __attribute__ ((no_split_stack, flatten)) {
}*/

void __gccgo_personality_v0(void) {
	//__go_print_string("In __gccgo_personality_v0\n");

};// __attribute__ ((no_split_stack, flatten));

// This should be done in Go, but it's not clear how to set the internals
// of a string in Go
String
__go_byte_array_to_string (const void* p, intgo len)
{
	String ret;
	ret.str = (const unsigned char *)p;
	ret.len = len;
	return ret;
}
