// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ServiceAck struct {
	LogPosition int64
	Timestamp   int64
	AckId       int64
	RelevantId  int64
	ServiceId   int32
}

func (s *ServiceAck) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, s.LogPosition); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.AckId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, s.RelevantId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, s.ServiceId); err != nil {
		return err
	}
	return nil
}

func (s *ServiceAck) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !s.LogPositionInActingVersion(actingVersion) {
		s.LogPosition = s.LogPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.LogPosition); err != nil {
			return err
		}
	}
	if !s.TimestampInActingVersion(actingVersion) {
		s.Timestamp = s.TimestampNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.Timestamp); err != nil {
			return err
		}
	}
	if !s.AckIdInActingVersion(actingVersion) {
		s.AckId = s.AckIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.AckId); err != nil {
			return err
		}
	}
	if !s.RelevantIdInActingVersion(actingVersion) {
		s.RelevantId = s.RelevantIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.RelevantId); err != nil {
			return err
		}
	}
	if !s.ServiceIdInActingVersion(actingVersion) {
		s.ServiceId = s.ServiceIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.ServiceId); err != nil {
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

func (s *ServiceAck) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.LogPositionInActingVersion(actingVersion) {
		if s.LogPosition < s.LogPositionMinValue() || s.LogPosition > s.LogPositionMaxValue() {
			return fmt.Errorf("Range check failed on s.LogPosition (%v < %v > %v)", s.LogPositionMinValue(), s.LogPosition, s.LogPositionMaxValue())
		}
	}
	if s.TimestampInActingVersion(actingVersion) {
		if s.Timestamp < s.TimestampMinValue() || s.Timestamp > s.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on s.Timestamp (%v < %v > %v)", s.TimestampMinValue(), s.Timestamp, s.TimestampMaxValue())
		}
	}
	if s.AckIdInActingVersion(actingVersion) {
		if s.AckId < s.AckIdMinValue() || s.AckId > s.AckIdMaxValue() {
			return fmt.Errorf("Range check failed on s.AckId (%v < %v > %v)", s.AckIdMinValue(), s.AckId, s.AckIdMaxValue())
		}
	}
	if s.RelevantIdInActingVersion(actingVersion) {
		if s.RelevantId < s.RelevantIdMinValue() || s.RelevantId > s.RelevantIdMaxValue() {
			return fmt.Errorf("Range check failed on s.RelevantId (%v < %v > %v)", s.RelevantIdMinValue(), s.RelevantId, s.RelevantIdMaxValue())
		}
	}
	if s.ServiceIdInActingVersion(actingVersion) {
		if s.ServiceId < s.ServiceIdMinValue() || s.ServiceId > s.ServiceIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ServiceId (%v < %v > %v)", s.ServiceIdMinValue(), s.ServiceId, s.ServiceIdMaxValue())
		}
	}
	return nil
}

func ServiceAckInit(s *ServiceAck) {
	return
}

func (*ServiceAck) SbeBlockLength() (blockLength uint16) {
	return 36
}

func (*ServiceAck) SbeTemplateId() (templateId uint16) {
	return 33
}

func (*ServiceAck) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*ServiceAck) SbeSchemaVersion() (schemaVersion uint16) {
	return 8
}

func (*ServiceAck) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ServiceAck) LogPositionId() uint16 {
	return 1
}

func (*ServiceAck) LogPositionSinceVersion() uint16 {
	return 0
}

func (s *ServiceAck) LogPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.LogPositionSinceVersion()
}

func (*ServiceAck) LogPositionDeprecated() uint16 {
	return 0
}

func (*ServiceAck) LogPositionMetaAttribute(meta int) string {
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

func (*ServiceAck) LogPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ServiceAck) LogPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*ServiceAck) LogPositionNullValue() int64 {
	return math.MinInt64
}

func (*ServiceAck) TimestampId() uint16 {
	return 2
}

func (*ServiceAck) TimestampSinceVersion() uint16 {
	return 0
}

func (s *ServiceAck) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TimestampSinceVersion()
}

func (*ServiceAck) TimestampDeprecated() uint16 {
	return 0
}

func (*ServiceAck) TimestampMetaAttribute(meta int) string {
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

func (*ServiceAck) TimestampMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ServiceAck) TimestampMaxValue() int64 {
	return math.MaxInt64
}

func (*ServiceAck) TimestampNullValue() int64 {
	return math.MinInt64
}

func (*ServiceAck) AckIdId() uint16 {
	return 3
}

func (*ServiceAck) AckIdSinceVersion() uint16 {
	return 0
}

func (s *ServiceAck) AckIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.AckIdSinceVersion()
}

func (*ServiceAck) AckIdDeprecated() uint16 {
	return 0
}

func (*ServiceAck) AckIdMetaAttribute(meta int) string {
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

func (*ServiceAck) AckIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ServiceAck) AckIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ServiceAck) AckIdNullValue() int64 {
	return math.MinInt64
}

func (*ServiceAck) RelevantIdId() uint16 {
	return 4
}

func (*ServiceAck) RelevantIdSinceVersion() uint16 {
	return 0
}

func (s *ServiceAck) RelevantIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.RelevantIdSinceVersion()
}

func (*ServiceAck) RelevantIdDeprecated() uint16 {
	return 0
}

func (*ServiceAck) RelevantIdMetaAttribute(meta int) string {
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

func (*ServiceAck) RelevantIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ServiceAck) RelevantIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ServiceAck) RelevantIdNullValue() int64 {
	return math.MinInt64
}

func (*ServiceAck) ServiceIdId() uint16 {
	return 5
}

func (*ServiceAck) ServiceIdSinceVersion() uint16 {
	return 0
}

func (s *ServiceAck) ServiceIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ServiceIdSinceVersion()
}

func (*ServiceAck) ServiceIdDeprecated() uint16 {
	return 0
}

func (*ServiceAck) ServiceIdMetaAttribute(meta int) string {
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

func (*ServiceAck) ServiceIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ServiceAck) ServiceIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ServiceAck) ServiceIdNullValue() int32 {
	return math.MinInt32
}
