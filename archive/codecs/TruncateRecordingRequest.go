// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type TruncateRecordingRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
	Position         int64
}

func (t *TruncateRecordingRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := t.RangeCheck(t.SbeSchemaVersion(), t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, t.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, t.Position); err != nil {
		return err
	}
	return nil
}

func (t *TruncateRecordingRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !t.ControlSessionIdInActingVersion(actingVersion) {
		t.ControlSessionId = t.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.ControlSessionId); err != nil {
			return err
		}
	}
	if !t.CorrelationIdInActingVersion(actingVersion) {
		t.CorrelationId = t.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.CorrelationId); err != nil {
			return err
		}
	}
	if !t.RecordingIdInActingVersion(actingVersion) {
		t.RecordingId = t.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.RecordingId); err != nil {
			return err
		}
	}
	if !t.PositionInActingVersion(actingVersion) {
		t.Position = t.PositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &t.Position); err != nil {
			return err
		}
	}
	if actingVersion > t.SbeSchemaVersion() && blockLength > t.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-t.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := t.RangeCheck(actingVersion, t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (t *TruncateRecordingRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if t.ControlSessionIdInActingVersion(actingVersion) {
		if t.ControlSessionId < t.ControlSessionIdMinValue() || t.ControlSessionId > t.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on t.ControlSessionId (%v < %v > %v)", t.ControlSessionIdMinValue(), t.ControlSessionId, t.ControlSessionIdMaxValue())
		}
	}
	if t.CorrelationIdInActingVersion(actingVersion) {
		if t.CorrelationId < t.CorrelationIdMinValue() || t.CorrelationId > t.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on t.CorrelationId (%v < %v > %v)", t.CorrelationIdMinValue(), t.CorrelationId, t.CorrelationIdMaxValue())
		}
	}
	if t.RecordingIdInActingVersion(actingVersion) {
		if t.RecordingId < t.RecordingIdMinValue() || t.RecordingId > t.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on t.RecordingId (%v < %v > %v)", t.RecordingIdMinValue(), t.RecordingId, t.RecordingIdMaxValue())
		}
	}
	if t.PositionInActingVersion(actingVersion) {
		if t.Position < t.PositionMinValue() || t.Position > t.PositionMaxValue() {
			return fmt.Errorf("Range check failed on t.Position (%v < %v > %v)", t.PositionMinValue(), t.Position, t.PositionMaxValue())
		}
	}
	return nil
}

func TruncateRecordingRequestInit(t *TruncateRecordingRequest) {
	return
}

func (*TruncateRecordingRequest) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*TruncateRecordingRequest) SbeTemplateId() (templateId uint16) {
	return 13
}

func (*TruncateRecordingRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*TruncateRecordingRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*TruncateRecordingRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*TruncateRecordingRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*TruncateRecordingRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (t *TruncateRecordingRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.ControlSessionIdSinceVersion()
}

func (*TruncateRecordingRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*TruncateRecordingRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*TruncateRecordingRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TruncateRecordingRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TruncateRecordingRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*TruncateRecordingRequest) CorrelationIdId() uint16 {
	return 2
}

func (*TruncateRecordingRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (t *TruncateRecordingRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.CorrelationIdSinceVersion()
}

func (*TruncateRecordingRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*TruncateRecordingRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*TruncateRecordingRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TruncateRecordingRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TruncateRecordingRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*TruncateRecordingRequest) RecordingIdId() uint16 {
	return 3
}

func (*TruncateRecordingRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (t *TruncateRecordingRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.RecordingIdSinceVersion()
}

func (*TruncateRecordingRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*TruncateRecordingRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*TruncateRecordingRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TruncateRecordingRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*TruncateRecordingRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*TruncateRecordingRequest) PositionId() uint16 {
	return 4
}

func (*TruncateRecordingRequest) PositionSinceVersion() uint16 {
	return 0
}

func (t *TruncateRecordingRequest) PositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.PositionSinceVersion()
}

func (*TruncateRecordingRequest) PositionDeprecated() uint16 {
	return 0
}

func (*TruncateRecordingRequest) PositionMetaAttribute(meta int) string {
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

func (*TruncateRecordingRequest) PositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*TruncateRecordingRequest) PositionMaxValue() int64 {
	return math.MaxInt64
}

func (*TruncateRecordingRequest) PositionNullValue() int64 {
	return math.MinInt64
}
