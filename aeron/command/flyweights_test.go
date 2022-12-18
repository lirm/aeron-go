// Copyright 2022 Steven Stern
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/stretchr/testify/suite"
	"testing"
)

type FlyweightsTestSuite struct {
	suite.Suite

	byteSlice []byte
	buffer    *atomic.Buffer
	flyweight CounterMessage
}

func (s *FlyweightsTestSuite) SetupSuite() {
	s.byteSlice = make([]byte, 128)
	s.buffer = atomic.MakeBuffer(s.byteSlice)
}

func (s *FlyweightsTestSuite) TestCounterMessageKeyBuffer() {
	offset := 24
	for i := 0; i <= offset; i++ {
		s.byteSlice[i] = byte(15)
	}
	s.flyweight.Wrap(s.buffer, offset)

	newBuffer := atomic.MakeBuffer(make([]byte, 16))
	s.flyweight.CopyKeyBuffer(newBuffer, 4, 8)

	s.Assert().EqualValues(8, s.flyweight.key.Length())
}

func (s *FlyweightsTestSuite) TestCounterMessageLabelBuffer() {
	offset := 40
	for i := 0; i <= offset; i++ {
		s.byteSlice[i] = byte(255)
	}
	s.flyweight.Wrap(s.buffer, offset)

	keyBuffer := atomic.MakeBuffer(make([]byte, 16))
	keyCopySize := int32(9)
	s.flyweight.CopyKeyBuffer(keyBuffer, 6, keyCopySize)

	labelBuffer := atomic.MakeBuffer(make([]byte, 32))
	labelCopySize := int32(21)
	s.flyweight.CopyLabelBuffer(labelBuffer, 2, labelCopySize)

	s.Assert().EqualValues(21, s.flyweight.label.Length())
	// Expected is sizes of: long, long, int, length field, keyCopySize
	s.Assert().EqualValues(8+8+4+4+keyCopySize, s.flyweight.getLabelOffset())
	// Above, plus label length and size
	s.Assert().EqualValues(8+8+4+4+keyCopySize+4+labelCopySize, s.flyweight.Size())
}

func TestFlyweights(t *testing.T) {
	suite.Run(t, new(FlyweightsTestSuite))
}
