// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type DetachSegmentsRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
	NewStartPosition int64
}

func (d *DetachSegmentsRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteInt64(_w, d.NewStartPosition); err != nil {
		return err
	}
	return nil
}

func (d *DetachSegmentsRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !d.NewStartPositionInActingVersion(actingVersion) {
		d.NewStartPosition = d.NewStartPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &d.NewStartPosition); err != nil {
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

func (d *DetachSegmentsRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if d.NewStartPositionInActingVersion(actingVersion) {
		if d.NewStartPosition < d.NewStartPositionMinValue() || d.NewStartPosition > d.NewStartPositionMaxValue() {
			return fmt.Errorf("Range check failed on d.NewStartPosition (%v < %v > %v)", d.NewStartPositionMinValue(), d.NewStartPosition, d.NewStartPositionMaxValue())
		}
	}
	return nil
}

func DetachSegmentsRequestInit(d *DetachSegmentsRequest) {
	return
}

func (*DetachSegmentsRequest) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*DetachSegmentsRequest) SbeTemplateId() (templateId uint16) {
	return 53
}

func (*DetachSegmentsRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*DetachSegmentsRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*DetachSegmentsRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*DetachSegmentsRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*DetachSegmentsRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (d *DetachSegmentsRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= d.ControlSessionIdSinceVersion()
}

func (*DetachSegmentsRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*DetachSegmentsRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*DetachSegmentsRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*DetachSegmentsRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*DetachSegmentsRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*DetachSegmentsRequest) CorrelationIdId() uint16 {
	return 2
}

func (*DetachSegmentsRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (d *DetachSegmentsRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= d.CorrelationIdSinceVersion()
}

func (*DetachSegmentsRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*DetachSegmentsRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*DetachSegmentsRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*DetachSegmentsRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*DetachSegmentsRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*DetachSegmentsRequest) RecordingIdId() uint16 {
	return 3
}

func (*DetachSegmentsRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (d *DetachSegmentsRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= d.RecordingIdSinceVersion()
}

func (*DetachSegmentsRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*DetachSegmentsRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*DetachSegmentsRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*DetachSegmentsRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*DetachSegmentsRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*DetachSegmentsRequest) NewStartPositionId() uint16 {
	return 4
}

func (*DetachSegmentsRequest) NewStartPositionSinceVersion() uint16 {
	return 0
}

func (d *DetachSegmentsRequest) NewStartPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= d.NewStartPositionSinceVersion()
}

func (*DetachSegmentsRequest) NewStartPositionDeprecated() uint16 {
	return 0
}

func (*DetachSegmentsRequest) NewStartPositionMetaAttribute(meta int) string {
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

func (*DetachSegmentsRequest) NewStartPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*DetachSegmentsRequest) NewStartPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*DetachSegmentsRequest) NewStartPositionNullValue() int64 {
	return math.MinInt64
}
