// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type StopRecordingRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	StreamId         int32
	Channel          []uint8
}

func (s *StopRecordingRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt32(_w, s.StreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(s.Channel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.Channel); err != nil {
		return err
	}
	return nil
}

func (s *StopRecordingRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !s.StreamIdInActingVersion(actingVersion) {
		s.StreamId = s.StreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &s.StreamId); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}

	if s.ChannelInActingVersion(actingVersion) {
		var ChannelLength uint32
		if err := _m.ReadUint32(_r, &ChannelLength); err != nil {
			return err
		}
		if cap(s.Channel) < int(ChannelLength) {
			s.Channel = make([]uint8, ChannelLength)
		}
		s.Channel = s.Channel[:ChannelLength]
		if err := _m.ReadBytes(_r, s.Channel); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *StopRecordingRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if s.StreamIdInActingVersion(actingVersion) {
		if s.StreamId < s.StreamIdMinValue() || s.StreamId > s.StreamIdMaxValue() {
			return fmt.Errorf("Range check failed on s.StreamId (%v < %v > %v)", s.StreamIdMinValue(), s.StreamId, s.StreamIdMaxValue())
		}
	}
	return nil
}

func StopRecordingRequestInit(s *StopRecordingRequest) {
	return
}

func (*StopRecordingRequest) SbeBlockLength() (blockLength uint16) {
	return 20
}

func (*StopRecordingRequest) SbeTemplateId() (templateId uint16) {
	return 5
}

func (*StopRecordingRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*StopRecordingRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*StopRecordingRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*StopRecordingRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*StopRecordingRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (s *StopRecordingRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ControlSessionIdSinceVersion()
}

func (*StopRecordingRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*StopRecordingRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*StopRecordingRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopRecordingRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopRecordingRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*StopRecordingRequest) CorrelationIdId() uint16 {
	return 2
}

func (*StopRecordingRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *StopRecordingRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*StopRecordingRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*StopRecordingRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*StopRecordingRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StopRecordingRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StopRecordingRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*StopRecordingRequest) StreamIdId() uint16 {
	return 3
}

func (*StopRecordingRequest) StreamIdSinceVersion() uint16 {
	return 0
}

func (s *StopRecordingRequest) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.StreamIdSinceVersion()
}

func (*StopRecordingRequest) StreamIdDeprecated() uint16 {
	return 0
}

func (*StopRecordingRequest) StreamIdMetaAttribute(meta int) string {
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

func (*StopRecordingRequest) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*StopRecordingRequest) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*StopRecordingRequest) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*StopRecordingRequest) ChannelMetaAttribute(meta int) string {
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

func (*StopRecordingRequest) ChannelSinceVersion() uint16 {
	return 0
}

func (s *StopRecordingRequest) ChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ChannelSinceVersion()
}

func (*StopRecordingRequest) ChannelDeprecated() uint16 {
	return 0
}

func (StopRecordingRequest) ChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (StopRecordingRequest) ChannelHeaderLength() uint64 {
	return 4
}
