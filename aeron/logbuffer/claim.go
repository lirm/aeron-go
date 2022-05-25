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
)

type Claim struct {
	buffer atomic.Buffer
}

func (c *Claim) Wrap(buf *atomic.Buffer, offset, length int32) {
	atomic.BoundsCheck(offset, length, buf.Capacity())
	ptr := unsafe.Pointer(uintptr(buf.Ptr()) + uintptr(offset))
	c.buffer.Wrap(ptr, length)
}

// The referenced buffer to be used.
//
// @return the referenced buffer to be used..
func (claim *Claim) Buffer() *atomic.Buffer {
	// TODO Should this be a pointer or a copy?
	return &claim.buffer
}

// offset in the buffer at which the claimed range begins.
//
// @return offset in the buffer at which the range begins.
func (claim *Claim) Offset() int32 {
	return DataFrameHeader.Length
}

// The length of the claimed range in the buffer.
//
// @return length of the range in the buffer.
func (claim *Claim) Length() int32 {
	return claim.buffer.Capacity() - DataFrameHeader.Length
}

// Get the value stored in the reserve space at the end of a data frame header.
//
// Note: The value is in {@link ByteOrder#LITTLE_ENDIAN} format.
//
// @return the value stored in the reserve space at the end of a data frame header.
// @see DataHeaderFlyweight
func (claim *Claim) ReservedValue() int64 {
	return claim.buffer.GetInt64(DataFrameHeader.ReservedValueFieldOffset)
}

// Write the provided value into the reserved space at the end of the data frame header.
//
// Note: The value will be written in {@link ByteOrder#LITTLE_ENDIAN} format.
//
// @param value to be stored in the reserve space at the end of a data frame header.
// @return this for fluent API semantics.
// @see DataHeaderFlyweight
func (claim *Claim) SetReservedValue(value int64) *Claim {
	claim.buffer.PutInt64(DataFrameHeader.ReservedValueFieldOffset, value)
	return claim
}

// Commit the message to the log buffer so that is it available to subscribers.
func (claim *Claim) Commit() {
	frameLength := claim.buffer.Capacity()

	claim.buffer.PutInt32Ordered(DataFrameHeader.FrameLengthFieldOffset, frameLength)
}

// Abort a claim of the message space to the log buffer so that the log can progress by ignoring this claim.
func (claim *Claim) Abort() {
	frameLength := claim.buffer.Capacity()

	claim.buffer.PutUInt16(DataFrameHeader.TypeFieldOffset, DataFrameHeader.TypePad)
	claim.buffer.PutInt32Ordered(DataFrameHeader.FrameLengthFieldOffset, frameLength)
}
