//go:build amd64

// TODO: VPCMPEQB and VPMOVMSKB are from SSE2 and always available in 64-bit
// VPBROADCASTB though is from AVX2, so we do need to detect with golang.org/x/sys/cpu

#include "textflag.h"

// NOSPLIT,NOFRAME?
// func matchMask(a uint8, b [16]uint8) uint16
TEXT ·matchMask(SB),$0-17
    VPBROADCASTB    a+0(FP),    X0
    VPCMPEQB        b+1(FP),    X0, X0
    VPMOVMSKB       X0,         AX
    MOVL            AX,         ret+24(FP)
    RET
