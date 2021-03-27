// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ListRecordingsForUriRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	FromRecordingId  int64
	RecordCount      int32
	StreamId         int32
	Channel          []uint8
}

func (l *ListRecordingsForUriRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := l.RangeCheck(l.SbeSchemaVersion(), l.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, l.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, l.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, l.FromRecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, l.RecordCount); err != nil {
		return err
	}
	if err := _m.WriteInt32(_w, l.StreamId); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(l.Channel))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, l.Channel); err != nil {
		return err
	}
	return nil
}

func (l *ListRecordingsForUriRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !l.ControlSessionIdInActingVersion(actingVersion) {
		l.ControlSessionId = l.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &l.ControlSessionId); err != nil {
			return err
		}
	}
	if !l.CorrelationIdInActingVersion(actingVersion) {
		l.CorrelationId = l.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &l.CorrelationId); err != nil {
			return err
		}
	}
	if !l.FromRecordingIdInActingVersion(actingVersion) {
		l.FromRecordingId = l.FromRecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &l.FromRecordingId); err != nil {
			return err
		}
	}
	if !l.RecordCountInActingVersion(actingVersion) {
		l.RecordCount = l.RecordCountNullValue()
	} else {
		if err := _m.ReadInt32(_r, &l.RecordCount); err != nil {
			return err
		}
	}
	if !l.StreamIdInActingVersion(actingVersion) {
		l.StreamId = l.StreamIdNullValue()
	} else {
		if err := _m.ReadInt32(_r, &l.StreamId); err != nil {
			return err
		}
	}
	if actingVersion > l.SbeSchemaVersion() && blockLength > l.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-l.SbeBlockLength()))
	}

	if l.ChannelInActingVersion(actingVersion) {
		var ChannelLength uint32
		if err := _m.ReadUint32(_r, &ChannelLength); err != nil {
			return err
		}
		if cap(l.Channel) < int(ChannelLength) {
			l.Channel = make([]uint8, ChannelLength)
		}
		l.Channel = l.Channel[:ChannelLength]
		if err := _m.ReadBytes(_r, l.Channel); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := l.RangeCheck(actingVersion, l.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (l *ListRecordingsForUriRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if l.ControlSessionIdInActingVersion(actingVersion) {
		if l.ControlSessionId < l.ControlSessionIdMinValue() || l.ControlSessionId > l.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on l.ControlSessionId (%v < %v > %v)", l.ControlSessionIdMinValue(), l.ControlSessionId, l.ControlSessionIdMaxValue())
		}
	}
	if l.CorrelationIdInActingVersion(actingVersion) {
		if l.CorrelationId < l.CorrelationIdMinValue() || l.CorrelationId > l.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on l.CorrelationId (%v < %v > %v)", l.CorrelationIdMinValue(), l.CorrelationId, l.CorrelationIdMaxValue())
		}
	}
	if l.FromRecordingIdInActingVersion(actingVersion) {
		if l.FromRecordingId < l.FromRecordingIdMinValue() || l.FromRecordingId > l.FromRecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on l.FromRecordingId (%v < %v > %v)", l.FromRecordingIdMinValue(), l.FromRecordingId, l.FromRecordingIdMaxValue())
		}
	}
	if l.RecordCountInActingVersion(actingVersion) {
		if l.RecordCount < l.RecordCountMinValue() || l.RecordCount > l.RecordCountMaxValue() {
			return fmt.Errorf("Range check failed on l.RecordCount (%v < %v > %v)", l.RecordCountMinValue(), l.RecordCount, l.RecordCountMaxValue())
		}
	}
	if l.StreamIdInActingVersion(actingVersion) {
		if l.StreamId < l.StreamIdMinValue() || l.StreamId > l.StreamIdMaxValue() {
			return fmt.Errorf("Range check failed on l.StreamId (%v < %v > %v)", l.StreamIdMinValue(), l.StreamId, l.StreamIdMaxValue())
		}
	}
	return nil
}

func ListRecordingsForUriRequestInit(l *ListRecordingsForUriRequest) {
	return
}

func (*ListRecordingsForUriRequest) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*ListRecordingsForUriRequest) SbeTemplateId() (templateId uint16) {
	return 9
}

func (*ListRecordingsForUriRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ListRecordingsForUriRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*ListRecordingsForUriRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ListRecordingsForUriRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*ListRecordingsForUriRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsForUriRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.ControlSessionIdSinceVersion()
}

func (*ListRecordingsForUriRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingsForUriRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ListRecordingsForUriRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingsForUriRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingsForUriRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingsForUriRequest) CorrelationIdId() uint16 {
	return 2
}

func (*ListRecordingsForUriRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsForUriRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.CorrelationIdSinceVersion()
}

func (*ListRecordingsForUriRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingsForUriRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*ListRecordingsForUriRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingsForUriRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingsForUriRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingsForUriRequest) FromRecordingIdId() uint16 {
	return 3
}

func (*ListRecordingsForUriRequest) FromRecordingIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsForUriRequest) FromRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.FromRecordingIdSinceVersion()
}

func (*ListRecordingsForUriRequest) FromRecordingIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingsForUriRequest) FromRecordingIdMetaAttribute(meta int) string {
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

func (*ListRecordingsForUriRequest) FromRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingsForUriRequest) FromRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingsForUriRequest) FromRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingsForUriRequest) RecordCountId() uint16 {
	return 4
}

func (*ListRecordingsForUriRequest) RecordCountSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsForUriRequest) RecordCountInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.RecordCountSinceVersion()
}

func (*ListRecordingsForUriRequest) RecordCountDeprecated() uint16 {
	return 0
}

func (*ListRecordingsForUriRequest) RecordCountMetaAttribute(meta int) string {
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

func (*ListRecordingsForUriRequest) RecordCountMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ListRecordingsForUriRequest) RecordCountMaxValue() int32 {
	return math.MaxInt32
}

func (*ListRecordingsForUriRequest) RecordCountNullValue() int32 {
	return math.MinInt32
}

func (*ListRecordingsForUriRequest) StreamIdId() uint16 {
	return 5
}

func (*ListRecordingsForUriRequest) StreamIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsForUriRequest) StreamIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.StreamIdSinceVersion()
}

func (*ListRecordingsForUriRequest) StreamIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingsForUriRequest) StreamIdMetaAttribute(meta int) string {
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

func (*ListRecordingsForUriRequest) StreamIdMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ListRecordingsForUriRequest) StreamIdMaxValue() int32 {
	return math.MaxInt32
}

func (*ListRecordingsForUriRequest) StreamIdNullValue() int32 {
	return math.MinInt32
}

func (*ListRecordingsForUriRequest) ChannelMetaAttribute(meta int) string {
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

func (*ListRecordingsForUriRequest) ChannelSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsForUriRequest) ChannelInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.ChannelSinceVersion()
}

func (*ListRecordingsForUriRequest) ChannelDeprecated() uint16 {
	return 0
}

func (ListRecordingsForUriRequest) ChannelCharacterEncoding() string {
	return "US-ASCII"
}

func (ListRecordingsForUriRequest) ChannelHeaderLength() uint64 {
	return 4
}
