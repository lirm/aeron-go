// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type AttachSegmentsRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
}

func (a *AttachSegmentsRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := a.RangeCheck(a.SbeSchemaVersion(), a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, a.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, a.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, a.RecordingId); err != nil {
		return err
	}
	return nil
}

func (a *AttachSegmentsRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !a.ControlSessionIdInActingVersion(actingVersion) {
		a.ControlSessionId = a.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.ControlSessionId); err != nil {
			return err
		}
	}
	if !a.CorrelationIdInActingVersion(actingVersion) {
		a.CorrelationId = a.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.CorrelationId); err != nil {
			return err
		}
	}
	if !a.RecordingIdInActingVersion(actingVersion) {
		a.RecordingId = a.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &a.RecordingId); err != nil {
			return err
		}
	}
	if actingVersion > a.SbeSchemaVersion() && blockLength > a.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-a.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := a.RangeCheck(actingVersion, a.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (a *AttachSegmentsRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if a.ControlSessionIdInActingVersion(actingVersion) {
		if a.ControlSessionId < a.ControlSessionIdMinValue() || a.ControlSessionId > a.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on a.ControlSessionId (%v < %v > %v)", a.ControlSessionIdMinValue(), a.ControlSessionId, a.ControlSessionIdMaxValue())
		}
	}
	if a.CorrelationIdInActingVersion(actingVersion) {
		if a.CorrelationId < a.CorrelationIdMinValue() || a.CorrelationId > a.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on a.CorrelationId (%v < %v > %v)", a.CorrelationIdMinValue(), a.CorrelationId, a.CorrelationIdMaxValue())
		}
	}
	if a.RecordingIdInActingVersion(actingVersion) {
		if a.RecordingId < a.RecordingIdMinValue() || a.RecordingId > a.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on a.RecordingId (%v < %v > %v)", a.RecordingIdMinValue(), a.RecordingId, a.RecordingIdMaxValue())
		}
	}
	return nil
}

func AttachSegmentsRequestInit(a *AttachSegmentsRequest) {
	return
}

func (*AttachSegmentsRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*AttachSegmentsRequest) SbeTemplateId() (templateId uint16) {
	return 56
}

func (*AttachSegmentsRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*AttachSegmentsRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*AttachSegmentsRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*AttachSegmentsRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*AttachSegmentsRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (a *AttachSegmentsRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.ControlSessionIdSinceVersion()
}

func (*AttachSegmentsRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*AttachSegmentsRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*AttachSegmentsRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AttachSegmentsRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AttachSegmentsRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*AttachSegmentsRequest) CorrelationIdId() uint16 {
	return 2
}

func (*AttachSegmentsRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (a *AttachSegmentsRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.CorrelationIdSinceVersion()
}

func (*AttachSegmentsRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*AttachSegmentsRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*AttachSegmentsRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AttachSegmentsRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AttachSegmentsRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*AttachSegmentsRequest) RecordingIdId() uint16 {
	return 3
}

func (*AttachSegmentsRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (a *AttachSegmentsRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= a.RecordingIdSinceVersion()
}

func (*AttachSegmentsRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*AttachSegmentsRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*AttachSegmentsRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*AttachSegmentsRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*AttachSegmentsRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}
