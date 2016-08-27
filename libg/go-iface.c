// Adapted from go-iface.goc in gccgo runtime to be pure C
// This is needed to be able to import "io"

#include <stdint.h>

#include "go-type.h"
#include "interface.h"
typedef struct __go_interface interface;
typedef struct __go_type_descriptor descriptor;

typedef struct{
	interface ret;
	uint8_t ok;
} ifaceI2I2Ret;
// Convert a non-empty interface to a non-empty interface.
ifaceI2I2Ret ifaceI2I2(descriptor *inter, interface i) {
	ifaceI2I2Ret r;
	if (i.__methods == 0) {
		r.ret.__methods = 0;
		r.ret.__object = 0;
		r.ok = 0;
	} else {
		r.ret.__methods = __go_convert_interface_2(inter,
							 i.__methods[0], 1);
		r.ret.__object = i.__object;
		r.ok = r.ret.__methods != 0;
	}
	return r;
}