package flyweight

import "github.com/lirm/aeron-go/aeron/buffer"

type Flyweight interface {
	Wrap(*buffer.Atomic, int) Flyweight
	Size() int
	SetSize(int)
}

type FWBase struct {
	size int
}

func (m *FWBase) Size() int {
	return m.size
}

func (m *FWBase) SetSize(size int) {
	m.size = size
}
