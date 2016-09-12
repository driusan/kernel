package ide

import "github.com/driusan/kernel/asm"

const (
	PrimaryDataPort = 0x1F0 + iota
	PrimaryFeaturesPort
	PrimarySectorCount
	PrimaryLBAlo
	PrimaryLBAmid
	PrimaryLBAhi
	PrimaryDriveSelect
	PrimaryCommand
	PrimaryStatus = PrimaryCommand
)

const (
	PrimaryControl         = 0x3F6
	PrimaryAlternateStatus = PrimaryControl
)

type DriveSelector byte

const (
	PrimaryDrive   = DriveSelector(0xA0)
	SecondaryDrive = DriveSelector(0xB0)
)

type DriveStatus byte

type IDEError string

func (i IDEError) Error() string {
	return string(i)
}

const (
	Identify = 0xEC
)

var InvalidDrive IDEError

func InitPkg() {
	InvalidDrive = IDEError("Invalid drive")
	lastSectorRead = DriveSector{}
}
func SelectDrive(drive DriveSelector) (DriveStatus, error) {
	if drive != PrimaryDrive && drive != SecondaryDrive {
		return 0, InvalidDrive
	}

	asm.OUTB(PrimaryDriveSelect, byte(drive))

	// 4 reads to give the drive time to warm up, one
	// read to get the value.
	asm.INB(PrimaryAlternateStatus)
	asm.INB(PrimaryAlternateStatus)
	asm.INB(PrimaryAlternateStatus)
	asm.INB(PrimaryAlternateStatus)

	status := asm.INB(PrimaryAlternateStatus)
	return DriveStatus(status), nil
}

// Masks for bits in the status byte
const (
	errstatus = 0x00

	dataready = 0x08
	busy      = 0x80
)

type IDEDrive struct {
	// Whether or not the drive is a master or a slave
	Drive DriveSelector

	MaxLBA28 uint32
	MaxLBA48 uint64
}

func IdentifyDrive(drive DriveSelector) (IDEDrive, error) {
	// Soft reset of the drive
	//asm.OUTB(PrimaryControl, 0x04)

	/*
		// Select the drive and wait 400ns
		asm.OUTB(PrimaryDriveSelect, 0xA0)
		asm.INB(PrimaryControl)
		asm.INB(PrimaryControl)
		asm.INB(PrimaryControl)
		asm.INB(PrimaryControl)

		cl := asm.INB(0x1F4)
		ch := asm.INB(0x1F5)
		print("Drive Signature ", cl, " ", ch)
	*/

	_, err := SelectDrive(drive)
	if err != nil {
		return IDEDrive{}, err
	}
	// disable interrupts from this drive
	/*
		asm.OUTB(PrimaryControl, 0x02)
		// wait 400ns
		asm.INB(PrimaryAlternateStatus)
		asm.INB(PrimaryAlternateStatus)
		asm.INB(PrimaryAlternateStatus)
		asm.INB(PrimaryAlternateStatus)
	*/

	asm.OUTB(PrimarySectorCount, 0)
	asm.OUTB(PrimaryLBAlo, 0)
	asm.OUTB(PrimaryLBAmid, 0)
	asm.OUTB(PrimaryLBAhi, 0)

	asm.OUTB(PrimaryCommand, Identify)

	status := asm.INB(PrimaryStatus)

	/*
		print("Status ", status)
		if status == 0 {
			return IDEDrive{}, InvalidDrive
		}*/

	for (status & busy) != 0 {
		status = asm.INB(PrimaryStatus)
	}

	lba := asm.INB(PrimaryLBAmid)
	if lba != 0 {
		return IDEDrive{}, InvalidDrive
	}
	lba = asm.INB(PrimaryLBAhi)
	if lba != 0 {
		return IDEDrive{}, InvalidDrive
	}

	for (status&dataready == 0) && (status&errstatus == 0) {
		status = asm.INB(PrimaryStatus)

	}

	var data uint16
	var maxLBA28 uint32
	var maxLBA48 uint64
	var noLBA48 bool
	var wordsPerLogicalSector uint32
	for i := 0; i < 256; i++ {
		data = asm.INW(PrimaryDataPort)
		switch i {
		case 83:
			if (data & (1 << 10)) != 0 {
				//println("LBA48 mode supported")
				noLBA48 = false
			} else {
				noLBA48 = true
			}
		case 88:
			//println("Active UDMA Mode", (data &0xFF00) >> 8, " Supported UDMA modes", data & 0x00FF)
		case 60:
			maxLBA28 |= uint32(data)
		case 61:
			maxLBA28 |= uint32(data) << 16
		case 100:
			maxLBA48 |= uint64(data)
		case 101:
			maxLBA48 |= uint64(data) << 16
		case 102:
			maxLBA48 |= uint64(data) << 32
		case 103:
			maxLBA48 |= uint64(data) << 48
		case 106:
			//println("Word 106 in identify device ", data)
			if (data & (1 << 12)) != 0 {
				println("WARNING: Drive block size is not 512 bytes")
			}
		case 117:
			wordsPerLogicalSector |= uint32(data)
		case 118:
			wordsPerLogicalSector |= uint32(data) << 16
		}
	}

	//println("Words per logical sector:", wordsPerLogicalSector)
	d := IDEDrive{
		Drive:    drive,
		MaxLBA28: maxLBA28,
		MaxLBA48: maxLBA48,
	}
	if noLBA48 {
		d.MaxLBA48 = 0
	}

	return d, nil
}
