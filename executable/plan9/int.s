# Define the syscall interrupt in asm, so that it will return
# with iret instead of ret

.text
.globl github_com_driusan_kernel_executable_plan9.exec
github_com_driusan_kernel_executable_plan9.exec:
	# MOV $0xC0000000, %esp
	# MOV $0, %eax
	PUSH $0
	PUSH $0
	CALL 0x20 
	# JMP $0x20

.text
.globl github_com_driusan_kernel_executable_plan9.p9int
github_com_driusan_kernel_executable_plan9.p9int:
	# Push all the registers onto the caller's parameters, so that we can
	# work with them in Go.
	cli
	pusha
	push %ds
	push %es
	push %fs
	push %gs
	mov $0x10, %ax
	mov %ax, %ds
	mov %ax, %es
	mov %ax, %fs
	mov %ax, %gs
	mov %esp, %eax
	push %eax
	mov $github_com_driusan_kernel_executable_plan9.Syscall, %eax
	call *%eax
	pop %eax
	pop %gs
	pop %fs
	pop %es
	pop %ds
	popa
	add $8, %esp
	iret

