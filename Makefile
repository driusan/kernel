AS=/home/driusan/opt/cross/bin/i686-elf-as
CC=/home/driusan/opt/cross/bin/i686-elf-gcc
GO=/home/driusan/opt/cross/bin/i686-elf-gccgo
LD=/home/driusan/opt/cross/bin/i686-elf-gcc

ASMOBJS=boot.o interrupts/interrupts.o asm/int.o descriptortables/dt.o
COBJS=libg/golang.o libg/go-type-error.o libg/go-type-identity.o libg/go-strcmp.o \
	libg/kernel.o libg/go-runtime-error.o libg/go-type-string.o \
	libg/go-type-interface.o \
	libg/go-typedesc-equal.o \
	libg/mem.o \
	libg/stubs.o \
	asm/inout.o \
	memory/cpaging.o interrupts/irq.o interrupts/isrs.o
LIBGPKGSRC=libg/print.go
GOSRC=kernel.go keyboard.go timer.go
ASMPKGSRC=asm/inout.go
PCIPKGSRC=pci/pci.go pci/class.go pci/header.go
INTERRUPTSPKGSRC=interrupts/isrs.go interrupts/irq.go
DTABLEPKGSRC=descriptortables/gdt.go descriptortables/idt.go
PS2PKGSRC=input/ps2/keyboard.go input/ps2/mouse.go
ACPIPKGSRC=acpi/find.go
IDEPKGSRC=ide/identify.go ide/drive.go
TERMINALPKGSRC=terminal/print.go terminal/terminal.go
MBRPKGSRC=mbr/mbr.go

all: myos.bin

clean:
	rm -f *.o myos.bin libg/*.o asm/*.o pci/*.o interrupts/*.o ide/*.o	

interrupts.o: ${INTERRUPTSPKGSRC} interrupts/irq.o interrupts/isrs.o descriptortables.o
	${GO} -I`go env GOPATH`/src -c interrupts/*.go -o interrupts.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/interrupts

terminal.o: ${TERMINALPKGSRC} 
	${GO} -I`go env GOPATH`/src -c terminal/*.go -o terminal.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/terminal

mbr.o: ${MBRPKGSRC} 
	${GO} -I`go env GOPATH`/src -c mbr/*.go -o mbr.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/mbr

libg.o: ${LIBGPKGSRC} terminal.o asm.o
	${GO} -I`go env GOPATH`/src -c libg/*.go -o libg.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/libg

descriptortables.o: ${DTABLEPKGSRC} descriptortables/dt.o
	${GO} -I`go env GOPATH`/src -c descriptortables/*.go -o descriptortables.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/descriptortables

asm.o: ${ASMPKGSRC} asm/inout.o asm/int.o
	${GO} -I`go env GOPATH`/src -c asm/*.go -o asm.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/asm

ide.o: ${IDEPKGSRC} asm.o
	${GO} -I`go env GOPATH`/src -c ide/*.go -o ide.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/ide

pci.o: ${PCIPKGSRC} terminal.o
	${GO} -I`go env GOPATH`/src -c pci/*.go -o pci.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/pci

memory.o: ${MEMPKGSRC} memory/cpaging.o
	${GO} -I`go env GOPATH`/src -c memory/*.go -o memory.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/memory

input/ps2.o: ${PS2PKGSRC}
	${GO} -I`go env GOPATH`/src -c input/ps2/*.go -o input/ps2.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/input/ps2

acpi.o: ${ACPIPKGSRC}
	${GO} -I`go env GOPATH`/src -c acpi/*.go -o acpi.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel/acpi

%.o: %.s
	${AS} $< -o $@

%.o: %.c  
	${CC} -c $< -o $@ -std=gnu99 -ffreestanding -fno-inline-small-functions -Wall -Wextra

kernel.o: $(GOSRC) asm.o pci.o interrupts.o descriptortables.o memory.o input/ps2.o acpi.o ide.o terminal.o mbr.o
	${GO} -I. -I`go env GOPATH`/src -c *.go -o kernel.o -Wall -Wextra -fgo-pkgpath=github.com/driusan/kernel

myos.bin: $(ASMOBJS) $(COBJS) kernel.o asm.o pci.o interrupts.o descriptortables.o memory.o input/ps2.o acpi.o ide.o libg.o terminal.o mbr.o
	${LD} -T linker.ld -o myos.bin -ffreestanding -nostdlib *.o libg/*.o asm/*.o interrupts/*.o memory/*.o descriptortables/*.o input/*.o -lgcc

run: myos.bin
	# qemu-system-x86_64 -m 4G -kernel myos.bin -d int -no-reboot 2>error
	qemu-system-x86_64 -m 4G -show-cursor -hda test.img -kernel myos.bin -no-reboot

