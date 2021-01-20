// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type StartPositionRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
}

func (s *StartPositionRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, s.RecordingId); err != nil {
		return err
	}
	return nil
}

func (s *StartPositionRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !s.RecordingIdInActingVersion(actingVersion) {
		s.RecordingId = s.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.RecordingId); err != nil {
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

func (s *StartPositionRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if s.RecordingIdInActingVersion(actingVersion) {
		if s.RecordingId < s.RecordingIdMinValue() || s.RecordingId > s.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on s.RecordingId (%v < %v > %v)", s.RecordingIdMinValue(), s.RecordingId, s.RecordingIdMaxValue())
		}
	}
	return nil
}

func StartPositionRequestInit(s *StartPositionRequest) {
	return
}

func (*StartPositionRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*StartPositionRequest) SbeTemplateId() (templateId uint16) {
	return 52
}

func (*StartPositionRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*StartPositionRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*StartPositionRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*StartPositionRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*StartPositionRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (s *StartPositionRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ControlSessionIdSinceVersion()
}

func (*StartPositionRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*StartPositionRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*StartPositionRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StartPositionRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StartPositionRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*StartPositionRequest) CorrelationIdId() uint16 {
	return 2
}

func (*StartPositionRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *StartPositionRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*StartPositionRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*StartPositionRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*StartPositionRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StartPositionRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StartPositionRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*StartPositionRequest) RecordingIdId() uint16 {
	return 3
}

func (*StartPositionRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (s *StartPositionRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.RecordingIdSinceVersion()
}

func (*StartPositionRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*StartPositionRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*StartPositionRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StartPositionRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StartPositionRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}
