package terminal

import (
	"unsafe"
)

// Prints a CString. A memory address that ends when it reaches a \0.
func PrintCString(s unsafe.Pointer) {
	var b byte = *(*byte)(s)
	for b != 0 {
		b = *(*byte)(s)
		PrintRune(rune(b))
		s = unsafe.Pointer(uintptr(s) + 1)
	}
}
func PrintRune(r rune) {
	if r < 256 {
		Term.PutChar(byte(r))
	} else {
		Term.PutChar(3)
	}
	/*	switch r {
		case
	}*/
}
func PrintHex(i uint64) {
	Term.PutChar('0')
	Term.PutChar('x')

	if i == 0 {
		Term.PutChar('0')
		return
	}

	foundByte := false
	for j := 15; j >= 0; j-- {
		mask := uint64(0xF << uint(j*4))
		thebyte := byte((i & mask) >> uint(j*4))
		if thebyte != 0 {
			foundByte = true
		}
		if foundByte {
			if thebyte < 10 {
				Term.PutChar(thebyte + '0')
			} else {
				Term.PutChar(thebyte + ('a' - 10))
			}
		}
	}

}

func PrintDec(i int64) {
	if i == 0 {
		Term.PutChar('0')
		return
	}
	// The highest int64_t is  18446744073709551616, a
	// 			12345678901234567890
	// 20 digit string. Since we don't have a malloc/free yet,
	// just use a 20 character array to store the string representation.
	// We need to do this, because the obvious algorithm counts it backwards
	// so we need to store an intermediary and then print the reverse.
	var c [20]byte

	if i < 0 {
		Term.PutChar('-')
		i = -i
	}

	digit := 0

	for i > 0 {
		digit++
		c[digit] = byte(i % 10)
		i /= 10
	}

	for ; digit > 0; digit-- {
		Term.PutChar(c[digit] + '0')
	}
}
