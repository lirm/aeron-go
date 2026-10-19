/*
Copyright 2016 Stanislav Liberman
Copyright 2020 Evan Wies

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

package util

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, uint32(512), SemanticVersionCompose(0, 2, 0))
}

func TestMemcpy(t *testing.T) {
	for length := 0; length <= 300; length++ {
		src := make([]byte, length+16)
		for i := range src {
			src[i] = byte(i*31 + 7)
		}
		for _, srcOffset := range []int{0, 1, 3, 8} {
			dest := make([]byte, length+32)
			destOffset := 5
			var srcPtr, destPtr uintptr
			if length > 0 {
				srcPtr = uintptr(unsafe.Pointer(&src[srcOffset]))
				destPtr = uintptr(unsafe.Pointer(&dest[destOffset]))
			}
			Memcpy(destPtr, srcPtr, int32(length))
			assert.Equal(t, src[srcOffset:srcOffset+length], dest[destOffset:destOffset+length],
				"length %d, src offset %d", length, srcOffset)
			for i := 0; i < destOffset; i++ {
				assert.Zero(t, dest[i], "length %d: byte %d before dest changed", length, i)
			}
			for i := destOffset + length; i < len(dest); i++ {
				assert.Zero(t, dest[i], "length %d: byte %d after dest changed", length, i)
			}
		}
	}
}

func BenchmarkMemcpy(b *testing.B) {
	for _, length := range []int{32, 100, 250, 1024} {
		b.Run(fmt.Sprintf("length=%d", length), func(b *testing.B) {
			src := make([]byte, length)
			dest := make([]byte, length)
			srcPtr := uintptr(unsafe.Pointer(&src[0]))
			destPtr := uintptr(unsafe.Pointer(&dest[0]))
			b.SetBytes(int64(length))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Memcpy(destPtr, srcPtr, int32(length))
			}
		})
	}
}
