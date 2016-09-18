package memory

import (
	"testing"
)

func TestTableEntryTranslation(t *testing.T) {
	tests := []struct {
		Address      uintptr
		Table, Entry uint16
	}{
		{0, 0, 0},
		{0x1000, 0, 1},
		{0xC0000000, 768, 0},
		{0xC0001000, 768, 1},
		{0xC000F000, 768, 15},
	}

	for i, test := range tests {
		tbl, e, err := getTableEntryForAddress(test.Address)
		if err != nil {
			t.Errorf("Unexpected error for %d: %v", i, err.Error())
		}
		if tbl != test.Table {
			t.Errorf("Invalid table for %d: got %v want %v", i, tbl, test.Table)
		}
		if e != test.Entry {
			t.Errorf("Invalid entry for %d: got %v want %v", i, e, test.Entry)
		}
	}

}
