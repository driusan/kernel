# GCCGo makes it nearly impossible to define functions
# without any package prefix, but it also tries to link
# to these functions. As a result, we define the symbol
# in assembly and then immediately jump to the prefixed
# symbol.
# Since the call pushed to the stack when calling these,
# and the function will RET at the end of it popping off
# the appropriate amounts, the stack should be maintained
# implicitly as long as the signatures match the memory
# sizes that would happen if these were defined in C.

.text
.globl __go_panic
__go_panic:
	JMP github_com_driusan_kernel_libg.GoPanic

.text
.globl __go_print_string
__go_print_string:
	JMP github_com_driusan_kernel_libg.GoPrintString

.text
.globl __go_print_nl
__go_print_nl:
	JMP github_com_driusan_kernel_libg.GoPrintNewline

.text
.globl __go_print_space
__go_print_space:
	JMP github_com_driusan_kernel_libg.GoPrintSpace

.text
.globl __go_print_int64
__go_print_int64:
	JMP github_com_driusan_kernel_libg.GoPrintInt64

.text
.globl __go_print_uint64
__go_print_uint64:
	JMP github_com_driusan_kernel_libg.GoPrintUint64


.text
.globl __go_print_pointer
__go_print_pointer:
	JMP github_com_driusan_kernel_libg.GoPrintPointer

#.text
#.globl runtime_panicstring
#runtime_panicstring:
#	JMP github_com_driusan_kernel_libg.GoRuntimePanicString


