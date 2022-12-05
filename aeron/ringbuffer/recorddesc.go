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

package rb

import (
	"fmt"

	"github.com/lirm/aeron-go/aeron/util"
)

var RecordDescriptor = struct {
	HeaderLength     int32
	RecordAlignment  int32
	PaddingMsgTypeID int32
	LengthOffset     int32
	TypeOffset       int32
}{
	util.SizeOfInt32 * 2,
	util.SizeOfInt32 * 2,
	-1,
	0,
	4,
}

func LengthOffset(recordOffset int32) int32 {
	return recordOffset
}

func TypeOffset(recordOffset int32) int32 {
	return recordOffset + util.SizeOfInt32
}

func EncodedMsgOffset(recordOffset int32) int32 {
	return recordOffset + RecordDescriptor.HeaderLength
}

func makeHeader(length, msgTypeID int32) int64 {
	return ((int64(msgTypeID) & 0xFFFFFFFF) << 32) | (int64(length) & 0xFFFFFFFF)
}

func checkMsgTypeID(msgTypeID int32) {
	if msgTypeID < 1 {
		panic(fmt.Sprintf("Message type id must be greater than zero, msgTypeId=%d", msgTypeID))
	}
}
