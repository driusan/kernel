package memory


//extern initialize_paging
func initPaging()

func InitializePaging() {
	initPaging()
}


 
type MultibootMemoryMap struct{
       Size uint32
       BaseAddr uint64
       Length uint64
       Memtype uint32
}
