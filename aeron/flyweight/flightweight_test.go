/*
Copyright 2016 Stanislav Liberman
Copyright 2022 Steven Stern

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
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	str1 = "Hello, World!"
	str2 = "I hope the alignment is right"
	str3 = "Goodbye!"
)

type FlyweightAlignmentTestSuite struct {
	suite.Suite

	shouldAlign bool
	stringBuf   *atomic.Buffer
}

type StringFly struct {
	FWBase

	s  StringField
	s2 StringField
	s3 StringField

	shouldWrap bool
}

func (m *StringFly) Wrap(buf *atomic.Buffer, offset int) Flyweight {
	pos := offset
	pos += m.s.Wrap(buf, pos, m, m.shouldWrap)
	pos += m.s2.Wrap(buf, pos, m, m.shouldWrap)
	pos += m.s3.Wrap(buf, pos, m, m.shouldWrap)
	m.SetSize(pos - offset)
	return m
}

// Only works for small strings
func (f *FlyweightAlignmentTestSuite) BuildStringFly(str string) []byte {
	bytes := []byte{byte(len(str)), 0, 0, 0}
	bytes = append(bytes, []byte(str)...)
	if f.shouldAlign {
		for len(bytes)%4 != 0 {
			bytes = append(bytes, 0)
		}
	}
	return bytes
}

func (f *FlyweightAlignmentTestSuite) SetupTest() {
	// This checks for a regression bug in StringField.Wrap.  The original bug was that if a StringField starts with an
	// unaligned offset, Wrap() returned a bad offset.  To check for regressions, we need a somewhat awkward
	// construction:
	//
	// str1: ends with an unaligned offset, setting up the bug
	// str2: triggers the bug, returning a bad offset
	// str3: uses the bad offset and corrupts the buffer (does not have to be a string)
	testBuffer := f.BuildStringFly(str1)
	testBuffer = append(testBuffer, f.BuildStringFly(str2)...)
	testBuffer = append(testBuffer, f.BuildStringFly(str3)...)
	f.stringBuf = atomic.MakeBuffer(testBuffer)

	sf := StringFly{shouldWrap: f.shouldAlign}
	sf.Wrap(f.stringBuf, 0)
}

// Return the possibly-aligned length of `str1`, plus 4 to account for the length field.
func (f *FlyweightAlignmentTestSuite) len(str string) int {
	if f.shouldAlign {
		return 4 + int(util.AlignInt32(int32(len(str)), 4))
	} else {
		return 4 + len(str)
	}
}

func (f *FlyweightAlignmentTestSuite) TestStringFlyweight() {
	sf := StringFly{shouldWrap: f.shouldAlign}
	sf.Wrap(f.stringBuf, 0)
	f.Assert().Equal(f.len(str1)+f.len(str2)+f.len(str3), sf.Size())
	f.Assert().Equal(sf.s.Get(), str1)
	f.Assert().Equal(sf.s2.Get(), str2)
	f.Assert().Equal(sf.s3.Get(), str3)
}

type FlyweightTestSuite struct {
	suite.Suite
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

func (f *FlyweightTestSuite) TestPaddingWrap() {
	buf := atomic.MakeBuffer(make([]byte, 256), 256)

	var fw PaddedFly
	fw.Wrap(buf, 0)

	f.Assert().Equal(fw.Size(), 192)
}

// Another regression test: Copy of PaddedFly, plus a tiny extra padding that should not be negative.
type MorePaddedFly struct {
	FWBase

	l1   Int64Field
	i1   Int32Field
	pad  Padding
	i2   Int32Field
	pad2 Padding
	i3   Int32Field
	pad3 Padding
}

func (m *MorePaddedFly) Wrap(buf *atomic.Buffer, offset int) Flyweight {
	pos := offset
	pos += m.l1.Wrap(buf, pos)
	pos += m.i1.Wrap(buf, pos)
	pos += m.pad.Wrap(buf, pos, 64, 64)
	pos += m.i2.Wrap(buf, pos)
	pos += m.pad2.Wrap(buf, pos, 128, 64)
	pos += m.i3.Wrap(buf, pos)
	pos += m.pad3.Wrap(buf, pos, 1, 64)

	m.SetSize(pos - offset)
	return m
}

func (f *FlyweightTestSuite) TestMinimumPaddingWrap() {
	buf := atomic.MakeBuffer(make([]byte, 512), 512)

	var fw MorePaddedFly
	fw.Wrap(buf, 0)

	f.Assert().Equal(fw.Size(), 256)
}

func (f *FlyweightTestSuite) TestLengthAndRawDataFieldCopyString() {
	testString := "Talos loves Aeron"
	buf := atomic.MakeBuffer(make([]byte, 128), 128)
	field := LengthAndRawDataField{}
	field.Wrap(buf, 0)
	field.CopyString(testString)
	f.Require().Equal(field.GetAsASCII(), testString)
	f.Require().Equal(field.GetAsBuffer().GetBytesArray(0, int32(len(testString))), []byte(testString))
	f.Require().EqualValues(field.Length(), len(testString))
}

func (f *FlyweightTestSuite) TestLengthAndRawDataFieldCopyBytes() {
	testString := "Talos loves Aeron"
	buf := atomic.MakeBuffer(make([]byte, 128), 128)
	field := LengthAndRawDataField{}
	field.Wrap(buf, 0)
	field.CopyBuffer(atomic.MakeBuffer([]byte(testString)), 0, int32(len(testString)))
	f.Require().Equal(field.GetAsASCII(), testString)
	f.Require().Equal(field.GetAsBuffer().GetBytesArray(0, int32(len(testString))), []byte(testString))
	f.Require().EqualValues(field.Length(), len(testString))
}

func (f *FlyweightTestSuite) TestLengthAndRawDataFieldGetFields() {
	testString := "Talos loves Aeron"
	bufferBytes := append([]byte{byte(len(testString)), 0, 0, 0},
		[]byte(testString)...)
	buf := atomic.MakeBuffer(bufferBytes)
	field := LengthAndRawDataField{}
	f.Require().Equal(field.Wrap(buf, 0), len(testString)+4)
	f.Require().Equal(field.GetAsASCII(), testString)
	f.Require().Equal(field.GetAsBuffer().GetBytesArray(0, int32(len(testString))), []byte(testString))
}

func TestFlyweight(t *testing.T) {
	suite.Run(t, &FlyweightAlignmentTestSuite{shouldAlign: true})
	suite.Run(t, &FlyweightAlignmentTestSuite{shouldAlign: false})
	suite.Run(t, new(FlyweightTestSuite))
}
