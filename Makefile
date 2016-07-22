AS=/home/driusan/opt/cross/bin/i686-elf-as
CC=/home/driusan/opt/cross/bin/i686-elf-gcc
GO=/home/driusan/opt/cross/bin/i686-elf-gccgo
LD=/home/driusan/opt/cross/bin/i686-elf-gcc

ASMOBJS=boot.o
COBJS=libg/golang.o libg/kernel.o
GOOBJS=kernel.o itoa.o

all: myos.bin

clean:
	rm -f *.o myos.bin libg/*.o	

%.o: %.s
	${AS} $< -o $@

%.o: %.c  
	${CC} -c $< -o $@ -std=gnu99 -ffreestanding -fno-inline-small-functions -Wall -Wextra

%.o: %.go
	${GO} -c $< -o $@ -Wall -Wextra -fgo-prefix=boot

myos.bin: $(ASMOBJS) $(COBJS) $(GOOBJS)
	${LD} -T linker.ld -o myos.bin -ffreestanding -nostdlib libg/*.o *.o -lgcc

run: myos.bin
	qemu-system-x86_64 -kernel myos.bin

