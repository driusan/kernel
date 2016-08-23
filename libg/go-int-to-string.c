// Manually transcribed without internet. TODO: just download this
// from the gccgo frontend and make sure it has the copyright

#include "runtime.h"
#include "arch.h"
#include "malloc.h"

String __go_int_to_string(intgo v) {
	char buf[4];
	int len;
	unsigned char *retdata;
	String ret;

	if (v < 0)
		v = 0xfffd;

	if (v <= 0x7f)
	{
	buf[0] = v;
	len = 1;
	} else if (v <= 0x7ff)
	{
		buf[0] = 0xc0 + (v >> 6);
		buf[1] = 0x80 + (v & 0x3f);
		len = 2;
	} else {
		// This should be enough for now.
	}

	retdata = runtime_mallocgc(len, 0, FlagNoScan);
	__builtin_memcpy(retdata, buf, len);
	ret.str = retdata;
	ret.len = len;
	return ret;
}
