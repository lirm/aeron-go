// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type PurgeRecordingRequest struct {
	ControlSessionId int64
	CorrelationId    int64
	RecordingId      int64
}

func (p *PurgeRecordingRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	return nil
}

func (p *PurgeRecordingRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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

func (p *PurgeRecordingRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	return nil
}

func PurgeRecordingRequestInit(p *PurgeRecordingRequest) {
	return
}

func (*PurgeRecordingRequest) SbeBlockLength() (blockLength uint16) {
	return 24
}

func (*PurgeRecordingRequest) SbeTemplateId() (templateId uint16) {
	return 104
}

func (*PurgeRecordingRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*PurgeRecordingRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 5
}

func (*PurgeRecordingRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*PurgeRecordingRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*PurgeRecordingRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (p *PurgeRecordingRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.ControlSessionIdSinceVersion()
}

func (*PurgeRecordingRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*PurgeRecordingRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*PurgeRecordingRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*PurgeRecordingRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*PurgeRecordingRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*PurgeRecordingRequest) CorrelationIdId() uint16 {
	return 2
}

func (*PurgeRecordingRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (p *PurgeRecordingRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.CorrelationIdSinceVersion()
}

func (*PurgeRecordingRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*PurgeRecordingRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*PurgeRecordingRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*PurgeRecordingRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*PurgeRecordingRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}

func (*PurgeRecordingRequest) RecordingIdId() uint16 {
	return 3
}

func (*PurgeRecordingRequest) RecordingIdSinceVersion() uint16 {
	return 0
}

func (p *PurgeRecordingRequest) RecordingIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= p.RecordingIdSinceVersion()
}

func (*PurgeRecordingRequest) RecordingIdDeprecated() uint16 {
	return 0
}

func (*PurgeRecordingRequest) RecordingIdMetaAttribute(meta int) string {
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

func (*PurgeRecordingRequest) RecordingIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*PurgeRecordingRequest) RecordingIdMaxValue() int64 {
	return math.MaxInt64
}

func (*PurgeRecordingRequest) RecordingIdNullValue() int64 {
	return math.MinInt64
}
