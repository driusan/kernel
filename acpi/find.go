// Package acpi handles parsing ACPI headers upon boot.
package acpi

import "unsafe"

var NoACPIFound, InvalidChecksum error

type ACPIError string

func (e ACPIError) Error() string {
	return string(e)
}

// FIXME: This calls __go_new before it's available because paging
// hasn't been initialized yet.
func InitPkg() {
	//NoACPIFound = ACPIError("No ACPI header found")
	//InvalidChecksum = ACPIError("Invalid checksum")
}

type RSDPtr struct {
	Signature [8]byte
	Checksum  uint8
	OEMID     [6]byte
	Revision  uint8

	// The RSDT is stored as a 32-bit pointer according to the ACPI spec.
	// This is unexported and not a pointer so that it will load and
	// serialize correctly regardless of the pointer size of the system,
	// ie we can't be sure that the compiler won't decide that a *RSDT is
	// a 64-bit pointer.
	// The GetRSDT method converts the int to a *RSDT
	rsdtAddress uint32
}

// BUG: This is just a stub right now.
type RSDT struct {
	Signature [4]byte
}

// Converts a 32 bit int found in the RSTP PTR header to a a pointer to
// an RSDT of the system-applicable size.
func (r RSDPtr) GetRSDT() (*RSDT, error) {
	return (*RSDT)(unsafe.Pointer(uintptr(r.rsdtAddress))), nil
}

// Find RSDP searches the memory for a valid RSDP header and returns a pointer
// to the RSDP Ptr defined by the ACPI spec.
func FindRSDP() (*RSDPtr, error) {
	var Desc *RSDPtr //[8]byte
	for addr := uint64(0xC00E0000); addr < 0xC00FFFFF; addr += 16 {
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
