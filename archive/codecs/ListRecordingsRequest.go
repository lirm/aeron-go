// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ListRecordingsRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	FromRecordingId  int64
	RecordCount      int32
}

func (l *ListRecordingsRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	return nil
}

func (l *ListRecordingsRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if actingVersion > l.SbeSchemaVersion() && blockLength > l.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-l.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := l.RangeCheck(actingVersion, l.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (l *ListRecordingsRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	return nil
}

func ListRecordingsRequestInit(l *ListRecordingsRequest) {
	return
}

func (*ListRecordingsRequest) SbeBlockLength() (blockLength uint16) {
	return 28
}

func (*ListRecordingsRequest) SbeTemplateId() (templateId uint16) {
	return 8
}

func (*ListRecordingsRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ListRecordingsRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*ListRecordingsRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ListRecordingsRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*ListRecordingsRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.ControlSessionIdSinceVersion()
}

func (*ListRecordingsRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingsRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ListRecordingsRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingsRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingsRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingsRequest) CorrelationIdId() uint16 {
	return 2
}

func (*ListRecordingsRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.CorrelationIdSinceVersion()
}

func (*ListRecordingsRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingsRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*ListRecordingsRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingsRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingsRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingsRequest) FromRecordingIdId() uint16 {
	return 3
}

func (*ListRecordingsRequest) FromRecordingIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsRequest) FromRecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.FromRecordingIdSinceVersion()
}

func (*ListRecordingsRequest) FromRecordingIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingsRequest) FromRecordingIdMetaAttribute(meta int) string {
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

func (*ListRecordingsRequest) FromRecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingsRequest) FromRecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingsRequest) FromRecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingsRequest) RecordCountId() uint16 {
	return 4
}

func (*ListRecordingsRequest) RecordCountSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingsRequest) RecordCountInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.RecordCountSinceVersion()
}

func (*ListRecordingsRequest) RecordCountDeprecated() uint16 {
	return 0
}

func (*ListRecordingsRequest) RecordCountMetaAttribute(meta int) string {
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

func (*ListRecordingsRequest) RecordCountMinValue() int32 {
	return math.MinInt32 + 1
}

func (*ListRecordingsRequest) RecordCountMaxValue() int32 {
	return math.MaxInt32
}

func (*ListRecordingsRequest) RecordCountNullValue() int32 {
	return math.MinInt32
}
