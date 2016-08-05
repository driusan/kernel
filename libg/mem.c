/* These are some functions that are needed by the implementations of
  various things taken from the gccgo frontend. */
#include <stddef.h>

int memcmp(const unsigned char *s1, const unsigned char *s2, size_t n) {
	size_t i;
	for(i =0; i < n; i++, s1++, s2++) {
		if(*s1 < *s2) {
			return -1;
		} else if(*s2 > *s1) {
			return 1;
		}
	}
	return 0;
}
void* memcpy(void *dst, const void *src, size_t n) {
	unsigned char *csrc = (unsigned char *)src;
	unsigned char *cdst = (unsigned char *)dst;
	size_t i;
	for (i = 0; i < n; i++) {
		cdst[i] = csrc[i];
	}
	return (void *)(cdst + i);

}