package plan9

const (
	Magic386 = ((((4 * 11) + 0) * 11) + 7)
)

// Simple helper to ensure data is interpreted in big endian format,
// regardless of architecture.
type BEuint32 [4]byte

func (b BEuint32) Uint32() uint32 {
	return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | uint32(b[3])

}

type ExecHeader struct {
	Magic           BEuint32
	Text, Data, Bss BEuint32
	Syms            BEuint32
	EntryPoint      BEuint32
	SPsz            BEuint32
	PCsz            BEuint32
}
