package flyweight

import (
	"github.com/lirm/aeron-go/aeron/buffers"
	"testing"
)

type StringFly struct {
	s      StringField
	length int
}

func (m *StringFly) Length() int {
	return m.length
}

func (m *StringFly) Wrap(buf *buffers.Atomic) *StringFly {
	offset := 0
	offset += m.s.wrap(buf, offset, &m.length)

	m.length = offset
	return m
}

func TestStringFlyweight(t *testing.T) {
	str := "Hello worlds!"
	buf := buffers.MakeAtomic(make([]byte, 128), 128)

	var fw StringFly
	fw.Wrap(buf)

	fw.s.Set(str)

	t.Logf("%v", fw)

	if 4+len(str) != fw.Length() {
		t.Error("Expected length", 4+len(str), "have", fw.Length())
	}

	if str != fw.s.Get() {
		t.Error("Got", fw.s.Get(), "instead of", str)
	}
}
