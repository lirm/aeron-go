// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type AddPassiveMember struct {
	CorrelationId   int64
	MemberEndpoints []uint8
}

func (a *AddPassiveMember) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := a.RangeCheck(a.SbeSchemaVersion(), a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, a.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(a.MemberEndpoints))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, a.MemberEndpoints); err != nil {
		return err
	}
	return nil
}

func (a *AddPassiveMember) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !a.CorrelationIdInActingVersion(actingVersion) {
		a.CorrelationId = a.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.CorrelationId); err != nil {
			return err
		}
	}
	if actingVersion > a.SbeSchemaVersion() && blockLength > a.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-a.SbeBlockLength()))
	}

	if a.MemberEndpointsInActingVersion(actingVersion) {
		var MemberEndpointsLength uint32
		if err := _m.ReadUint32(_r, &MemberEndpointsLength); err != nil {
			return err
		}
		if cap(a.MemberEndpoints) < int(MemberEndpointsLength) {
			a.MemberEndpoints = make([]uint8, MemberEndpointsLength)
		}
		a.MemberEndpoints = a.MemberEndpoints[:MemberEndpointsLength]
		if err := _m.ReadBytes(_r, a.MemberEndpoints); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := a.RangeCheck(actingVersion, a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (a *AddPassiveMember) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if a.CorrelationIdInActingVersion(actingVersion) {
		if a.CorrelationId < a.CorrelationIdMinValue() || a.CorrelationId > a.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on a.CorrelationId (%v < %v > %v)", a.CorrelationIdMinValue(), a.CorrelationId, a.CorrelationIdMaxValue())
		}
	}
	for idx, ch := range a.MemberEndpoints {
		if ch > 127 {
			return fmt.Errorf("a.MemberEndpoints[%d]=%d failed ASCII validation", idx, ch)
		}
	}
	return nil
}

func AddPassiveMemberInit(a *AddPassiveMember) {
	return
}

func (*AddPassiveMember) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*AddPassiveMember) SbeTemplateId() (templateId uint16) {
	return 70
}

func (*AddPassiveMember) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*AddPassiveMember) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*AddPassiveMember) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*AddPassiveMember) CorrelationIdId() uint16 {
	return 1
}

func (*AddPassiveMember) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (a *AddPassiveMember) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.CorrelationIdSinceVersion()
}

func (*AddPassiveMember) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*AddPassiveMember) CorrelationIdMetaAttribute(meta int) string {
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

func (*AddPassiveMember) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AddPassiveMember) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AddPassiveMember) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*AddPassiveMember) MemberEndpointsMetaAttribute(meta int) string {
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

func (*AddPassiveMember) MemberEndpointsSinceVersion() uint16 {
	return 0
}

func (a *AddPassiveMember) MemberEndpointsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.MemberEndpointsSinceVersion()
}

func (*AddPassiveMember) MemberEndpointsDeprecated() uint16 {
	return 0
}

func (AddPassiveMember) MemberEndpointsCharacterEncoding() string {
	return "US-ASCII"
}

func (AddPassiveMember) MemberEndpointsHeaderLength() uint64 {
	return 4
}
