package pci

import (
	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/terminal"
	//"fmt"
)

type VendorID uint16
type DeviceID uint16

type Device struct {
	BusID    uint8
	Slot     uint8
	Function uint8

	Vendor VendorID
	Device DeviceID
}

type PCIError string

func (p PCIError) Error() string {
	return string(p)
}

var NoDevice PCIError
var Invalid PCIError

// global variable initializations don't work in freestanding mode, so we
// need to manually initialize the package with a function
func InitPkg() {
	NoDevice = PCIError("No such device")
	Invalid = PCIError("Invalid parameter")
}

// Prints all devices found on the PCI Bus.
// TODO: Dynamically create the slice instead of having a fixed
// array of the maximum size once Malloc/Free are implemented
func EnumerateDevices() { //[256 * 32]Device{
	var header HeaderType
	var class Class
	var err error

	multipleBuses := true

	for b := 0; b < 256; b++ {
		busDevices := EnumerateBus(uint8(b))

		for i, d := range busDevices {
			if d.Vendor != 0xFFFF {
				header, err = d.GetHeaderType()
				if err != nil {
					print(err.Error())
					continue
				}

				if header.IsMultifunction() {
					for f := uint8(0); f < 8; f++ {
						err := (&d).Probe(f)
						if err == nil {
							print(i, " ")
							terminal.PrintHex(uint64(d.Vendor))
							print(" ")
							terminal.PrintHex(uint64(d.Device))
							print(" ", f, " ")

							class, err = d.GetClass(f)
							print(class.String())

							print("\n")
						}

					}

				} else {
					print(i, " ")
					terminal.PrintHex(uint64(d.Vendor))
					print(" ")
					terminal.PrintHex(uint64(d.Device))

					class, err = d.GetClass(0)
					if err == nil {
						print(" ", class.String())
					}
					print("\n")

				}

			}
			// TODO: Make this smarter. It should check the class
			// to see if it's a PCI-to-PCI bridge. If we don't find
			// any, there's no more buses to scan.
			// BUG(driusan): This check is completely wrong and
			// doesn't do what it should.
			if i == 0 && b == 0 && !header.IsMultifunction() {
				multipleBuses = false
			}
		}
		if !multipleBuses {
			break
		}
	}
	//return devices
}

func EnumerateBus(busNum uint8) (devices [32]Device) {
	for device := uint8(0); device < 32; device++ {
		d := &devices[device]
		d.BusID = busNum
		d.Slot = device

		err := d.Probe(0)

		if err != nil {
			// This should have been handled by d.Probe()
			// but for some reason doesn't work, so needs
			// to be explicitly propagated.
			d.Vendor = 0xFFFF
		}
	}
	return devices
}

func (d Device) ReadWord(fnc, offset uint8) uint16 {
	var address uint32
	address = uint32(d.BusID)<<16 |
		uint32(d.Slot)<<11 |
		uint32(fnc)<<8 |
		uint32(offset&0xfc) | uint32(0x80000000)
	asm.OUTL(0xCF8, address)
	return (uint16(asm.INL(0xCFC) >> (uint16(offset&2) * 8)))
}

func (d *Device) Probe(fnc uint8) error {
	if fnc > 8 {
		return Invalid
	}
	vendor := d.ReadWord(fnc, 0)
	if vendor != 0xFFFF {
		d.Vendor = VendorID(vendor)
		d.Device = DeviceID(d.ReadWord(fnc, 2))
		return nil
	}
	d.Vendor = 0xFFFF
	d.Device = 0
	return NoDevice
}

func (d *Device) GetClass(fnc uint8) (Class, error) {
	if d.Vendor == 0xFFFF {
		return 0xFFFF, NoDevice
	}
	classAndSub := d.ReadWord(fnc, 10)
	return Class(classAndSub), nil
}

func (d Device) GetHeaderType() (HeaderType, error) {
	if d.Vendor == 0xFFFF {
		return 0, NoDevice
	}
	bistAndHeader := d.ReadWord(0, 14)
	return HeaderType(bistAndHeader & 0xFF), nil
}
