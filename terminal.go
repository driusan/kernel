package kernel

// TODO: Port this to another package, and add vesa support and support for
// more than just printing text.
import "unsafe"

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

func init() {
	Term = &Terminal{}
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
	if c == '\n' {
		t.Column = 0
		if t.Row < VGA_HEIGHT-1 {
			t.Row++
		} else {
			/* scroll everything up 1 row */
			for y := uint8(1); y < VGA_HEIGHT; y++ {
				for x := uint8(0); x < VGA_WIDTH; x++ {
					var idx uint16 = uint16(y*VGA_WIDTH + x)
					setbuffer(t, idx-VGA_WIDTH, getbuffer(t, idx))
				}
			}
			/* clear the last row. */
			for x := uint16(0); x < VGA_WIDTH; x++ {
				t.PutEntryAt(' ', t.Color, x, VGA_HEIGHT-1)
			}
		}
		return
	}

	t.PutEntryAt(c, t.Color, t.Column, t.Row)
	t.Column++
	if t.Column >= VGA_WIDTH {
		t.Column = 0
		t.Row++
		if t.Row >= VGA_HEIGHT {
			t.Row = 0
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