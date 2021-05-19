// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type StopRecordingSubscriptionRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	SubscriptionId   int64
}

func (s *StopRecordingSubscriptionRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, s.SubscriptionId); err != nil {
		return err
	}
	return nil
}

func (s *StopRecordingSubscriptionRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !s.SubscriptionIdInActingVersion(actingVersion) {
		s.SubscriptionId = s.SubscriptionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.SubscriptionId); err != nil {
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

func (s *StopRecordingSubscriptionRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if s.SubscriptionIdInActingVersion(actingVersion) {
		if s.SubscriptionId < s.SubscriptionIdMinValue() || s.SubscriptionId > s.SubscriptionIdMaxValue() {
			return fmt.Errorf("Range check failed on s.SubscriptionId (%v < %v > %v)", s.SubscriptionIdMinValue(), s.SubscriptionId, s.SubscriptionIdMaxValue())
		}
	}
	return nil
}

func StopRecordingSubscriptionRequestInit(s *StopRecordingSubscriptionRequest) {
	return
}

func (*StopRecordingSubscriptionRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*StopRecordingSubscriptionRequest) SbeTemplateId() (templateId uint16) {
	return 14
}

func (*StopRecordingSubscriptionRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*StopRecordingSubscriptionRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*StopRecordingSubscriptionRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*StopRecordingSubscriptionRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*StopRecordingSubscriptionRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (s *StopRecordingSubscriptionRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ControlSessionIdSinceVersion()
}

func (*StopRecordingSubscriptionRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*StopRecordingSubscriptionRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*StopRecordingSubscriptionRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopRecordingSubscriptionRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopRecordingSubscriptionRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*StopRecordingSubscriptionRequest) CorrelationIdId() uint16 {
	return 2
}

func (*StopRecordingSubscriptionRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *StopRecordingSubscriptionRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*StopRecordingSubscriptionRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*StopRecordingSubscriptionRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*StopRecordingSubscriptionRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopRecordingSubscriptionRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopRecordingSubscriptionRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*StopRecordingSubscriptionRequest) SubscriptionIdId() uint16 {
	return 3
}

func (*StopRecordingSubscriptionRequest) SubscriptionIdSinceVersion() uint16 {
	return 0
}

func (s *StopRecordingSubscriptionRequest) SubscriptionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.SubscriptionIdSinceVersion()
}

func (*StopRecordingSubscriptionRequest) SubscriptionIdDeprecated() uint16 {
	return 0
}

func (*StopRecordingSubscriptionRequest) SubscriptionIdMetaAttribute(meta int) string {
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

func (*StopRecordingSubscriptionRequest) SubscriptionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopRecordingSubscriptionRequest) SubscriptionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopRecordingSubscriptionRequest) SubscriptionIdNullValue() int64 {
	return math.MinInt64
}
