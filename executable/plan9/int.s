# Define the syscall interrupt in asm, so that it will return
# with iret instead of ret

.text
.globl github_com_driusan_kernel_executable_plan9.p9int
github_com_driusan_kernel_executable_plan9.p9int:
	cli
	cmp $8, %eax
	je github_com_driusan_kernel_executable_plan9._exits
	cmp $51, %eax
	je github_com_driusan_kernel_executable_plan9._pwrite
retsyscall:
	iret

.text
.globl github_com_driusan_kernel_executable_plan9._pwrite
github_com_driusan_kernel_executable_plan9._pwrite:
	# The parameters are located at 16(%esp), but each time
	# we push, the value of %esp changes by 4, so this is a little weird.
	# It's basically pushing (16-32)(%esp) onto the stack, to ensure
	# that when we get to the IRET, we haven't screwed up the return pointer
	# PUSH 16(%esp)
	PUSH 32(%esp)
	# PUSH 20(%esp)
	PUSH 32(%esp)
	# PUSH 24(%esp)
	PUSH 32(%esp)
	# PUSH 28(%esp)
	PUSH 32(%esp)
	# PUSH 32(%esp)
	PUSH 32(%esp)
	CALL github_com_driusan_kernel_executable_plan9.Pwrite
	ADD $20, %esp
	IRET

.text
.globl github_com_driusan_kernel_executable_plan9._exits
github_com_driusan_kernel_executable_plan9._exits:
	CALL github_com_driusan_kernel_executable_plan9.Exits
ihalt:
	HLT
	JMP ihalt
#	JMP retsyscall

