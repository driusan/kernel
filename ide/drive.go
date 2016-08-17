package ide

import (
	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/interrupts"
)

const (
	ReadSectors = 0x20
)

/*
type IDEDrive struct{
	// Whether or not the drive is a master or a slave
	Drive DriveSelector

	MaxLBA28 uint32
	MaxLBA48 uint64
}*/

type DriveSector struct {
	Data   [256]uint16
	Sector uint64
}

var LastSectorRead DriveSector

func PrimaryDriveHandler(r *interrupts.Registers) {
	//print("Reading data")
	for i := 0; i < 256; i++ {
		LastSectorRead.Data[i] = asm.INW(PrimaryDataPort)
		/*

			if LastSectorRead.Data[i] != 0 {
				println("Word ", i, LastSectorRead.Data[i])
			}
				LastSectorRead[i*2] = byte((data >> 8) & 0xFF)
				LastSectorRead[(i*2) + 1] = byte((data >> 8) & 0xFF)
				LastSectorRead = */
	}

	//print("Read data")
	//print("Word 216 ", LastSectorRead.Data[216], " from ", LastSectorRead.Sector)
	//status := asm.INB(PrimaryStatus)
	asm.INB(PrimaryStatus)
	//print("Waiting for drive to not be busy")
	//for (status & busy) != 0 {
	//	status = asm.INB(PrimaryStatus)
	//}
	//print("Status after interrupt ", status)
}

// for some reason this doesn't work as a method, even though it should
// and methods work other places
//func (d IDEDrive) ReadLBA(lba uint64) error {
func ReadLBA(d IDEDrive, lba uint64) error {
	//print("Reading LBA ", lba)

	if d.Drive != PrimaryDrive {
		print("Drive ", d.Drive)
		return InvalidDrive
	}

	if lba > uint64(d.MaxLBA28) {
		println("Can not read sector")
		// 48 bit not yet supported
		return InvalidDrive
	}

	// assume the master drive. TODO: Handle slave drives
	asm.OUTB(PrimaryDriveSelect, 0xE0|byte((lba>>24)&0xF))
	asm.OUTB(PrimarySectorCount, 1)
	asm.OUTB(PrimaryLBAlo, byte(lba&0xFF))
	asm.OUTB(PrimaryLBAmid, byte((lba>>8)&0xFF))
	asm.OUTB(PrimaryLBAhi, byte((lba>>16)&0xFF))
	asm.OUTB(PrimaryCommand, ReadSectors)

	LastSectorRead.Sector = lba
	//println("Leaving ReadLBA ", lba)
	return nil
}
