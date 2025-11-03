GLOBL COLLIDINGSYMBOL<>(SB), 8, $32

// func subtract(x, y int64)
TEXT ·subtract(SB),$0-24
    MOVQ x+0(FP), BX
	MOVQ y+8(FP), BP
    SUBQ BP, BX
    MOVQ BX, ret+16(FP)
    RET
