# Declare constants for the multiboot header.
.set ALIGN,    1<<0             # align loaded modules on page boundaries
.set MEMINFO,  1<<1             # provide memory map
.set FLAGS,    ALIGN | MEMINFO  # this is the Multiboot 'flag' field
.set MAGIC,    0x1BADB002       # 'magic number' lets bootloader find the header
.set CHECKSUM, -(MAGIC + FLAGS) # checksum of above, to prove we are multiboot

# Declare a multiboot header that marks the program as a kernel.
.section .multiboot
.align 4
.long MAGIC
.long FLAGS
.long CHECKSUM

# Allocate the initial stack.
.section .bootstrap_stack, "aw", @nobits
stack_bottom:
.skip 16384 # 16 KiB
stack_top:

# Preallocate pages used for paging. Don't hard-code addresses and assume they
# are available, as the bootloader might have loaded its multiboot structures or
# modules there. This lets the bootloader know it must avoid the addresses.
.section .bss, "aw", @nobits
	.align 4096
boot_page_directory:
	.skip 4096
boot_page_table:
	.skip 4096
	.skip 4096

# Further page tables may be required if the kernel grows beyond 3 MiB.

# The kernel entry point.
.section .text
.global _start
.type _start, @function
_start:
	# Map into the boot_page_table's physical address
	movl $(boot_page_table - 0xC0000000), %edi

	# Map 2048 4K pages (=8MB)
	movl $0, %esi
	movl $2048, %ecx

1:
	# Map from address physical address 0 
	cmpl $0, %esi
	jl 2f
	# cmpl $(_kernel_end - 0xC0000000), %esi
	cmpl $0x00800000, %esi
	jge 3f

	# Map physical address as "present, writable". Note that this maps
	# .text and .rodata as writable. Mind security and map them as non-writable.
	movl %esi, %edx
	orl $0x003, %edx
	movl %edx, (%edi)

2:
	# Size of page is 4096 bytes.
	addl $4096, %esi
	# Size of entries in boot_pagetab1 is 4 bytes.
	addl $4, %edi
	# Loop to the next entry if we haven't finished.
	loop 1b

3:
	# The page table is used at both page directory entry 0 (virtually from 0x0
	# to 0x3FFFFF) (thus identity mapping the kernel) and page directory entry
	# 768 (virtually from 0xC0000000 to 0xC03FFFFF) (thus mapping it in the
	# higher half). The kernel is identity mapped because enabling paging does
	# not change the next instruction, which continues to be physical. The CPU
	# would instead page fault if there was no identity mapping.

	# Map the page table to both virtual addresses 0x00000000 and 0xC0000000.
	movl $(boot_page_table - 0xC0000000 + 0x003), boot_page_directory - 0xC0000000 + 0
	movl $(boot_page_table - 0xC0000000 + 0x003 + 4096), boot_page_directory - 0xC0000000 + 1*4
	movl $(boot_page_table - 0xC0000000 + 0x003), boot_page_directory - 0xC0000000 + 768 * 4
	movl $(boot_page_table - 0xC0000000 + 0x003 + 4096), boot_page_directory - 0xC0000000 + 769 * 4

	# Set cr3 to the address of the boot_page_directory.
	movl $(boot_page_directory - 0xC0000000), %ecx
	movl %ecx, %cr3

	# Enable paging and the write-protect bit.
	movl %cr0, %ecx
	orl $0x80010000, %ecx
	movl %ecx, %cr0

	# Jump to higher half with an absolute jump. 
	lea 4f, %ecx
	jmp *%ecx

4:
postjump:
	# At this point, paging is fully set up and enabled.

	# Don't unmap the identity paging yet, we still need it to access
	# the multiboot header.
	# Unmap the identity mapping as it is now unnecessary. 
	#movl $0, boot_page_directory + 0
	# Reload crc3 to force a TLB flush so the changes to take effect.
	#movl %cr3, %ecx
	#movl %ecx, %cr3

	# Set up the stack.
	mov $stack_top, %esp

	# Enter the high-level kernel.
	push %ebx
	push %eax
	call github_com_driusan_kernel.KernelMain

	# Infinite loop if the system has nothing more to do.
	cli
1:	hlt
	jmp 1b

# Define a way for us to halt the system manually in case of a panic
.globl github_com_driusan_kernel.Halt
.globl halt
github_com_driusan_kernel.Halt:
halt:
	cli
	hlt
	jmp halt


