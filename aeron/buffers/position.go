package buffers

type Position struct {
	Buffer *Atomic
	Id     int32
	Offset int32
}

func (pos *Position) Get() int64 {
	return pos.Buffer.GetInt64(pos.Offset)
}

func (pos *Position) Set(value int64) {
	pos.Buffer.PutInt64(pos.Offset, value)
}
