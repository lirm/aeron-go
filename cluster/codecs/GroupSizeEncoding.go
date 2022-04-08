// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"math"
)

type GroupSizeEncoding struct {
	BlockLength uint16
	NumInGroup  uint16
}

func (g *GroupSizeEncoding) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteUint16(_w, g.BlockLength); err != nil {
		return err
	}
	if err := _m.WriteUint16(_w, g.NumInGroup); err != nil {
		return err
	}
	return nil
}

func (g *GroupSizeEncoding) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if !g.BlockLengthInActingVersion(actingVersion) {
		g.BlockLength = g.BlockLengthNullValue()
	} else {
		if err := _m.ReadUint16(_r, &g.BlockLength); err != nil {
			return err
		}
	}
	if !g.NumInGroupInActingVersion(actingVersion) {
		g.NumInGroup = g.NumInGroupNullValue()
	} else {
		if err := _m.ReadUint16(_r, &g.NumInGroup); err != nil {
			return err
		}
	}
	return nil
}

func (g *GroupSizeEncoding) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if g.BlockLengthInActingVersion(actingVersion) {
		if g.BlockLength < g.BlockLengthMinValue() || g.BlockLength > g.BlockLengthMaxValue() {
			return fmt.Errorf("Range check failed on g.BlockLength (%v < %v > %v)", g.BlockLengthMinValue(), g.BlockLength, g.BlockLengthMaxValue())
		}
	}
	if g.NumInGroupInActingVersion(actingVersion) {
		if g.NumInGroup < g.NumInGroupMinValue() || g.NumInGroup > g.NumInGroupMaxValue() {
			return fmt.Errorf("Range check failed on g.NumInGroup (%v < %v > %v)", g.NumInGroupMinValue(), g.NumInGroup, g.NumInGroupMaxValue())
		}
	}
	return nil
}

func GroupSizeEncodingInit(g *GroupSizeEncoding) {
	return
}

func (*GroupSizeEncoding) EncodedLength() int64 {
	return 4
}

func (*GroupSizeEncoding) BlockLengthMinValue() uint16 {
	return 0
}

func (*GroupSizeEncoding) BlockLengthMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*GroupSizeEncoding) BlockLengthNullValue() uint16 {
	return math.MaxUint16
}

func (*GroupSizeEncoding) BlockLengthSinceVersion() uint16 {
	return 0
}

func (g *GroupSizeEncoding) BlockLengthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= g.BlockLengthSinceVersion()
}

func (*GroupSizeEncoding) BlockLengthDeprecated() uint16 {
	return 0
}

func (*GroupSizeEncoding) NumInGroupMinValue() uint16 {
	return 0
}

func (*GroupSizeEncoding) NumInGroupMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*GroupSizeEncoding) NumInGroupNullValue() uint16 {
	return math.MaxUint16
}

func (*GroupSizeEncoding) NumInGroupSinceVersion() uint16 {
	return 0
}

func (g *GroupSizeEncoding) NumInGroupInActingVersion(actingVersion uint16) bool {
	return actingVersion >= g.NumInGroupSinceVersion()
}

func (*GroupSizeEncoding) NumInGroupDeprecated() uint16 {
	return 0
}
