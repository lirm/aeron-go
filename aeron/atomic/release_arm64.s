#include "textflag.h"

// func storeInt32Release(addr *int32, value int32)
TEXT ·storeInt32Release(SB), NOSPLIT, $0-12
	MOVD addr+0(FP), R0
	MOVW value+8(FP), R1
	STLRW R1, (R0)
	RET

// func storeInt32StoreFence(addr *int32, value int32)
// DMB ISHST: later stores do not become visible before this one.
TEXT ·storeInt32StoreFence(SB), NOSPLIT, $0-12
	MOVD addr+0(FP), R0
	MOVW value+8(FP), R1
	MOVW R1, (R0)
	DMB $0xa
	RET
