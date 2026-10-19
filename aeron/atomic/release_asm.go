//go:build amd64 || arm64

package atomic

// storeInt32Release is a release store, like putIntOrdered in Agrona: earlier stores become visible
// before this one. On amd64 (TSO) it is a plain MOV, on arm64 it is STLR. Unlike atomic.StoreInt32
// (XCHG on amd64, a full barrier) the CPU does not wait for its store buffer to drain. The compiler
// does not inline calls to assembly functions and does not move memory accesses across them.
//
//go:noescape
func storeInt32Release(addr *int32, value int32)

// storeInt32StoreFence is a store followed by a StoreStore barrier, like putLongOrdered plus
// storeFence in the Java HeaderWriter: later stores do not become visible before this one. On amd64
// it is a MOV, on arm64 it is STR followed by DMB ISHST.
//
//go:noescape
func storeInt32StoreFence(addr *int32, value int32)
