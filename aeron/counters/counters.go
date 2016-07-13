package counters

import (
	"github.com/lirm/aeron-go/aeron/util"
)

/*
   struct CounterValueDefn
   {
       std::int64_t counterValue;
       std::int8_t pad1[(2 * util::BitUtil::CACHE_LINE_LENGTH) - sizeof(std::int64_t)];
   };

   struct CounterMetaDataDefn
   {
       std::int32_t state;
       std::int32_t typeId;
       std::int8_t key[(2 * util::BitUtil::CACHE_LINE_LENGTH) - (2 * sizeof(std::int32_t))];
       std::int32_t labelLength;
       std::int8_t label[(2 * util::BitUtil::CACHE_LINE_LENGTH) - sizeof(std::int32_t)];
   };
*/

var ReaderConsts = struct {
	RECORD_UNUSED    int32
	RECORD_ALLOCATED int32
	RECORD_RECLAIMED int32

	COUNTER_LENGTH      int32
	METADATA_LENGTH     int32
	KEY_OFFSET          int32
	LABEL_LENGTH_OFFSET int32

	MAX_LABEL_LENGTH int32
	MAX_KEY_LENGTH   int32
}{
	0,
	1,
	-1,

	2 * util.CACHE_LINE_LENGTH,
	4 * util.CACHE_LINE_LENGTH,
	8,
	2 * util.CACHE_LINE_LENGTH,

	2*int32(util.CACHE_LINE_LENGTH) - 4,
	2*int32(util.CACHE_LINE_LENGTH) - 8,
}

func Offset(id int32) int32 {
	return id * ReaderConsts.COUNTER_LENGTH
}
