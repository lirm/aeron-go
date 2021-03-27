// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type DeleteDetachedSegmentsRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
}

func (d *DeleteDetachedSegmentsRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := d.RangeCheck(d.SbeSchemaVersion(), d.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, d.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, d.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, d.RecordingId); err != nil {
		return err
	}
	return nil
}

func (d *DeleteDetachedSegmentsRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !d.ControlSessionIdInActingVersion(actingVersion) {
		d.ControlSessionId = d.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &d.ControlSessionId); err != nil {
			return err
		}
	}
	if !d.CorrelationIdInActingVersion(actingVersion) {
		d.CorrelationId = d.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &d.CorrelationId); err != nil {
			return err
		}
	}
	if !d.RecordingIdInActingVersion(actingVersion) {
		d.RecordingId = d.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &d.RecordingId); err != nil {
			return err
		}
	}
	if actingVersion > d.SbeSchemaVersion() && blockLength > d.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-d.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := d.RangeCheck(actingVersion, d.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (d *DeleteDetachedSegmentsRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if d.ControlSessionIdInActingVersion(actingVersion) {
		if d.ControlSessionId < d.ControlSessionIdMinValue() || d.ControlSessionId > d.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on d.ControlSessionId (%v < %v > %v)", d.ControlSessionIdMinValue(), d.ControlSessionId, d.ControlSessionIdMaxValue())
		}
	}
	if d.CorrelationIdInActingVersion(actingVersion) {
		if d.CorrelationId < d.CorrelationIdMinValue() || d.CorrelationId > d.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on d.CorrelationId (%v < %v > %v)", d.CorrelationIdMinValue(), d.CorrelationId, d.CorrelationIdMaxValue())
		}
	}
	if d.RecordingIdInActingVersion(actingVersion) {
		if d.RecordingId < d.RecordingIdMinValue() || d.RecordingId > d.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on d.RecordingId (%v < %v > %v)", d.RecordingIdMinValue(), d.RecordingId, d.RecordingIdMaxValue())
		}
	}
	return nil
}

func DeleteDetachedSegmentsRequestInit(d *DeleteDetachedSegmentsRequest) {
	return
}

func (*DeleteDetachedSegmentsRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*DeleteDetachedSegmentsRequest) SbeTemplateId() (templateId uint16) {
	return 54
}

func (*DeleteDetachedSegmentsRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*DeleteDetachedSegmentsRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*DeleteDetachedSegmentsRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*DeleteDetachedSegmentsRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*DeleteDetachedSegmentsRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (d *DeleteDetachedSegmentsRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= d.ControlSessionIdSinceVersion()
}

func (*DeleteDetachedSegmentsRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*DeleteDetachedSegmentsRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*DeleteDetachedSegmentsRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*DeleteDetachedSegmentsRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*DeleteDetachedSegmentsRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*DeleteDetachedSegmentsRequest) CorrelationIdId() uint16 {
	return 2
}

func (*DeleteDetachedSegmentsRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (d *DeleteDetachedSegmentsRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= d.CorrelationIdSinceVersion()
}

func (*DeleteDetachedSegmentsRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*DeleteDetachedSegmentsRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*DeleteDetachedSegmentsRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*DeleteDetachedSegmentsRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*DeleteDetachedSegmentsRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*DeleteDetachedSegmentsRequest) RecordingIdId() uint16 {
	return 3
}

func (*DeleteDetachedSegmentsRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (d *DeleteDetachedSegmentsRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= d.RecordingIdSinceVersion()
}

func (*DeleteDetachedSegmentsRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*DeleteDetachedSegmentsRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*DeleteDetachedSegmentsRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*DeleteDetachedSegmentsRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*DeleteDetachedSegmentsRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}
