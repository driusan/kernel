.text
.globl gdt_flush
gdt_flush:
	lgdt github_com_driusan_kernel_descriptortables.GDTPtr
	mov $0x10, %ax
	mov %ax, %ds
	mov %ax, %es
	mov %ax, %fs
	mov %ax, %gs
	mov %ax, %ss
	ljmp $0x08, $flush2
flush2:
	ret



.globl idt_load
idt_load:
	lidt github_com_driusan_kernel_descriptortables.IDTPtr
	ret

