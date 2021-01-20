// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RecordingDescriptorHeader struct {
	Length   int32
	State    RecordingStateEnum
	Checksum int32
	Reserved int8
}

func (r *RecordingDescriptorHeader) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt32(_w, r.Length); err != nil {
		return err
	}
	if err := r.State.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, r.Checksum); err != nil {
		return err
	}

	for i := 0; i < 19; i++ {
		if err := _m.WriteUint8(_w, uint8(0)); err != nil {
			return err
		}
	}
	if err := _m.WriteInt8(_w, r.Reserved); err != nil {
		return err
	}
	return nil
}

func (r *RecordingDescriptorHeader) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.LengthInActingVersion(actingVersion) {
		r.Length = r.LengthNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.Length); err != nil {
			return err
		}
	}
	if r.StateInActingVersion(actingVersion) {
		if err := r.State.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !r.ChecksumInActingVersion(actingVersion) {
		r.Checksum = r.ChecksumNullValue()
	} else {
		if err := _m.ReadInt32(_r, &r.Checksum); err != nil {
			return err
		}
	}
	io.CopyN(ioutil.Discard, _r, 19)
	if !r.ReservedInActingVersion(actingVersion) {
		r.Reserved = r.ReservedNullValue()
	} else {
		if err := _m.ReadInt8(_r, &r.Reserved); err != nil {
			return err
		}
	}
	if actingVersion > r.SbeSchemaVersion() && blockLength > r.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-r.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := r.RangeCheck(actingVersion, r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (r *RecordingDescriptorHeader) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.LengthInActingVersion(actingVersion) {
		if r.Length < r.LengthMinValue() || r.Length > r.LengthMaxValue() {
			return fmt.Errorf("Range check failed on r.Length (%v < %v > %v)", r.LengthMinValue(), r.Length, r.LengthMaxValue())
		}
	}
	if err := r.State.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if r.ChecksumInActingVersion(actingVersion) {
		if r.Checksum < r.ChecksumMinValue() || r.Checksum > r.ChecksumMaxValue() {
			return fmt.Errorf("Range check failed on r.Checksum (%v < %v > %v)", r.ChecksumMinValue(), r.Checksum, r.ChecksumMaxValue())
		}
	}
	if r.ReservedInActingVersion(actingVersion) {
		if r.Reserved < r.ReservedMinValue() || r.Reserved > r.ReservedMaxValue() {
			return fmt.Errorf("Range check failed on r.Reserved (%v < %v > %v)", r.ReservedMinValue(), r.Reserved, r.ReservedMaxValue())
		}
	}
	return nil
}

func RecordingDescriptorHeaderInit(r *RecordingDescriptorHeader) {
	return
}

func (*RecordingDescriptorHeader) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*RecordingDescriptorHeader) SbeTemplateId() (templateId uint16) {
	return 21
}

func (*RecordingDescriptorHeader) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*RecordingDescriptorHeader) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*RecordingDescriptorHeader) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RecordingDescriptorHeader) LengthId() uint16 {
	return 1
}

func (*RecordingDescriptorHeader) LengthSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptorHeader) LengthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.LengthSinceVersion()
}

func (*RecordingDescriptorHeader) LengthDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptorHeader) LengthMetaAttribute(meta int) string {
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

func (*RecordingDescriptorHeader) LengthMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingDescriptorHeader) LengthMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingDescriptorHeader) LengthNullValue() int32 {
	return math.MinInt32
}

func (*RecordingDescriptorHeader) StateId() uint16 {
	return 2
}

func (*RecordingDescriptorHeader) StateSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptorHeader) StateInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.StateSinceVersion()
}

func (*RecordingDescriptorHeader) StateDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptorHeader) StateMetaAttribute(meta int) string {
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

func (*RecordingDescriptorHeader) ChecksumId() uint16 {
	return 4
}

func (*RecordingDescriptorHeader) ChecksumSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptorHeader) ChecksumInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ChecksumSinceVersion()
}

func (*RecordingDescriptorHeader) ChecksumDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptorHeader) ChecksumMetaAttribute(meta int) string {
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

func (*RecordingDescriptorHeader) ChecksumMinValue() int32 {
	return math.MinInt32 + 1
}

func (*RecordingDescriptorHeader) ChecksumMaxValue() int32 {
	return math.MaxInt32
}

func (*RecordingDescriptorHeader) ChecksumNullValue() int32 {
	return math.MinInt32
}

func (*RecordingDescriptorHeader) ReservedId() uint16 {
	return 3
}

func (*RecordingDescriptorHeader) ReservedSinceVersion() uint16 {
	return 0
}

func (r *RecordingDescriptorHeader) ReservedInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ReservedSinceVersion()
}

func (*RecordingDescriptorHeader) ReservedDeprecated() uint16 {
	return 0
}

func (*RecordingDescriptorHeader) ReservedMetaAttribute(meta int) string {
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

func (*RecordingDescriptorHeader) ReservedMinValue() int8 {
	return math.MinInt8 + 1
}

func (*RecordingDescriptorHeader) ReservedMaxValue() int8 {
	return math.MaxInt8
}

func (*RecordingDescriptorHeader) ReservedNullValue() int8 {
	return math.MinInt8
}
