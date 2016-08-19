package memory

import (
	"unsafe"
	//"github.com/driusan/kernel/asm"
	//"github.com/driusan/kernel/terminal"
)

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
	return uint32(uintptr(unsafe.Pointer(pt)))
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

func InitializePaging(MMapAddr, MMapLength uintptr) {
	pd := getPageDirectory()
	table := getPageTable(0)
	println("Table 0 address", GetTableAddress(table))

	var i uint32
	// Mark all pages as readwrite, but not present.
	for i = 0; i < 1024; i++ {
		pd[i] = PageReadWrite
	}

	// Start by identity mapping the first page and mark it as present,
	// regardless of what the boot loader told us.
	for i = 0; i < 1024; i++ {
		table[i] = (i * 0x1000) | PagePresent | PageReadWrite
	}
	pd[0] = GetTableAddress(table) | PagePresent | PageReadWrite

	// Now identity map the rest of the memory that the multiboot loader
	// told us about.
	// TODO: Map the kernel into the higher portion of memory and locate
	//       the page table there.
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
			//println(mmap.Length, " memory of type", mmap.Memtype, " at ", mmap.BaseAddr, "(Size:", mmap.Size, ")")
		}

		offset += unsafe.Sizeof(*MultibootMemoryMap)
		offset += uintptr(mmap.Size)
	}
	loadPageDirectory(pd)
	enablePaging()

	// Assume the whole kernel is in the first 3 pages (12MB) and mark
	// it as allocated for Malloc.
	// (One page for the page table, one page for the kernel code, and
	// one page for good measure.)
	// TODO: Make this smarter.
	pagesAllocated.Set(0, true)
	pagesAllocated.Set(1, true)
	pagesAllocated.Set(2, true)

}
