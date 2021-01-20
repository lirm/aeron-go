// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type PurgeSegmentsRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
	NewStartPosition int64
}

func (p *PurgeSegmentsRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := p.RangeCheck(p.SbeSchemaVersion(), p.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, p.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, p.CorrelationId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, p.RecordingId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, p.NewStartPosition); err != nil {
		return err
	}
	return nil
}

func (p *PurgeSegmentsRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !p.ControlSessionIdInActingVersion(actingVersion) {
		p.ControlSessionId = p.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &p.ControlSessionId); err != nil {
			return err
		}
	}
	if !p.CorrelationIdInActingVersion(actingVersion) {
		p.CorrelationId = p.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &p.CorrelationId); err != nil {
			return err
		}
	}
	if !p.RecordingIdInActingVersion(actingVersion) {
		p.RecordingId = p.RecordingIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &p.RecordingId); err != nil {
			return err
		}
	}
	if !p.NewStartPositionInActingVersion(actingVersion) {
		p.NewStartPosition = p.NewStartPositionNullValue()
	} else {
		if err := _m.ReadInt64(_r, &p.NewStartPosition); err != nil {
			return err
		}
	}
	if actingVersion > p.SbeSchemaVersion() && blockLength > p.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-p.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := p.RangeCheck(actingVersion, p.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (p *PurgeSegmentsRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if p.ControlSessionIdInActingVersion(actingVersion) {
		if p.ControlSessionId < p.ControlSessionIdMinValue() || p.ControlSessionId > p.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on p.ControlSessionId (%v < %v > %v)", p.ControlSessionIdMinValue(), p.ControlSessionId, p.ControlSessionIdMaxValue())
		}
	}
	if p.CorrelationIdInActingVersion(actingVersion) {
		if p.CorrelationId < p.CorrelationIdMinValue() || p.CorrelationId > p.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on p.CorrelationId (%v < %v > %v)", p.CorrelationIdMinValue(), p.CorrelationId, p.CorrelationIdMaxValue())
		}
	}
	if p.RecordingIdInActingVersion(actingVersion) {
		if p.RecordingId < p.RecordingIdMinValue() || p.RecordingId > p.RecordingIdMaxValue() {
			return fmt.Errorf("Range check failed on p.RecordingId (%v < %v > %v)", p.RecordingIdMinValue(), p.RecordingId, p.RecordingIdMaxValue())
		}
	}
	if p.NewStartPositionInActingVersion(actingVersion) {
		if p.NewStartPosition < p.NewStartPositionMinValue() || p.NewStartPosition > p.NewStartPositionMaxValue() {
			return fmt.Errorf("Range check failed on p.NewStartPosition (%v < %v > %v)", p.NewStartPositionMinValue(), p.NewStartPosition, p.NewStartPositionMaxValue())
		}
	}
	return nil
}

func PurgeSegmentsRequestInit(p *PurgeSegmentsRequest) {
	return
}

func (*PurgeSegmentsRequest) SbeBlockLength() (blockLength uint16) {
	return 32
}

func (*PurgeSegmentsRequest) SbeTemplateId() (templateId uint16) {
	return 55
}

func (*PurgeSegmentsRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*PurgeSegmentsRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*PurgeSegmentsRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*PurgeSegmentsRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*PurgeSegmentsRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (p *PurgeSegmentsRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.ControlSessionIdSinceVersion()
}

func (*PurgeSegmentsRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*PurgeSegmentsRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*PurgeSegmentsRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*PurgeSegmentsRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*PurgeSegmentsRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*PurgeSegmentsRequest) CorrelationIdId() uint16 {
	return 2
}

func (*PurgeSegmentsRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (p *PurgeSegmentsRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.CorrelationIdSinceVersion()
}

func (*PurgeSegmentsRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*PurgeSegmentsRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*PurgeSegmentsRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*PurgeSegmentsRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*PurgeSegmentsRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*PurgeSegmentsRequest) RecordingIdId() uint16 {
	return 3
}

func (*PurgeSegmentsRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (p *PurgeSegmentsRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.RecordingIdSinceVersion()
}

func (*PurgeSegmentsRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*PurgeSegmentsRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*PurgeSegmentsRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*PurgeSegmentsRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*PurgeSegmentsRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}

func (*PurgeSegmentsRequest) NewStartPositionId() uint16 {
	return 4
}

func (*PurgeSegmentsRequest) NewStartPositionSinceVersion() uint16 {
	return 0
}

func (p *PurgeSegmentsRequest) NewStartPositionInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.NewStartPositionSinceVersion()
}

func (*PurgeSegmentsRequest) NewStartPositionDeprecated() uint16 {
	return 0
}

func (*PurgeSegmentsRequest) NewStartPositionMetaAttribute(meta int) string {
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

func (*PurgeSegmentsRequest) NewStartPositionMinValue() int64 {
	return math.MinInt64 + 1
}

func (*PurgeSegmentsRequest) NewStartPositionMaxValue() int64 {
	return math.MaxInt64
}

func (*PurgeSegmentsRequest) NewStartPositionNullValue() int64 {
	return math.MinInt64
}
