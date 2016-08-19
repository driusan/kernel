package memory

type eightBools byte

const (
	gB = (1024 * 1024 * 1024)
)

// Stores a bit map of if physical pages are currently in use.
type pageAllocTable [(4 * gB / 4096) / 8]eightBools

//type pageAllocTable [32768]eightBools

var pagesAllocated pageAllocTable

func (pat pageAllocTable) isPageAllocated(page int) bool {
	entry := page / 8
	bit := uint8(page % 8)

	return pat[entry].Get(bit)
}

func (pat pageAllocTable) Set(page int, allocated bool) {
	entry := page / 8
	bit := uint8(page % 8)

	pat[entry].Set(bit, allocated)
}

func (eb eightBools) Get(num uint8) bool {
	return eb&(1<<(num-1)) != 0
}

func (eb *eightBools) Set(num uint8, val bool) {
	if val {
		*eb |= eightBools(1 << (num))
	} else {
		*eb &= eightBools(1 << num)
	}
}
