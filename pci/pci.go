package pci

import "asm"

//extern printhex
func printhex(int64)

type Device struct {
	BusID uint8
	DevID uint8
}

// Prints all devices found on the PCI Bus.
// TODO: Return a struct representing them instead.
// Malloc/Free are required for returning a slice.
func EnumerateDevices() { //[]Device {
	var bus, device uint8
	//var devices []Device

	for bus = 0; bus <= 255; bus++ {
		for device = 0; device < 32; device++ {
			d := checkDevice(bus, device)
			if d != 0 {
				printhex(int64(d))
				print("\n")
			}
		}

		// without this, the counter will overflow and get into
		// an infinite loop, but if we change the comparison to <
		// we could technically miss the last device
		if bus == 255 {
			break
		}
	}
	//return devices
}

func checkDevice(bus, device uint8) uint32 {
	vendorID := getVendorID(bus, device)
	if vendorID == 0xFFFF {
		return 0
	}
	return uint32(vendorID)
}

func PCIConfigReadRegister(bus, slot, fnc, offset uint8) uint32 {
	var address uint32
	address = uint32(bus)<<16 |
		uint32(slot)<<11 |
		uint32(fnc)<<8 |
		uint32(offset&0xfc) | uint32(0x80000000)
	asm.OUTL(0xCF8, address)
	return (uint32(asm.INL(0xCFC) >> (uint32(offset&2) * 8)))
}

func getVendorID(bus, slot uint8) uint32 {
	if vendor := PCIConfigReadRegister(bus, slot, 0, 0); vendor != 0xFFFFFFFF {
		return vendor
	}
	return 0
}
