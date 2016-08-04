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

package counters

import (
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/op/go-logging"
	"log"
)

var logger = logging.MustGetLogger("counters")

var ReaderConsts = struct {
	RecordUnused    int32
	RecordAllocated int32
	RecordReclaimed int32

	CounterLen     int32
	MetadataLen    int32
	KeyOffset      int32
	LabelLenOffset int32

	MaxLabelLen int32
	MaxKeyLen   int32
}{
	0,
	1,
	-1,

	2 * util.CacheLineLength,
	4 * util.CacheLineLength,
	8,
	2 * util.CacheLineLength,

	2*int32(util.CacheLineLength) - 4,
	2*int32(util.CacheLineLength) - 8,
}

func MapFile(filename string) *memmap.File {

	logger.Debugf("Trying to map file: %s", filename)
	cncBuffer, err := memmap.MapExisting(filename, 0, 0)
	if err != nil {
		log.Fatalf("Failed to map the file %s with %s", filename, err.Error())
	}

	cncVer := CncVersion(cncBuffer)
	logger.Debugf("Mapped %s for ver %d", filename, cncVer)

	if CurrentCncVersion != cncVer {
		log.Fatalf("aeron cnc file version not understood: version=%d", cncVer)
	}

	return cncBuffer
}
