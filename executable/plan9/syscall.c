// defined in idt.go
extern void idt_set_gate(unsigned char num, unsigned long base, unsigned short sel, unsigned char flags) __asm__("github_com_driusan_kernel_descriptortables.IDTSetGate");

void syscall(void) __asm__("github_com_driusan_kernel_executable_plan9.p9int");

void installInt() {
	idt_set_gate(64, (unsigned) syscall, 0x08, 0x8E);
}