package fat

import (
	"github.com/driusan/kernel/filesystem"
	"github.com/driusan/kernel/ide"
)

// Names and extensions are space padded, so look for the
// first space. We should probably be using the long filename
// instead, but for now this is better than nothing.
func trimFatName(name []byte, ext []byte) filesystem.Path {
	for i, c := range name {
		switch c {
		case 0, ' ':
			name = name[0:i]
			goto namedone
		}
	}
namedone:
	for i, c := range ext {
		switch c {
		case 0, ' ':
			ext = ext[0:i]
			goto extdone
		}
	}
extdone:
	if len(ext) > 0 {
		return filesystem.Path(name) + "." + filesystem.Path(ext)
	}
	return filesystem.Path(name)
}

type fat32File struct {
	ClusterStart uint32
	FileSize     uint32
	IsDir        bool

	fs      *FatFS
	seekPos int64
}

func (f *fat32File) Read(buf []byte) (n int, err error) {
	if len(buf) <= 0 {
		// If the buffer size is 0, we can't do anything.
		return 0, filesystem.FilesystemError("Can not Read into nil buffer")
	}
	if f.seekPos >= int64(f.FileSize) {
		// The reader should have stopped at the first EOF, but if they
		// didn't just keep giving them EOFs
		return 0, filesystem.EOF
	}

	bytesRead := 0
	clusters, err := f.fs.getClusterChain(f.ClusterStart)

	clusterPos := int64(0)
	clusterSize := int64(f.fs.biosParameterBlock.SectorsPerCluster) * int64(f.fs.biosParameterBlock.BytesPerSector)
	startedReading := false
	for _, c := range clusters {
		// Skip over clusters that are less than seekPos, we read forward
		// not backwards
		if clusterPos+clusterSize < f.seekPos {
			clusterPos += clusterSize
			continue
		}

		sectorPos := clusterPos
		sectorSize := int64(f.fs.biosParameterBlock.BytesPerSector)
		startPos := int64(0)
		for s := uint8(0); s < f.fs.biosParameterBlock.SectorsPerCluster; s++ {
			// Skip over sectors that are less than seekPos.
			// This isn't very cleanly written, it should be cleaned up
			if sectorPos+sectorSize < f.seekPos {
				sectorPos += sectorSize
				continue
			} else {
				// Calculate the index into the sector that matches
				// seekPos
				if !startedReading {
					startPos = f.seekPos - sectorPos
					startedReading = true
				} else {
					// we've already started reading, so
					// start from the beginning of the sector
					startPos = 0
				}
			}

			// Read the sector from the hard drive
			sector, err := ide.ReadLBA(f.fs.Drive, f.fs.ClusterToLBA(uint64(c))+uint64(s))
			if err != nil {
				return bytesRead, err
			}

			// Copy from startPos into the appropriate place in the buffer
			copy(buf[bytesRead:], sector.Data[startPos:])
			bytesRead += len(sector.Data[startPos:])
			if bytesRead >= len(buf) {
				// We've filled up the buffer
				f.seekPos += int64(len(buf))
				return len(buf), nil
			} else if f.seekPos+int64(bytesRead) >= int64(f.FileSize) {
				// We've read the entire file
				initSeek := f.seekPos
				f.seekPos = int64(f.FileSize)
				return int(f.FileSize) - int(initSeek), filesystem.EOF
			}
		}
	}

	f.seekPos += int64(bytesRead)
	return bytesRead, nil
}

func (f fat32File) Write([]byte) (n int, err error) {
	return 0, filesystem.FilesystemError("Writing files not yet implemented")
}

func (f fat32File) ReadByte() (byte, error) {
	return 0, filesystem.FilesystemError("ReadByte not yet implemented for FAT")
}

func (f fat32File) WriteByte(byte) error {
	return filesystem.FilesystemError("Writing files not yet implemented")
}
func (f fat32File) WriteRune(rune) error {
	return filesystem.FilesystemError("Writing files not yet implemented")
}
func (f fat32File) Seek(offset int64, whence int) (int64, error) {
	var attemptedFinal int64
	switch whence {
	case 0: // io.SeekStart
		attemptedFinal = offset
	case 1: // io.SeekCurrent
		attemptedFinal = int64(f.seekPos) + offset
	case 2: // io.SeekEnd
		attemptedFinal = int64(f.FileSize) + offset
	default:
		return 0, filesystem.FilesystemError("Invalid seek usage")
	}
	if attemptedFinal < 0 {
		return 0, filesystem.FilesystemError("Attempt to Seek before start of file")
	} else if attemptedFinal > int64(f.FileSize) {
		f.seekPos = int64(f.FileSize)
	} else {
		f.seekPos = attemptedFinal
	}
	return f.seekPos, nil
}

func (f fat32File) AsFATDirectory() (fat32Directory, error) {
	if f.IsDir {
		return f.fs.readDir(f.ClusterStart)
	} else {
		return nil, filesystem.FilesystemError("File is not a directory.")
	}
}
func (f fat32File) AsDirectory() (filesystem.Directory, error) {
	if f.IsDir {
		return f.fs.readDir(f.ClusterStart)
	} else {
		return nil, filesystem.FilesystemError("File is not a directory.")
	}
}

func (f fat32File) Close() error {
	return nil
}

func (f fat32File) IsDirectory() bool {
	return f.IsDir
}
