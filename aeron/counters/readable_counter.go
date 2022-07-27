package counters

import "fmt"

type ReadableCounter struct {
	Reader         *Reader
	CounterId      int32
	metaDataOffset int32
	valueOffset    int32
}

func NewReadableCounter(reader *Reader, counterId int32) (*ReadableCounter, error) {
	if counterId < 0 || counterId >= int32(reader.maxCounterID) {
		return nil, fmt.Errorf("counterId=%d maxCounterId=%d", counterId, reader.maxCounterID)
	}
	metaDataOffset := counterId * MetadataLength
	valueOffset := counterId * CounterLength
	return &ReadableCounter{reader, counterId, metaDataOffset, valueOffset}, nil
}

func (rc *ReadableCounter) State() int32 {
	return rc.Reader.metaData.GetInt32Volatile(rc.metaDataOffset)
}

func (rc *ReadableCounter) Label() string {
	return rc.Reader.labelValue(rc.metaDataOffset)
}

func (rc *ReadableCounter) Get() int64 {
	return rc.Reader.values.GetInt64Volatile(rc.valueOffset)
}

func (rc *ReadableCounter) GetWeak() int64 {
	return rc.Reader.values.GetInt64(rc.valueOffset)
}
