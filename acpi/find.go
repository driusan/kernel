package acpi

import "unsafe"

var NoACPIFound, InvalidChecksum error

type ACPIError string
func (e ACPIError) Error() string {
	return string(e)
}
func InitPkg() {
	NoACPIFound = ACPIError("No ACPI header found")
	InvalidChecksum = ACPIError("Invalid checksum")
}

type RSDPtr struct {
	Signature   [8]byte
	Checksum    uint8
	OEMID       [6]byte
	Revision    uint8

	// The RSDT is stored as a 32-bit pointer according to the ACPI spec.
	// This is unexported and not a pointer so that it will load and
	// serialize correctly regardless of the pointer size of the system,
	// ie we can't be sure that the compiler won't decide that a *RSDT is
	// a 64-bit pointer.
	// The GetRSDT method converts the int to 
	rsdtAddress uint32
}

// BUG: This is just a stub right now.
type RSDT struct {
	Signature [4]byte
}
func (r RSDPtr) GetRSDT() (*RSDT, error) {
	return (*RSDT)(unsafe.Pointer(uintptr(r.rsdtAddress))), nil
}

func FindRSDP() (*RSDPtr, error) {
	var Desc *RSDPtr //[8]byte
	for addr := 0xE0000; addr < 0xFFFFF; addr += 16 {
		Desc = ((*RSDPtr)(unsafe.Pointer(uintptr(addr))))
		if Desc.Signature[0] == 'R' && Desc.Signature[1] == 'S' && Desc.Signature[2] == 'D' && Desc.Signature[3] == ' ' &&
			Desc.Signature[4] == 'P' && Desc.Signature[5] == 'T' && Desc.Signature[6] == 'R' && Desc.Signature[7] == ' ' {
			break
		} else {
			Desc = nil
		}
	}
	if Desc == nil {
		return nil, NoACPIFound
	}

	// TODO: Validate checksum here
	return Desc, nil
}
