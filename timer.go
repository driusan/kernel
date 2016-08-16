package kernel

import "github.com/driusan/kernel/interrupts"

var ticks uint64

func TimerHandler(r *interrupts.Registers) {
	ticks++
	if ticks%18 == 0 {
		//println("Approximately one second")
	}
}
