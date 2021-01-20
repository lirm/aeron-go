// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type StartRecordingRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	StreamId         int32
	SourceLocation   SourceLocationEnum
	Channel          []uint8
}

func (s *StartRecordingRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := s.SourceLocation.Encode(_m, _w); err != nil {
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

func (s *StartRecordingRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if s.SourceLocationInActingVersion(actingVersion) {
		if err := s.SourceLocation.Decode(_m, _r, actingVersion); err != nil {
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

func (s *StartRecordingRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if err := s.SourceLocation.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func StartRecordingRequestInit(s *StartRecordingRequest) {
	return
}

func (*StartRecordingRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*StartRecordingRequest) SbeTemplateId() (templateId uint16) {
	return 4
}

func (*StartRecordingRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*StartRecordingRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*StartRecordingRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*StartRecordingRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*StartRecordingRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (s *StartRecordingRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ControlSessionIdSinceVersion()
}

func (*StartRecordingRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*StartRecordingRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*StartRecordingRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StartRecordingRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StartRecordingRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*StartRecordingRequest) CorrelationIdId() uint16 {
	return 2
}

func (*StartRecordingRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (s *StartRecordingRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.CorrelationIdSinceVersion()
}

func (*StartRecordingRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*StartRecordingRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*StartRecordingRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*StartRecordingRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*StartRecordingRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*StartRecordingRequest) StreamIdId() uint16 {
	return 3
}

func (*StartRecordingRequest) StreamIdSinceVersion() uint16 {
	return 0
}

func (s *StartRecordingRequest) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.StreamIdSinceVersion()
}

func (*StartRecordingRequest) StreamIdDeprecated() uint16 {
	return 0
}

func (*StartRecordingRequest) StreamIdMetaAttribute(meta int) string {
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

func (*StartRecordingRequest) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*StartRecordingRequest) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*StartRecordingRequest) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*StartRecordingRequest) SourceLocationId() uint16 {
	return 4
}

func (*StartRecordingRequest) SourceLocationSinceVersion() uint16 {
	return 0
}

func (s *StartRecordingRequest) SourceLocationInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.SourceLocationSinceVersion()
}

func (*StartRecordingRequest) SourceLocationDeprecated() uint16 {
	return 0
}

func (*StartRecordingRequest) SourceLocationMetaAttribute(meta int) string {
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

func (*StartRecordingRequest) ChannelMetaAttribute(meta int) string {
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

func (*StartRecordingRequest) ChannelSinceVersion() uint16 {
	return 0
}

func (s *StartRecordingRequest) ChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ChannelSinceVersion()
}

func (*StartRecordingRequest) ChannelDeprecated() uint16 {
	return 0
}

func (StartRecordingRequest) ChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (StartRecordingRequest) ChannelHeaderLength() uint64 {
	return 4
}
