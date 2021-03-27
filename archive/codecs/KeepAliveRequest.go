// Generated SBE (Simple Binary Encoding) message codec

package codecs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type KeepAliveRequest struct {
	ControlSessionId int64
	CorrelationId    int64
}

func (k *KeepAliveRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := k.RangeCheck(k.SbeSchemaVersion(), k.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := _m.WriteInt64(_w, k.ControlSessionId); err != nil {
		return err
	}
	if err := _m.WriteInt64(_w, k.CorrelationId); err != nil {
		return err
	}
	return nil
}

func (k *KeepAliveRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if !k.ControlSessionIdInActingVersion(actingVersion) {
		k.ControlSessionId = k.ControlSessionIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &k.ControlSessionId); err != nil {
			return err
		}
	}
	if !k.CorrelationIdInActingVersion(actingVersion) {
		k.CorrelationId = k.CorrelationIdNullValue()
	} else {
		if err := _m.ReadInt64(_r, &k.CorrelationId); err != nil {
			return err
		}
	}
	if actingVersion > k.SbeSchemaVersion() && blockLength > k.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-k.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := k.RangeCheck(actingVersion, k.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (k *KeepAliveRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if k.ControlSessionIdInActingVersion(actingVersion) {
		if k.ControlSessionId < k.ControlSessionIdMinValue() || k.ControlSessionId > k.ControlSessionIdMaxValue() {
			return fmt.Errorf("Range check failed on k.ControlSessionId (%v < %v > %v)", k.ControlSessionIdMinValue(), k.ControlSessionId, k.ControlSessionIdMaxValue())
		}
	}
	if k.CorrelationIdInActingVersion(actingVersion) {
		if k.CorrelationId < k.CorrelationIdMinValue() || k.CorrelationId > k.CorrelationIdMaxValue() {
			return fmt.Errorf("Range check failed on k.CorrelationId (%v < %v > %v)", k.CorrelationIdMinValue(), k.CorrelationId, k.CorrelationIdMaxValue())
		}
	}
	return nil
}

func KeepAliveRequestInit(k *KeepAliveRequest) {
	return
}

func (*KeepAliveRequest) SbeBlockLength() (blockLength uint16) {
	return 16
}

func (*KeepAliveRequest) SbeTemplateId() (templateId uint16) {
	return 61
}

func (*KeepAliveRequest) SbeSchemaId() (schemaId uint16) {
	return 101
}

func (*KeepAliveRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 6
}

func (*KeepAliveRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*KeepAliveRequest) ControlSessionIdId() uint16 {
	return 1
}

func (*KeepAliveRequest) ControlSessionIdSinceVersion() uint16 {
	return 0
}

func (k *KeepAliveRequest) ControlSessionIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= k.ControlSessionIdSinceVersion()
}

func (*KeepAliveRequest) ControlSessionIdDeprecated() uint16 {
	return 0
}

func (*KeepAliveRequest) ControlSessionIdMetaAttribute(meta int) string {
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

func (*KeepAliveRequest) ControlSessionIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*KeepAliveRequest) ControlSessionIdMaxValue() int64 {
	return math.MaxInt64
}

func (*KeepAliveRequest) ControlSessionIdNullValue() int64 {
	return math.MinInt64
}

func (*KeepAliveRequest) CorrelationIdId() uint16 {
	return 2
}

func (*KeepAliveRequest) CorrelationIdSinceVersion() uint16 {
	return 0
}

func (k *KeepAliveRequest) CorrelationIdInActingVersion(actingVersion uint16) bool {
	return actingVersion >= k.CorrelationIdSinceVersion()
}

func (*KeepAliveRequest) CorrelationIdDeprecated() uint16 {
	return 0
}

func (*KeepAliveRequest) CorrelationIdMetaAttribute(meta int) string {
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

func (*KeepAliveRequest) CorrelationIdMinValue() int64 {
	return math.MinInt64 + 1
}

func (*KeepAliveRequest) CorrelationIdMaxValue() int64 {
	return math.MaxInt64
}

func (*KeepAliveRequest) CorrelationIdNullValue() int64 {
	return math.MinInt64
}
