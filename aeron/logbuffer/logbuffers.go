package logbuffer

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"unsafe"
)

type LogBuffers struct {
	memoryMappedFiles []*MemoryMappedFile
	buffers           [PARTITION_COUNT + 1]buffers.Atomic
}

func Wrap(fileName string) *LogBuffers {
	buffers := new(LogBuffers)

	logLength := GetFileSize(fileName)
	termLength := computeTermLength(logLength)

	checkTermLength(termLength)

	if logLength < Descriptor.MAX_SINGLE_MAPPING_SIZE {
		mmap, err := MapExistingMemoryMappedFile(fileName, 0, 0)
		if err != nil {
			panic(err)
		}

		buffers.memoryMappedFiles = [](*MemoryMappedFile){mmap}
		basePtr := uintptr(buffers.memoryMappedFiles[0].GetMemoryPtr())
		for i := 0; i < PARTITION_COUNT; i++ {
			ptr := unsafe.Pointer(basePtr + uintptr(int64(i)*termLength))
			buffers.buffers[i].Wrap(ptr, int32(termLength))
		}

		ptr := unsafe.Pointer(basePtr + uintptr(logLength-int64(Descriptor.LOG_META_DATA_LENGTH)))
		buffers.buffers[PARTITION_COUNT].Wrap(ptr, Descriptor.LOG_META_DATA_LENGTH)
	} else {
		buffers.memoryMappedFiles = make([](*MemoryMappedFile), PARTITION_COUNT+1)
		metaDataSectionOffset := termLength * int64(PARTITION_COUNT)
		metaDataSectionLength := int(logLength - metaDataSectionOffset)

		mmap, err := MapExistingMemoryMappedFile(fileName, metaDataSectionOffset, metaDataSectionLength)
		if err != nil {
			panic("Failed to map the log buffer")
		}
		buffers.memoryMappedFiles = append(buffers.memoryMappedFiles, mmap)

		for i := 0; i < PARTITION_COUNT; i++ {
			// one map for each term
			mmap, err := MapExistingMemoryMappedFile(fileName, int64(i)*termLength, int(termLength))
			if err != nil {
				panic("Failed to map the log buffer")
			}
			buffers.memoryMappedFiles = append(buffers.memoryMappedFiles, mmap)

			basePtr := buffers.memoryMappedFiles[i+1].GetMemoryPtr()

			buffers.buffers[i].Wrap(basePtr, int32(termLength))
		}
		metaDataBasePtr := buffers.memoryMappedFiles[0].GetMemoryPtr()
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

func (logBuffers *LogBuffers) Buffer(index int) *buffers.Atomic {
	return &logBuffers.buffers[index]
}
