package memory

import (
	"unsafe"
	//"github.com/driusan/kernel/terminal"
)

const (
	PagePresent = 1 << iota
	PageReadWrite
	PageUserspace
	PageWriteThrough
	PageCacheDisaled
	PageAccessed
	_ // What is bit 7? 
	PageIs4MB
	PageGlobal
)

type PageDirectory *[1024]uint32
type PageTable *[1024]uint32

func GetAddress(pt PageTable) uint32{
	return uint32(uintptr(unsafe.Pointer(pt))) 
}
//extern getPageDirectory
func getPageDirectory() PageDirectory

//extern getInitialPageTable
func getInitialPageTable() PageTable

//extern initialize_paging
func initPaging() *uint32

//extern loadPageDirectory
func loadPageDirectory(PageDirectory)

//extern enablePaging
func enablePaging()

type MultibootMemoryMap struct{
       Size uint32
       BaseAddr uint64
       Length uint64
       Memtype uint32
}

func InitializePaging() {
	pd := getPageDirectory()
	table1 := getInitialPageTable()

	var i uint32
	// Mark all pages as readwrite, but not present.
	for i = 0; i < 1024; i++ {
		pd[i] = PageReadWrite
	}

	for i = 0; i < 1024; i++ {
		table1[i] = (i*0x1000) | PagePresent | PageReadWrite;
	}

	pd[0] = GetAddress(table1) | PagePresent | PageReadWrite

	loadPageDirectory(pd)
	enablePaging();
}

