package plan9

import (
	_ "C"
	"io"
	"unsafe"

	"github.com/driusan/kernel/interrupts"
	"github.com/driusan/kernel/terminal"
)

// Valid syscalls from Plan 9
const (
	SysR1 = iota // I'm not sure what this is.
	_            // Was: errstr, removed from Plan 9
	Bind
	Chdir
	Close
	Dup
	Alarm
	Exec
	ExitS
	_ // Was: fsession, removed from Plan 9
	FAuth
	_ // Was: fstat, removed from Plan 9
	SegBRK
	_ // Was: mount, removed from Plan 9
	Open
	_ // Was: read, removed from Plan9 (replaced by pread)
	OSeek
	Sleep
	_ // Was: stat, replaced with other stat syscall
	RFork
	_ // Was: write, replaced with pwrite
	Pipe
	Create
	Fd2Path
	Brk_
	Remove
	_ // was: wstat
	_ // was: fwstat
	Notify
	Noted
	SegAttach
	SegDetach
	SegFree
	SegFlush
	RendezVous
	Unmount
	_ // Was: wait
	SemAcquire
	SemRelease
	Seek
	FVersion
	ErrStr
	Stat
	FStat
	WStat
	FWStat
	Mount
	AWait
	PRead
	PWrite
	TSemAcquire
	_ // was: nsec, removed from 9front. Still used in Plan 9 Go?
)

//extern installInt
func installIntC()

func InstallSyscallInterrupt() {
	installIntC()
}

func _pwrite()

func Syscall(r *interrupts.Registers) {
	// TODO: Look these up in the real Plan9 src to make sure
	// they're correct. These are the interrupts used by
	// sys_plan9_386.s in the Go runtime.
	switch r.Eax {
	case 0: // SYSR1 (what is this?)
	case 2: // unused (was _ERRSTR)
	case 4: // close()
	case ExitS: // exits()
		println("Should exits")
	case Open: // open()
		println("Should open")
	case 17: // sleep()
	case 19: // rfork
	case 24: // brk_
	case 28: // notify
	case 29: // noted
	case 37: // semacquire
	case 38: // semrelease
	case 39: // seek()
	case 41: // errstr(int8 *buf, int32 len)
		println("Should errstr")
	case 50: // pread()
		println("Should pread")
	case 51: // pwrite()
		println("Should pwrite")
		_pwrite()
	case 52: // tsemacquire
	case 53: // nsec?
	}
	println("In Plan9 syscall")
}

// Implements the syscall:
//
//	long pwrite(int fd, void *buf, long nbytes, vlong offset)
//
// This can't have a signature that's more idiomatic Go such as
// 	func PWrite(fd FileDescriptor, []buf, int64 offset) (n, error)
// because the syscalls in Plan 9's ABI are defined in terms of
// C.
func Pwrite(fd int, buf unsafe.Pointer, nbytes int32, offset int64) int32 {
	if fd > len(activeProc.FDs) {
		println("Invalid file descriptor")
		return 0
	}

	// Convert *buf into a []byte for compatibility with the io.Writer
	// interface
	b := make([]byte, nbytes, nbytes)
	for i := uintptr(0); i < uintptr(nbytes); i++ {
		b[i] = *(*byte)(unsafe.Pointer(uintptr(buf) + uintptr(i)))
	}

	// TODO: Add a Mutex, to make sure Seek + Write is atomic. (Is it
	// necessary? This is an interrupt handler, it shouldn't be interrupted
	// before the IRET call.)
	if offset > 0 {
		n, err := activeProc.FDs[fd].Seek(offset, io.SeekStart)
		if err != nil {
			println("Error seeking to offset: ", err.Error())
			return 0
		}
		if n != offset {
			println("Warning: Seek returned ", n, " not ", offset)
		}
	}
	n, err := activeProc.FDs[fd].Write(b)
	if err != nil {
		println(err.Error())
		return int32(n)
	}
	return int32(n)
}

type CString uintptr

func Exits(s CString) {
	if s == 0 {
		return
	}
	// This should actually put it in a Waitmsg for the parent, but we're
	// not that advanced yet. See exits(2).
	terminal.PrintCString(unsafe.Pointer(s))
	println("There should be an error printed above.")
}
