// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ExtendRecordingRequest2 struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
	StreamId         int32
	SourceLocation   SourceLocationEnum
	AutoStop         BooleanTypeEnum
	Channel          []uint8
}

func (e *ExtendRecordingRequest2) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := e.AutoStop.Encode(_m, _w); err != nil {
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

func (e *ExtendRecordingRequest2) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if e.AutoStopInActingVersion(actingVersion) {
		if err := e.AutoStop.Decode(_m, _r, actingVersion); err != nil {
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

func (e *ExtendRecordingRequest2) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if err := e.AutoStop.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func ExtendRecordingRequest2Init(e *ExtendRecordingRequest2) {
	return
}

func (*ExtendRecordingRequest2) SbeBlockLength() (blockLength uint16) {
	return 36
}

func (*ExtendRecordingRequest2) SbeTemplateId() (templateId uint16) {
	return 64
}

func (*ExtendRecordingRequest2) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ExtendRecordingRequest2) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*ExtendRecordingRequest2) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ExtendRecordingRequest2) ControlSessionIdId() uint16 {
	return 1
}

func (*ExtendRecordingRequest2) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest2) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.ControlSessionIdSinceVersion()
}

func (*ExtendRecordingRequest2) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest2) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest2) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ExtendRecordingRequest2) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ExtendRecordingRequest2) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ExtendRecordingRequest2) CorrelationIdId() uint16 {
	return 2
}

func (*ExtendRecordingRequest2) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest2) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.CorrelationIdSinceVersion()
}

func (*ExtendRecordingRequest2) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest2) CorrelationIdMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest2) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ExtendRecordingRequest2) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ExtendRecordingRequest2) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ExtendRecordingRequest2) RecordingIdId() uint16 {
	return 3
}

func (*ExtendRecordingRequest2) RecordingIdSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest2) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.RecordingIdSinceVersion()
}

func (*ExtendRecordingRequest2) RecordingIdDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest2) RecordingIdMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest2) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ExtendRecordingRequest2) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ExtendRecordingRequest2) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*ExtendRecordingRequest2) StreamIdId() uint16 {
	return 4
}

func (*ExtendRecordingRequest2) StreamIdSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest2) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.StreamIdSinceVersion()
}

func (*ExtendRecordingRequest2) StreamIdDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest2) StreamIdMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest2) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ExtendRecordingRequest2) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ExtendRecordingRequest2) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ExtendRecordingRequest2) SourceLocationId() uint16 {
	return 5
}

func (*ExtendRecordingRequest2) SourceLocationSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest2) SourceLocationInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.SourceLocationSinceVersion()
}

func (*ExtendRecordingRequest2) SourceLocationDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest2) SourceLocationMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest2) AutoStopId() uint16 {
	return 6
}

func (*ExtendRecordingRequest2) AutoStopSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest2) AutoStopInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.AutoStopSinceVersion()
}

func (*ExtendRecordingRequest2) AutoStopDeprecated() uint16 {
	return 0
}

func (*ExtendRecordingRequest2) AutoStopMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest2) ChannelMetaAttribute(meta int) string {
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

func (*ExtendRecordingRequest2) ChannelSinceVersion() uint16 {
	return 0
}

func (e *ExtendRecordingRequest2) ChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.ChannelSinceVersion()
}

func (*ExtendRecordingRequest2) ChannelDeprecated() uint16 {
	return 0
}

func (ExtendRecordingRequest2) ChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ExtendRecordingRequest2) ChannelHeaderLength() uint64 {
	return 4
}
