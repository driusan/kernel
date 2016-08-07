AS=/home/driusan/opt/cross/bin/i686-elf-as
CC=/home/driusan/opt/cross/bin/i686-elf-gcc
GO=/home/driusan/opt/cross/bin/i686-elf-gccgo
LD=/home/driusan/opt/cross/bin/i686-elf-gcc

ASMOBJS=boot.o interrupts/interrupts.o asm/int.s
COBJS=libg/golang.o libg/go-type-error.o libg/go-type-identity.o libg/go-strcmp.o \
	libg/kernel.o libg/go-runtime-error.o libg/go-type-string.o \
	libg/go-type-interface.o \
	libg/go-typedesc-equal.o \
	libg/mem.o \
	asm/inout.o \
	memory/cpaging.o interrupts/irq.o interrupts/isrs.o
GOSRC=kernel.go keyboard.go timer.go
ASMPKGSRC=asm/inout.go
PCIPKGSRC=pci/pci.go
INTERRUPTSPKGSRC=interrupts/isrs.go interrupts/irq.go
DTABLEPKGSRC=descriptortables/gdt.go descriptortables/idt.go

all: myos.bin

clean:
	rm -f *.o myos.bin libg/*.o asm/*.o pci/*.o interrupts/*.o	

interrupts.o: ${INTERRUPTSPKGSRC} interrupts/irq.o interrupts/isrs.o descriptortables.o
	${GO} -c interrupts/*.go -o interrupts.o -Wall -Wextra -fgo-prefix=boot

descriptortables.o: ${DTABLEPKGSRC}
	${GO} -c descriptortables/*.go -o descriptortables.o -Wall -Wextra -fgo-prefix=boot

asm.o: ${ASMPKGSRC} asm/inout.o asm/int.o
	${GO} -c asm/*.go -o asm.o -Wall -Wextra -fgo-prefix=boot

pci.o: ${PCIPKGSRC}
	${GO} -c pci/*.go -o pci.o -Wall -Wextra -fgo-prefix=boot

memory.o: ${MEMPKGSRC} memory/cpaging.o
	${GO} -c memory/*.go -o memory.o -Wall -Wextra -fgo-prefix=boot

%.o: %.s
	${AS} $< -o $@

%.o: %.c  
	${CC} -c $< -o $@ -std=gnu99 -ffreestanding -fno-inline-small-functions -Wall -Wextra

kernel.o: $(GOSRC) asm.o pci.o interrupts.o descriptortables.o memory.o
	${GO} -c *.go -o kernel.o -Wall -Wextra -fgo-prefix=boot

myos.bin: $(ASMOBJS) $(COBJS) kernel.o asm.o pci.o interrupts.o descriptortables.o memory.o
	${LD} -T linker.ld -o myos.bin -ffreestanding -nostdlib *.o libg/*.o asm/*.o interrupts/*.o memory/*.o -lgcc

run: myos.bin
	# qemu-system-x86_64 -m 4G -kernel myos.bin -d int -no-reboot 2>error
	qemu-system-x86_64 -m 4G -kernel myos.bin -no-reboot

