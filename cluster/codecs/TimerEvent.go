// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type TimerEvent struct {
	LeadershipTermId int64
	CorrelationId    int64
	Timestamp        int64
}

func (t *TimerEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := t.RangeCheck(t.SbeSchemaVersion(), t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, t.LeadershipTermId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.Timestamp); err != nil {
		return err
	}
	return nil
}

func (t *TimerEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !t.LeadershipTermIdInActingVersion(actingVersion) {
		t.LeadershipTermId = t.LeadershipTermIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.LeadershipTermId); err != nil {
			return err
		}
	}
	if !t.CorrelationIdInActingVersion(actingVersion) {
		t.CorrelationId = t.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.CorrelationId); err != nil {
			return err
		}
	}
	if !t.TimestampInActingVersion(actingVersion) {
		t.Timestamp = t.TimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.Timestamp); err != nil {
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

func (t *TimerEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if t.LeadershipTermIdInActingVersion(actingVersion) {
		if t.LeadershipTermId < t.LeadershipTermIdMinValue() || t.LeadershipTermId > t.LeadershipTermIdMaxValue() {
			return fmt.Errorf("Range check failed on t.LeadershipTermId (%v < %v > %v)", t.LeadershipTermIdMinValue(), t.LeadershipTermId, t.LeadershipTermIdMaxValue())
		}
	}
	if t.CorrelationIdInActingVersion(actingVersion) {
		if t.CorrelationId < t.CorrelationIdMinValue() || t.CorrelationId > t.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on t.CorrelationId (%v < %v > %v)", t.CorrelationIdMinValue(), t.CorrelationId, t.CorrelationIdMaxValue())
		}
	}
	if t.TimestampInActingVersion(actingVersion) {
		if t.Timestamp < t.TimestampMinValue() || t.Timestamp > t.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on t.Timestamp (%v < %v > %v)", t.TimestampMinValue(), t.Timestamp, t.TimestampMaxValue())
		}
	}
	return nil
}

func TimerEventInit(t *TimerEvent) {
	return
}

func (*TimerEvent) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*TimerEvent) SbeTemplateId() (templateId uint16) {
	return 20
}

func (*TimerEvent) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*TimerEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*TimerEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*TimerEvent) LeadershipTermIdId() uint16 {
	return 1
}

func (*TimerEvent) LeadershipTermIdSinceVersion() uint16 {
	return 0
}

func (t *TimerEvent) LeadershipTermIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.LeadershipTermIdSinceVersion()
}

func (*TimerEvent) LeadershipTermIdDeprecated() uint16 {
	return 0
}

func (*TimerEvent) LeadershipTermIdMetaAttribute(meta int) string {
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

func (*TimerEvent) LeadershipTermIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TimerEvent) LeadershipTermIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TimerEvent) LeadershipTermIdNullValue() int64 {
	return math.MinInt64
}

func (*TimerEvent) CorrelationIdId() uint16 {
	return 2
}

func (*TimerEvent) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (t *TimerEvent) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.CorrelationIdSinceVersion()
}

func (*TimerEvent) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*TimerEvent) CorrelationIdMetaAttribute(meta int) string {
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

func (*TimerEvent) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TimerEvent) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TimerEvent) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*TimerEvent) TimestampId() uint16 {
	return 3
}

func (*TimerEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (t *TimerEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.TimestampSinceVersion()
}

func (*TimerEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*TimerEvent) TimestampMetaAttribute(meta int) string {
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

func (*TimerEvent) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TimerEvent) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*TimerEvent) TimestampNullValue() int64 {
	return math.MinInt64
}
