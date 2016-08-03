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
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"unsafe"
)

var Descriptor = struct {
	CncFile               string
	CncVersion            int32
	VersionAndMetaDataLen int
}{
	"cnc.dat",
	5,
	int(util.CacheLineLength * 2),
}

var MetadataOffsets = struct {
	cncVersion     uintptr
	driverBufLen   uintptr
	clientsBufLen  uintptr
	metaDataBufLen uintptr
	valuesBufLen   uintptr
	livelinessTo   uintptr
	errorLogBufLen uintptr
}{
	0,
	4,
	8,
	12,
	16,
	20,
	28,
}

type MetaDataDefn struct {
	cncVersion                  int32
	toDriverBufferLength        int32
	toClientsBufferLength       int32
	counterMetadataBufferLength int32
	counterValuesBufferLength   int32
	clientLivenessTimeout       int64
	errorLogBufferLength        int32
}

func CreateToDriverBuffer(cncFile *memmap.File) *atomic.Buffer {

	toDriverBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), MetadataOffsets.driverBufLen)

	ptr := uintptr(cncFile.GetMemoryPtr()) + uintptr(Descriptor.VersionAndMetaDataLen)

	return atomic.MakeBuffer(unsafe.Pointer(ptr), toDriverBufferLength)
}

func CreateToClientsBuffer(cncFile *memmap.File) *atomic.Buffer {
	metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())
	// log.Printf("createToClientsBuffer: MetaData: %v\n", metaData)

	toDriverBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), MetadataOffsets.driverBufLen)
	ptr := uintptr(cncFile.GetMemoryPtr()) + uintptr(Descriptor.VersionAndMetaDataLen) + uintptr(toDriverBufferLength)

	return atomic.MakeBuffer(unsafe.Pointer(ptr), metaData.toClientsBufferLength)
}

func CreateCounterMetadataBuffer(cncFile *memmap.File) *atomic.Buffer {
	metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())
	// log.Printf("createCounterMetadataBuffer: MetaData: %v\n", metaData)

	ptr := uintptr(cncFile.GetMemoryPtr()) + uintptr(Descriptor.VersionAndMetaDataLen) + uintptr(metaData.toDriverBufferLength) + uintptr(metaData.toClientsBufferLength)

	return atomic.MakeBuffer(unsafe.Pointer(ptr), metaData.counterMetadataBufferLength)
}

func CreateCounterValuesBuffer(cncFile *memmap.File) *atomic.Buffer {
	metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())
	// log.Printf("createCounterValuesBuffer: MetaData: %v\n", metaData)

	ptr := uintptr(cncFile.GetMemoryPtr()) +
		uintptr(Descriptor.VersionAndMetaDataLen) +
		uintptr(metaData.toDriverBufferLength) +
		uintptr(metaData.toClientsBufferLength) +
		uintptr(metaData.counterMetadataBufferLength)

	return atomic.MakeBuffer(unsafe.Pointer(ptr), metaData.counterValuesBufferLength)
}

func CreateErrorLogBuffer(cncFile *memmap.File) *atomic.Buffer {
	// metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())
	// log.Printf("createErrorLogBuffer: MetaData: %v\n", metaData)

	toDriverBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), MetadataOffsets.driverBufLen)
	toClientsBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), MetadataOffsets.clientsBufLen)
	counterMetadataBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), MetadataOffsets.metaDataBufLen)
	counterValuesBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), MetadataOffsets.valuesBufLen)
	errorLogBufferLength := readInt32FromPointer(cncFile.GetMemoryPtr(), MetadataOffsets.errorLogBufLen)

	ptr := uintptr(cncFile.GetMemoryPtr()) +
		uintptr(Descriptor.VersionAndMetaDataLen) +
		uintptr(toDriverBufferLength) +
		uintptr(toClientsBufferLength) +
		uintptr(counterMetadataBufferLength) +
		uintptr(counterValuesBufferLength)

	return atomic.MakeBuffer(unsafe.Pointer(ptr), errorLogBufferLength)
}

func CncVersion(cncFile *memmap.File) int32 {
	metaData := (*MetaDataDefn)(cncFile.GetMemoryPtr())

	return metaData.cncVersion
}

func ClientLivenessTimeout(cncFile *memmap.File) int64 {
	ptr := uintptr(cncFile.GetMemoryPtr()) + MetadataOffsets.livelinessTo

	return *(*int64)(unsafe.Pointer(ptr))
}

func readInt32FromPointer(basePtr unsafe.Pointer, offset uintptr) int32 {
	ptr := uintptr(basePtr) + offset
	return *(*int32)(unsafe.Pointer(ptr))
}
