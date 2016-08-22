package terminal

import "C"

import (
	"unsafe"
)

// If we import C, the GCCGO cross-compiler claims we're not using it.
// If we import it as _ "C", go test claims we can't rename it since it's
// the real go tool chain.
// This is a hack so that gmake, go test ./... and go fmt ./...
// all work. C.Nothing is a struct{}, so hopefully this gets
// optimized away by the compiler.
type cNoop C.Nothing

// TODO: Add vesa support and support for more than just printing text.

const (
	VGA_WIDTH  = 80
	VGA_HEIGHT = 25
)

const (
	COLOR_BLACK = iota
	COLOR_BLUE
	COLOR_GREEN
	COLOR_CYAN
	COLOR_RED
	COLOR_MAGENTA
	COLOR_BROWN
	COLOR_LIGHT_GREY
	COLOR_DARK_GREY
	COLOR_LIGHT_BLUE
	COLOR_LIGHT_GREEN
	COLOR_LIGHT_CYAN
	COLOR_LIGHT_RED
	COLOR_LIGHT_MAGENTA
	COLOR_LIGHT_BROWN
	COLOR_WHITE
)

type Color uint8
type VgaCharacter uint16
type Terminal struct {
	Row    uint16
	Column uint16
	Color  Color
	Buffer *uint16
}

var Term *Terminal

func InitPkg() {
	Term = &Terminal{}
	InitializeTerminal()
}
func MakeColor(fg, bg Color) Color {
	return fg | bg<<4
}

func MakeVgaEntry(chr byte, color Color) VgaCharacter {
	return VgaCharacter(chr) | VgaCharacter(color)<<8
}

func (t *Terminal) PutEntryAt(c byte, color Color, x, y uint16) {
	idx := y*VGA_WIDTH + x
	setbuffer(t, idx, MakeVgaEntry(c, color))
}

func (t *Terminal) PutChar(c byte) {
	switch c {
	case '\n':
		if t.Row < VGA_HEIGHT-1 {
			t.Row++
		} else {
			// Scroll everything up by 1 row
			for y := uint16(1); y <= VGA_HEIGHT; y++ {
				for x := uint16(0); x < VGA_WIDTH; x++ {
					idx := uint16(y*VGA_WIDTH + x)
					val := getbuffer(t, idx)
					setbuffer(t, idx-VGA_WIDTH, val)
					//setbuffer(t, idx-VGA_WIDTH, MakeVgaEntry('x', COLOR_BLUE))
				}
			}
			// Clear the last row.
			for x := uint16(0); x < VGA_WIDTH; x++ {
				t.PutEntryAt(' ', t.Color, x, VGA_HEIGHT-1)
			}
		}
		t.Column = 0
	case '\t':
		spaces := 8 - (t.Column % 8)
		t.Column += spaces
		if t.Column >= VGA_WIDTH {
			t.PutChar('\n')
		} else {
			for ; spaces > 0; spaces-- {
				t.PutEntryAt(' ', t.Color, t.Column-spaces, t.Row)
			}
		}
	default:

		if t.Column < VGA_WIDTH {
			t.PutEntryAt(c, t.Color, t.Column, t.Row)
			t.Column++
		}

		if t.Column >= VGA_WIDTH {
			if t.Row < VGA_HEIGHT-1 {
				t.Row++
			} else {
				// Scroll everything up by 1 row
				for y := uint16(1); y <= VGA_HEIGHT; y++ {
					for x := uint16(0); x < VGA_WIDTH; x++ {
						idx := uint16(y*VGA_WIDTH + x)
						val := getbuffer(t, idx)
						setbuffer(t, idx-VGA_WIDTH, val)
						//setbuffer(t, idx-VGA_WIDTH, MakeVgaEntry('x', COLOR_BLUE))
					}
				}
				// Clear the last row.
				for x := uint16(0); x < VGA_WIDTH; x++ {
					t.PutEntryAt(' ', t.Color, x, VGA_HEIGHT-1)
				}
			}
			t.Column = 0

		}
	}
}

//extern setbuffer
func setbuffer(t *Terminal, idx uint16, val VgaCharacter)

//extern getbuffer
func getbuffer(t *Terminal, idx uint16) VgaCharacter

func InitializeTerminal() {
	Term.Row = 0
	Term.Column = 0
	Term.Color = MakeColor(COLOR_LIGHT_GREY, COLOR_BLACK)
	Term.Buffer = (*uint16)(unsafe.Pointer(uintptr(0xB8000)))
	for y := uint16(0); y < VGA_HEIGHT; y++ {
		for x := uint16(0); x < VGA_WIDTH; x++ {
			Term.PutEntryAt(' ', Term.Color, x, y)
		}
	}
}
