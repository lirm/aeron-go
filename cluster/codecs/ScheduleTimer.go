// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ScheduleTimer struct {
	CorrelationId int64
	Deadline      int64
}

func (s *ScheduleTimer) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.Deadline); err != nil {
		return err
	}
	return nil
}

func (s *ScheduleTimer) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.CorrelationIdInActingVersion(actingVersion) {
		s.CorrelationId = s.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.CorrelationId); err != nil {
			return err
		}
	}
	if !s.DeadlineInActingVersion(actingVersion) {
		s.Deadline = s.DeadlineNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.Deadline); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *ScheduleTimer) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.CorrelationIdInActingVersion(actingVersion) {
		if s.CorrelationId < s.CorrelationIdMinValue() || s.CorrelationId > s.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on s.CorrelationId (%v < %v > %v)", s.CorrelationIdMinValue(), s.CorrelationId, s.CorrelationIdMaxValue())
		}
	}
	if s.DeadlineInActingVersion(actingVersion) {
		if s.Deadline < s.DeadlineMinValue() || s.Deadline > s.DeadlineMaxValue() {
			return fmt.Errorf("Range check failed on s.Deadline (%v < %v > %v)", s.DeadlineMinValue(), s.Deadline, s.DeadlineMaxValue())
		}
	}
	return nil
}

func ScheduleTimerInit(s *ScheduleTimer) {
	return
}

func (*ScheduleTimer) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*ScheduleTimer) SbeTemplateId() (templateId uint16) {
	return 31
}

func (*ScheduleTimer) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ScheduleTimer) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ScheduleTimer) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ScheduleTimer) CorrelationIdId() uint16 {
	return 1
}

func (*ScheduleTimer) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *ScheduleTimer) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*ScheduleTimer) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ScheduleTimer) CorrelationIdMetaAttribute(meta int) string {
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

func (*ScheduleTimer) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ScheduleTimer) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ScheduleTimer) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ScheduleTimer) DeadlineId() uint16 {
	return 2
}

func (*ScheduleTimer) DeadlineSinceVersion() uint16 {
	return 0
}

func (s *ScheduleTimer) DeadlineInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.DeadlineSinceVersion()
}

func (*ScheduleTimer) DeadlineDeprecated() uint16 {
	return 0
}

func (*ScheduleTimer) DeadlineMetaAttribute(meta int) string {
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

func (*ScheduleTimer) DeadlineMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ScheduleTimer) DeadlineMaxValue() int64 {
	return math.MaxInt64
}

func (*ScheduleTimer) DeadlineNullValue() int64 {
	return math.MinInt64
}
