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

var DataFrameHeader = struct {
	FrameLengthFieldOffset   int32
	VersionFieldOffset       int32
	FlagsFieldOffset         int32
	TypeFieldOffset          int32
	TermOffsetFieldOffset    int32
	SessionIDFieldOffset     int32
	StreamIDFieldOffset      int32
	TermIDFieldOffset        int32
	ReservedValueFieldOffset int32
	DataOffset               int32

	Length int32

	TypePad   uint16
	TypeData  uint16
	TypeNAK   uint16
	TypeSM    uint16
	TypeErr   uint16
	TypeSetup uint16
	TypeExt   uint16

	CurrentVersion int8
}{
	0,
	4,
	5,
	6,
	8,
	12,
	16,
	20,
	24,
	32,

	32,

	0x00,
	0x01,
	0x02,
	0x03,
	0x04,
	0x05,
	0xFFFF,

	0x0,
}
