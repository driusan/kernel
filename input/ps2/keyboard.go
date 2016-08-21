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
	case 8:
		return '7', nil
	case 9:
		return '8', nil
	case 10:
		return '9', nil
	case 11:
		return '0', nil
	case 12:
		return '-', nil
	case 13:
		return '=', nil
	case 16:
		return 'q', nil
	case 17:
		return 'w', nil
	case 18:
		return 'e', nil
	case 19:
		return 'r', nil
	case 20:
		return 't', nil
	case 21:
		return 'y', nil
	case 22:
		return 'u', nil
	case 23:
		return 'i', nil
	case 24:
		return 'o', nil
	case 25:
		return 'p', nil
	case 28:
		return '\n', nil
	case 30:
		return 'a', nil
	case 31:
		return 's', nil
	case 32:
		return 'd', nil
	case 33:
		return 'f', nil
	case 34:
		return 'g', nil
	case 35:
		return 'h', nil
	case 36:
		return 'j', nil
	case 37:
		return 'k', nil
	case 38:
		return 'l', nil
	case 44:
		return 'z', nil
	case 45:
		return 'x', nil
	case 46:
		return 'c', nil
	case 47:
		return 'v', nil
	case 48:
		return 'b', nil
	case 49:
		return 'n', nil
	case 50:
		return 'm', nil
	case 51:
		return ',', nil
	case 52:
		return '.', nil
	case 53:
		return '/', nil
	case 57:
		return ' ', nil
	default:

		return scancode, InvalidKey
	}
}
