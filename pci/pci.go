package pci

import (
	"asm"
	//"fmt"
)

type VendorID uint16
type DeviceID uint16

//extern printhex
func printhex(int64)

type Device struct {
	BusID uint8
	Slot  uint8

	Vendor VendorID
	Device DeviceID
}

type PCIError string

func (p PCIError) Error() string {
	return string(p)
}

var NoDevice PCIError

// global variable initializations don't work in freestanding mode, so we
// need to manually initialize the package with a function
func InitPkg() {
	NoDevice = PCIError("No such device")
}

// Prints all devices found on the PCI Bus.
// TODO: Dynamically create the slice instead of having a fixed
// array of the maximum size once Malloc/Free are implemented
func EnumerateDevices() { //[256 * 32]Device{
	//var devices [256*32]Device
	for b := 0; b < 256; b++ {
		busDevices := EnumerateBus(uint8(b))

		//println("Finished bus ")
		//printhex(int64(bus))
		//print("\n")
		for i, d := range busDevices {
			if d.Vendor != 0xFFFF {
				println(b, i, d.Vendor)
			}
			//devices[(bus*32)+i] = d
		}
	}
	//return devices
}

func EnumerateBus(busNum uint8) (devices [32]Device) {
	for device := uint8(0); device < 32; device++ {
		d := &devices[device]
		d.BusID = busNum
		d.Slot = device

		err := d.Probe()

		if err == nil {
			print(device, " ")
			printhex(int64(devices[device].Vendor))
			print(" ")
			printhex(int64(devices[device].Device))
			print("\n")
		} else {
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

func (d *Device) Probe() error {
	vendor := d.ReadWord(0, 0)
	if vendor != 0xFFFF {
		d.Vendor = VendorID(vendor)
		d.Device = DeviceID(d.ReadWord(0, 2))
		//(vendor >> 8) & 0xFFFF)
		return nil
	}
	d.Vendor = 0xFFFF
	d.Device = 0
	return NoDevice
}
