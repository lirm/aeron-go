package buffers

import "unsafe"

type Claim struct {
	buffer Atomic
}

func (c *Claim) Wrap(buf *Atomic, offset, length int32) {
	buf.BoundsCheck(offset, length)
	ptr := unsafe.Pointer(uintptr(buf.Ptr()) + uintptr(offset))
	c.buffer.Wrap(ptr, length)
}
