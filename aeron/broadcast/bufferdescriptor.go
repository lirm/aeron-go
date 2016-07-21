package broadcast

import (
	"github.com/lirm/aeron-go/aeron/util"
	"log"
)

var BufferDescriptor = struct {
	TAIL_INTENT_COUNTER_OFFSET int32
	TAIL_COUNTER_OFFSET        int32
	LATEST_COUNTER_OFFSET      int32
	TRAILER_LENGTH             int32
}{
	0,
	util.SIZEOF_INT64,
	util.SIZEOF_INT64 * 2,
	util.CACHE_LINE_LENGTH * 2,
}

func CheckCapacity(capacity int32) {
	if !util.IsPowerOfTwo(capacity) {
		log.Fatalf("Capacity must be a positive power of 2 + TRAILER_LENGTH: capacity=%d", capacity)
	}
}
