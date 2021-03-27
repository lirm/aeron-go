// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ExtendRecordingRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
	StreamId         int32
	SourceLocation   SourceLocationEnum
	Channel          []uint8
}

func (e *ExtendRecordingRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := e.RangeCheck(e.SbeSchemaVersion(), e.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, e.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, e.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, e.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, e.StreamId); err != nil {
		return err
	}
	if err := e.SourceLocation.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(e.Channel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, e.Channel); err != nil {
		return err
	}
	return nil
}

func (e *ExtendRecordingRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !e.ControlSessionIdInActingVersion(actingVersion) {
		e.ControlSessionId = e.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &e.ControlSessionId); err != nil {
			return err
		}
	}
	if !e.CorrelationIdInActingVersion(actingVersion) {
		e.CorrelationId = e.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &e.CorrelationId); err != nil {
			return err
		}
	}
	if !e.RecordingIdInActingVersion(actingVersion) {
		e.RecordingId = e.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &e.RecordingId); err != nil {
			return err
		}
	}
	if !e.StreamIdInActingVersion(actingVersion) {
		e.StreamId = e.StreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &e.StreamId); err != nil {
			return err
		}
	}
	if e.SourceLocationInActingVersion(actingVersion) {
		if err := e.SourceLocation.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if actingVersion > e.SbeSchemaVersion() && blockLength > e.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-e.SbeBlockLength()))
	}

	if e.ChannelInActingVersion(actingVersion) {
		var ChannelLength uint32
		if err := _m.ReadUint32(_r, &ChannelLength); err != nil {
			return err
		}
		if cap(e.Channel) < int(ChannelLength) {
			e.Channel = make([]uint8, ChannelLength)
		}
		e.Channel = e.Channel[:ChannelLength]
		if err := _m.ReadBytes(_r, e.Channel); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := e.RangeCheck(actingVersion, e.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (e *ExtendRecordingRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if e.ControlSessionIdInActingVersion(actingVersion) {
		if e.ControlSessionId < e.ControlSessionIdMinValue() || e.ControlSessionId > e.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on e.ControlSessionId (%v < %v > %v)", e.ControlSessionIdMinValue(), e.ControlSessionId, e.ControlSessionIdMaxValue())
		}
	}
	if e.CorrelationIdInActingVersion(actingVersion) {
		if e.CorrelationId < e.CorrelationIdMinValue() || e.CorrelationId > e.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on e.CorrelationId (%v < %v > %v)", e.CorrelationIdMinValue(), e.CorrelationId, e.CorrelationIdMaxValue())
		}
	}
	if e.RecordingIdInActingVersion(actingVersion) {
		if e.RecordingId < e.RecordingIdMinValue() || e.RecordingId > e.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on e.RecordingId (%v < %v > %v)", e.RecordingIdMinValue(), e.RecordingId, e.RecordingIdMaxValue())
		}
	}
	if e.StreamIdInActingVersion(actingVersion) {
		if e.StreamId < e.StreamIdMinValue() || e.StreamId > e.StreamIdMaxValue() {
			return fmt.Errorf("Range check failed on e.StreamId (%v < %v > %v)", e.StreamIdMinValue(), e.StreamId, e.StreamIdMaxValue())
		}
	}
	if err := e.SourceLocation.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func ExtendRecordingRequestInit(e *ExtendRecordingRequest) {
	return
}

func (*ExtendRecordingRequest) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*ExtendRecordingRequest) SbeTemplateId() (templateId uint16) {
	return 11
}

func (*ExtendRecordingRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ExtendRecordingRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*ExtendRecordingRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ExtendRecordingRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*ExtendRecordingRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.ControlSessionIdSinceVersion()
}

func (*ExtendRecordingRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ExtendRecordingRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ExtendRecordingRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ExtendRecordingRequest) CorrelationIdId() uint16 {
	return 2
}

func (*ExtendRecordingRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.CorrelationIdSinceVersion()
}

func (*ExtendRecordingRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ExtendRecordingRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ExtendRecordingRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ExtendRecordingRequest) RecordingIdId() uint16 {
	return 3
}

func (*ExtendRecordingRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.RecordingIdSinceVersion()
}

func (*ExtendRecordingRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ExtendRecordingRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ExtendRecordingRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*ExtendRecordingRequest) StreamIdId() uint16 {
	return 4
}

func (*ExtendRecordingRequest) StreamIdSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.StreamIdSinceVersion()
}

func (*ExtendRecordingRequest) StreamIdDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest) StreamIdMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ExtendRecordingRequest) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ExtendRecordingRequest) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ExtendRecordingRequest) SourceLocationId() uint16 {
	return 5
}

func (*ExtendRecordingRequest) SourceLocationSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest) SourceLocationInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.SourceLocationSinceVersion()
}

func (*ExtendRecordingRequest) SourceLocationDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest) SourceLocationMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest) ChannelMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest) ChannelSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest) ChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.ChannelSinceVersion()
}

func (*ExtendRecordingRequest) ChannelDeprecated() uint16 {
	return 0
}

func (ExtendRecordingRequest) ChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ExtendRecordingRequest) ChannelHeaderLength() uint64 {
	return 4
}
