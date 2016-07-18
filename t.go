package kernel

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

type Terminal struct {
	Row    uint16
	Column uint16
	Color  uint8
	Buffer *uint16
}

//extern make_color
func MakeColor(fg, bg uint8) uint8

//extern make_vgaentry
func make_vgaentry(chr uint8, color uint8) uint16

//extern reset_terminal_buffer
func resetbuffer(t *Terminal)

func InitializeTerminal(t *Terminal) {
	t.Row = 0
	t.Column = 0
	t.Color = MakeColor(COLOR_LIGHT_GREY, COLOR_BLACK)
	resetbuffer(t)

}
