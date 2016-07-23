AS=/home/driusan/opt/cross/bin/i686-elf-as
CC=/home/driusan/opt/cross/bin/i686-elf-gcc
GO=/home/driusan/opt/cross/bin/i686-elf-gccgo
LD=/home/driusan/opt/cross/bin/i686-elf-gcc

ASMOBJS=boot.o
COBJS=libg/golang.o libg/kernel.o cpaging.o

GOSRC=itoa.go kernel.go #test.go #//paging.go

all: myos.bin

clean:
	rm -f *.o myos.bin libg/*.o	

%.o: %.s
	${AS} $< -o $@

%.o: %.c  
	${CC} -c $< -o $@ -std=gnu99 -ffreestanding -fno-inline-small-functions -Wall -Wextra -O0

kernel.o: $(GOSRC)
	${GO} -c *.go -o kernel.o -Wall -Wextra -fgo-prefix=boot

myos.bin: $(ASMOBJS) $(COBJS) kernel.o
	${LD} -T linker.ld -o myos.bin -ffreestanding -nostdlib libg/*.o *.o -lgcc -O0

run: myos.bin
	qemu-system-x86_64 -kernel myos.bin

