.text
.globl STI
.globl github_com_driusan_kernel_asm.STI
STI:
github_com_driusan_kernel_asm.STI:
	sti
	ret

.text
.globl CLI
.globl github_com_driusan_kernel_asm.CLI
CLI:
github_com_driusan_kernel_asm.CLI:
	cli
	ret


.text
.globl HLT
.globl github_com_driusan_kernel_asm.HLT
HLT:
github_com_driusan_kernel_asm.HLT:
	hlt
	ret
