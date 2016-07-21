package counters

import (
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/op/go-logging"
	"log"
)

var logger = logging.MustGetLogger("counters")

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

func MapCncFile(filename string) *memmap.File {

	logger.Debugf("Trying to map file: %s", filename)
	cncBuffer, err := memmap.MapExisting(filename, 0, 0)
	if err != nil {
		log.Fatalf("Failed to map the file %s with %s", filename, err.Error())
	}

	cncVer := CncVersion(cncBuffer)
	logger.Debugf("Mapped %s for ver %d", filename, cncVer)

	if Descriptor.CNC_VERSION != cncVer {
		log.Fatalf("aeron cnc file version not understood: version=%d", cncVer)
	}

	return cncBuffer
}
