// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type TerminationPosition struct {
	LeadershipTermId int64
	LogPosition      int64
}

func (t *TerminationPosition) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	return nil
}

func (t *TerminationPosition) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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

func (t *TerminationPosition) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	return nil
}

func TerminationPositionInit(t *TerminationPosition) {
	return
}

func (*TerminationPosition) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*TerminationPosition) SbeTemplateId() (templateId uint16) {
	return 75
}

func (*TerminationPosition) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*TerminationPosition) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*TerminationPosition) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*TerminationPosition) LeadershipTermIdId() uint16 {
	return 1
}

func (*TerminationPosition) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (t *TerminationPosition) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.LeadershipTermIdSinceVersion()
}

func (*TerminationPosition) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*TerminationPosition) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*TerminationPosition) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TerminationPosition) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TerminationPosition) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*TerminationPosition) LogPositionId() uint16 {
	return 2
}

func (*TerminationPosition) LogPositionSinceVersion() uint16 {
	return 0
}

func (t *TerminationPosition) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.LogPositionSinceVersion()
}

func (*TerminationPosition) LogPositionDeprecated() uint16 {
	return 0
}

func (*TerminationPosition) LogPositionMetaAttribute(meta int) string {
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

func (*TerminationPosition) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TerminationPosition) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*TerminationPosition) LogPositionNullValue() int64 {
	return math.MinInt64
}
