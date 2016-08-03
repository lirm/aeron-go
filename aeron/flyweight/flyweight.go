package flyweight

import "github.com/lirm/aeron-go/aeron/atomic"

type Flyweight interface {
	Wrap(*atomic.Buffer, int) Flyweight
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
