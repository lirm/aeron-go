package logbuffer

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"unsafe"
)

type Header struct {
	buffer              buffers.Atomic
	offset              int32
	initialTermId       int32
	positionBitsToShift int32
}

func (hdr *Header) Wrap(ptr unsafe.Pointer, length int32) *Header {
	hdr.buffer.Wrap(ptr, length)
	return hdr
}

func (hdr *Header) SetOffset(offset int32) *Header {
	hdr.offset = offset
	return hdr
}

func (hdr *Header) SetInitialTermId(initialTermId int32) *Header {
	hdr.initialTermId = initialTermId
	return hdr
}

func (hdr *Header) SetPositionBitsToShift(positionBitsToShift int32) *Header {
	hdr.positionBitsToShift = positionBitsToShift
	return hdr
}
