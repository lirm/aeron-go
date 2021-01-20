// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type StopReplicationRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	ReplicationId    int64
}

func (s *StopReplicationRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, s.ReplicationId); err != nil {
		return err
	}
	return nil
}

func (s *StopReplicationRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !s.ReplicationIdInActingVersion(actingVersion) {
		s.ReplicationId = s.ReplicationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &s.ReplicationId); err != nil {
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

func (s *StopReplicationRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if s.ReplicationIdInActingVersion(actingVersion) {
		if s.ReplicationId < s.ReplicationIdMinValue() || s.ReplicationId > s.ReplicationIdMaxValue() {
			return fmt.Errorf("Range check failed on s.ReplicationId (%v < %v > %v)", s.ReplicationIdMinValue(), s.ReplicationId, s.ReplicationIdMaxValue())
		}
	}
	return nil
}

func StopReplicationRequestInit(s *StopReplicationRequest) {
	return
}

func (*StopReplicationRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*StopReplicationRequest) SbeTemplateId() (templateId uint16) {
	return 51
}

func (*StopReplicationRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*StopReplicationRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*StopReplicationRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*StopReplicationRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*StopReplicationRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (s *StopReplicationRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ControlSessionIdSinceVersion()
}

func (*StopReplicationRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*StopReplicationRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*StopReplicationRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopReplicationRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopReplicationRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*StopReplicationRequest) CorrelationIdId() uint16 {
	return 2
}

func (*StopReplicationRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *StopReplicationRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*StopReplicationRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*StopReplicationRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*StopReplicationRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopReplicationRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopReplicationRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*StopReplicationRequest) ReplicationIdId() uint16 {
	return 3
}

func (*StopReplicationRequest) ReplicationIdSinceVersion() uint16 {
	return 0
}

func (s *StopReplicationRequest) ReplicationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ReplicationIdSinceVersion()
}

func (*StopReplicationRequest) ReplicationIdDeprecated() uint16 {
	return 0
}

func (*StopReplicationRequest) ReplicationIdMetaAttribute(meta int) string {
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

func (*StopReplicationRequest) ReplicationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopReplicationRequest) ReplicationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopReplicationRequest) ReplicationIdNullValue() int64 {
	return math.MinInt64
}
