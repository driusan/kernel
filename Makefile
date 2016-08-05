AS=/home/driusan/opt/cross/bin/i686-elf-as
CC=/home/driusan/opt/cross/bin/i686-elf-gcc
GO=/home/driusan/opt/cross/bin/i686-elf-gccgo
LD=/home/driusan/opt/cross/bin/i686-elf-gcc

ASMOBJS=boot.o interrupts.o
COBJS=libg/golang.o libg/go-type-error.o libg/go-type-identity.o libg/go-strcmp.o \
	libg/kernel.o libg/go-runtime-error.o libg/go-type-string.o \
	libg/mem.o \
	asm/inout.o \
	cpaging.o irq.o isrs.o
GOSRC=kernel.go gdt.go idt.go isrs.go irq.go keyboard.go timer.go
ASMPKGSRC=asm/inout.go
PCIPKGSRC=asm/pci.go

all: myos.bin

clean:
	rm -f *.o myos.bin libg/*.o	

asm.o: ${ASMPKGSRC} asm/inout.o
	${GO} -c asm/*.go -o asm.o -Wall -Wextra -fgo-prefix=boot
pci.o: ${ASMPKGSRC}
	${GO} -c pci/*.go -o pci.o -Wall -Wextra -fgo-prefix=boot

%.o: %.s
	${AS} $< -o $@

%.o: %.c  
	${CC} -c $< -o $@ -std=gnu99 -ffreestanding -fno-inline-small-functions -Wall -Wextra

kernel.o: $(GOSRC) asm.o pci.o
	${GO} -c *.go -o kernel.o -Wall -Wextra -fgo-prefix=boot

myos.bin: $(ASMOBJS) $(COBJS) kernel.o asm.o pci.o
	${LD} -T linker.ld -o myos.bin -ffreestanding -nostdlib *.o libg/*.o asm/*.o -lgcc

run: myos.bin
	# qemu-system-x86_64 -m 4G -kernel myos.bin -d int -no-reboot 2>error
	qemu-system-x86_64 -m 4G -kernel myos.bin -no-reboot

