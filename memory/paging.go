package memory

import _ "C"

import (
	"unsafe"

	"github.com/driusan/kernel/asm"
)

var isInitialized bool

func IsPagingInitialized() bool {
	return isInitialized
}

const (
	PagePresent = 1 << iota
	PageReadWrite
	PageUserspace
	PageWriteThrough
	PageCacheDisabled
	PageAccessed
	_
	PageIs4MB
	PageGlobal
)

type PageDirectory *[1024]uint32
type PageTableEntry *[1024]uint32

func GetTableAddress(pt PageTableEntry) uint32 {
	// The linker linked all the symbols in virtual address space, but
	// paging needs to use the physical address, so we subtract 3GB
	// from the pointer.
	return uint32(uintptr(unsafe.Pointer(pt))) - 0xC0000000
}

// maps an address to the page table and entry in that table which
// corresponds to that address
func getTableEntryForAddress(a uintptr) (uint16, uint16, error) {
	if a%PageSize != 0 {
		return 0, 0, MemoryError("Address is not page aligned")
	}
	tbl := a / (1024 * PageSize)
	e := (a / PageSize) % 1024
	return uint16(tbl), uint16(e), nil
}

//extern getPageDirectory
func getPageDirectory() PageDirectory

//extern getPageTable
func getPageTable(uint16) PageTableEntry

//extern initialize_paging
func initPaging() *uint32

//extern loadPageDirectory
func loadPageDirectory(PageDirectory)

//extern enablePaging
func enablePaging()

type MultibootMemoryMap struct {
	Size     uint32
	BaseAddr uint64
	Length   uint64
	Memtype  uint32
}

// Denotes a single mapping from a physical address to a
// virtual address space.
type MMapEntry struct {
	// The physical address to be mapped. Must be page aligned.
	PAddr PhysicalAddress
	// The virtual address. The zero value for VAddr means
	// "the next available address" when loading a map.
	VAddr VirtualAddress

	// The size in bytes of the memory to be mapped.
	Length uint
}

func GetPhysicalAddress(addr VirtualAddress) (PhysicalAddress, error) {
	t, e, err := getTableEntryForAddress(uintptr(addr))
	if err != nil {
		return 0, err
	}

	table := getPageTable(t)
	adr := uintptr(table[e]) & 0xFFFFF000
	return PhysicalAddress(adr), nil
}

// Updates the page table so that for each segment in segments maps to
// the physical address to the corresponding virtual address.
func LoadMap(segments ...*MMapEntry) error {
	var startAddr VirtualAddress

	for _, me := range segments {
		if me.PAddr%PageSize != 0 {
			return MemoryError("Physical address is not page aligned.")
		}
		if me.VAddr%PageSize != 0 {
			return MemoryError("Virtual address is not page aligned.")
		}

		if me.VAddr == 0 {
			me.VAddr = startAddr
		}

		for addr := me.VAddr; addr <= (me.VAddr + VirtualAddress(me.Length)); addr += PageSize {
			if addr >= 0xC0000000 {
				return MemoryError("Can not remap kernel address space.")
			}

			t, e, err := getTableEntryForAddress(uintptr(addr))
			if err != nil {
				return err
			}
			table := getPageTable(t)

			offset := uint32(addr - me.VAddr)
			table[e] = (uint32(me.PAddr) + (offset)) | PagePresent | PageReadWrite
			asm.INVLPG(uintptr(addr))
			startAddr = addr + PageSize
		}
	}
	return nil
}
func InitializePaging(MMapAddr, MMapLength uintptr) {
	pd := getPageDirectory()
	table := getPageTable(0)

	var i uint32
	// Mark all pages as readwrite, but not present.
	for i = 0; i < 1024; i++ {
		pd[i] = PageReadWrite
	}

	// Start by identity mapping the first page and mark it as present,
	// regardless of what the boot loader told us.
	// TODO: Make a proper frame page allocator instead of this hack which
	// makes all memory except for kernel space identity mapped.

	for i = 0; i < 1024; i++ {
		table[i] = (i * 0x1000) | PagePresent | PageReadWrite
	}
	pd[0] = GetTableAddress(table) | PagePresent | PageReadWrite

	// Now identity map the rest of the memory that the multiboot loader
	// told us about.
	var mmap *MultibootMemoryMap

	// Now mark anything above the first MB that the multiboot boot loader
	// told us about as present.
	for offset := uintptr(0); offset < MMapLength; {
		mmap = (*MultibootMemoryMap)(unsafe.Pointer(uintptr(MMapAddr) + offset))
		//i++
		if mmap.Memtype == 1 && mmap.BaseAddr >= (1024*1024) {
			println(mmap.Length, " of available RAM at ", mmap.BaseAddr, "(Size:", mmap.Size, ")")
			// TODO: Verify this math.
			startDirIdx := uint16(mmap.BaseAddr / (4096 * 1024))
			startPageIdx := uint16((mmap.BaseAddr / 4096) % (1024))
			sizeDirEntries := uint16(mmap.Length / (4096 * 1024))
			lastIdx := uint16((mmap.Length / 4096) % 1024)
			if startDirIdx+sizeDirEntries >= 1024 {
				println("Warning: can't access memory above 4GB with 4KB pages")
				println("Losing", mmap.Length, " bytes of memory.")
			} else {
				var startIdx, endIdx uint16
				var pageIdx uint16
				for pageTableIdx := startDirIdx; pageTableIdx <= startDirIdx+sizeDirEntries; pageTableIdx++ {
					switch pageTableIdx {
					case startDirIdx:
						startIdx = startPageIdx
						endIdx = 1024
					case startDirIdx + sizeDirEntries:
						startIdx = 0
						endIdx = lastIdx
					default:
						startIdx = 0
						endIdx = 1024
					}
					table = getPageTable(pageTableIdx)
					for pageIdx = startIdx; pageIdx < endIdx; pageIdx++ {
						table[pageIdx] = (uint32(pageIdx)*0x1000 + uint32(pageTableIdx)*(4096*1024)) | PagePresent | PageReadWrite
					}
					pd[pageTableIdx] = GetTableAddress(table) | PagePresent | PageReadWrite
				}
			}
		} else {
			// Reserved memory.
			// println(mmap.Length, " memory of type", mmap.Memtype, " at ", mmap.BaseAddr, "(Size:", mmap.Size, ")")
		}

		offset += unsafe.Sizeof(mmap) // *MultibootMemoryMap)
		offset += uintptr(mmap.Size)
	}

	// Now, map to page table entries from 0xC0000000 for the kernel.
	i = 768
	table = getPageTable(uint16(i))

	for pageIdx := 0; pageIdx < 1024; pageIdx++ {
		table[pageIdx] = (uint32(pageIdx) * 0x1000) | PagePresent | PageReadWrite
	}
	pd[i] = GetTableAddress(table) | PagePresent | PageReadWrite

	i = 769
	table = getPageTable(uint16(i))
	for pageIdx := 0; pageIdx < 1024; pageIdx++ {
		table[pageIdx] = (uint32(pageIdx)*0x1000 + (4096 * 1024)) | PagePresent | PageReadWrite
	}
	pd[i] = GetTableAddress(table) | PagePresent | PageReadWrite

	loadPageDirectory(pd)
	enablePaging()

	afterPagingInit()
}
