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
	FRAME_LENGTH_FIELD_OFFSET   int32
	VERSION_FIELD_OFFSET        int32
	FLAGS_FIELD_OFFSET          int32
	TYPE_FIELD_OFFSET           int32
	TERM_OFFSET_FIELD_OFFSET    int32
	SESSION_ID_FIELD_OFFSET     int32
	STREAM_ID_FIELD_OFFSET      int32
	TERM_ID_FIELD_OFFSET        int32
	RESERVED_VALUE_FIELD_OFFSET int32
	DATA_OFFSET                 int32

	LENGTH int32

	HDR_TYPE_PAD   uint16
	HDR_TYPE_DATA  uint16
	HDR_TYPE_NAK   uint16
	HDR_TYPE_SM    uint16
	HDR_TYPE_ERR   uint16
	HDR_TYPE_SETUP uint16
	HDR_TYPE_EXT   uint16

	CURRENT_VERSION int8
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
