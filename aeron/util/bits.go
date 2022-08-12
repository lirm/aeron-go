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

import (
	"fmt"
	"unsafe"
)

var i32 int32
var i64 int64

const (
	// CacheLineLength is a constant for the size of a CPU cache line
	CacheLineLength int32 = 64

	// SizeOfInt32 is a constant for the size of int32. Ha. Just for Clarity
	SizeOfInt32 int32 = int32(unsafe.Sizeof(i32))

	// SizeOfInt64 is a constant for the size of int64
	SizeOfInt64 int32 = int32(unsafe.Sizeof(i64))
)

// AlignInt32 will return a number rounded up to the alignment boundary
func AlignInt32(value, alignment int32) int32 {
	return (value + (alignment - 1)) & ^(alignment - 1)
}

// NumberOfTrailingZeroes is HD recipe for determining the number of leading zeros on 32 bit integer
func NumberOfTrailingZeroes(value uint32) uint8 {
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

// FastMod3 is HD recipe for faster division by 3 for 32 bit integers
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

// IsPowerOfTwo checks that the argument number is a power of two
func IsPowerOfTwo(value int64) bool {
	return value > 0 && ((value & (^value + 1)) == value)
}

// Memcpy will copy length bytes from pointer src to dest
//go:nocheckptr
func Memcpy(dest uintptr, src uintptr, length int32) {
	var i int32

	// batches of 8
	i8 := length >> 3
	for ; i < i8; i += 8 {
		destPtr := unsafe.Pointer(dest + uintptr(i))
		srcPtr := unsafe.Pointer(src + uintptr(i))

		*(*uint64)(destPtr) = *(*uint64)(srcPtr)
	}

	// batches of 4
	i4 := (length - i) >> 2
	for ; i < i4; i += 4 {
		destPtr := unsafe.Pointer(dest + uintptr(i))
		srcPtr := unsafe.Pointer(src + uintptr(i))

		*(*uint32)(destPtr) = *(*uint32)(srcPtr)
	}

	// remainder
	for ; i < length; i++ {
		destPtr := unsafe.Pointer(dest + uintptr(i))
		srcPtr := unsafe.Pointer(src + uintptr(i))

		*(*int8)(destPtr) = *(*int8)(srcPtr)
	}
}

func MemPrint(ptr uintptr, len int) string {
	var output string

	for i := 0; i < len; i += 1 {
		ptr := unsafe.Pointer(ptr + uintptr(i))
		output += fmt.Sprintf("%02x ", *(*int8)(ptr))
	}

	return output
}

func Print(bytes []byte) {
	for i, b := range bytes {
		if i > 0 && i%16 == 0 && i%32 != 0 {
			fmt.Print(" :  ")
		}
		if i > 0 && i%32 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("%02x ", b)
	}
	fmt.Print("\n")
}
