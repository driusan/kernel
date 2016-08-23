package fat

import (
	"github.com/driusan/kernel/filesystem"
)

type fat32Directory []fat32DirectoryEntry

func (f fat32Directory) Read(output []byte) (n int, err error) {
	return 0, filesystem.FilesystemError("Directory should not be opened as a file.")
}
func (f fat32Directory) Write([]byte) (n int, err error) {
	return 0, filesystem.FilesystemError("Directory should not be opened as a file.")
}

func (f fat32Directory) ReadByte() (byte, error) {
	return 0, filesystem.FilesystemError("Directory should not be opened as a file.")
}

func (f fat32Directory) WriteByte(byte) error {
	return filesystem.FilesystemError("Directory should not be opened as a file.")
}
func (f fat32Directory) WriteRune(rune) error {
	return filesystem.FilesystemError("Directory should not be opened as a file.")
}

func (f fat32Directory) AsDirectory() (filesystem.Directory, error) {
	return f, nil
}

func (f fat32Directory) IsDirectory() bool {
	return true
}

func (f fat32Directory) Close() error {
	return nil
}
func (f fat32Directory) Seek(offset int64, whence int) (int64, error) {
	return 0, filesystem.FilesystemError("Directory should not be opened as a file.")
}

func (dir fat32Directory) OpenFile(name filesystem.Path) (fat32File, error) {
	for _, f := range dir {
		var sname filesystem.Path
		if f.lfn == "" {
			sname = trimFatName(f.Name[:], f.Ext[:])
		} else {
			sname = filesystem.Path(f.lfn)
		}
		if sname == name {
			file := fat32File{
				ClusterStart: uint32(f.ClusterHigh)<<16 | uint32(f.ClusterLow),
				FileSize:     f.FileSize,
				fs:           f.fs,
			}
			if f.Attrib&0x10 != 0 {
				file.IsDir = true
			}

			return file, nil
		}
	}
	return fat32File{}, filesystem.FilesystemError("File not found")
}

func (dir fat32Directory) Files() map[string]filesystem.File {
	m := make(map[string]filesystem.File, 0)

	for _, f := range dir {
		var sname string
		if f.lfn == "" {
			sname = string(trimFatName(f.Name[:], f.Ext[:]))
		} else {
			sname = f.lfn
		}
		if len(sname) > 0 {
			file := fat32File{
				ClusterStart: uint32(f.ClusterHigh)<<16 | uint32(f.ClusterLow),
				FileSize:     f.FileSize,
				fs:           f.fs,
			}
			if f.Attrib&0x10 != 0 {
				file.IsDir = true
			}
			m[sname] = &file
		}
	}
	return m
}
