package fat

import (
	"github.com/driusan/kernel/filesystem"
	"github.com/driusan/kernel/ide"
)

var Fat filesystem.Filesystem

type FatFS struct {
	LBAStart uint64
	LBASize  uint64

	// This should be an interface that isn't tied to IDE, but for now
	// that's all we've got
	Drive              ide.IDEDrive
	biosParameterBlock fatBPB

	// Only one of these should be populated, the other nil
	fat16    *fatEBR
	fat32EBR *fat32EBR
}

// Gets the first useable cluster in a FAT filesystem.
func (f FatFS) FirstUseableCluster() uint64 {
	bpb := &f.biosParameterBlock
	if bpb.TotalSectors != 0 {
		// FAT12/16, use SectorsPerFAT from BPB
		return f.LBAStart + uint64(bpb.ReservedSectorCount) + (uint64(bpb.FileAllocationTableCount) * uint64(bpb.SectorsPerFat))
	} else {
		// FAT32, use SectorsPerFAT from EBR
		return f.LBAStart + uint64(bpb.ReservedSectorCount) + (uint64(bpb.FileAllocationTableCount) * uint64(f.fat32EBR.SectorsPerFAT))
	}
}

func (f FatFS) ClusterToLBA(c uint64) uint64 {
	return ((uint64(c) - 2) * uint64(f.biosParameterBlock.SectorsPerCluster)) + f.FirstUseableCluster() //first_data_sector;
}

// Common Partition Boot Block for all FAT subtypes
// All values must be stored and serialized in little endian format.
type fatBPB struct {
	JumpByte                 [3]byte
	OEM                      [8]byte
	BytesPerSector           uint16
	SectorsPerCluster        uint8
	ReservedSectorCount      uint16
	FileAllocationTableCount uint8
	DirectoryEntriesCount    uint16
	TotalSectors             uint16 // 0 means > 65536. Use Bytes 32-35 instead
	MediaDescriptorType      uint8
	SectorsPerFat            uint16 // FAT12/FAT16 only
	SectorsPerTrack          uint16
	StorageHeadersCount      uint16
	HiddenSectorsCount       uint32
	LargeSectorCount         uint32
}

// Extended Boot Record for FAT12 and FAT16
// All values must be stored and serialized in little endian format.
type fatEBR struct {
	DriveNumber      uint8
	Flags            uint8
	Signature        uint8
	VolumeID         [4]byte
	VolumeLabel      [11]byte
	SystemIdentifier [8]byte
	_                [448]byte
	BootableSig      [2]byte // should be 0xAA55

}

// Extended Boot Record for FAT32
// All values must be stored and serialized in little endian format.
type fat32EBR struct {
	SectorsPerFAT    uint32
	Flags            uint16
	FatVersion       [2]byte
	RootCluster      uint32
	FSInfoSector     uint16
	BackupBootSector uint16
	_                [12]byte
	DriveNumber      uint8
	NTFlags          uint8
	Signature        uint8
	VolumeID         [4]byte
	VolumeLabel      [11]byte
	SystemIdentifier [8]byte
	_                [420]byte
	BootableSig      uint16 // should be 0xAA55
}

type fat32DirectoryEntry struct {
	Name       [8]byte
	Ext        [3]byte
	Attrib     byte
	UserAttrib byte

	Undelete    byte
	CreateTime  uint16
	CreateDate  uint16
	AccessDate  uint16
	ClusterHigh uint16

	ModifiedTime uint16
	ModifiedDate uint16

	ClusterLow uint16
	FileSize   uint32

	// Below this line in the struct isn't part of the FAT32
	// directory structure definition, but is used internally
	fs *FatFS
	// long file name, if applicable
	lfn string
}

func (f FatFS) Type() string {
	return "FAT Filesystem"
}

func (f *FatFS) Initialize() error {
	firstPartitionBlock, err := ide.ReadLBA(f.Drive, f.LBAStart)
	if err != nil {
		return err
	}
	data := firstPartitionBlock.Data
	bpb := &f.biosParameterBlock
	for i := 0; i < 3; i++ {
		bpb.JumpByte[i] = data[i]

	}
	for i := 0; i < 8; i++ {
		bpb.OEM[i] = data[3+i]

	}
	bpb.BytesPerSector = (uint16(data[12]) << 8) | uint16(data[11])
	bpb.SectorsPerCluster = data[13]
	bpb.ReservedSectorCount = uint16(data[15])<<8 | uint16(data[14])
	bpb.FileAllocationTableCount = data[16]
	bpb.DirectoryEntriesCount = uint16(data[18])<<8 | uint16(data[17])
	bpb.TotalSectors = uint16(data[20])<<8 | uint16(data[19])
	bpb.MediaDescriptorType = uint8(data[21])
	bpb.SectorsPerFat = uint16(data[23])<<8 | uint16(data[22])
	bpb.SectorsPerTrack = uint16(data[25])<<8 | uint16(data[24])
	bpb.StorageHeadersCount = uint16(data[27])<<8 | uint16(data[26])
	bpb.HiddenSectorsCount = uint32(data[31])<<24 | uint32(data[30])<<16 | uint32(data[29])<<8 | uint32(data[28])
	bpb.LargeSectorCount = uint32(data[35])<<24 | uint32(data[34])<<16 | uint32(data[33])<<8 | uint32(data[32])

	if bpb.TotalSectors != 0 {
		return filesystem.FilesystemError("FAT12 and FAT16 not implemented, only FAT32")
	} else {
		f.fat32EBR = &fat32EBR{
			SectorsPerFAT:    uint32(data[39])<<24 | uint32(data[38])<<16 | uint32(data[37])<<8 | uint32(data[36]),
			Flags:            uint16(data[41])<<8 | uint16(data[40]),
			FatVersion:       [2]byte{data[42], data[43]},
			RootCluster:      uint32(data[47])<<24 | uint32(data[46])<<16 | uint32(data[45])<<8 | uint32(data[44]),
			FSInfoSector:     uint16(data[49])<<8 | uint16(data[48]),
			BackupBootSector: uint16(data[51])<<8 | uint16(data[50]),
			DriveNumber:      data[64],
			NTFlags:          data[65],
			Signature:        data[66],
			VolumeID:         [4]byte{data[67], data[68], data[69], data[70]},
			VolumeLabel: [11]byte{
				data[71],
				data[72], data[73], data[74], data[75],
				data[76], data[77], data[78], data[79], data[80], data[81],
			},
			SystemIdentifier: [8]byte{data[82], data[83], data[84], data[85], data[86], data[87], data[88], data[89]},
			BootableSig:      uint16(data[511])<<8 | uint16(data[510]),
		}
	}

	return nil
}

func (f FatFS) Open(name filesystem.Path) (filesystem.File, error) {
	curDir, err := f.readDir(f.fat32EBR.RootCluster)

	switch name {
	case "", "/":
		return curDir, err
	}

	pieces := filesystem.SplitPath(name)
	for i, p := range pieces {
		filePiece, err := curDir.OpenFile(p)
		if err != nil {
			return nil, err
		}
		if i == len(pieces)-1 {
			return &filePiece, nil
		}

		curDir, err = filePiece.AsFATDirectory()
		if err != nil {
			return nil, err
		}
	}
	return nil, filesystem.FilesystemError("File not found")
}

// The LFN comes before the file entry and may span multiple entries, so
// we need to keep track of what the file name constructed so far is.
var lastLFN string

func (f FatFS) readDir(clusterStart uint32) (files fat32Directory, err error) {
	clusters, err := f.getClusterChain(uint32(clusterStart))
	if err != nil {
		return nil, err
	}

	// BUG(driusan): append() currently breaks if it needs to call memmove.
	// This should be debugged and fixed. For now, we just allocate a 32
	// entry slice with nothing in it so that it doesn't crash on the first
	// append. This will still crash on directories with more than 32 entries
	// unless the memmove problem is fixed.
	files = make([]fat32DirectoryEntry, 0, 32)

	for _, c := range clusters {
		for s := uint64(0); s < uint64(f.biosParameterBlock.SectorsPerCluster); s++ {
			sector, err := ide.ReadLBA(f.Drive, f.ClusterToLBA(uint64(c))+uint64(s))
			if err != nil {
				return files, err
			}
			for i := 0; i < 512; i += 32 {
				switch sector.Data[i] {
				case 0:
					return files, err
				case 0xE5:
					continue
				default:
					// This could probably done more efficiently with an unsafe.Pointer,
					// but I'm not confident enough in the runtime as implemented in
					// the kernel to not get corrupted by using unsafe, or to deal
					// with the endianness appropriately.

					file := fat32DirectoryEntry{}
					for j := 0; j < 8; j++ {
						file.Name[j] = sector.Data[i+j]
					}
					for j := 0; j < 3; j++ {
						file.Ext[j] = sector.Data[i+j+8]
					}
					file.Attrib = sector.Data[i+11]
					file.UserAttrib = sector.Data[i+12]
					file.Undelete = sector.Data[i+13]
					file.CreateTime = uint16(sector.Data[i+15])<<8 | uint16(sector.Data[i+14])
					file.CreateDate = uint16(sector.Data[i+17])<<8 | uint16(sector.Data[i+16])
					file.AccessDate = uint16(sector.Data[i+19])<<8 | uint16(sector.Data[i+18])
					file.ClusterHigh = uint16(sector.Data[i+21])<<8 | uint16(sector.Data[i+20])
					file.ModifiedTime = uint16(sector.Data[i+23])<<8 | uint16(sector.Data[i+22])
					file.ModifiedDate = uint16(sector.Data[i+25])<<8 | uint16(sector.Data[i+24])
					file.ClusterLow = uint16(sector.Data[i+27])<<8 | uint16(sector.Data[i+26])
					file.FileSize = uint32(sector.Data[i+31])<<24 | uint32(sector.Data[i+30])<<16 | uint32(sector.Data[i+29])<<8 | uint32(sector.Data[i+28])
					if file.Attrib == 0x0F {
						name, err := LongFileName(sector.Data[i : i+32]).Decode()
						if err != nil {
							print(err.Error())
							continue
							//return files, err
						}
						lastLFN = name + lastLFN
						// skip long file name entries for now
						continue
					}

					if lastLFN != "" {
						file.lfn = lastLFN
						lastLFN = ""
					}
					file.fs = &f
					files = append(files, file)
				}

			}
		}
	}
	return
}

func (f FatFS) getClusterChain(firstcluster uint32) ([]uint32, error) {
	clusters := make([]uint32, 0, 32)
	for cluster := firstcluster; ; {
		FatSector := f.LBAStart + uint64(f.biosParameterBlock.ReservedSectorCount) + uint64((cluster*4)/512)
		FatOffset := (cluster * 4) % 512

		sector, err := ide.ReadLBA(f.Drive, FatSector)
		if err != nil {
			return nil, err
		}

		cchain := uint32(sector.Data[FatOffset+3])<<24 | uint32(sector.Data[FatOffset+2])<<16 | uint32(sector.Data[FatOffset+1])<<8 | uint32(sector.Data[FatOffset])&0x0FFFFFFF

		clusters = append(clusters, cluster)
		if cchain == 0 || ((cchain & 0x0FFFFFFF) >= 0x0FFFFFF8) {
			break
		}
		cluster = cchain
	}
	return clusters, nil
}
