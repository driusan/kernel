package memory

type eightBools byte

type PageNumber int
type PageSpan int

const (
	gB = (1024 * 1024 * 1024)
)

const PageSize = 4096

type MemoryError string

func (m MemoryError) Error() string {
	return string(m)
}

var InvalidUsage error
var NoMemory error

var largePages map[PageNumber]PageSpan

func InitPkg() {
}
func afterPagingInit() {
	// Initialize the allocation map.
	for i, _ := range pagesAllocated {
		pagesAllocated[i] = 0
	}

	// Assume the whole kernel is in the first 4 page tables (16MB) and mark
	// it as allocated for Malloc.
	// (One page for the page table, one page for the kernel code, and
	// one page for good measure.)
	// TODO: Make this smarter.
	for i := 0; i < 4*1024; i++ {
		(&pagesAllocated).Set(PageNumber(i), true)

	}

	// It's now safe to call __go_new
	isInitialized = true

	// This can't be done in InitPkg because initPaging hasn't marked
	// the appropriate places as allocated, and make() will call Malloc
	largePages = make(map[PageNumber]PageSpan)

	// These also call __go_new, so the allocation map needs to be initialized
	InvalidUsage = MemoryError("Invalid usage")
	NoMemory = MemoryError("No memory available")

}

// Stores a bit map of if physical pages are currently in use.
type pageAllocTable [(4 * gB / PageSize) / 8]eightBools

//type pageAllocTable [32768]eightBools

var pagesAllocated pageAllocTable

func (pat *pageAllocTable) isPageAllocated(page PageNumber) bool {
	entry := page / 8
	bit := uint8(page % 8)

	return pat[entry].Get(bit)
}

func (pat *pageAllocTable) Set(page PageNumber, allocated bool) {
	entry := page / 8
	bit := uint8(page % 8)

	(&pat[entry]).Set(bit, allocated)
}

func (eb eightBools) Get(num uint8) bool {
	return eb&(1<<(num)) != 0
}

func (eb *eightBools) Set(num uint8, val bool) {
	if val {
		*eb |= eightBools(1 << (num))
	} else {
		*eb &= (255 - (1 << (num)))
	}
}

// This is a very stupid/simple malloc implementation which always allocates at
// least one page, and keeps track of which pages allocated in a large bitmap.
func Malloc(amt uint) (uintptr, error) {
	var numPages PageSpan = PageSpan((amt-1)/PageSize) + 1

	for i := PageNumber(0); i < (4 * gB / PageSize); i++ {
		if pagesAllocated.isPageAllocated(i) == false {
			for j := 0; j < int(numPages); j++ {
				if (i + PageNumber(j)) >= (4 * gB / PageSize) {
					return 0, NoMemory
				}

				if pagesAllocated.isPageAllocated(i + PageNumber(j)) {
					// There was an interruption in the span, so move
					// to the next possible span
					i = i + PageNumber(j)
					goto next
				}
			}
			if numPages > 1 {
				largePages[i] = numPages
			}
			for j := i; j < i+PageNumber(numPages); j++ {
				pagesAllocated.Set(j, true)

			}
			return uintptr(i * PageSize), nil
		}
	next:
	}
	return 0, NoMemory

}

// This is a very stupid and insecure Free implementation to go along with the
// stupid Malloc implementation
func Free(addr uintptr) error {
	pageStart := PageNumber(addr / PageSize)
	if !pagesAllocated.isPageAllocated(pageStart) {
		return InvalidUsage
	}

	if pageSpan, ok := largePages[pageStart]; ok {
		for j := pageStart; j < pageStart+PageNumber(pageSpan); j++ {
			pagesAllocated.Set(j, false)

		}
		return nil
	}
	pagesAllocated.Set(pageStart, false)
	return nil
}
