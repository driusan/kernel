package fat

import (
	"unicode/utf16"

	"github.com/driusan/kernel/filesystem"
)

type LongFileName []byte

func (lfn LongFileName) Decode() (name string, err error) {
	if len(lfn) != 32 {
		return "", filesystem.FilesystemError("Invalid usage")
	}
	if lfn[11] != 0x0F {
		return "", filesystem.FilesystemError("Not a long file name entry")
	}
	if lfn[26] != 0 || lfn[27] != 0 {
		return "", filesystem.FilesystemError("Not a long file name entry")
	}

	// TODO: This should be done 2 uint16s at a time, so that it handles
	// surrogate pairs
	var chr [1]uint16
	for i := 1; i < 11; i += 2 {
		chr[0] = uint16(lfn[i+1])<<8 | uint16(lfn[i])
		if chr[0] == 0 {
			return
		}
		asRune := utf16.Decode(chr[:])
		name += string(asRune[0])
	}
	for i := 14; i < 26; i += 2 {
		chr[0] = uint16(lfn[i+1])<<8 | uint16(lfn[i])
		if chr[0] == 0 || chr[0] == 0xFF {
			return
		}
		asRune := utf16.Decode(chr[:])
		name += string(asRune[0])
	}
	for i := 28; i < 32; i += 2 {
		chr[0] = uint16(lfn[i+1])<<8 | uint16(lfn[i])
		if chr[0] == 0 || chr[0] == 0xFF {
			return
		}
		asRune := utf16.Decode(chr[:])
		name += string(asRune[0])
	}
	return
}
