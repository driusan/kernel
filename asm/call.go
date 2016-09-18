package asm

import "unsafe"

//extern call
func call(addr uintptr)

// CALL will call a function at address addr, such that a RET instruction will
// return control to the current location
func CALL(addr unsafe.Pointer) {
	println("In Go: calling ", addr)
	call(uintptr(addr))
}

//extern jmp
func jmp(addr uintptr)

func JMP(addr unsafe.Pointer) {
	jmp(uintptr(addr))
}

//extern invlpg
func invlpg(addr uintptr)

// Invalidates the page table cache for a virtual address, so that the next
// time it's accessed the CPU will reread the page table from main memory
func INVLPG(vaddr uintptr) {
	invlpg(vaddr)
}
