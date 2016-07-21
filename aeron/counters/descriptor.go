package counters

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"unsafe"
)

var Descriptor = struct {
	CNC_FILE                     string
	CNC_VERSION                  int32
	VERSION_AND_META_DATA_LENGTH int
}{
	"cnc.dat",
	5,
	int(util.CACHE_LINE_LENGTH * 2),
}

const (
	METADATA_OFFSET_CNCVERSION     uintptr = 0
	METADATA_OFFSET_DRIVERBUFLEN   uintptr = 4
	METADATA_OFFSET_CLIENTSBUFLEN  uintptr = 8
	METADATA_OFFSET_METADATABUFLEN uintptr = 12
	METADATA_OFFSET_VALUESBUFLEN   uintptr = 16
	METADATA_OFFSET_LIVELINESSTO   uintptr = 20
	METADATA_OFFSET_ERRORLOGBUFLEN uintptr = 28
)

type MetaDataDefn struct {
	cncVersion                  int32
	toDriverBufferLength        int32
	toClientsBufferLength       int32
	counterMetadataBufferLength int32
	counterValuesBufferLength   int32
	clientLivenessTimeout       int64
	errorLogBufferLength        int32
}

func CreateToDriverBuffer(cncFile *memmap.File) *buffers.Atomic {

	toDriverBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), METADATA_OFFSET_DRIVERBUFLEN)

	ptr := uintptr(cncFile.GetMemoryPtr()) + uintptr(Descriptor.VERSION_AND_META_DATA_LENGTH)

	return buffers.MakeAtomic(unsafe.Pointer(ptr), toDriverBufferLength)
}

func CreateToClientsBuffer(cncFile *memmap.File) *buffers.Atomic {
	metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())
	// log.Printf("createToClientsBuffer: MetaData: %v\n", metaData)

	toDriverBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), METADATA_OFFSET_DRIVERBUFLEN)
	ptr := uintptr(cncFile.GetMemoryPtr()) + uintptr(Descriptor.VERSION_AND_META_DATA_LENGTH) + uintptr(toDriverBufferLength)

	return buffers.MakeAtomic(unsafe.Pointer(ptr), metaData.toClientsBufferLength)
}

func CreateCounterMetadataBuffer(cncFile *memmap.File) *buffers.Atomic {
	metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())
	// log.Printf("createCounterMetadataBuffer: MetaData: %v\n", metaData)

	ptr := uintptr(cncFile.GetMemoryPtr()) + uintptr(Descriptor.VERSION_AND_META_DATA_LENGTH) + uintptr(metaData.toDriverBufferLength) + uintptr(metaData.toClientsBufferLength)

	return buffers.MakeAtomic(unsafe.Pointer(ptr), metaData.counterMetadataBufferLength)
}

func CreateCounterValuesBuffer(cncFile *memmap.File) *buffers.Atomic {
	metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())
	// log.Printf("createCounterValuesBuffer: MetaData: %v\n", metaData)

	ptr := uintptr(cncFile.GetMemoryPtr()) +
		uintptr(Descriptor.VERSION_AND_META_DATA_LENGTH) +
		uintptr(metaData.toDriverBufferLength) +
		uintptr(metaData.toClientsBufferLength) +
		uintptr(metaData.counterMetadataBufferLength)

	return buffers.MakeAtomic(unsafe.Pointer(ptr), metaData.counterValuesBufferLength)
}

func CreateErrorLogBuffer(cncFile *memmap.File) *buffers.Atomic {
	// metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())
	// log.Printf("createErrorLogBuffer: MetaData: %v\n", metaData)

	toDriverBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), METADATA_OFFSET_DRIVERBUFLEN)
	toClientsBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), METADATA_OFFSET_CLIENTSBUFLEN)
	counterMetadataBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), METADATA_OFFSET_METADATABUFLEN)
	counterValuesBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), METADATA_OFFSET_VALUESBUFLEN)
	errorLogBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), METADATA_OFFSET_ERRORLOGBUFLEN)

	ptr := uintptr(cncFile.GetMemoryPtr()) +
		uintptr(Descriptor.VERSION_AND_META_DATA_LENGTH) +
		uintptr(toDriverBufferLength) +
		uintptr(toClientsBufferLength) +
		uintptr(counterMetadataBufferLength) +
		uintptr(counterValuesBufferLength)

	return buffers.MakeAtomic(unsafe.Pointer(ptr), errorLogBufferLength)
}

func CncVersion(cncFile *memmap.File) int32 {
	metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())

	return metaData.cncVersion
}

func ClientLivenessTimeout(cncFile *memmap.File) int64 {
	ptr := uintptr(cncFile.GetMemoryPtr()) + METADATA_OFFSET_LIVELINESSTO

	return *(*int64)(unsafe.Pointer(ptr))
}

func readInt32FromPointer(basePtr unsafe.Pointer, offset uintptr) int32 {
	ptr := uintptr(basePtr) + offset
	return *(*int32)(unsafe.Pointer(ptr))
}
