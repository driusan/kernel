package shell

import (
	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/process"
)

func Run() {
	namespace := process.NewNamespace()
	
	print(namespace.Test())
	println("Entering the shell.")
	for {
		asm.HLT()
	}
	println("Leaving the shell.")
}