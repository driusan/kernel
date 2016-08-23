// Copyright 2009, 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted from string.goc for libg in the kernel



#include "runtime.h"
#include "go-string.h"
enum {
	Runeself = 0x80
};
#define charntorune(pv, str, len) __go_get_rune(str, len, pv)

intgo stringiter(String s, int k) {
	int32 l;

	if(k >= s.len) {
		// end of iteration
		return 0;
	}

	l = s.str[k];
	if(l < Runeself) {
		return k+1;
	}

	// multi-char rune
	return k + charntorune(&l, s.str+k, s.len-k);
}

typedef struct{
	intgo K;
	int32_t V;
} stringiter2Ret;
stringiter2Ret stringiter2(String s, int k) {
	stringiter2Ret r;
	if(k >= s.len) {
		// retk=0 is end of iteration
		r.K = 0;
		r.V = 0;
		goto out;
	}

	r.V = s.str[k];
	if(r.V < Runeself) {
		r.K = k+1;
		goto out;
	}

	// multi-char rune
	r.K = k + charntorune(&r.V, s.str+k, s.len-k);

out:
	return r;
}