package buffers

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/util"
)

var RecordDescriptor = struct {
	HEADER_LENGTH       int32
	RECORD_ALIGNMENT    int32
	PADDING_MSG_TYPE_ID int32
	LENGTH_OFFSET       int32
	TYPE_OFFSET         int32
}{
	util.SIZEOF_INT32 * 2,
	util.SIZEOF_INT32 * 2,
	-1,
	0,
	4,
}

func LengthOffset(recordOffset int32) int32 {
	return recordOffset
}

func TypeOffset(recordOffset int32) int32 {
	return recordOffset + util.SIZEOF_INT32
}

func EncodedMsgOffset(recordOffset int32) int32 {
	return recordOffset + RecordDescriptor.HEADER_LENGTH
}

func MakeHeader(length, msgTypeId int32) int64 {
	return ((int64(msgTypeId) & 0xFFFFFFFF) << 32) | (int64(length) & 0xFFFFFFFF)
}

func RecordLength(header int64) int32 {
	return int32(header)
}

func MessageTypeId(header int64) int32 {
	return int32(header >> 32)
}

func CheckMsgTypeId(msgTypeId int32) {
	if msgTypeId < 1 {
		panic(fmt.Sprintf("Message type id must be greater than zero, msgTypeId=%d", msgTypeId))
	}
}
