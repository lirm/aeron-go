package buffers

import "github.com/lirm/aeron-go/aeron/util"

type Handler func(int32, *Atomic, int32, int32)

var RingBufferDescriptor = struct {
	TAIL_POSITION_OFFSET       int32
	HEAD_CACHE_POSITION_OFFSET int32
	HEAD_POSITION_OFFSET       int32
	CORRELATION_COUNTER_OFFSET int32
	CONSUMER_HEARTBEAT_OFFSET  int32
	TRAILER_LENGTH             int32
}{
	util.CACHE_LINE_LENGTH * 2,
	util.CACHE_LINE_LENGTH * 4,
	util.CACHE_LINE_LENGTH * 6,
	util.CACHE_LINE_LENGTH * 8,
	util.CACHE_LINE_LENGTH * 10,
	util.CACHE_LINE_LENGTH * 12,
}
