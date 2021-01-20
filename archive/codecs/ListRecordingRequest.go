// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type ListRecordingRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
}

func (l *ListRecordingRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, l.RecordingId); err != nil {
		return err
	}
	return nil
}

func (l *ListRecordingRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !l.RecordingIdInActingVersion(actingVersion) {
		l.RecordingId = l.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &l.RecordingId); err != nil {
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

func (l *ListRecordingRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if l.RecordingIdInActingVersion(actingVersion) {
		if l.RecordingId < l.RecordingIdMinValue() || l.RecordingId > l.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on l.RecordingId (%v < %v > %v)", l.RecordingIdMinValue(), l.RecordingId, l.RecordingIdMaxValue())
		}
	}
	return nil
}

func ListRecordingRequestInit(l *ListRecordingRequest) {
	return
}

func (*ListRecordingRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*ListRecordingRequest) SbeTemplateId() (templateId uint16) {
	return 10
}

func (*ListRecordingRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*ListRecordingRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*ListRecordingRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*ListRecordingRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*ListRecordingRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.ControlSessionIdSinceVersion()
}

func (*ListRecordingRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*ListRecordingRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingRequest) CorrelationIdId() uint16 {
	return 2
}

func (*ListRecordingRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.CorrelationIdSinceVersion()
}

func (*ListRecordingRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*ListRecordingRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*ListRecordingRequest) RecordingIdId() uint16 {
	return 3
}

func (*ListRecordingRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (l *ListRecordingRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.RecordingIdSinceVersion()
}

func (*ListRecordingRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*ListRecordingRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*ListRecordingRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*ListRecordingRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*ListRecordingRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}
