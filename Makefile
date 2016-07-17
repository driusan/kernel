AS=/home/driusan/opt/cross/bin/i686-elf-as
CC=/home/driusan/opt/cross/bin/i686-elf-gcc
LD=/home/driusan/opt/cross/bin/i686-elf-gcc

all: boot.o kernel.o myos.bin

clean:
	rm boot.o kernel.o myos.bin	

boot.o:
	${AS} boot.s -o boot.o

kernel.o:
	${CC} -c kernel.c -o kernel.o -std=gnu99 -ffreestanding -O2 -Wall -Wextra

myos.bin:
	${LD} -T linker.ld -o myos.bin -ffreestanding -O2 -nostdlib boot.o kernel.o -lgcc

run: myos.bin
	qemu-system-x86_64 -kernel myos.bin

