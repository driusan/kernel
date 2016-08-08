package ps2

import (
	"asm"
	"interrupts"
)

//extern printhex
func printhex(int64)

type PS2Error string

func (p PS2Error) Error() string {
	return string(p)
}

var SendFailure error
var TooManyRetries error
var UnknownError error

// Acknowledgement responses to PS2 commands
const (
	Success = byte(0xFA)
	Fail    = byte(0xFC)
	Resend  = byte(0xFE)
)

func InitPkg() {
	TooManyRetries = PS2Error("Too many retries sending PS2 command")
	SendFailure = PS2Error("Failure sending PS2 command")
	UnknownError = PS2Error("Unknown error")
}

func outputReady() bool {
	status := asm.INB(0x64)
	if status&2 == 0 {
		return true
	}
	return false
}

func waitOutput() {
	for {
		if outputReady() == true {
			return
		}
	}
}

func inputReady() bool {
	status := asm.INB(0x64)
	if status&1 == 1 {
		return true
	}
	return false
}

func waitInput() {
	for {
		if inputReady() == true {
			return
		}
	}
}

func readPS2Port() byte {
	waitInput()
	return asm.INB(0x60)

}

func KeyboardHandler(r *interrupts.Registers) {
	scancode := readPS2Port()
	if scancode&0x80 != 0 {
		println("Released key", scancode)
		// key has been released
	} else {
		println("Pressed key")
		// key has been pressed
	}
}
