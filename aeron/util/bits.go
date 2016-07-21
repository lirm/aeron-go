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
