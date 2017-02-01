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
	# The caller pushed the parameters onto the stack, but then
	# the INT call pushed the return address, return code segment selector,
	# and Eflags onto the stack too.
	# We need to fiddle with things so that the call to pwrite has
	# the correct parameters. We start by re-pushing the params, so
	# that the call doesn't mess up the return address from the interrupt.
	#
	# The parameters are located at 16(%esp), but each time
	# we push, the value of %esp changes by 4, so this is a little weird.
	# It's basically pushing (16-32)(%esp) onto the stack, to ensure
	# that when we get to the IRET, we haven't screwed up the return pointer.
	PUSH %ebp
	# PUSH 32(%esp)
	PUSH 36(%esp)
	# PUSH 28(%esp)
	PUSH 36(%esp)
	# PUSH 24(%esp)
	PUSH 36(%esp)
	# PUSH 20(%esp)
	PUSH 36(%esp)
	# PUSH 16(%esp)
	PUSH 36(%esp)
	CALL github_com_driusan_kernel_executable_plan9.Pwrite

	# Now, the return value is in %eax, but our fake parameters are
	# shadowing the IRET return, so clean them up.
	ADD $20, %esp
	POP %ebp
	IRET

.text
.globl github_com_driusan_kernel_executable_plan9._exits
github_com_driusan_kernel_executable_plan9._exits:
	# The interrupt pushed the return address on the stack. Since this is
	# an exits syscall, we don't really care, we just want the parameter
	# that the caller pushed on the stack before the INT instruction. Get
	# rid of the stuff that INT thought was important.
	ADD $16, %esp
	CALL github_com_driusan_kernel_executable_plan9.Exits
ihalt:
	# Exits should never return. If it does, just hang.
	HLT
	JMP ihalt

