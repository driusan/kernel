package memory

import (
	"unsafe"
)

func Move(dst, src unsafe.Pointer, len int) (unsafe.Pointer, error) {
	if uintptr(dst) > uintptr(src) {
		for pos := len - 1; pos >= 0; pos-- {
			*(*byte)(unsafe.Pointer(uintptr(dst) + uintptr(pos))) = *(*byte)(unsafe.Pointer(uintptr(src) + uintptr(pos)))
			if pos == 0 {
				return dst, nil
			}
		}
	} else if uintptr(dst) < uintptr(src) {
		for pos := 0; pos < len; pos++ {
			*(*byte)(unsafe.Pointer(uintptr(dst) + uintptr(pos))) = *(*byte)(unsafe.Pointer(uintptr(src) + uintptr(pos)))
		}
	}
	return dst, nil
}
