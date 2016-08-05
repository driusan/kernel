// Package asm provides wrappers around assembly instructions to be
// called with Go syntax.
package asm

// INB executes an intel INB assembly instruction. Returns a byte from
// the input device at port.
func INB(port uint16) byte {
	// for some reason if this isn't a separate call, GCCGO strips it
	// out of the object file entirely. It's probably something I'm
	// missing.
	return inb(port)
}

//extern inb
func inb(port uint16) byte

// OUTB executes the Intel OUTB assembly instruction. Sends a byte to
// the output device at port.
func OUTB(port uint16, data byte) {
	outb(port, data)
}

//
//extern outb
func outb(port uint16, data byte)

// INL executes the Intel INL assembly instruction. Returns a uint32
// from the output device at port.
func INL(port uint16) uint32 {
	return inl(port)
}

//extern inl
func inl(port uint16) uint32

// OUTL executes the Intel OUTL assembly instruction. Writes a uint32
// to the output device at port.
func OUTL(port uint16, data uint32) {
	outl(port, data)
}

//extern outl
func outl(port uint16, data uint32)
