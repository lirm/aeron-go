#include "textflag.h"

// func storeInt32Release(addr *int32, value int32)
TEXT ·storeInt32Release(SB), NOSPLIT, $0-12
	MOVQ addr+0(FP), AX
	MOVL value+8(FP), BX
	MOVL BX, 0(AX)
	RET

// func storeInt32StoreFence(addr *int32, value int32)
// On amd64 (TSO) stores are already visible in program order, so a MOV is enough.
TEXT ·storeInt32StoreFence(SB), NOSPLIT, $0-12
	MOVQ addr+0(FP), AX
	MOVL value+8(FP), BX
	MOVL BX, 0(AX)
	RET
