.text
.globl loadPageDirectory
loadPageDirectory:
	push %ebp
	mov %esp, %ebp
	mov 8(%esp), %eax
	# Convert the pointer that was passed from virtual address space
	# to physical before putting it in register cr3.
	sub $0xC0000000, %eax
	mov %eax, %cr3
	mov %ebp, %esp
	pop %ebp
	ret

.text
.globl enablePaging
enablePaging:
	push %ebp
	mov %esp, %ebp
	mov %cr0, %eax
	or $0x80000000, %eax
	mov %eax, %cr0
	mov %ebp, %esp
	pop %ebp
	ret
