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

.text
.globl __go_alloc
__go_alloc:
	JMP github_com_driusan_kernel_libg.GoAlloc

.text
.globl __go_free
__go_free:
	JMP github_com_driusan_kernel_libg.GoFree

.text
.globl malloc
malloc:
	JMP github_com_driusan_kernel_libg.GoAlloc

.text
.globl free
free:
	JMP github_com_driusan_kernel_libg.GoFree

.text
.globl __go_print_bool
__go_print_bool:
	JMP github_com_driusan_kernel_libg.GoPrintBool

.text
.globl pagingInitialized
pagingInitialized:
	JMP github_com_driusan_kernel_memory.IsPagingInitialized

#.text
#.globl runtime_panicstring
#runtime_panicstring:
#	JMP github_com_driusan_kernel_libg.GoRuntimePanicString

# This is used internal by the Go runtime, to resolve
# if x, ok :=  map[idx]; ok { .. } type statements
# It's assumed as part of runtime, but we don't have the runtime available
# need to figure out how to compile .goc files with gccgo before we can
# use the version from gccgo frontend..
.text
.globl runtime.mapaccess2
runtime.mapaccess2:
	JMP mapaccess2

# Similarly used by for x, y := range map
.text
.globl runtime.mapiterinit
runtime.mapiterinit:
	JMP __go_mapiterinit
.text
.globl runtime.mapiter2
runtime.mapiter2:
	JMP __go_mapiter2

.text
.globl runtime.mapiter1
runtime.mapiter1:
	JMP __go_mapiter1

.text
.globl runtime.mapiternext
runtime.mapiternext:
	JMP __go_mapiternext

.text
.globl runtime.stringiter
runtime.stringiter:
	JMP stringiter
.text
.globl runtime.stringiter2
runtime.stringiter2:
	JMP stringiter2


