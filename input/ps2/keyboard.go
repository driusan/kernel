package ps2

import (
	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/filesystem"
	"github.com/driusan/kernel/interrupts"
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
var InvalidKey error

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
	InvalidKey = PS2Error("Invalid key")
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
		//println("Released key", scancode)
		// key has been released
	} else {
		char, err := keymap(scancode)
		if err != nil {
			println("Unknown key", scancode)
			return
		}
		filesystem.Cons.SendByte(char)
		//println("Pressed key")
		// key has been pressed
	}
}

var modState keyMods

type keyMods byte

// Translates a scancode to an ASCII character
// TODO: Make this an interface that can handle different keyboard maps
//       (from the filesystem?)
// TODO: Change this from a byte to a rune?
func keymap(scancode byte) (byte, error) {
	switch scancode {
	case 2:
		return '1', nil
	case 3:
		return '2', nil
	case 4:
		return '3', nil
	case 5:
		return '4', nil
	case 6:
		return '5', nil
	case 7:
		return '6', nil

	default:

		return scancode, InvalidKey
	}
}
