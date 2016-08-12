package ps2

import (
	"asm"
	"interrupts"
)

// the packet size of the current protocol
var packetSize uint8

// the current packet we're receiving
var curPacketIdx uint8

// the complete current packet. Normal PS2 is 3 bytes, Intellimouse is 4,
// so reserve an array of 4 bytes just in case. This may need to increase
// if multitouch drivers are ever added.
var curPacket [4]byte

func sendPS2AuxCmd(cmd byte) error {
	// try to send
	for i := 0; i < 4; i++ {
		waitOutput()
		asm.OUTB(0x64, 0xD4)
		waitOutput()
		asm.OUTB(0x60, cmd)

		waitInput()
		resp := asm.INB(0x60)
		/*
				switch statements result in a compiler error.
				right now. Something needs to be fixed, probably
				related to us linking in freestanding mode.
			switch resp {
				case PS2Success:
					return nil
				case PS2Fail:
					return SendFailure
				case PS2Resend:
					continue
				default:
					return UnknownError
			}
		*/
		if resp == Success {
			return nil
		} else if resp == Fail {
			return SendFailure
		} else if resp == Resend {
			continue
		}
	}
	return TooManyRetries
}

func EnableMouse() error {
	// check if there is a mouse port
	waitOutput()
	asm.OUTB(0x64, 0xA9)

	waitInput()
	b := asm.INB(0x60)
	if b != 0 {
		return PS2Error("No PS2 mouse port exists")
	}
	packetSize = 3

	// enable the mouse port
	waitOutput()
	asm.OUTB(0x64, 0xA8)

	// get the status byte so we can enable interrupt 12
	waitOutput()
	asm.OUTB(0x64, 0x20)
	waitInput()
	b = asm.INB(0x60)

	// and enable IRQ12 for the mouse
	b |= (1 << 1)
	waitOutput()
	asm.OUTB(0x64, 0x60)
	waitOutput()
	asm.OUTB(0x60, b)

	// Reset the mouse
	sendPS2AuxCmd(0xFF)
	readPS2Port()
	// and enable streaming
	sendPS2AuxCmd(0xF4)
	readPS2Port()

	return nil
}

// Handles an interrupt from the mouse on IRQ 12
func MouseHandler(r *interrupts.Registers) {
	b := readPS2Port()
	curPacket[curPacketIdx] = b

	curPacketIdx++

	// Received a complete packet, so process it, and then reset
	// everything for the next packet.
	if curPacketIdx >= packetSize {
		processPacket(curPacket)
		curPacketIdx = 0
		for i := 0; i < len(curPacket); i++ {
			curPacket[i] = 0
		}
	}
}

// This is a stub. It doesn't do anything exciting other than
// print how much X and Y have changed by. We should be using
// an interface instead.
func processPacket(packet [4]byte) {
	var x, y int16
	if packet[0]&0x10 != 0 {
		x = -int16(byte(255) - packet[1])
	} else {
		x = int16(packet[1])
	}
	if packet[0]&0x20 != 0 {
		y = -int16(byte(255) - packet[2])
	} else {
		y = int16(packet[2])
	}
	if x != 0 || y != 0 {
		println("Mouse moved by (X, Y)=(", x, ",", y, ")")
	}
}
