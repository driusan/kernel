package pci

type HeaderType byte

func (ht HeaderType) String() string {
	switch ht {
	case 0x00:
		return "General PCI device"
	case 0x01:
		return "PCI-to-PCI bridge"
	case 0x02:
		return "CardBus bridge"
	case 0x80:
		return "Multifunction general PCI device"
	case 0x81:
		return "Multifunction PCI-to-PCI bridge"
	case 0x82:
		return "Multifunction CardBus bridge"
	}
	if ht.IsMultifunction() {
		return "Unknown multifunction PCI header type"
	}
	return "Unknown PCI header type"
}

func (ht HeaderType) IsMultifunction() bool {
	return ht >= 0x80
}
