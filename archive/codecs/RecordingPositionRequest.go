// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type RecordingPositionRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
}

func (r *RecordingPositionRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := r.RangeCheck(r.SbeSchemaVersion(), r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, r.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, r.RecordingId); err != nil {
		return err
	}
	return nil
}

func (r *RecordingPositionRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !r.ControlSessionIdInActingVersion(actingVersion) {
		r.ControlSessionId = r.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.ControlSessionId); err != nil {
			return err
		}
	}
	if !r.CorrelationIdInActingVersion(actingVersion) {
		r.CorrelationId = r.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.CorrelationId); err != nil {
			return err
		}
	}
	if !r.RecordingIdInActingVersion(actingVersion) {
		r.RecordingId = r.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &r.RecordingId); err != nil {
			return err
		}
	}
	if actingVersion > r.SbeSchemaVersion() && blockLength > r.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-r.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := r.RangeCheck(actingVersion, r.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (r *RecordingPositionRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if r.ControlSessionIdInActingVersion(actingVersion) {
		if r.ControlSessionId < r.ControlSessionIdMinValue() || r.ControlSessionId > r.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on r.ControlSessionId (%v < %v > %v)", r.ControlSessionIdMinValue(), r.ControlSessionId, r.ControlSessionIdMaxValue())
		}
	}
	if r.CorrelationIdInActingVersion(actingVersion) {
		if r.CorrelationId < r.CorrelationIdMinValue() || r.CorrelationId > r.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on r.CorrelationId (%v < %v > %v)", r.CorrelationIdMinValue(), r.CorrelationId, r.CorrelationIdMaxValue())
		}
	}
	if r.RecordingIdInActingVersion(actingVersion) {
		if r.RecordingId < r.RecordingIdMinValue() || r.RecordingId > r.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on r.RecordingId (%v < %v > %v)", r.RecordingIdMinValue(), r.RecordingId, r.RecordingIdMaxValue())
		}
	}
	return nil
}

func RecordingPositionRequestInit(r *RecordingPositionRequest) {
	return
}

func (*RecordingPositionRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*RecordingPositionRequest) SbeTemplateId() (templateId uint16) {
	return 12
}

func (*RecordingPositionRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*RecordingPositionRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*RecordingPositionRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*RecordingPositionRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*RecordingPositionRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingPositionRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.ControlSessionIdSinceVersion()
}

func (*RecordingPositionRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*RecordingPositionRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*RecordingPositionRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingPositionRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingPositionRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingPositionRequest) CorrelationIdId() uint16 {
	return 2
}

func (*RecordingPositionRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingPositionRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.CorrelationIdSinceVersion()
}

func (*RecordingPositionRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*RecordingPositionRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*RecordingPositionRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingPositionRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingPositionRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*RecordingPositionRequest) RecordingIdId() uint16 {
	return 3
}

func (*RecordingPositionRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (r *RecordingPositionRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= r.RecordingIdSinceVersion()
}

func (*RecordingPositionRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*RecordingPositionRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*RecordingPositionRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*RecordingPositionRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*RecordingPositionRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}
