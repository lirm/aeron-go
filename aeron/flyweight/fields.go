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

package flyweight

//#include <string.h>
import "C"

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"unsafe"
)

type Field interface {
	Wrap(buffer *buffers.Atomic, offset int)
	Get() interface{}
}

type Int32Field struct {
	offset unsafe.Pointer
}

func (fld *Int32Field) Wrap(buffer *buffers.Atomic, offset int) int {
	buffer.BoundsCheck(int32(offset), 4)

	fld.offset = unsafe.Pointer(uintptr(buffer.Ptr()) + uintptr(offset))
	return 4
}

func (fld *Int32Field) Get() int32 {
	return *(*int32)(fld.offset)
}

func (fld *Int32Field) Set(value int32) {
	*(*int32)(fld.offset) = value
}

type Int64Field struct {
	offset unsafe.Pointer
}

func (fld *Int64Field) Wrap(buffer *buffers.Atomic, offset int) int {
	buffer.BoundsCheck(int32(offset), 8)

	fld.offset = unsafe.Pointer(uintptr(buffer.Ptr()) + uintptr(offset))
	return 8
}

func (fld *Int64Field) Get() int64 {
	return *(*int64)(fld.offset)
}

func (fld *Int64Field) Set(value int64) {
	*(*int64)(fld.offset) = value
}

type StringField struct {
	lenOffset  unsafe.Pointer
	dataOffset unsafe.Pointer
	length     *int
}

func (fld *StringField) Wrap(buffer *buffers.Atomic, offset int, length *int) int {
	buffer.BoundsCheck(int32(offset), 4)

	fld.lenOffset = unsafe.Pointer(uintptr(buffer.Ptr()) + uintptr(offset))
	len := *(*int32)(fld.lenOffset)

	buffer.BoundsCheck(int32(offset)+4, len)

	fld.length = length
	fld.dataOffset = unsafe.Pointer(uintptr(buffer.Ptr()) + uintptr(offset+4))
	return 4 + int(len)
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

	//log.Printf("length: %d, srcUptr: %v, destPtr: %v", length, srcUptr, fld.dataOffset)

	C.memcpy(fld.dataOffset, unsafe.Pointer(srcUptr), C.size_t(length))

	*fld.length -= int(prevLen)
	*fld.length += int(length)
}
