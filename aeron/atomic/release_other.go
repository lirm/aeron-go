//go:build !amd64 && !arm64

package atomic

import "sync/atomic"

// On other architectures fall back to atomic.StoreInt32, which is at least as strong as a release
// store and a StoreStore barrier.

func storeInt32Release(addr *int32, value int32) {
	atomic.StoreInt32(addr, value)
}

func storeInt32StoreFence(addr *int32, value int32) {
	atomic.StoreInt32(addr, value)
}
