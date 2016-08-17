package asm

// Executes a CLI assembly instruction to disable interrupts
func CLI()

// Executes an STI assembly instruction to enable interrupts
func STI()

// Executes an HLT assembly instruction, which causes the CPU to
// wait for an interrupt to happen
func HLT()
