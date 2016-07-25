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
	"github.com/lirm/aeron-go/aeron/buffers"
	"testing"
)

type StringFly struct {
	FWBase

	s StringField
}

func (m *StringFly) Wrap(buf *buffers.Atomic, offset int) Flyweight {
	pos := offset
	pos += m.s.Wrap(buf, pos, m)
	m.SetSize(pos - offset)
	return m
}

func TestStringFlyweight(t *testing.T) {
	str := "Hello worlds!"
	buf := buffers.MakeAtomic(make([]byte, 128), 128)

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
