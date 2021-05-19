// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type CatalogHeader struct {
	Version         int32
	Length          int32
	NextRecordingId int64
	Alignment       int32
	Reserved        int8
}

func (c *CatalogHeader) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt32(_w, c.Version); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.Length); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, c.NextRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, c.Alignment); err != nil {
		return err
	}

	for i := 0; i < 11; i++ {
		if err := _m.WriteUint8(_w, uint8(0)); err != nil {
			return err
		}
	}
	if err := _m.WriteInt8(_w, c.Reserved); err != nil {
		return err
	}
	return nil
}

func (c *CatalogHeader) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !c.VersionInActingVersion(actingVersion) {
		c.Version = c.VersionNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.Version); err != nil {
			return err
		}
	}
	if !c.LengthInActingVersion(actingVersion) {
		c.Length = c.LengthNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.Length); err != nil {
			return err
		}
	}
	if !c.NextRecordingIdInActingVersion(actingVersion) {
		c.NextRecordingId = c.NextRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &c.NextRecordingId); err != nil {
			return err
		}
	}
	if !c.AlignmentInActingVersion(actingVersion) {
		c.Alignment = c.AlignmentNullValue()
	} else {
		if err := _m.ReadInt32(_r, &c.Alignment); err != nil {
			return err
		}
	}
	io.CopyN(ioutil.Discard, _r, 11)
	if !c.ReservedInActingVersion(actingVersion) {
		c.Reserved = c.ReservedNullValue()
	} else {
		if err := _m.ReadInt8(_r, &c.Reserved); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := c.RangeCheck(actingVersion, c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (c *CatalogHeader) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if c.VersionInActingVersion(actingVersion) {
		if c.Version < c.VersionMinValue() || c.Version > c.VersionMaxValue() {
			return fmt.Errorf("Range check failed on c.Version (%v < %v > %v)", c.VersionMinValue(), c.Version, c.VersionMaxValue())
		}
	}
	if c.LengthInActingVersion(actingVersion) {
		if c.Length < c.LengthMinValue() || c.Length > c.LengthMaxValue() {
			return fmt.Errorf("Range check failed on c.Length (%v < %v > %v)", c.LengthMinValue(), c.Length, c.LengthMaxValue())
		}
	}
	if c.NextRecordingIdInActingVersion(actingVersion) {
		if c.NextRecordingId < c.NextRecordingIdMinValue() || c.NextRecordingId > c.NextRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on c.NextRecordingId (%v < %v > %v)", c.NextRecordingIdMinValue(), c.NextRecordingId, c.NextRecordingIdMaxValue())
		}
	}
	if c.AlignmentInActingVersion(actingVersion) {
		if c.Alignment < c.AlignmentMinValue() || c.Alignment > c.AlignmentMaxValue() {
			return fmt.Errorf("Range check failed on c.Alignment (%v < %v > %v)", c.AlignmentMinValue(), c.Alignment, c.AlignmentMaxValue())
		}
	}
	if c.ReservedInActingVersion(actingVersion) {
		if c.Reserved < c.ReservedMinValue() || c.Reserved > c.ReservedMaxValue() {
			return fmt.Errorf("Range check failed on c.Reserved (%v < %v > %v)", c.ReservedMinValue(), c.Reserved, c.ReservedMaxValue())
		}
	}
	return nil
}

func CatalogHeaderInit(c *CatalogHeader) {
	return
}

func (*CatalogHeader) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*CatalogHeader) SbeTemplateId() (templateId uint16) {
	return 20
}

func (*CatalogHeader) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*CatalogHeader) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*CatalogHeader) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*CatalogHeader) VersionId() uint16 {
	return 1
}

func (*CatalogHeader) VersionSinceVersion() uint16 {
	return 0
}

func (c *CatalogHeader) VersionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.VersionSinceVersion()
}

func (*CatalogHeader) VersionDeprecated() uint16 {
	return 0
}

func (*CatalogHeader) VersionMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*CatalogHeader) VersionMinValue() int32 {
	return math.MinInt32 + 1
}

func (*CatalogHeader) VersionMaxValue() int32 {
	return math.MaxInt32
}

func (*CatalogHeader) VersionNullValue() int32 {
	return math.MinInt32
}

func (*CatalogHeader) LengthId() uint16 {
	return 2
}

func (*CatalogHeader) LengthSinceVersion() uint16 {
	return 0
}

func (c *CatalogHeader) LengthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.LengthSinceVersion()
}

func (*CatalogHeader) LengthDeprecated() uint16 {
	return 0
}

func (*CatalogHeader) LengthMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*CatalogHeader) LengthMinValue() int32 {
	return math.MinInt32 + 1
}

func (*CatalogHeader) LengthMaxValue() int32 {
	return math.MaxInt32
}

func (*CatalogHeader) LengthNullValue() int32 {
	return math.MinInt32
}

func (*CatalogHeader) NextRecordingIdId() uint16 {
	return 3
}

func (*CatalogHeader) NextRecordingIdSinceVersion() uint16 {
	return 0
}

func (c *CatalogHeader) NextRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.NextRecordingIdSinceVersion()
}

func (*CatalogHeader) NextRecordingIdDeprecated() uint16 {
	return 0
}

func (*CatalogHeader) NextRecordingIdMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*CatalogHeader) NextRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*CatalogHeader) NextRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*CatalogHeader) NextRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*CatalogHeader) AlignmentId() uint16 {
	return 4
}

func (*CatalogHeader) AlignmentSinceVersion() uint16 {
	return 0
}

func (c *CatalogHeader) AlignmentInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.AlignmentSinceVersion()
}

func (*CatalogHeader) AlignmentDeprecated() uint16 {
	return 0
}

func (*CatalogHeader) AlignmentMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*CatalogHeader) AlignmentMinValue() int32 {
	return math.MinInt32 + 1
}

func (*CatalogHeader) AlignmentMaxValue() int32 {
	return math.MaxInt32
}

func (*CatalogHeader) AlignmentNullValue() int32 {
	return math.MinInt32
}

func (*CatalogHeader) ReservedId() uint16 {
	return 5
}

func (*CatalogHeader) ReservedSinceVersion() uint16 {
	return 0
}

func (c *CatalogHeader) ReservedInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ReservedSinceVersion()
}

func (*CatalogHeader) ReservedDeprecated() uint16 {
	return 0
}

func (*CatalogHeader) ReservedMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*CatalogHeader) ReservedMinValue() int8 {
	return math.MinInt8 + 1
}

func (*CatalogHeader) ReservedMaxValue() int8 {
	return math.MaxInt8
}

func (*CatalogHeader) ReservedNullValue() int8 {
	return math.MinInt8
}
