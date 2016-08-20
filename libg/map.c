// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "map.h"

typedef struct __go_map Hmap;

/* Access a value in a map, returning a value and a presence indicator.  */

uint8_t mapaccess2(MapType* t, Hmap *h, byte* key, byte* val) {
	byte *mapval;
	size_t valsize;

	mapval = __go_map_index(h, key, 0);
	valsize = t->__val_type->__size;
	if (mapval == 0) {
		__builtin_memset(val, 0, valsize);
		return 0;
	} else {
		__builtin_memcpy(val, mapval, valsize);
		return 1;
	}
}

