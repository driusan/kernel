// This package only exists to get rid of errors about import file 'C'
// not existing when using a gccgo elf cross-compiler. It does nothing,
// but without it either go fmt/go test will complain about using C without
// import "C" or i686-elf-gccgo will complain about not being able to find
// package C.
package C

type Nothing struct{}
