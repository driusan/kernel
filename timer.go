package kernel

var ticks uint64
func TimerHandler(r *Registers) {
	ticks++
	if ticks % 18 == 0 {
		println("Approximately one second")
	}
}
