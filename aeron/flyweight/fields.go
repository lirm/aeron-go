/*
Copyright 2016-2018 Stanislav Liberman
Copyright 2022 Steven Stern

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

package flyweight

import (
	syncatomic "sync/atomic"
	"unsafe"

	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/util"
)

// Field is the interface for a field in a flyweight wrapper. It expects a preallocated buffer and offset into it, as
// arguments.
type Field interface {
	Wrap(buffer *atomic.Buffer, offset int)
	Get() interface{}
}

// Int32Field is int32 field for flyweight
type Int32Field struct {
	offset unsafe.Pointer
}

func (fld *Int32Field) Wrap(buffer *atomic.Buffer, offset int) int {
	atomic.BoundsCheck(int32(offset), 4, buffer.Capacity())

	fld.offset = unsafe.Pointer(uintptr(buffer.Ptr()) + uintptr(offset))
	return 4
}

func (fld *Int32Field) Get() int32 {
	return *(*int32)(fld.offset)
}

func (fld *Int32Field) Set(value int32) {
	*(*int32)(fld.offset) = value
}

func (fld *Int32Field) CAS(curValue, newValue int32) bool {
	n := syncatomic.CompareAndSwapInt32((*int32)(fld.offset), curValue, newValue)
	return n
}

// Int64Field is int64 field for flyweight
type Int64Field struct {
	offset unsafe.Pointer
}

func (fld *Int64Field) Wrap(buffer *atomic.Buffer, offset int) int {
	atomic.BoundsCheck(int32(offset), 8, buffer.Capacity())

	fld.offset = unsafe.Pointer(uintptr(buffer.Ptr()) + uintptr(offset))
	return 8
}

func (fld *Int64Field) Get() int64 {
	return *(*int64)(fld.offset)
}

func (fld *Int64Field) Set(value int64) {
	*(*int64)(fld.offset) = value
}

func (fld *Int64Field) GetAndAddInt64(value int64) int64 {
	n := syncatomic.AddInt64((*int64)(fld.offset), value)
	return n - value
}

func (fld *Int64Field) CAS(curValue, newValue int64) bool {
	n := syncatomic.CompareAndSwapInt64((*int64)(fld.offset), curValue, newValue)
	return n
}

// StringField is string field for flyweight
type StringField struct {
	lenOffset  unsafe.Pointer
	dataOffset unsafe.Pointer
	fly        Flyweight
}

func (fld *StringField) Wrap(buffer *atomic.Buffer, offset int, fly Flyweight, align bool) int {

	off := int32(offset)
	offsetAdjustment := 0
	if align {
		off = util.AlignInt32(int32(offset), 4)
		offsetAdjustment = int(off) - offset
	}

	atomic.BoundsCheck(off, 4, buffer.Capacity())

	fld.lenOffset = unsafe.Pointer(uintptr(buffer.Ptr()) + uintptr(off))
	l := *(*int32)(fld.lenOffset)

	atomic.BoundsCheck(off+4, l, buffer.Capacity())

	fld.fly = fly
	fld.dataOffset = unsafe.Pointer(uintptr(buffer.Ptr()) + uintptr(off+4))
	return 4 + int(l) + offsetAdjustment
}

func (fld *StringField) Get() string {
	length := *(*int32)(fld.lenOffset)

	bArr := make([]byte, length)
	for ix := 0; ix < int(length); ix++ {
		uptr := unsafe.Pointer(uintptr(fld.dataOffset) + uintptr(ix))
		bArr[ix] = *(*uint8)(uptr)
	}

	return string(bArr)
}

func (fld *StringField) Set(value string) {
	length := int32(len(value))
	prevLen := *(*int32)(fld.lenOffset)
	*(*int32)(fld.lenOffset) = length

	bArr := []byte(value)
	srcUptr := unsafe.Pointer(&bArr[0])

	util.Memcpy(uintptr(fld.dataOffset), uintptr(srcUptr), length)

	size := fld.fly.Size()
	size -= int(prevLen)
	size += int(length)
	fld.fly.SetSize(size)
}

type RawDataField struct {
	buf atomic.Buffer
}

func (f *RawDataField) Wrap(buffer *atomic.Buffer, offset int, length int32) int {
	ptr := uintptr(buffer.Ptr()) + uintptr(offset)
	f.buf.Wrap(unsafe.Pointer(ptr), length)

	return int(length)
}

func (f *RawDataField) Get() *atomic.Buffer {
	return &f.buf
}

type Padding struct {
	raw RawDataField
}

// Wrap for padding takes size to pass this particular position to.  That size will be rounded down to match alignment.
func (f *Padding) Wrap(buffer *atomic.Buffer, offset int, size int32, alignment int32) int {
	maxLength := int32(offset) + size
	newLen := maxLength - maxLength%alignment - int32(offset)
	if newLen < 0 {
		newLen += alignment
	}

	return f.raw.Wrap(buffer, offset, newLen)
}

func (f *Padding) Get() *atomic.Buffer {
	return f.raw.Get()
}

type LengthAndRawDataField struct {
	lenOffset unsafe.Pointer
	buf       atomic.Buffer
}

func (fld *LengthAndRawDataField) Length() int32 {
	return *(*int32)(fld.lenOffset)
}

func (fld *LengthAndRawDataField) SetLength(length int32) {
	*(*int32)(fld.lenOffset) = length
}

func (fld *LengthAndRawDataField) Wrap(buffer *atomic.Buffer, rawOffset int) int {
	offset := util.AlignInt32(int32(rawOffset), 4)
	offsetAdjustment := int(offset) - rawOffset

	atomic.BoundsCheck(offset, 4, buffer.Capacity())
	fld.lenOffset = unsafe.Pointer(uintptr(buffer.Ptr()) + uintptr(offset))

	atomic.BoundsCheck(offset+4, fld.Length(), buffer.Capacity())
	ptr := uintptr(buffer.Ptr()) + uintptr(offset+4)
	fld.buf.Wrap(unsafe.Pointer(ptr), buffer.Capacity()-offset)

	return 4 + int(fld.Length()) + offsetAdjustment
}

func (fld *LengthAndRawDataField) CopyBuffer(buffer *atomic.Buffer, offset int32, length int32) {
	fld.SetLength(length)
	fld.buf.PutBytes(0, buffer, offset, length)
}

func (fld *LengthAndRawDataField) GetAsBuffer() *atomic.Buffer {
	return &fld.buf
}

func (fld *LengthAndRawDataField) GetAsASCII() string {
	bArr := make([]byte, fld.Length())
	fld.buf.GetBytes(0, bArr)
	return string(bArr)
}
