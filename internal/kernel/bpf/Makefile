CLANG ?= clang

all:
	$(CLANG) \
		-target bpf \
		-O2 \
		-g \
		-c sockmap.bpf.c \
		-o sockmap.bpf.o


clean:
	rm -f sockmap.bpf.o