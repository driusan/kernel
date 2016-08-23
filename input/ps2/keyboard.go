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

var isShifted bool

func KeyboardHandler(r *interrupts.Registers) {
	scancode := readPS2Port()
	if scancode&0x80 != 0 {
		// key has been released
		switch scancode {
		case 170, 182:
			isShifted = false
		default:
			// println("Released key", scancode)

		}
	} else {
		switch scancode {
		case 42, 54:
			isShifted = true
		default:
			char, err := keymap(scancode)
			if err != nil {
				println("Unknown key", scancode)
				return
			}
			filesystem.Cons.SendByte(char)
		}
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
		if isShifted {
			return 'Q', nil
		} else {
			return 'q', nil
		}
	case 17:
		if isShifted {
			return 'W', nil
		} else {
			return 'w', nil
		}
	case 18:
		if isShifted {
			return 'E', nil
		} else {
			return 'e', nil
		}
	case 19:
		if isShifted {
			return 'R', nil
		} else {
			return 'r', nil
		}
	case 20:
		if isShifted {
			return 'T', nil
		} else {
			return 't', nil
		}
	case 21:
		if isShifted {
			return 'Y', nil
		} else {
			return 'y', nil
		}
	case 22:
		if isShifted {
			return 'U', nil
		} else {
			return 'u', nil
		}
	case 23:
		if isShifted {
			return 'I', nil
		} else {
			return 'i', nil
		}
	case 24:
		if isShifted {
			return 'O', nil
		} else {
			return 'o', nil
		}
	case 25:
		if isShifted {
			return 'P', nil
		} else {
			return 'p', nil
		}
	case 28:
		return '\n', nil
	case 30:
		if isShifted {
			return 'A', nil
		} else {
			return 'a', nil
		}
	case 31:
		if isShifted {
			return 'S', nil
		} else {
			return 's', nil
		}
	case 32:
		if isShifted {
			return 'D', nil
		} else {
			return 'd', nil
		}
	case 33:
		if isShifted {
			return 'F', nil
		} else {
			return 'f', nil
		}
	case 34:
		if isShifted {
			return 'G', nil
		} else {
			return 'g', nil
		}
	case 35:
		if isShifted {
			return 'H', nil
		} else {
			return 'h', nil
		}
	case 36:
		if isShifted {
			return 'J', nil
		} else {
			return 'j', nil
		}
	case 37:
		if isShifted {
			return 'K', nil
		} else {
			return 'k', nil
		}
	case 38:
		if isShifted {
			return 'L', nil
		} else {
			return 'l', nil
		}
	case 44:
		if isShifted {
			return 'Z', nil
		} else {
			return 'z', nil
		}
	case 45:
		if isShifted {
			return 'X', nil
		} else {
			return 'x', nil
		}
	case 46:
		if isShifted {
			return 'C', nil
		} else {
			return 'c', nil
		}
	case 47:
		if isShifted {
			return 'V', nil
		} else {
			return 'v', nil
		}
	case 48:
		if isShifted {
			return 'B', nil
		} else {
			return 'b', nil
		}
	case 49:
		if isShifted {
			return 'N', nil
		} else {
			return 'n', nil
		}
	case 50:
		if isShifted {
			return 'M', nil
		} else {
			return 'm', nil
		}
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
