package kernel

import "interrupts"

var ticks uint64

func TimerHandler(r *interrupts.Registers) {
	ticks++
	if ticks%18 == 0 {
		//println("Approximately one second")
	}
}
