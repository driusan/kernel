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

func ExtractPartitions(hdsector [512]byte) *Partitions {
	return (*Partitions)(unsafe.Pointer(&hdsector[446]))
}
