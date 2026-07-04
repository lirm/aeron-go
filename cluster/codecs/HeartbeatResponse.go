// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type HeartbeatResponse struct {
	CorrelationId int64
}

func (h *HeartbeatResponse) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := h.RangeCheck(h.SbeSchemaVersion(), h.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, h.CorrelationId); err != nil {
		return err
	}
	return nil
}

func (h *HeartbeatResponse) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !h.CorrelationIdInActingVersion(actingVersion) {
		h.CorrelationId = h.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &h.CorrelationId); err != nil {
			return err
		}
	}
	if actingVersion > h.SbeSchemaVersion() && blockLength > h.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-h.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := h.RangeCheck(actingVersion, h.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (h *HeartbeatResponse) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if h.CorrelationIdInActingVersion(actingVersion) {
		if h.CorrelationId < h.CorrelationIdMinValue() || h.CorrelationId > h.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on h.CorrelationId (%v < %v > %v)", h.CorrelationIdMinValue(), h.CorrelationId, h.CorrelationIdMaxValue())
		}
	}
	return nil
}

func HeartbeatResponseInit(h *HeartbeatResponse) {
	return
}

func (*HeartbeatResponse) SbeBlockLength() (blockLength uint16) {
	return 8
}

func (*HeartbeatResponse) SbeTemplateId() (templateId uint16) {
	return 80
}

func (*HeartbeatResponse) SbeSchemaId() (schemaId uint16) {
	return 111
}

func (*HeartbeatResponse) SbeSchemaVersion() (schemaVersion uint16) {
	return 16
}

func (*HeartbeatResponse) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*HeartbeatResponse) SbeSemanticVersion() (semanticVersion string) {
	return "5.4"
}

func (*HeartbeatResponse) CorrelationIdId() uint16 {
	return 1
}

func (*HeartbeatResponse) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (h *HeartbeatResponse) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= h.CorrelationIdSinceVersion()
}

func (*HeartbeatResponse) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*HeartbeatResponse) CorrelationIdMetaAttribute(meta int) string {
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

func (*HeartbeatResponse) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*HeartbeatResponse) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*HeartbeatResponse) CorrelationIdNullValue() int64 {
	return math.MinInt64
}
