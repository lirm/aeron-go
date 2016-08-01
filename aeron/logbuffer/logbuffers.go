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

import (
	"github.com/lirm/aeron-go/aeron/buffer"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/op/go-logging"
	"unsafe"
)

var logger = logging.MustGetLogger("logbuffers")

type LogBuffers struct {
	mmapFiles []*memmap.File
	buffers   [PARTITION_COUNT + 1]buffer.Atomic
}

func Wrap(fileName string) *LogBuffers {
	buffers := new(LogBuffers)

	logLength := memmap.GetFileSize(fileName)
	termLength := computeTermLength(logLength)

	checkTermLength(termLength)

	if logLength < Descriptor.MAX_SINGLE_MAPPING_SIZE {
		mmap, err := memmap.MapExisting(fileName, 0, 0)
		if err != nil {
			panic(err)
		}

		buffers.mmapFiles = [](*memmap.File){mmap}
		basePtr := uintptr(buffers.mmapFiles[0].GetMemoryPtr())
		for i := 0; i < PARTITION_COUNT; i++ {
			ptr := unsafe.Pointer(basePtr + uintptr(int64(i)*termLength))
			buffers.buffers[i].Wrap(ptr, int32(termLength))
		}

		ptr := unsafe.Pointer(basePtr + uintptr(logLength-int64(Descriptor.LOG_META_DATA_LENGTH)))
		buffers.buffers[PARTITION_COUNT].Wrap(ptr, Descriptor.LOG_META_DATA_LENGTH)
	} else {
		buffers.mmapFiles = make([](*memmap.File), PARTITION_COUNT+1)
		metaDataSectionOffset := termLength * int64(PARTITION_COUNT)
		metaDataSectionLength := int(logLength - metaDataSectionOffset)

		mmap, err := memmap.MapExisting(fileName, metaDataSectionOffset, metaDataSectionLength)
		if err != nil {
			panic("Failed to map the log buffer")
		}
		buffers.mmapFiles = append(buffers.mmapFiles, mmap)

		for i := 0; i < PARTITION_COUNT; i++ {
			// one map for each term
			mmap, err := memmap.MapExisting(fileName, int64(i)*termLength, int(termLength))
			if err != nil {
				panic("Failed to map the log buffer")
			}
			buffers.mmapFiles = append(buffers.mmapFiles, mmap)

			basePtr := buffers.mmapFiles[i+1].GetMemoryPtr()

			buffers.buffers[i].Wrap(basePtr, int32(termLength))
		}
		metaDataBasePtr := buffers.mmapFiles[0].GetMemoryPtr()
		buffers.buffers[PARTITION_COUNT].Wrap(metaDataBasePtr, Descriptor.LOG_META_DATA_LENGTH)
	}

	//log.Printf("WrapLogBuffers: activePartitionIndex:%d, initialTermId:%d, mtuLength:%d, timeOfLastStatusMessage:%d; LogBuffer: %v",
	//	ActivePartitionIndex(&buffers.buffers[PARTITION_COUNT]),
	//	InitialTermId(&buffers.buffers[PARTITION_COUNT]),
	//	MtuLength(&buffers.buffers[PARTITION_COUNT]),
	//	TimeOfLastStatusMessage(&buffers.buffers[PARTITION_COUNT]),
	//	buffers)

	return buffers
}

func (logBuffers *LogBuffers) Buffer(index int) *buffer.Atomic {
	return &logBuffers.buffers[index]
}

func (logBuffers *LogBuffers) Close() error {
	logger.Debug("Closing logBuffers")
	// TODO accumulate errors
	var err error
	for _, mmap := range logBuffers.mmapFiles {
		err = mmap.Close()
	}
	logBuffers.mmapFiles = nil
	return err
}
