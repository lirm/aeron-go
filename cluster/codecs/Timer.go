// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type Timer struct {
	CorrelationId int64
	Deadline      int64
}

func (t *Timer) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := t.RangeCheck(t.SbeSchemaVersion(), t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, t.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.Deadline); err != nil {
		return err
	}
	return nil
}

func (t *Timer) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !t.CorrelationIdInActingVersion(actingVersion) {
		t.CorrelationId = t.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.CorrelationId); err != nil {
			return err
		}
	}
	if !t.DeadlineInActingVersion(actingVersion) {
		t.Deadline = t.DeadlineNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.Deadline); err != nil {
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

func (t *Timer) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if t.CorrelationIdInActingVersion(actingVersion) {
		if t.CorrelationId < t.CorrelationIdMinValue() || t.CorrelationId > t.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on t.CorrelationId (%v < %v > %v)", t.CorrelationIdMinValue(), t.CorrelationId, t.CorrelationIdMaxValue())
		}
	}
	if t.DeadlineInActingVersion(actingVersion) {
		if t.Deadline < t.DeadlineMinValue() || t.Deadline > t.DeadlineMaxValue() {
			return fmt.Errorf("Range check failed on t.Deadline (%v < %v > %v)", t.DeadlineMinValue(), t.Deadline, t.DeadlineMaxValue())
		}
	}
	return nil
}

func TimerInit(t *Timer) {
	return
}

func (*Timer) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*Timer) SbeTemplateId() (templateId uint16) {
	return 104
}

func (*Timer) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*Timer) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*Timer) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*Timer) CorrelationIdId() uint16 {
	return 1
}

func (*Timer) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (t *Timer) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.CorrelationIdSinceVersion()
}

func (*Timer) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*Timer) CorrelationIdMetaAttribute(meta int) string {
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

func (*Timer) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*Timer) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*Timer) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*Timer) DeadlineId() uint16 {
	return 2
}

func (*Timer) DeadlineSinceVersion() uint16 {
	return 0
}

func (t *Timer) DeadlineInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.DeadlineSinceVersion()
}

func (*Timer) DeadlineDeprecated() uint16 {
	return 0
}

func (*Timer) DeadlineMetaAttribute(meta int) string {
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

func (*Timer) DeadlineMinValue() int64 {
	return math.MinInt64 + 1
}

func (*Timer) DeadlineMaxValue() int64 {
	return math.MaxInt64
}

func (*Timer) DeadlineNullValue() int64 {
	return math.MinInt64
}
