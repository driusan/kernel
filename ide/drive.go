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
	Data   [512]byte
	Sector uint64
}

var LastSectorRead DriveSector

func PrimaryDriveHandler(r *interrupts.Registers) {
	//print("Reading data")
	var data uint16
	for i := 0; i < 256; i++ {
		data = asm.INW(PrimaryDataPort)
		LastSectorRead.Data[(2*i)+1] = byte((data >> 8) & 0xFF)
		LastSectorRead.Data[(2*i)+0] = byte(data & 0xFF)
	}
	asm.INB(PrimaryStatus)
}

// for some reason this doesn't work as a method, even though it should
// and methods work other places
//func (d IDEDrive) ReadLBA(lba uint64) error {
func ReadLBA(d IDEDrive, lba uint64) (DriveSector, error) {
	//print("Reading LBA ", lba)

	if d.Drive != PrimaryDrive {
		print("Drive ", d.Drive)
		return DriveSector{}, InvalidDrive
	}

	if lba > uint64(d.MaxLBA28) {
		println("Can not read sector")
		// 48 bit not yet supported
		return DriveSector{}, InvalidDrive
	}

	// assume the master drive. TODO: Handle slave drives
	asm.OUTB(PrimaryDriveSelect, 0xE0|byte((lba>>24)&0xF))
	asm.OUTB(PrimarySectorCount, 1)
	asm.OUTB(PrimaryLBAlo, byte(lba&0xFF))
	asm.OUTB(PrimaryLBAmid, byte((lba>>8)&0xFF))
	asm.OUTB(PrimaryLBAhi, byte((lba>>16)&0xFF))
	asm.OUTB(PrimaryCommand, ReadSectors)

	LastSectorRead.Sector = lba
	asm.HLT()
	// TODO: Verify that it was the right type of interrupt that woke
	// us up
	return LastSectorRead, nil
	//println("Leaving ReadLBA ", lba)
	//return nil
}
