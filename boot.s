# Declare constants for the multiboot header.
.set ALIGN,    1<<0             # align loaded modules on page boundaries
.set MEMINFO,  1<<1             # provide memory map
.set FLAGS,    ALIGN | MEMINFO  # this is the Multiboot 'flag' field
.set MAGIC,    0x1BADB002       # 'magic number' lets bootloader find the header
.set CHECKSUM, -(MAGIC + FLAGS) # checksum of above, to prove we are multiboot

# Declare a multiboot header that marks the program as a kernel. These are magic
# values that are documented in the multiboot standard. The bootloader will
# search for this signature in the first 8 KiB of the kernel file, aligned at a
# 32-bit boundary. The signature is in its own section so the header can be
# forced to be within the first 8 KiB of the kernel file.
.section .multiboot
.align 4
.long MAGIC
.long FLAGS
.long CHECKSUM

# The multiboot standard does not define the value of the stack pointer register
# (esp) and it is up to the kernel to provide a stack. This allocates room for a
# small stack by creating a symbol at the bottom of it, then allocating 16384
# bytes for it, and finally creating a symbol at the top. The stack grows
# downwards on x86. The stack is in its own section so it can be marked nobits,
# which means the kernel file is smaller because it does not contain an
# uninitialized stack.
.section .bootstrap_stack, "aw", @nobits
stack_bottom:
.skip 16384 # 16 KiB
stack_top:

# The linker script specifies _start as the entry point to the kernel and the
# bootloader will jump to this position once the kernel has been loaded. It
# doesn't make sense to return from this function as the bootloader is gone.
.section .text
.global _start
.type _start, @function
_start:
	cli
	# The bootloader has loaded us into 32-bit protected mode on a x86 machine.
	# Interrupts are disabled. Paging is disabled. The processor state is as
	# defined in the multiboot standard. The kernel has full control of the CPU.
	# The kernel can only make use of hardware features and any code it provides
	# as part of itself. There's no printf function, unless the kernel provides
	# its own <stdio.h> header and a printf implementation. There are no
	# security restrictions, no safeguards, no debugging mechanisms, only what
	# the kernel provides itself. It has absolute and complete power over the
	# machine.

	# To set up a stack, we set the esp register to point to the top of our
	# stack (as it grows downwards on x86 systems). This is necessarily done in
	# assembly as languages such as C cannot function without a stack.
	mov $stack_top, %esp

	# This is a good place to initialize crucial processor state before the
	# high-level kernel is entered. It's best to minimize the early environment
	# where crucial features are offline. Note that the processor is not fully
	# initialized yet: Features such as floating point instructions and
	# instruction set extensions are not initialized yet. The GDT should be
	# loaded here. Paging should be enabled here. C++ features such as global
	# constructors and exceptions will require runtime support to work as well.
	# Enter the high-level kernel.

	# Add the multiboot info onto the stack as the first parameter to KernelMain
	leal -4(%esp), %esp
	movl %ebx, (%esp)
	call boot.kernel.KernelMain
	# If the system has nothing more to do, put the computer into an infinite
	# loop. To do that:
	# 1) Disable interrupts with cli (clear interrupt enable in eflags). They
	#    are already disabled by the bootloader, so this is not needed. Mind
	#    that you might later enable interrupts and return from kernel_main
	#    (which is sort of nonsensical to do).
	# 2) Wait for the next interrupt to arrive with hlt (halt instruction).
	#    Since they are disabled, this will lock up the computer.
	# 3) Jump to the hlt instruction if it ever wakes up due to a
	#    non-maskable interrupt occurring or due to system management mode.
	cli
1:	hlt
	jmp 1b

# Set the size of the _start symbol to the current location '.' minus its start.
# This is useful when debugging or when you implement call tracing.
.size _start, . - _start

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


.text
.globl loadPageDirectory
loadPageDirectory:
	push %ebp
	mov %esp, %ebp
	mov 8(%esp), %eax
	mov %eax, %cr3
	mov %ebp, %esp
	pop %ebp
	ret

.text
.globl gdt_flush
gdt_flush:
	lgdt boot.kernel.Gp
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
	lidt idtp
	ret

.text
.globl boot.kernel.Halt
.globl halt
boot.kernel.Halt:
halt:
	hlt
	jmp boot.kernel.Halt

.text 
.globl isr0
.globl isr1
.globl isr2
.globl isr3
.globl isr4
.globl isr5
.globl isr6
.globl isr7
.globl isr8
.globl isr9
.globl isr10
.globl isr11
.globl isr12
.globl isr13
.globl isr14
.globl isr15
.globl isr16
.globl isr17
.globl isr18
.globl isr19
.globl isr20
.globl isr21
.globl isr22
.globl isr23
.globl isr24
.globl isr25
.globl isr26
.globl isr27
.globl isr28
.globl isr29
.globl isr30
.globl isr31
isr0:
	cli
	push $0
	push $0
	jmp isr_common_stub
isr1:
	cli
	push $0
	push $1
	jmp isr_common_stub
isr2:
	cli
	push $0
	push $2
	jmp isr_common_stub
isr3:
	cli
	push $0
	push $3
	jmp isr_common_stub
isr4:
	cli
	push $0
	push $4
	jmp isr_common_stub
isr5:
	cli
	push $0
	push $5
	jmp isr_common_stub
isr6:
	cli
	push $0
	push $6
	jmp isr_common_stub
isr7:
	cli
	push $0
	push $7
	jmp isr_common_stub
isr8:
	cli
	push $8
	jmp isr_common_stub
isr9:
	cli
	push $0
	push $9
	jmp isr_common_stub
isr10:
	cli
	push $10
	jmp isr_common_stub
isr11:
	cli
	push $11
	jmp isr_common_stub
isr12:
	cli
	push $12
	jmp isr_common_stub
isr13:
	cli
	push $13
	jmp isr_common_stub
isr14:
	cli
	push $14
	jmp isr_common_stub
isr15:
	cli
	push $0
	push $15
	jmp isr_common_stub
isr16:
	cli
	push $0
	push $16
	jmp isr_common_stub
isr17:
	cli
	push $0
	push $17
	jmp isr_common_stub
isr18:
	cli
	push $0
	push $18
	jmp isr_common_stub
isr19:
	cli
	push $0
	push $19
	jmp isr_common_stub
isr20:
	cli
	push $0
	push $20
	jmp isr_common_stub
isr21:
	cli
	push $0
	push $21
	jmp isr_common_stub
isr22:
	cli
	push $0
	push $22
	jmp isr_common_stub
isr23:
	cli
	push $0
	push $23
	jmp isr_common_stub
isr24:
	cli
	push $0
	push $24
	jmp isr_common_stub
isr25:
	cli
	push $0
	push $25
	jmp isr_common_stub
isr26:
	cli
	push $0
	push $26
	jmp isr_common_stub
isr27:
	cli
	push $0
	push $27
	jmp isr_common_stub
isr28:
	cli
	push $0
	push $28
	jmp isr_common_stub
isr29:
	cli
	push $0
	push $29
	jmp isr_common_stub
isr30:
	cli
	push $0
	push $30
	jmp isr_common_stub
isr31:
	cli
	push $0
	push $31
	jmp isr_common_stub

.text
.globl isr_common_stub
isr_common_stub:
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
	call fault_handler # mov fault_handler, %eax
	# call %eax
	pop %eax
	pop %gs
	pop %fs
	pop %es
	pop %ds
	popa
	add $8, %esp
	iret

.text
.globl irq0
.globl irq1
.globl irq2
.globl irq3
.globl irq4
.globl irq5
.globl irq6
.globl irq7
.globl irq8
.globl irq9
.globl irq10
.globl irq11
.globl irq12
.globl irq13
.globl irq14
.globl irq15

irq0:
	cli
	push $0
	push $32
	hlt
	jmp irq_common_stub
irq1:
	cli
	push $0
	push $33
	jmp irq_common_stub
irq2:
	cli
	push $0
	push $34
	jmp irq_common_stub
irq3:
	cli
	push $0
	push $35
	jmp irq_common_stub
irq4:
	cli
	push $0
	push $36
	jmp irq_common_stub
irq5:
	cli
	push $0
	push $37
	jmp irq_common_stub
irq6:
	cli
	push $0
	push $38
	jmp irq_common_stub
irq7:
	cli
	push $0
	push $39
	jmp irq_common_stub
irq8:
	cli
	push $0
	push $40
	jmp irq_common_stub
irq9:
	cli
	push $0
	push $41
	jmp irq_common_stub
irq10:
	cli
	push $0
	push $42
	jmp irq_common_stub
irq11:
	cli
	push $0
	push $43
	jmp irq_common_stub
irq12:
	cli
	push $0
	push $44
	jmp irq_common_stub
irq13:
	cli
	push $0
	push $45
	jmp irq_common_stub
irq14:
	cli
	push $0
	push $46
	jmp irq_common_stub
irq15:
	cli
	push $0
	push $47
	jmp irq_common_stub

irq_common_stub:
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
	mov %ax, %ds
	mov %esp, %eax
	push %eax
	mov irq_handler, %eax
	call *%eax
	pop %eax
	pop %gs
	pop %fs
	pop %es
	pop %ds
	popa
	add $8, %esp
	iret

.text
.globl Halt
Halt:
	hlt
	ret
