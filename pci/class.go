package pci

type Class uint16

func (c Class) GetClass() byte {
	return byte((c >> 8) & 0xFF)
}
func (c Class) GetSubclass() byte {
	return byte(c & 0xFF)
}

func (c Class) getBridgeClass() string {
	subclass := c.GetSubclass()
	switch subclass {
	case 0x00:
		return "Host/PCI bridge"
	case 0x01:
		return "PCI/ISA bridge"
	case 0x02:
		return "PCI/EISA bridge"
	case 0x03:
		return "PCI/Micro Channel bridge"
	case 0x04:
		return "PCI/PCI bridge"
	case 0x05:
		return "PCI/PCMCIA bridge"
	case 0x06:
		return "PCI/NuBus bridge"
	case 0x07:
		return "PCI/CardBus bridge"
	case 0x08:
		return "RACEway bridge"
	case 0x80:
		return "Other bridge type"
	}
	return "Unknown Bridge subclass"

}

func (c Class) getMassStorageClass() string {
	subclass := c.GetSubclass()
	switch subclass {
	case 0x00:
		return "SCSI controller"
	case 0x01:
		return "IDE controller"
	case 0x02:
		return "Floppy disk controller"
	case 0x03:
		return "IPI controller"
	case 0x04:
		return "RAID controller"
	case 0x80:
		return "Other mass storage controller"
	}
	return "Unknown mass storage controller"

}

func (c Class) getDisplayClass() string {
	subclass := c.GetSubclass()
	switch subclass {
	case 0x00:
		// This depends on the IF register that we don't have
		// access to from here.
		return "VGA compatible controller"
		//return "8514-compatible controller"
	case 0x01:
		return "XGA controller"
	case 0x02:
		return "3D controller"
	case 0x80:
		return "Other display controller"
	}
	return "Unknown display controller"
}

func (c Class) getNetworkClass() string {
	subclass := c.GetSubclass()
	switch subclass {
	case 0x00:
		return "Ethernet controller"
	case 0x01:
		return "Token ring controller"
	case 0x02:
		return "FDDI controller"
	case 0x03:
		return "ATM controller"
	case 0x04:
		return "ISDN controller"
	case 0x80:
		return "Other network controller"
	}
	return "Unknown network controller"
}

func (c Class) String() string {
	class := c.GetClass()
	switch class {
	case 0x00:
		return "Unknown class type"
	case 0x01:
		return c.getMassStorageClass()
	case 0x02:
		return c.getNetworkClass()
	case 0x03:
		return c.getDisplayClass()
	case 0x04:
		return "Multimedia device"
	case 0x05:
		return "Memory controller"
	case 0x06:
		return c.getBridgeClass()
	}
	return "Not yet implemented class type"
}
