package memory

import (
	"testing"
)

func TestEightBools(t *testing.T) {
	var b eightBools

	for i := uint8(0); i < 8; i++ {
		if b.Get(i) != false {
			t.Errorf("Bit %d is not false upon initialization", i)
		}
	}
	b.Set(3, true)
	if val := byte(b); val != 0x08 {
		t.Errorf("Got %v expected 0x08", val)
	}
	b.Set(7, true)
	if val := byte(b); val != 0x88 {
		t.Errorf("Got %v expected 0x88", val)
	}

	b.Set(0, true)
	if val := byte(b); val != 0x89 {
		t.Errorf("Got %x expected 0x89", val)
	}
	b.Set(1, true)
	b.Set(2, true)
	if val := byte(b); val != 0x8F {
		t.Errorf("Got %x expected 0x8F", val)
	}
	b.Set(4, true)
	b.Set(5, true)
	b.Set(6, true)
	if val := byte(b); val != 0xFF {
		t.Errorf("Got %x expected 0xFF", val)
	}

	b.Set(4, true)
	if val := byte(b); val != 0xFF {
		t.Errorf("Got %x expected 0xFF", val)
	}
	b.Set(3, false)
	if val := byte(b); val != 0xF7 {
		t.Errorf("Got %x expected 0xF7", val)
	}
	b.Set(0, false)
	if val := byte(b); val != 0xF6 {
		t.Errorf("Got %x expected 0xF6", val)
	}

}
