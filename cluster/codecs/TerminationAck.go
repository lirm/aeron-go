// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type TerminationAck struct {
	LeadershipTermId int64
	LogPosition      int64
	MemberId         int32
}

func (t *TerminationAck) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := t.RangeCheck(t.SbeSchemaVersion(), t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, t.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, t.MemberId); err != nil {
		return err
	}
	return nil
}

func (t *TerminationAck) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !t.LeadershipTermIdInActingVersion(actingVersion) {
		t.LeadershipTermId = t.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.LeadershipTermId); err != nil {
			return err
		}
	}
	if !t.LogPositionInActingVersion(actingVersion) {
		t.LogPosition = t.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.LogPosition); err != nil {
			return err
		}
	}
	if !t.MemberIdInActingVersion(actingVersion) {
		t.MemberId = t.MemberIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &t.MemberId); err != nil {
			return err
		}
	}
	if actingVersion > t.SbeSchemaVersion() && blockLength > t.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-t.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := t.RangeCheck(actingVersion, t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (t *TerminationAck) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if t.LeadershipTermIdInActingVersion(actingVersion) {
		if t.LeadershipTermId < t.LeadershipTermIdMinValue() || t.LeadershipTermId > t.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on t.LeadershipTermId (%v < %v > %v)", t.LeadershipTermIdMinValue(), t.LeadershipTermId, t.LeadershipTermIdMaxValue())
		}
	}
	if t.LogPositionInActingVersion(actingVersion) {
		if t.LogPosition < t.LogPositionMinValue() || t.LogPosition > t.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on t.LogPosition (%v < %v > %v)", t.LogPositionMinValue(), t.LogPosition, t.LogPositionMaxValue())
		}
	}
	if t.MemberIdInActingVersion(actingVersion) {
		if t.MemberId < t.MemberIdMinValue() || t.MemberId > t.MemberIdMaxValue() {
			return fmt.Errorf("Range check failed on t.MemberId (%v < %v > %v)", t.MemberIdMinValue(), t.MemberId, t.MemberIdMaxValue())
		}
	}
	return nil
}

func TerminationAckInit(t *TerminationAck) {
	return
}

func (*TerminationAck) SbeBlockLength() (blockLength uint16) {
	return 20
}

func (*TerminationAck) SbeTemplateId() (templateId uint16) {
	return 76
}

func (*TerminationAck) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*TerminationAck) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*TerminationAck) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*TerminationAck) LeadershipTermIdId() uint16 {
	return 1
}

func (*TerminationAck) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (t *TerminationAck) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.LeadershipTermIdSinceVersion()
}

func (*TerminationAck) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*TerminationAck) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*TerminationAck) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TerminationAck) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TerminationAck) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*TerminationAck) LogPositionId() uint16 {
	return 2
}

func (*TerminationAck) LogPositionSinceVersion() uint16 {
	return 0
}

func (t *TerminationAck) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.LogPositionSinceVersion()
}

func (*TerminationAck) LogPositionDeprecated() uint16 {
	return 0
}

func (*TerminationAck) LogPositionMetaAttribute(meta int) string {
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

func (*TerminationAck) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TerminationAck) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*TerminationAck) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*TerminationAck) MemberIdId() uint16 {
	return 3
}

func (*TerminationAck) MemberIdSinceVersion() uint16 {
	return 0
}

func (t *TerminationAck) MemberIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.MemberIdSinceVersion()
}

func (*TerminationAck) MemberIdDeprecated() uint16 {
	return 0
}

func (*TerminationAck) MemberIdMetaAttribute(meta int) string {
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

func (*TerminationAck) MemberIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*TerminationAck) MemberIdMaxValue() int32 {
	return math.MaxInt32
}

func (*TerminationAck) MemberIdNullValue() int32 {
	return math.MinInt32
}
