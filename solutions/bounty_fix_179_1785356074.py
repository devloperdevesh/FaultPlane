//go:build amd64

#include "textflag.h"

// func cldemote(ptr unsafe.Pointer)
// Moves the cache line containing 'ptr' from L1/L2 cache to L3/LLC.
TEXT ·cldemote(SB), NOSPLIT, $0-8
    MOVQ ptr+0(FP), AX
    // Byte sequence for CLDEMOTE (AX) -> 0x0F 0x1C 0x00
    // ModRM: mod=00 (direct indirect [AX]), reg=000 (extension), rm=000 (AX)
    BYTE $0x0F; BYTE $0x1C; BYTE $0x00
    RET