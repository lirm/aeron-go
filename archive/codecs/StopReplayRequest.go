// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type StopReplayRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	ReplaySessionId  int64
}

func (s *StopReplayRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.ReplaySessionId); err != nil {
		return err
	}
	return nil
}

func (s *StopReplayRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.ControlSessionIdInActingVersion(actingVersion) {
		s.ControlSessionId = s.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.ControlSessionId); err != nil {
			return err
		}
	}
	if !s.CorrelationIdInActingVersion(actingVersion) {
		s.CorrelationId = s.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.CorrelationId); err != nil {
			return err
		}
	}
	if !s.ReplaySessionIdInActingVersion(actingVersion) {
		s.ReplaySessionId = s.ReplaySessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.ReplaySessionId); err != nil {
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

func (s *StopReplayRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.ControlSessionIdInActingVersion(actingVersion) {
		if s.ControlSessionId < s.ControlSessionIdMinValue() || s.ControlSessionId > s.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ControlSessionId (%v < %v > %v)", s.ControlSessionIdMinValue(), s.ControlSessionId, s.ControlSessionIdMaxValue())
		}
	}
	if s.CorrelationIdInActingVersion(actingVersion) {
		if s.CorrelationId < s.CorrelationIdMinValue() || s.CorrelationId > s.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on s.CorrelationId (%v < %v > %v)", s.CorrelationIdMinValue(), s.CorrelationId, s.CorrelationIdMaxValue())
		}
	}
	if s.ReplaySessionIdInActingVersion(actingVersion) {
		if s.ReplaySessionId < s.ReplaySessionIdMinValue() || s.ReplaySessionId > s.ReplaySessionIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ReplaySessionId (%v < %v > %v)", s.ReplaySessionIdMinValue(), s.ReplaySessionId, s.ReplaySessionIdMaxValue())
		}
	}
	return nil
}

func StopReplayRequestInit(s *StopReplayRequest) {
	return
}

func (*StopReplayRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*StopReplayRequest) SbeTemplateId() (templateId uint16) {
	return 7
}

func (*StopReplayRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*StopReplayRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*StopReplayRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*StopReplayRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*StopReplayRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (s *StopReplayRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ControlSessionIdSinceVersion()
}

func (*StopReplayRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*StopReplayRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*StopReplayRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopReplayRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopReplayRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*StopReplayRequest) CorrelationIdId() uint16 {
	return 2
}

func (*StopReplayRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *StopReplayRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*StopReplayRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*StopReplayRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*StopReplayRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopReplayRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopReplayRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*StopReplayRequest) ReplaySessionIdId() uint16 {
	return 3
}

func (*StopReplayRequest) ReplaySessionIdSinceVersion() uint16 {
	return 0
}

func (s *StopReplayRequest) ReplaySessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ReplaySessionIdSinceVersion()
}

func (*StopReplayRequest) ReplaySessionIdDeprecated() uint16 {
	return 0
}

func (*StopReplayRequest) ReplaySessionIdMetaAttribute(meta int) string {
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

func (*StopReplayRequest) ReplaySessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopReplayRequest) ReplaySessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopReplayRequest) ReplaySessionIdNullValue() int64 {
	return math.MinInt64
}
