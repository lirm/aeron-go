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

package util

import "unsafe"

var i32 int32
var i64 int64

const (
	CACHE_LINE_LENGTH int32 = 64
	SIZEOF_INT32      int32 = int32(unsafe.Sizeof(i32))
	SIZEOF_INT64      int32 = int32(unsafe.Sizeof(i64))
)

func AlignInt32(value, alignment int32) int32 {
	return (value + (alignment - 1)) & ^(alignment - 1)
}

func NumberOfTrailingZeroes(value int32) uint8 {
	table := [32]uint8{
		0, 1, 2, 24, 3, 19, 6, 25,
		22, 4, 20, 10, 16, 7, 12, 26,
		31, 23, 18, 5, 21, 9, 15, 11,
		30, 17, 8, 14, 29, 13, 28, 27}

	if value == 0 {
		return 32
	}

	value = (value & -value) * 0x04D7651F

	return table[value>>27]
}

func FastMod3(value uint64) int32 {

	table := [62]int32{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2,
		0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2,
		0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2,
		0, 1, 2, 0, 1, 2, 0, 1}

	value = (value >> 16) + (value & 0xFFFF) // Max 0x1FFFE.
	value = (value >> 8) + (value & 0x00FF)  // Max 0x2FD.
	value = (value >> 4) + (value & 0x000F)  // Max 0x3D.
	return table[value]
}

func IsPowerOfTwo(value int32) bool {
	return value > 0 && ((value & (^value + 1)) == value)
}
