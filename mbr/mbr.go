package mbr

import "unsafe"

type Partitions [4]PartitionTableEntry

type PartitionTableEntry struct {
	Active                  byte // 0 = no, 0x80 = yes
	StartingHead            byte
	StartSectorAndCylinder  uint16 // lower 10 bits are cylinder, upper are sector
	PartitionType           byte
	EndingHead              byte
	EndingSectorAndCylinder uint16 // same encoding as start
	LBAStart                uint32
	LBASize                 uint32
}

func (pte PartitionTableEntry) Type() string {
	switch pte.PartitionType {
	case 0:
		return "Unused"
	case 0x05:
		return "Extended"
	case 0xb, 0xc:
		return "FAT32"
	case 0x83:
		return "EXT2"
	default:
		print("Unknown", pte.PartitionType)
		return ""

	}
}
func ExtractPartitions(hdsector [512]byte) *Partitions {
	return (*Partitions)(unsafe.Pointer(&hdsector[446]))
}
