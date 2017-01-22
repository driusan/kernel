package plan9

import (
	"github.com/driusan/kernel/interrupts"

	_ "C"
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
	_ // was: nsec, removed from Plan 9
)

//extern installInt
func installIntC()

func InstallSyscallInterrupt() {
	installIntC()
}

func Syscall(r *interrupts.Registers) {
	// TODO: Look these up in the real Plan9 src to make sure
	// they're correct. These are the interrupts used by
	// sys_plan9_386.s in the Go runtime.
	switch r.Eax {
	case 0: // SYSR1 (what is this?)
	case 2: // unused (was _ERRSTR)
	case 4: // close()
	case 8: // exits()
		println("Should exits")
	case 14: // open()
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
	case 52: // tsemacquire
	case 53: // nsec?
	}
	println("In Plan9 syscall")
}
