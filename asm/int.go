package asm

// typedef struct{} Nothing;
import "C"

// If we import C, the GCCGO cross-compiler claims we're not using it.
// If we import it as _ "C", go test claims we can't rename it since it's
// the real go tool chain.
// This is a hack so that gmake, go test ./... and go fmt ./...
// all work. C.Nothing is a struct{}, so hopefully this gets
// optimized away by the compiler.
type cNoop C.Nothing

// Executes a CLI assembly instruction to disable interrupts
func CLI()

// Executes an STI assembly instruction to enable interrupts
func STI()

// Executes an HLT assembly instruction, which causes the CPU to
// wait for an interrupt to happen
func HLT()
