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
#.skip 131072 # 128 KiB
#.skip 655360 # 640KiB oughta be enough for anyone
stack_top:

# The linker script specifies _start as the entry point to the kernel and the
# bootloader will jump to this position once the kernel has been loaded. It
# doesn't make sense to return from this function as the bootloader is gone.
.section .text
.global _start
.type _start, @function
_start:
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

	# TODO: Setup paging and map the kernel into a higher memory address before
	# calling KernelMain

	# Add the multiboot info onto the stack as the first parameter to KernelMain
	push %ebx
	push %eax
	call github_com_driusan_kernel.KernelMain
	cli
1:	hlt
	jmp 1b

# Set the size of the _start symbol to the current location '.' minus its start.
# This is useful when debugging or when you implement call tracing.
.size _start, . - _start
.text

# Define a way for us to halt the system manually in case of a panic
.globl github_com_driusan_kernel.Halt
.globl halt
github_com_driusan_kernel.Halt:
halt:
	cli
	hlt
	jmp halt


