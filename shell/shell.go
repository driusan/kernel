package shell

import (
	"github.com/driusan/kernel/asm"
	//"github.com/driusan/kernel/process"
)

func Run() {
	//ns := process.NewNamespace()

	//_, err := ns.Open("/dev/cons")
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//prompt := []byte{ '>', ' '}
	println("Entering the shell.")
	for {
		//cons.Write(prompt)

		asm.HLT()
	}
	println("Leaving the shell.")
}
