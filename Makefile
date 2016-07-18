AS=/home/driusan/opt/cross/bin/i686-elf-as
CC=/home/driusan/opt/cross/bin/i686-elf-gcc
GO=/home/driusan/opt/cross/bin/i686-elf-gccgo
LD=/home/driusan/opt/cross/bin/i686-elf-gcc

all: boot.o kernel.o t.o myos.bin

clean:
	rm -f *.o myos.bin	

boot.o: boot.s
	${AS} boot.s -o boot.o

kernel.o: kernel.c
	${CC} -c kernel.c -o kernel.o -std=gnu99 -ffreestanding -fno-inline-small-functions -O0 -Wall -Wextra -fdump-go-spec=headers

t.o: t.go
	${GO} -c t.go -o t.o -O0 -Wall -Wextra -fgo-prefix=boot

myos.bin: kernel.o boot.o t.o
	${LD} -T linker.ld -o myos.bin -ffreestanding -O0 -nostdlib boot.o kernel.o t.o -lgcc

run: myos.bin
	qemu-system-x86_64 -kernel myos.bin

