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

package flyweight

import (
	"github.com/corymonroe-coinbase/aeron-go/aeron/atomic"
	"testing"
)

type StringFly struct {
	FWBase

	s StringField
}

func (m *StringFly) Wrap(buf *atomic.Buffer, offset int) Flyweight {
	pos := offset
	pos += m.s.Wrap(buf, pos, m, false)
	m.SetSize(pos - offset)
	return m
}

func TestStringFlyweight(t *testing.T) {
	str := "Hello worlds!"
	buf := atomic.MakeBuffer(make([]byte, 128), 128)

	// TODO Test aligned reads

	var fw StringFly
	fw.Wrap(buf, 0)

	fw.s.Set(str)

	t.Logf("%v", fw)

	if 4+len(str) != fw.Size() {
		t.Error("Expected length", 4+len(str), "have", fw.Size())
	}

	if str != fw.s.Get() {
		t.Error("Got", fw.s.Get(), "instead of", str)
	}
}

type PaddedFly struct {
	FWBase

	l1   Int64Field
	i1   Int32Field
	pad  Padding
	i2   Int32Field
	pad2 Padding
}

func (m *PaddedFly) Wrap(buf *atomic.Buffer, offset int) Flyweight {
	pos := offset
	pos += m.l1.Wrap(buf, pos)
	pos += m.i1.Wrap(buf, pos)
	pos += m.pad.Wrap(buf, pos, 64, 64)
	pos += m.i2.Wrap(buf, pos)
	pos += m.pad2.Wrap(buf, pos, 128, 64)

	m.SetSize(pos - offset)
	return m
}

func TestPadding_Wrap(t *testing.T) {
	buf := atomic.MakeBuffer(make([]byte, 256), 256)

	var fw PaddedFly
	fw.Wrap(buf, 0)

	if fw.Size() != 192 {
		t.Logf("fw size: %d", fw.size)
		t.Fail()
	}
}
