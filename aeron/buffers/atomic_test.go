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

package buffers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

var buffer *Atomic

func TestMakeAtomicBuffer(t *testing.T) {
	var b *Atomic

	b = MakeAtomic()
	t.Logf("buf: %v", b)

	b = MakeAtomic(nil)
	t.Logf("buf: %v", b)

	arr := make([]byte, 64)

	b = MakeAtomic(arr)
	t.Logf("buf: %v", b)

	b = MakeAtomic(arr, 64)
	t.Logf("buf: %v", b)

	b = MakeAtomic(unsafe.Pointer(&arr[0]))
	t.Logf("buf: %v", b)

	b = MakeAtomic(unsafe.Pointer(&arr[0]), 64)
	t.Logf("buf: %v", b)
}

func TestInit(t *testing.T) {
	buffer = new(Atomic)
	assert.Equal(t, 0, buffer.Capacity())

	bytes := make([]byte, 32)
	bufLen := len(bytes)

	buffer.Wrap(unsafe.Pointer(&bytes), bufLen)
	assert.Equal(t, bufLen, buffer.Capacity())
}

func TestGetAndAddInt64(t *testing.T) {
	buffer = MakeAtomic(make([]byte, 32), 32)
	buffer.Fill(0)

	assert.Equal(t, int64(0), buffer.GetAndAddInt64(0, 7))
	assert.Equal(t, int64(0), buffer.GetAndAddInt64(8, 7))
	assert.Equal(t, int64(0), buffer.GetAndAddInt64(16, 7))
	assert.Equal(t, int64(0), buffer.GetAndAddInt64(24, 7))

	assert.Equal(t, int64(7), buffer.GetAndAddInt64(0, 7))
	assert.Equal(t, int64(7), buffer.GetAndAddInt64(8, 7))
	assert.Equal(t, int64(7), buffer.GetAndAddInt64(16, 7))
	assert.Equal(t, int64(7), buffer.GetAndAddInt64(24, 7))

	assert.Equal(t, int64(14), buffer.GetAndAddInt64(0, 7))
	assert.Equal(t, int64(14), buffer.GetAndAddInt64(8, 7))
	assert.Equal(t, int64(14), buffer.GetAndAddInt64(16, 7))
	assert.Equal(t, int64(14), buffer.GetAndAddInt64(24, 7))
}

func TestPutInt64Ordered(t *testing.T) {
	buffer = MakeAtomic(make([]byte, 32), 32)
	buffer.Fill(0)

	buffer.PutInt64Ordered(1, 31415)
	assert.Equal(t, int64(31415), buffer.GetInt64Volatile(1))
	assert.NotEqual(t, int64(31415), buffer.GetInt64Volatile(2))
	assert.NotEqual(t, int64(31415), buffer.GetInt64Volatile(0))
}

func TestWrap(t *testing.T) {
	bytes := make([]byte, 32)
	ptr := unsafe.Pointer(&bytes[0])
	buffer = MakeAtomic(bytes, 32)
	buffer.Fill(0)
	t.Logf("buf: %v", buffer)

	buffer.PutInt64Ordered(1, 31415)
	assert.Equal(t, int64(31415), buffer.GetInt64Volatile(1))
	assert.NotEqual(t, int64(31415), buffer.GetInt64Volatile(2))
	assert.NotEqual(t, int64(31415), buffer.GetInt64Volatile(0))

	newPtr := unsafe.Pointer(uintptr(ptr) + uintptr(1))
	t.Logf("Old pointer: %v; new pointer: %v", ptr, newPtr)
	var newLen int32 = 31
	buffer.Wrap(newPtr, newLen)
	assert.Equal(t, newLen, buffer.Capacity())

	assert.NotEqual(t, int64(31415), buffer.GetInt64Volatile(1))
	assert.Equal(t, int64(31415), buffer.GetInt64Volatile(0))
}
