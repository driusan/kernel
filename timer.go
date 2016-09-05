package kernel

import "github.com/driusan/kernel/interrupts"

var ticks uint64

// TimerHandler handlers an interrupt from the PIT.
// It currently only increments the tick counter.
func TimerHandler(r *interrupts.Registers) {
	ticks++
	if ticks%18 == 0 {
		//println("Approximately one second")
	}
}
