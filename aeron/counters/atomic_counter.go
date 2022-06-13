package counters

import "fmt"

type AtomicCounter struct {
	Reader         *Reader
	CounterId      int32
	metaDataOffset int32
	valueOffset    int32
}

func NewAtomicCounter(reader *Reader, counterId int32) (*AtomicCounter, error) {
	if counterId < 0 || counterId >= int32(reader.maxCounterID) {
		return nil, fmt.Errorf("counterId=%d maxCounterId=%d", counterId, reader.maxCounterID)
	}
	metaDataOffset := counterId * MetadataLength
	valueOffset := counterId * CounterLength
	return &AtomicCounter{reader, counterId, metaDataOffset, valueOffset}, nil
}

func (ac *AtomicCounter) State() int32 {
	return ac.Reader.metaData.GetInt32Volatile(ac.metaDataOffset)
}

func (ac *AtomicCounter) Label() string {
	return ac.Reader.labelValue(ac.metaDataOffset)
}

func (ac *AtomicCounter) Get() int64 {
	return ac.Reader.values.GetInt64Volatile(ac.valueOffset)
}

func (ac *AtomicCounter) GetWeak() int64 {
	return ac.Reader.values.GetInt64(ac.valueOffset)
}

func (ac *AtomicCounter) Set(value int64) {
	ac.Reader.values.PutInt64Ordered(ac.valueOffset, value)
}

func (ac *AtomicCounter) SetWeak(value int64) {
	ac.Reader.values.PutInt64(ac.valueOffset, value)
}
