package asm

//extern call
func call(addr uintptr)

// CALL will call a function at address addr, such that a RET instruction will
// return control to the current location
func CALL(addr uintptr) {
	println("In Go: calling ", addr)
	call(addr)
}
