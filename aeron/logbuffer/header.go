/*
Copyright 2016 Stanislav Liberman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logbuffer

import (
	"unsafe"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/util"
)

type Header struct {
	buffer              atomic.Buffer
	offset              int32
	initialTermID       int32
	positionBitsToShift int32
}

//go:norace
func (hdr *Header) Wrap(ptr unsafe.Pointer, length int32) *Header {
	hdr.buffer.Wrap(ptr, length)
	return hdr
}

// Position calculates the current position to which the image has advanced on
// reading this message.
func (hdr *Header) Position() int64 {
	resultingOffset := util.AlignInt32(hdr.Offset()+hdr.FrameLength(), FrameAlignment)
	return computePosition(hdr.TermId(), resultingOffset, hdr.positionBitsToShift, hdr.InitialTermId())
}

func (hdr *Header) Offset() int32 {
	return hdr.offset
}

func (hdr *Header) Flags() uint8 {
	return GetFlags(&hdr.buffer, hdr.offset)
}

func (hdr *Header) FrameLength() int32 {
	return GetFrameLength(&hdr.buffer, hdr.offset)
}

func (hdr *Header) TermId() int32 {
	return GetTermId(&hdr.buffer, hdr.offset)
}

func (hdr *Header) SessionId() int32 {
	return GetSessionId(&hdr.buffer, hdr.offset)
}

func (hdr *Header) StreamId() int32 {
	return GetStreamId(&hdr.buffer, hdr.offset)
}

func (hdr *Header) GetReservedValue() int64 {
	return GetReservedValue(&hdr.buffer, hdr.offset)
}

func (hdr *Header) SetOffset(offset int32) *Header {
	hdr.offset = offset
	return hdr
}

func (hdr *Header) InitialTermId() int32 {
	return hdr.initialTermID
}

func (hdr *Header) SetInitialTermID(initialTermID int32) *Header {
	hdr.initialTermID = initialTermID
	return hdr
}

func (hdr *Header) SetPositionBitsToShift(positionBitsToShift int32) *Header {
	hdr.positionBitsToShift = positionBitsToShift
	return hdr
}

func (hdr *Header) SetReservedValue(reservedValue int64) *Header {
	hdr.buffer.PutInt64(reservedValueOffset(hdr.offset), reservedValue)
	return hdr
}

func (hdr *Header) SetSessionId(value int32) *Header {
	hdr.buffer.PutInt32(sessionIdOffset(hdr.offset), value)
	return hdr
}

// computePosition computes the current position in absolute number of bytes.
func computePosition(activeTermId int32, termOffset int32, positionBitsToShift int32, initialTermId int32) int64 {
	termCount := activeTermId - initialTermId // copes with negative activeTermId on rollover
	return (int64(termCount) << uint32(positionBitsToShift)) + int64(termOffset)
}
