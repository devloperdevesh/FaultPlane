CLANG ?= clang

all:
	clang \
		-target bpf \
		-O2 \
		-g \
		-I/usr/include \
		-I/usr/include/x86_64-linux-gnu \
		-c sockmap.bpf.c \
		-o sockmap.bpf.o

clean:
	rm -f sockmap.bpf.o
